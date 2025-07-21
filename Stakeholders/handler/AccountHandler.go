package handler

import (
	"encoding/json"
	"net/http"

	"stakeholders.com/model"
	"stakeholders.com/service"
)

type AccountHandler struct {
	AccountService *service.AccountService
}

func (handler *AccountHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	accounts, err := handler.AccountService.FindAll()

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(accounts)
}

func (handler *AccountHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var account model.Account

	err := json.NewDecoder(req.Body).Decode(&account)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !account.IsValid() {
		println("Account is not valid")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if account.Role != "guide" && account.Role != "tourist" {
		println("Account role is not valid")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.AccountService.Create(&account)

	if err != nil {
		println("Error while creating a new account")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}
