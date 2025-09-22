package repo

import (
	"Shopping/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourPurchaseTokenRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *TourPurchaseTokenRepository) FindAllByAccountId(accountId string) ([]model.TourPurchaseToken, error) {
	var purchaseTokens []model.TourPurchaseToken
	dbResult := repo.DatabaseConnection.
		Select("id", "account_id", "tour_id").
		Where("account_id = ?", accountId).
		Find(&purchaseTokens)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return purchaseTokens, nil
}

func (repo *TourPurchaseTokenRepository) Create(purchaseToken *model.TourPurchaseToken) error {
	if purchaseToken.Id == "" {
		purchaseToken.Id = uuid.New().String()
	}

	dbResult := repo.DatabaseConnection.Create(purchaseToken)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	fmt.Println("Rows affected: ", dbResult.RowsAffected)
	fmt.Printf("Created shopping cart: %+v\n", purchaseToken)
	return nil
}
