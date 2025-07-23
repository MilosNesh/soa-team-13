package service

import (
	"errors"

	"gorm.io/gorm"
	"stakeholders.com/model"
	"stakeholders.com/repo"
)

type AccountService struct {
	AccountRepo *repo.AccountRepository
}

func (service *AccountService) FindAll() ([]model.Account, error) {
	accounts, err := service.AccountRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (service *AccountService) Create(account *model.Account) error {
	err := service.AccountRepo.Create(account)

	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) FindAccount(accountId string) (bool, error) {
	err := service.AccountRepo.FindAccount(accountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
