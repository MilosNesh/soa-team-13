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

	for i := 0; i < 10; i++ {
		driver, err = neo4j.NewDriverWithContext("bolt://neo4j:7687", neo4j.BasicAuth("neo4j", "password123", ""))
		if err != nil {
			log.Println("Neo4j connection error, retrying:", err)
		} else {
			err = driver.VerifyConnectivity(ctx)
			if err == nil {
				log.Println("Connected to Neo4j")
				break
			}
			log.Println("Neo4j connectivity check failed, retrying:", err)
		}
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to Neo4j after retries:", err)
	}
	return driver
}

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
