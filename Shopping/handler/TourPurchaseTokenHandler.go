package handler

import (
	"Shopping/model"
	"Shopping/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type TourPurchaseTokenHandler struct {
	TourPurchaseTokenService *service.TourPurchaseTokenService
}

func (handler *TourPurchaseTokenHandler) GetAllByAccountId(writer http.ResponseWriter, req *http.Request) {
	var accountId string = req.Header.Get("X-Account-Id")
	log.Println("X-Account-Id =", accountId)
	if accountId == "" {
		accountId = req.URL.Query().Get("accountId")
		if accountId == "" {
			http.Error(writer, "missing account id", http.StatusUnauthorized)
			return
		}
	}

	purchaseTokens, err := handler.TourPurchaseTokenService.FindAllByAccountId(accountId)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(purchaseTokens)
}

func (handler *TourPurchaseTokenHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var purchaseToken model.TourPurchaseToken

	err := json.NewDecoder(req.Body).Decode(&purchaseToken)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	purchaseToken.Id = uuid.New().String()

	created, err := handler.TourPurchaseTokenService.Create(&purchaseToken)
	if err != nil {
		println("Error while creating a new purchase token")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(created)
}
