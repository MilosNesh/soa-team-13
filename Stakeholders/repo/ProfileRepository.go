package repo

import (
	"errors"

	"gorm.io/gorm"
	"stakeholders.com/model"
)

type ProfileRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *ProfileRepository) FindByAccountId(accountId string) (*model.Profile, error) {
	var profile model.Profile
	result := repo.DatabaseConnection.First(&profile, "account_id = ?", accountId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

func (repo *ProfileRepository) Save(profile *model.Profile) error {
	return repo.DatabaseConnection.Save(profile).Error
}

func (repo *ProfileRepository) UpdateBalance(accountId string, delta float64) error {
	dbResult := repo.DatabaseConnection.
		Model(&model.Profile{}).
		Where("account_id = ?", accountId).
		Update("balance", gorm.Expr("balance + ?", delta))

	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return errors.New("profile not found")
	}
	return nil
}
