package handler

import (
	"Shopping/model"
	"Shopping/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ShoppingCartHandler struct {
	ShoppingCartService *service.ShoppingCartService
}

func (handler *ShoppingCartHandler) GetAll(writer http.ResponseWriter, req *http.Request) {

	shoppingCarts, err := handler.ShoppingCartService.FindAll()

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(shoppingCarts)
}

func (handler *ShoppingCartHandler) GetOrCreate(writer http.ResponseWriter, req *http.Request) {
	var accountId string = req.Header.Get("X-Account-Id")
	log.Println("X-Account-Id =", accountId)
	if accountId == "" {
		accountId = req.URL.Query().Get("accountId")
		if accountId == "" {
			http.Error(writer, "missing account id", http.StatusUnauthorized)
			return
		}
	}

	shoppingCart, err := handler.ShoppingCartService.FindOrCreate(accountId)

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(shoppingCart)
}

func (handler *ShoppingCartHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var shoppingCart model.ShoppingCart

	err := json.NewDecoder(req.Body).Decode(&shoppingCart)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	shoppingCart.Id = uuid.New().String()

	err = handler.ShoppingCartService.Create(&shoppingCart)

	if err != nil {
		println("Error while creating a new shopping cart")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *ShoppingCartHandler) Update(writer http.ResponseWriter, req *http.Request) {
	var shoppingCart model.ShoppingCart

	err := json.NewDecoder(req.Body).Decode(&shoppingCart)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	cart, err := handler.ShoppingCartService.Update(&shoppingCart)

	if err != nil {
		println("Error while updating shopping cart")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cart)
}

func (handler *ShoppingCartHandler) AddItem(writer http.ResponseWriter, req *http.Request) {
	var accountId string = req.Header.Get("X-Account-Id")
	log.Println("X-Account-Id =", accountId)
	if accountId == "" {
		accountId = req.URL.Query().Get("accountId")
		if accountId == "" {
			http.Error(writer, "missing account id", http.StatusUnauthorized)
			return
		}
	}

	var orderItem model.OrderItem
	err := json.NewDecoder(req.Body).Decode(&orderItem)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	err = handler.ShoppingCartService.AddItem(accountId, &orderItem)

	if err != nil {
		println("Error while updating shopping cart")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(true)
}

func (handler *ShoppingCartHandler) Checkout(writer http.ResponseWriter, req *http.Request) {
	var accountId string = req.Header.Get("X-Account-Id")
	log.Println("X-Account-Id =", accountId)
	if accountId == "" {
		accountId = req.URL.Query().Get("accountId")
		if accountId == "" {
			http.Error(writer, "missing account id", http.StatusUnauthorized)
			return
		}
	}

	purchaseTokens, err := handler.ShoppingCartService.Checkout(accountId)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(purchaseTokens)
}
