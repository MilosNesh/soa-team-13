package handler

import (
	"blog_project/dto"
	"blog_project/model"
	"blog_project/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UploadDir = "/uploads"

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
			writer.WriteHeader(http.StatusNotFound)
		} else if err == context.DeadlineExceeded {
			log.Printf("Pretraga bloga sa ID '%s' je predugo trajala (timeout).", id)
			writer.WriteHeader(http.StatusGatewayTimeout)
		} else {
			log.Printf("Greška prilikom dohvatanja bloga sa ID '%s': %v", id, err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}

func (handler *BlogHandler) GetAll(writer http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blogs, err := handler.BlogService.FindAllBlogs(ctx)

	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Printf("Greška prilikom dohvatanja svih blogova: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (handler *BlogHandler) Create(writer http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("X-Account-Id")
	if userID == "" {
		http.Error(writer, "Autorizacija neuspešna: User ID (X-Account-Id) nije pronađen.", http.StatusUnauthorized)
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Greška prilikom parsiranja forme: %v", err)
		http.Error(writer, "Greška prilikom parsiranja forme: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := req.PostFormValue("title")
	content := req.PostFormValue("content")

	if title == "" || content == "" {
		http.Error(writer, "Naslov i sadržaj su obavezna polja.", http.StatusBadRequest)
		return
	}

	newBlog := model.Blog{
		Title:     title,
		Content:   content,
		UserId:    userID,
		CreatedAt: time.Now(),
	}

	file, fileHeader, err := req.FormFile("image")

	if err == nil {
		defer file.Close()

		ext := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "blog_img", ext)
		filePath := filepath.Join(UploadDir, fileName)

		if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
			os.Mkdir(UploadDir, os.ModePerm)
		}

		// Čuvanje fajla na disku
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("Greška pri kreiranju fajla na serveru: %v", err)
			http.Error(writer, "Greška prilikom čuvanja fajla: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Greška pri kopiranju fajla: %v", err)
			http.Error(writer, "Greška prilikom kopiranja fajla: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newBlog.ImageUrl = fileName
	} else if err != http.ErrMissingFile {
		log.Printf("Greška pri obradi fajla: %v", err)
		http.Error(writer, "Greška pri obradi fajla: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = handler.BlogService.Create(ctx, &newBlog)
	if err != nil {
		log.Printf("Greška prilikom kreiranja novog bloga: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(newBlog)
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

func (handler *BlogHandler) AddComment(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	blogIdStr := vars["blogId"]

	authorId := req.Header.Get("X-Account-Id")
	if authorId == "" {
		http.Error(writer, "Autorizacija neuspešna: User ID (X-Account-Id) nije pronađen.", http.StatusUnauthorized)
		return
	}

	blogId, err := primitive.ObjectIDFromHex(blogIdStr)
	if err != nil {
		http.Error(writer, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var dto dto.AddCommentDto
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment := model.Comment{
		BlogId:   blogId,
		AuthorId: authorId,
		Content:  dto.Content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = handler.BlogService.AddComment(ctx, &comment)

	if err != nil {
		if err.Error() == "blog not found" || err.Error() == "author not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(comment)
}
