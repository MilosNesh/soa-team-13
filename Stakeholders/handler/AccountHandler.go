package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"stakeholders.com/dto"
	"stakeholders.com/model"
	"stakeholders.com/service"
)

type AccountHandler struct {
	AccountService *service.AccountService
}

func (handler *AccountHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	var role string = req.Header.Get("X-Account-Role")

	if role != "admin" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	account.Id = uuid.New().String()

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

func (handler *AccountHandler) FindAccount(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	accountId := vars["accountId"]

	exists, err := handler.AccountService.FindAccount(accountId)

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "Account not found", http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Account exists"))
}

func (handler *AccountHandler) Login(writer http.ResponseWriter, req *http.Request) {
	var loginDetails dto.LoginDetailsDto
	err := json.NewDecoder(req.Body).Decode(&loginDetails)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	str := handler.AccountService.Login(&loginDetails)
	if str == "Nema" {
		http.Error(writer, "Account not found", http.StatusNotFound)
		return
	}
	if str == "Lozinka" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if str == "Blok" {
		writer.WriteHeader((http.StatusForbidden))
		return
	}

	if str == "Token" {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"token": str})
}

func (handler *AccountHandler) Block(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	accountId := vars["accountId"]

	//
	var role string = req.Header.Get("X-Account-Role")

	if role != "admin" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	admin := &model.Account{Role: role}

	log.Printf("Blockkkk: accountId=%s, role=%q", accountId, role)

	err := handler.AccountService.BlockAccount(admin, accountId)
	if err != nil {
		http.Error(writer, "Only admin can block account", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Account blocked successfully"))
}

func (handler *AccountHandler) ParseToken(writer http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claims *dto.Claims
	claims = handler.AccountService.ParseToken(parts[1])
	if claims == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(map[string]interface{}{
		"username": claims.Username,
		"role":     claims.Role,
		"id":       claims.Subject,
	})
	if err != nil {
		log.Println("Error encoding JSON response:", err)
	}

}
