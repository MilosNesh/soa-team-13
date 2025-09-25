package main

import (
	"Shopping/handler"
	"Shopping/model"
	"Shopping/repo"
	"Shopping/service"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	saga "github.com/MilosNesh/soa-team-13/common/saga/messaging"
	natsmsg "github.com/MilosNesh/soa-team-13/common/saga/messaging/nats"
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
	router.HandleFunc("/shopping/checkout-saga/", handler.ShoppingCartHandler.CheckoutWithSaga).Methods("POST")

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

	cfg := loadConfig()

	//Repozitorijumi i servisi

	tourPurchaseTokenRepo := &repo.TourPurchaseTokenRepository{DatabaseConnection: database}
	tourPurchaseTokenService := &service.TourPurchaseTokenService{TourPurchaseTokenRepo: tourPurchaseTokenRepo}

	shoppingCartRepo := &repo.ShoppingCartRepository{DatabaseConnection: database}
	shoppingCartService := &service.ShoppingCartService{
		ShoppingCartRepo:         shoppingCartRepo,
		TourPurchaseTokenService: tourPurchaseTokenService,
	}

	//HTTP handleri
	tourPurchaseTokenHandler := handler.TourPurchaseTokenHandler{TourPurchaseTokenService: tourPurchaseTokenService}
	shoppingCartHandler := handler.ShoppingCartHandler{ShoppingCartService: shoppingCartService}
	httpHandlers := &handler.ShoppingHandler{
		ShoppingCartHandler:      shoppingCartHandler,
		TourPurchaseTokenHandler: tourPurchaseTokenHandler,
	}

	//SAGA
	handlerPublisher := mustPublisher(cfg, cfg.ReplySubject)
	handlerSubscriber := mustSubscriber(cfg, cfg.CommandSubject, cfg.QueueGroup)

	// Create orchestrator for starting SAGA
	orchestratorPublisher := mustPublisher(cfg, cfg.CommandSubject)
	orchestratorSubscriber := mustSubscriber(cfg, cfg.ReplySubject, "orchestrator")

	orchestrator, err := service.NewPurchaseCartOrchestrator(orchestratorPublisher, orchestratorSubscriber)
	if err != nil {
		log.Fatalf("Failed to init PurchaseCartOrchestrator: %v", err)
	}

	_, err = handler.NewPurchaseCartHandler(shoppingCartService, handlerPublisher, handlerSubscriber)
	if err != nil {
		log.Fatalf("Failed to init PurchaseCartHandler: %v", err)
	}

	// Store orchestrator for use in checkout
	shoppingCartService.SetOrchestrator(orchestrator)

	startServer(httpHandlers)
}

type config struct {
	NatsHost string
	NatsPort string
	NatsUser string
	NatsPass string

	CommandSubject string
	ReplySubject   string
	QueueGroup     string
}

func loadConfig() config {
	c := config{
		NatsHost:       getenv("NATS_HOST", "nats-server"),
		NatsPort:       getenv("NATS_PORT", "4222"),
		NatsUser:       getenv("NATS_USER", ""),
		NatsPass:       getenv("NATS_PASS", ""),
		CommandSubject: getenv("PURCHASE_COMMAND_SUBJECT", "purchase.checkout.command"),
		ReplySubject:   getenv("PURCHASE_REPLY_SUBJECT", "purchase.checkout.reply"),
		QueueGroup:     getenv("NATS_QUEUE_GROUP", "shopping_service"),
	}
	return c
}

func getenv(k, def string) string {
	if v := getenvReal(k); v != "" {
		return v
	}
	return def
}

// zameni implementacijom za svoj runtime (os.Getenv)
func getenvReal(k string) string { return os.Getenv(k) }

func mustPublisher(cfg config, subject string) saga.Publisher {
	pub, err := natsmsg.NewNATSPublisher(cfg.NatsHost, cfg.NatsPort, cfg.NatsUser, cfg.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return pub
}

func mustSubscriber(cfg config, subject, queueGroup string) saga.Subscriber {
	sub, err := natsmsg.NewNATSSubscriber(cfg.NatsHost, cfg.NatsPort, cfg.NatsUser, cfg.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return sub
}
