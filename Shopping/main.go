package main

import (
	"Shopping/handler"
	"Shopping/model"
	"Shopping/repo"
	"Shopping/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connectionStr := "host=shopping_database user=postgres password=super dbname=shopping port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.ShoppingCart{}, &model.OrderItem{}, &model.TourPurchaseToken{})

	return database
}

type Server struct {
	shoppingCartService *service.ShoppingCartService
}

func startServer(handler *handler.ShoppingHandler) {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/shopping/", handler.ShoppingCartHandler.GetOrCreate).Methods("GET")
	router.HandleFunc("/shopping/all", handler.ShoppingCartHandler.GetAll).Methods("GET")
	//router.HandleFunc("/shopping/{accountId}", handler.ShoppingCartHandler.FindByAccountId).Methods("GET")
	router.HandleFunc("/shopping/", handler.ShoppingCartHandler.Create).Methods("POST")
	router.HandleFunc("/shopping/", handler.ShoppingCartHandler.Update).Methods("PUT")

	router.HandleFunc("/shopping/orderItems/", handler.ShoppingCartHandler.AddItem).Methods("POST")

	router.HandleFunc("/shopping/checkout/", handler.ShoppingCartHandler.Checkout).Methods("POST")

	router.HandleFunc("/shopping/purchaseTokens/", handler.TourPurchaseTokenHandler.GetAllByAccountId).Methods("GET")
	router.HandleFunc("/shopping/purchaseTokens/", handler.TourPurchaseTokenHandler.Create).Methods("POST")

	println("Server started...")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func main() {
	database := initDB()

	if database == nil {
		print("Failed to connect to DB")
		return
	}

	tourPurchaseTokenRepo := &repo.TourPurchaseTokenRepository{DatabaseConnection: database}
	tourPurchaseTokenService := &service.TourPurchaseTokenService{TourPurchaseTokenRepo: tourPurchaseTokenRepo}
	tourPurchaseTokenHandler := handler.TourPurchaseTokenHandler{TourPurchaseTokenService: tourPurchaseTokenService}

	shoppingCartRepo := &repo.ShoppingCartRepository{DatabaseConnection: database}
	shoppingCartService := &service.ShoppingCartService{
		ShoppingCartRepo:         shoppingCartRepo,
		TourPurchaseTokenService: tourPurchaseTokenService,
	}
	shoppingCartHandler := handler.ShoppingCartHandler{ShoppingCartService: shoppingCartService}

	handler := &handler.ShoppingHandler{
		ShoppingCartHandler:      shoppingCartHandler,
		TourPurchaseTokenHandler: tourPurchaseTokenHandler,
	}

	startServer(handler)
}
