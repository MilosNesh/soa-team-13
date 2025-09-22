package service

import (
	"Shopping/model"
	"Shopping/repo"
	"errors"

	"gorm.io/gorm"
)

type ShoppingCartService struct {
	ShoppingCartRepo *repo.ShoppingCartRepository
}

func (service *ShoppingCartService) FindAll() ([]model.ShoppingCart, error) {
	shoppingCarts, err := service.ShoppingCartRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return shoppingCarts, nil
}

func (service *ShoppingCartService) FindByAccountId(accountId string) (bool, error) {
	err := service.ShoppingCartRepo.FindByAccountId(accountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *ShoppingCartService) Create(shoppingCart *model.ShoppingCart) error {
	err := service.ShoppingCartRepo.Create(shoppingCart)

	if err != nil {
		return err
	}
	return nil
}

func (service *ShoppingCartService) Update(shoppingCart *model.ShoppingCart) error {
	err := service.ShoppingCartRepo.Update(shoppingCart)

	if err != nil {
		return err
	}
	return nil
}
