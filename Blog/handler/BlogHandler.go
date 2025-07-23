package handler

import (
	"blog_project/dto"
	"blog_project/model"
	"blog_project/service"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogHandler struct {
	BlogService *service.BlogService
}

func (handler *BlogHandler) Get(writer http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	log.Printf("Blog sa id-em %s", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blog, err := handler.BlogService.FindBlogById(ctx, id)

	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Blog sa ID '%s' nije pronađen.", id)
			writer.WriteHeader(http.StatusNotFound) // Vrati 404 ako blog ne postoji
		} else if err == context.DeadlineExceeded {
			log.Printf("Pretraga bloga sa ID '%s' je predugo trajala (timeout).", id)
			writer.WriteHeader(http.StatusGatewayTimeout) // Vrati 504 za timeout
		} else {
			log.Printf("Greška prilikom dohvatanja bloga sa ID '%s': %v", id, err)
			writer.WriteHeader(http.StatusInternalServerError) // Za ostale greške
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}

func (handler *BlogHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var blog model.Blog
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = handler.BlogService.Create(ctx, &blog)
	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(blog)
}

func (handler *BlogHandler) HandleLike(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	blogIdStr := vars["blogId"]

	blogId, err := primitive.ObjectIDFromHex(blogIdStr)
	if err != nil {
		http.Error(writer, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var request dto.LikeDto
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil || request.AccountId == "" {
		http.Error(writer, "Invalid request body or missing username", http.StatusBadRequest)
		return
	}

	liked, err := handler.BlogService.HandleLike(req.Context(), blogId, request.AccountId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	if liked {
		writer.Write([]byte("like added"))
	} else {
		writer.Write([]byte("like removed"))
	}
}
