package service

import (
	"errors"

	"gorm.io/gorm"
	"stakeholders.com/dto"
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

func (service *AccountService) FindAccountById(accountId string) (*model.Account, error) {
	account, err := service.AccountRepo.FindById(accountId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (service *AccountService) Login(loginDetails *dto.LoginDetailsDto) string {
	str, _ := service.AccountRepo.Login(loginDetails)

	return str
}

func (s *AccountService) GetUsernameById(accountId string) (string, error) {
	return s.AccountRepo.GetUsernameById(accountId)
}

func (service *AccountService) BlockAccount(admin *model.Account, userId string) error {
	if admin == nil || admin.Role != "admin" {
		return errors.New("forbidden: admin only")
	}
	return service.AccountRepo.BlockAccount(userId)
}

func (service *AccountService) ParseToken(token string) *dto.Claims {
	return service.AccountRepo.ParseToken(token)
}
