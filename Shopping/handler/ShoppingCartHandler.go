package handler

import (
	"Shopping/model"
	"Shopping/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (handler *ShoppingCartHandler) FindByAccountId(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	accountId := vars["accountId"]

	exists, err := handler.ShoppingCartService.FindByAccountId(accountId)

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
	if !exists {
		http.Error(writer, "Shopping cart not found", http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Shopping cart Exists"))
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

	err = handler.ShoppingCartService.Update(&shoppingCart)

	if err != nil {
		println("Error while updating shopping cart")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
