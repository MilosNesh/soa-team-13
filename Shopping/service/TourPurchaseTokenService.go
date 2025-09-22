package service

import (
	"Shopping/model"
	"Shopping/repo"
)

type TourPurchaseTokenService struct {
	TourPurchaseTokenRepo *repo.TourPurchaseTokenRepository
}

func (service *TourPurchaseTokenService) FindAllByAccountId(accountId string) ([]model.TourPurchaseToken, error) {
	purchaseTokens, err := service.TourPurchaseTokenRepo.FindAllByAccountId(accountId)

	if err != nil {
		return nil, err
	}
	return purchaseTokens, nil
}

func (service *TourPurchaseTokenService) Create(purchaseToken *model.TourPurchaseToken) error {
	err := service.TourPurchaseTokenRepo.Create(purchaseToken)

	if err != nil {
		return err
	}
	return nil
}
