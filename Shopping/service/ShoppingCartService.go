package service

import (
	"Shopping/model"
	"Shopping/repo"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShoppingCartService struct {
	ShoppingCartRepo         *repo.ShoppingCartRepository
	TourPurchaseTokenService *TourPurchaseTokenService
}

func (service *ShoppingCartService) FindAll() ([]model.ShoppingCart, error) {
	shoppingCarts, err := service.ShoppingCartRepo.FindAll()

	if err != nil {
		return nil, err
	}
	return shoppingCarts, nil
}

func (service *ShoppingCartService) FindOrCreate(accountId string) (*model.ShoppingCart, error) {
	shoppingCart, err := service.ShoppingCartRepo.FindOrCreate(accountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return shoppingCart, nil
}

func (service *ShoppingCartService) Create(shoppingCart *model.ShoppingCart) error {
	err := service.ShoppingCartRepo.Create(shoppingCart)

	if err != nil {
		return err
	}
	return nil
}

func (service *ShoppingCartService) Update(shoppingCart *model.ShoppingCart) (*model.ShoppingCart, error) {
	shoppingCart, err := service.ShoppingCartRepo.Update(shoppingCart)

	if err != nil {
		return nil, err
	}
	return shoppingCart, nil
}

func (service *ShoppingCartService) AddItem(accountId string, orderItem *model.OrderItem) error {

	shoppingCart, err := service.ShoppingCartRepo.FindOrCreate(accountId)
	if err != nil {
		return err
	}

	orderItem.ShoppingCartId = shoppingCart.Id

	err = service.ShoppingCartRepo.AddItem(orderItem)

	if err != nil {
		return err
	}
	return nil
}

func (service *ShoppingCartService) Checkout(accountId string) ([]model.TourPurchaseToken, error) {
	shoppingCart, err := service.ShoppingCartRepo.FindOrCreate(accountId)
	if err != nil {
		return nil, err
	}

	if len(shoppingCart.Items) == 0 {
		return nil, fmt.Errorf("korpa je prazna")
	}

	var purchaseTokens []model.TourPurchaseToken

	for _, item := range shoppingCart.Items {

		token := model.TourPurchaseToken{
			Id:        uuid.NewString(),
			AccountID: accountId,
			TourId:    item.TourId,
		}

		if _, err := service.TourPurchaseTokenService.Create(&token); err != nil {
			return nil, err
		}

		purchaseTokens = append(purchaseTokens, token)

		filtered := make([]model.OrderItem, 0, len(shoppingCart.Items))
		for _, it := range shoppingCart.Items {
			if it.Id != item.Id {
				filtered = append(filtered, it)
			}
		}
		shoppingCart.Items = filtered
	}

	if _, err := service.ShoppingCartRepo.Update(shoppingCart); err != nil {
		return nil, err
	}

	return purchaseTokens, nil
}
