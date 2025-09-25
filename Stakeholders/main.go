package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"stakeholders.com/handler"
	"stakeholders.com/model"
	"stakeholders.com/proto/stakeholders"
	"stakeholders.com/repo"
	"stakeholders.com/service"

	saga "github.com/MilosNesh/soa-team-13/common/saga/messaging"
	natsmsg "github.com/MilosNesh/soa-team-13/common/saga/messaging/nats"
)

func initDB() *gorm.DB {
	connectionStr := "host=database user=postgres password=super dbname=stakeholders port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Account{}, &model.Profile{})

	ddl := `
	CREATE OR REPLACE FUNCTION create_profile_for_account() RETURNS trigger AS $$
	BEGIN
	IF lower(coalesce(NEW.role, '')) = 'admin' THEN
		RETURN NEW;
	END IF;

	INSERT INTO profiles (account_id, name, surname, profile_picture, bio, motto)
	VALUES (NEW.id, '', '', '', '', '')
	ON CONFLICT (account_id) DO NOTHING;

	RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS trg_create_profile ON accounts;

	CREATE TRIGGER trg_create_profile
	AFTER INSERT ON accounts
	FOR EACH ROW
	EXECUTE FUNCTION create_profile_for_account();
	`
	if err := database.Exec(ddl).Error; err != nil {
		panic(err)
	}

	newID := uuid.New()
	newID2 := uuid.New()

	database.Exec("INSERT INTO accounts (id, username, password, email, role) VALUES (?, ?, ?, ?, ?)", newID, "mika", "mika123", "mika@gmail.com", "admin")
	database.Exec("INSERT INTO accounts (id, username, password, email, role) VALUES (?, ?, ?, ?, ?)", newID2, "zika", "zika123", "zika@gmail.com", "admin")

	return database
}

func startGRPCServer(accountService *service.AccountService) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	// Bootstrap gRPC server.
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Bootstrap gRPC service server and respond to request.
	accountHandler := handler.AccountGrpcHandler{Service: accountService}
	stakeholders.RegisterStakeholdersServiceServer(grpcServer, accountHandler)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	grpcServer.Stop()
}

func startServer(handler *handler.StakeholdersHandler) {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/accounts/", handler.AccountHandler.GetAll).Methods("GET")
	router.HandleFunc("/accounts/", handler.AccountHandler.Create).Methods("POST")
	router.HandleFunc("/accounts/doesExists/{accountId}", handler.AccountHandler.FindAccount).Methods("GET")
	router.HandleFunc("/accounts/login", handler.AccountHandler.Login).Methods("POST")
	router.HandleFunc("/accounts/{accountId}/block", handler.AccountHandler.Block).Methods("POST")
	router.HandleFunc("/accounts/parsetoken", handler.AccountHandler.ParseToken).Methods("GET")

	router.HandleFunc("/accounts/profiles/{accountId}", handler.ProfileHandler.FindByAccountId).Methods("GET")
	router.HandleFunc("/accounts/profiles/", handler.ProfileHandler.UpdateProfile).Methods("PUT")

	println("Server started...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	database := initDB()

	if database == nil {
		print("Falied to connect to DB")
		return
	}

	cfg := loadConfig()

	accountRepo := &repo.AccountRepository{DatabaseConnection: database}
	accountService := &service.AccountService{AccountRepo: accountRepo}
	accountHandler := handler.AccountHandler{AccountService: accountService}

	// handler := &handler.StakeholdersHandler{AccountHandler: handler.AccountHandler{AccountService: accountService}}

	profileRepo := &repo.ProfileRepository{DatabaseConnection: database}
	profileService := &service.ProfileService{ProfileRepo: profileRepo}
	profileHandler := handler.ProfileHandler{ProfileService: profileService}

	// profileHandler := &handler.StakeholdersHandler{ProfileHandler: handler.ProfileHandler {profileService: profileService}}

	go startGRPCServer(accountService)

	httpHandlers := &handler.StakeholdersHandler{
		AccountHandler: accountHandler,
		ProfileHandler: profileHandler,
	}

	//SAGA
	replyPublisher := mustPublisher(cfg, cfg.ReplySubject)
	commandSubscriber := mustSubscriber(cfg, cfg.CommandSubject, cfg.QueueGroup)

	_, err := handler.NewPurchaseStakeholdersHandler(profileService, replyPublisher, commandSubscriber)
	if err != nil {
		log.Fatal(err)
	}

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
		QueueGroup:     getenv("NATS_QUEUE_GROUP", "stakeholders_service"),
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
