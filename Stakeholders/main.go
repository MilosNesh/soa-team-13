package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"stakeholders.com/handler"
	"stakeholders.com/model"
	"stakeholders.com/repo"
	"stakeholders.com/service"
)

func initDB() *gorm.DB {
	connectionStr := "host=database user=postgres password=super dbname=stakeholders port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Account{})
	database.Exec("INSERT INTO accounts VALUES ('mika', 'mika123', 'mika@gmail.com', 'admin')")
	database.Exec("INSERT INTO accounts VALUES ('zika', 'zika123', 'zika@gmail.com', 'admin')")
	return database
}

func startServer(handler *handler.StakeholdersHandler) {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/accounts/", handler.AccountHandler.GetAll).Methods("GET")
	router.HandleFunc("/accounts/", handler.AccountHandler.Create).Methods("POST")

	println("Server started...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	database := initDB()

	if database == nil {
		print("Falied to connect to DB")
		return
	}

	accountRepo := &repo.AccountRepository{DatabaseConnection: database}
	accountService := &service.AccountService{AccountRepo: accountRepo}
	handler := &handler.StakeholdersHandler{AccountHandler: handler.AccountHandler{AccountService: accountService}}

	startServer(handler)
}
