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

func (repo *ShoppingCartRepository) FindOrCreate(accountId string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart

	err := repo.DatabaseConnection.
		Where(&model.ShoppingCart{AccountID: accountId}).
		Attrs(&model.ShoppingCart{Id: uuid.NewString()}).
		FirstOrCreate(&cart).Error
	if err != nil {
		return nil, err
	}

	if err := repo.DatabaseConnection.
		Preload("Items").
		First(&cart, "id = ?", cart.Id).Error; err != nil {
		return nil, err
	}
	return &cart, nil
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

func (repo *ShoppingCartRepository) Update(cart *model.ShoppingCart) (*model.ShoppingCart, error) {
	if cart == nil || cart.Id == "" {
		return nil, fmt.Errorf("shoppingCart.Id is empty")
	}

	var out *model.ShoppingCart

	err := repo.DatabaseConnection.Transaction(func(tx *gorm.DB) error {
		// 1) Učitaj postojeći
		var existing model.ShoppingCart
		if err := tx.Preload("Items").
			First(&existing, "id = ?", cart.Id).Error; err != nil {
			return err
		}

		// 2) Update polja na korpi (po potrebi)
		if err := tx.Model(&existing).
			Updates(map[string]any{
				"account_id": cart.AccountID,
			}).Error; err != nil {
			return err
		}

		// 3) Upsert stavki (uvek postavi FK); skupljamo ID-eve koje zadržavamo
		keep := make([]string, 0, len(cart.Items))
		for i := range cart.Items {
			it := &cart.Items[i]
			if it.Id == "" {
				it.Id = uuid.NewString()
			}
			it.ShoppingCartId = existing.Id
			keep = append(keep, it.Id)

			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "price", "tour_id", "shopping_cart_id"}),
			}).Create(it).Error; err != nil {
				return err
			}
		}

		// 4) Obriši one koje više nisu u listi (umesto NULL FK!)
		if len(keep) == 0 {
			if err := tx.Where("shopping_cart_id = ?", existing.Id).
				Delete(&model.OrderItem{}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Where("shopping_cart_id = ? AND id NOT IN ?", existing.Id, keep).
				Delete(&model.OrderItem{}).Error; err != nil {
				return err
			}
		}

		// 5) (Opcionalno) Izračunaj total u bazi ili u kodu
		//    Ako imaš polje TotalPrice u modelu:
		// total := int64(0)
		// for _, it := range cart.Items { total += it.Price }
		// if err := tx.Model(&existing).Update("total_price", total).Error; err != nil { return err }

		// 6) Vrati svež objekat
		var fresh model.ShoppingCart
		if err := tx.Preload("Items").
			First(&fresh, "id = ?", existing.Id).Error; err != nil {
			return err
		}
		out = &fresh
		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}

func (repo *ShoppingCartRepository) AddItem(orderItem *model.OrderItem) error {
	if orderItem.Id == "" {
		orderItem.Id = uuid.New().String()
	}
	if orderItem.ShoppingCartId == "" {
		return fmt.Errorf("missing ShoppingCartId")
	}
	if orderItem.TourId == 0 {
		return fmt.Errorf("missing TourId")
	}

	dbResult := repo.DatabaseConnection.Create(orderItem)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	fmt.Println("Rows affected: ", dbResult.RowsAffected)
	fmt.Printf("Created shopping cart: %+v\n", orderItem)
	return nil
}

func (repo *ShoppingCartRepository) ClearCart(cartID string) error {

	dbResult := repo.DatabaseConnection.
		Where("shopping_cart_id = ?", cartID).
		Delete(&model.OrderItem{})

	if dbResult.Error != nil {
		return dbResult.Error
	}

	fmt.Println("Rows affected: ", dbResult.RowsAffected)
	fmt.Printf("Deleted items from shopping cart")
	return nil
}
