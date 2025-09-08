package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"followers.com/handler"
	"followers.com/repo"
	"followers.com/service"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func initDB() neo4j.DriverWithContext {
	ctx := context.Background()
	var driver neo4j.DriverWithContext
	var err error

	// Retry loop dok Neo4j ne bude dostupan
	for i := 0; i < 10; i++ { // 10 pokušaja
		driver, err = neo4j.NewDriverWithContext("bolt://neo4j:7687", neo4j.BasicAuth("neo4j", "password123", ""))
		if err != nil {
			log.Println("Neo4j connection error, retrying:", err)
		} else {
			// Test connection
			err = driver.VerifyConnectivity(ctx)
			if err == nil {
				log.Println("Connected to Neo4j")
				break
			}
			log.Println("Neo4j connectivity check failed, retrying:", err)
		}
		// Čekanje pre sledećeg pokušaja
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to Neo4j after retries:", err)
	}

	// Dodavanje sample account-a
	// createSampleAccounts(ctx, driver)

	return driver
}

// func createSampleAccounts(ctx context.Context, driver neo4j.DriverWithContext) {
// 	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
// 	defer session.Close(ctx)

// 	accounts := []struct {
// 		ID       string
// 		Username string
// 	}{
// 		{"7ac9fcb3-5785-4f57-ab5f-bb7f9adff513", "mika"},
// 		{"8a8dbf45-80ee-45cc-a3eb-971f657b7b1f", "zika"},
// 	}

// 	for _, acc := range accounts {
// 		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
// 			query := `
// 				MERGE (a:Account {id: $id})
// 				SET a.username = $username
// 				RETURN a
// 			`
// 			_, err := tx.Run(ctx, query, map[string]any{
// 				"id":       acc.ID,
// 				"username": acc.Username,
// 			})
// 			return nil, err
// 		})
// 		if err != nil {
// 			log.Println("Error creating account:", err)
// 		}
// 	}
// 	log.Println("Sample accounts created in Neo4j")
// }

func main() {
	driver := initDB()
	defer driver.Close(context.TODO())

	stakeholderService := &service.StakeholderService{
		BaseURL: "http://stakeholders:8080/",
		Client:  &http.Client{},
	}

	followRepo := &repo.FollowRepository{Driver: driver}
	followService := &service.FollowService{FollowRepo: followRepo, StakeholderService: stakeholderService}
	followHandler := &handler.FollowHandler{FollowService: followService}

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")

	router.HandleFunc("/follow/{followerId}/{followedId}", followHandler.FollowUser).Methods("POST")

	log.Println("Followers service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
