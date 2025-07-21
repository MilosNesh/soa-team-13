package repo

import (
	"gorm.io/gorm"
	"stakeholders.com/model"
)

type AccountRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *AccountRepository) FindAll() ([]model.Account, error) {
	var accounts []model.Account
	dbResult := repo.DatabaseConnection.Select("username", "email", "role").Find(&accounts)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return accounts, nil
}

func (repo *AccountRepository) Create(account *model.Account) error {
	dbResult := repo.DatabaseConnection.Create(account)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
