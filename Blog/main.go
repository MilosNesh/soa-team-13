package main

import (
	"blog_project/handler"
	"blog_project/model"
	"blog_project/repo"
	"blog_project/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Uvek pozovi cancel() da oslobodiš resurse

	clientOptions := options.Client().ApplyURI("mongodb://mongo_db:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Greška prilikom konekcije na MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Greška prilikom pingovanja MongoDB-a: %v", err)
	}

	fmt.Println("Uspešno povezan na MongoDB!")

	collection := client.Database("blogdb").Collection("blogs")

	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Upozorenje: Greška prilikom brisanja postojećih blogova (možda ne postoje): %v", err)
	} else {
		fmt.Println("Obrisani svi prethodni blogovi (ako su postojali).")
	}

	blog1 := model.Blog{
		Title:     "Moj prvi blog post",
		Content:   "Ovo je sadrzaj prvog blog posta.",
		CreatedAt: time.Now(),
		ImageUrl:  "https://example.com/image1.jpg",
	}

	blog2 := model.Blog{
		Title:     "Drugi post o GoLangu",
		Content:   "Detaljno objašnjenje GoLang Context paketa.",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		ImageUrl:  "https://example.com/image2.jpg",
	}

	blog3 := model.Blog{
		Title:     "Treći post, bez slike",
		Content:   "Sta je pistac bajo, ker ili zivotinja",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}

	blogsToInsert := []interface{}{blog1, blog2, blog3}

	insertResult, err := collection.InsertMany(ctx, blogsToInsert)
	if err != nil {
		log.Fatalf("Greška prilikom ubacivanja početnih podataka: %v", err)
	}

	fmt.Printf("Uspešno ubaceno %d početnih blog postova. ID-evi: %v\n", len(insertResult.InsertedIDs), insertResult.InsertedIDs)

	return client
}

func startServer(handler *handler.BlogHandler) {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/blogs/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/blogs/", handler.GetAll).Methods("GET")
	router.HandleFunc("/blogs/", handler.Create).Methods("POST")
	router.HandleFunc("/blogs/{blogId}/like", handler.HandleLike).Methods("POST")
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("/uploads"))))

	log.Println("Server started on port :8081...")
	log.Fatal(http.ListenAndServe(":8081", router)) // stakeholders slusa na 8080
}

func main() {
	mongoClient := initDB()

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Greška prilikom zatvaranja MongoDB konekcije: %v", err)
		}
		fmt.Println("MongoDB konekcija zatvorena.")
	}()

	blogsCollection := mongoClient.Database("blogdb").Collection("blogs")

	stakeholderService := &service.StakeholderService{
		BaseURL: "http://stakeholders:8080/",
		Client:  &http.Client{},
	}
	blogRepo := &repo.BlogRepository{Collection: blogsCollection}

	blogService := &service.BlogService{BlogRepo: blogRepo, StakeholderService: stakeholderService}

	blogHandler := &handler.BlogHandler{BlogService: blogService}

	startServer(blogHandler)
}
