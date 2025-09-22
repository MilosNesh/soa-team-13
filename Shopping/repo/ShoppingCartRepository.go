package repo

import (
	"Shopping/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShoppingCartRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *ShoppingCartRepository) FindAll() ([]model.ShoppingCart, error) {
	var shoppingCarts []model.ShoppingCart
	dbResult := repo.DatabaseConnection.
		Model(&model.ShoppingCart{}).
		Select("shopping_carts.id", "shopping_carts.account_id").
		Preload("Items").
		Find(&shoppingCarts)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return shoppingCarts, nil
}

func (repo *ShoppingCartRepository) FindByAccountId(accountId string) error {
	var shoppingCart model.ShoppingCart
	dbResult := repo.DatabaseConnection.First(&shoppingCart, "id = ?", accountId)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *ShoppingCartRepository) Create(shoppingCart *model.ShoppingCart) error {
	if shoppingCart.Id == "" {
		shoppingCart.Id = uuid.New().String()
	}

	for i := range shoppingCart.Items {
		if shoppingCart.Items[i].Id == "" {
			shoppingCart.Items[i].Id = uuid.NewString()
		}
		shoppingCart.Items[i].ShoppingCartId = shoppingCart.Id
	}

	dbResult := repo.DatabaseConnection.Create(shoppingCart)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	fmt.Println("Rows affected: ", dbResult.RowsAffected)
	fmt.Printf("Created shopping cart: %+v\n", shoppingCart)
	return nil
}

func (repo *ShoppingCartRepository) Update(shoppingCart *model.ShoppingCart) error {
	if shoppingCart.Id == "" {
		return fmt.Errorf("shoppingCart.Id is empty")
	}

	return repo.DatabaseConnection.Transaction(func(tx *gorm.DB) error {

		var existing model.ShoppingCart
		if err := tx.Preload("Items").
			First(&existing, "id = ?", shoppingCart.Id).Error; err != nil {
			return err
		}

		if err := tx.Model(&existing).
			Updates(map[string]any{
				"account_id": shoppingCart.AccountID,
			}).Error; err != nil {
			return err
		}

		for i := range shoppingCart.Items {
			if shoppingCart.Items[i].Id == "" {
				shoppingCart.Items[i].Id = uuid.NewString()
			}
			shoppingCart.Items[i].ShoppingCartId = existing.Id
		}

		if err := tx.Model(&existing).
			Association("Items").
			Replace(shoppingCart.Items); err != nil {
			return err
		}

		if len(shoppingCart.Items) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&shoppingCart.Items).Error; err != nil {
				return err
			}
		}

		ids := make([]string, 0, len(shoppingCart.Items))
		for _, it := range shoppingCart.Items {
			ids = append(ids, it.Id)
		}
		q := tx.Where("shopping_cart_id = ?", existing.Id)
		if len(ids) > 0 {
			q = q.Where("id NOT IN ?", ids)
		}
		if err := q.Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}
