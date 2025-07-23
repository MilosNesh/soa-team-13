package repo

import (
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
