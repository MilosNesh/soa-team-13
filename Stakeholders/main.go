package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
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

	accountRepo := &repo.AccountRepository{DatabaseConnection: database}
	accountService := &service.AccountService{AccountRepo: accountRepo}
	accountHandler := handler.AccountHandler{AccountService: accountService}

	// handler := &handler.StakeholdersHandler{AccountHandler: handler.AccountHandler{AccountService: accountService}}

	profileRepo := &repo.ProfileRepository{DatabaseConnection: database}
	profileService := &service.ProfileService{ProfileRepo: profileRepo}
	profileHandler := handler.ProfileHandler{ProfileService: profileService}

	// profileHandler := &handler.StakeholdersHandler{ProfileHandler: handler.ProfileHandler {profileService: profileService}}

	handler := &handler.StakeholdersHandler{
		AccountHandler: accountHandler,
		ProfileHandler: profileHandler,
	}

	startServer(handler)
}
