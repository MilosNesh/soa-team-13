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

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const UploadDir = "/uploads"

type BlogHandler struct {
	BlogService *service.BlogService
	Tracer      *sdktrace.TracerProvider
	ServiceName string
}

func (handler *BlogHandler) Get(writer http.ResponseWriter, req *http.Request) {
	traceContext, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(req.Context(), "blog-get-by-id")
	defer func() { span.End() }()

	id := mux.Vars(req)["id"]
	span.AddEvent(fmt.Sprintf("Fetching blog with ID: %s", id))
	log.Printf("Blog sa id-em %s", id)

	ctx, cancel := context.WithTimeout(traceContext, 5*time.Second)
	defer cancel()

	blog, err := handler.BlogService.FindBlogById(ctx, id)

	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		span.RecordError(err, trace.WithAttributes())
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
	traceContext, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(req.Context(), "blog-get-all")
	defer func() { span.End() }()

	requesterId := req.Header.Get("X-Account-Id")

	if requesterId == "" {
		http.Error(writer, "Unauthorized: X-Account-Id required", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(traceContext, 5*time.Second)
	defer cancel()

	blogs, err := handler.BlogService.FindAllBlogs(ctx, requesterId)

	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		span.RecordError(err, trace.WithAttributes())
		log.Printf("Greška prilikom dohvatanja svih blogova: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (handler *BlogHandler) Create(writer http.ResponseWriter, req *http.Request) {
	traceContext, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(req.Context(), "blog-create")
	defer func() { span.End() }()

	userID := req.Header.Get("X-Account-Id")
	if userID == "" {
		http.Error(writer, "Autorizacija neuspešna: User ID (X-Account-Id) nije pronađen.", http.StatusUnauthorized)
		return
	}

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		span.RecordError(err, trace.WithAttributes())
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

		dst, err := os.Create(filePath)
		if err != nil {
			span.RecordError(err, trace.WithAttributes())
			log.Printf("Greška pri kreiranju fajla na serveru: %v", err)
			http.Error(writer, "Greška prilikom čuvanja fajla: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			span.RecordError(err, trace.WithAttributes())
			log.Printf("Greška pri kopiranju fajla: %v", err)
			http.Error(writer, "Greška prilikom kopiranja fajla: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newBlog.ImageUrl = fileName
	} else if err != http.ErrMissingFile {
		span.RecordError(err, trace.WithAttributes())
		log.Printf("Greška pri obradi fajla: %v", err)
		http.Error(writer, "Greška pri obradi fajla: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(traceContext, 5*time.Second)
	defer cancel()

	err = handler.BlogService.Create(ctx, &newBlog)
	if err != nil {
		span.RecordError(err, trace.WithAttributes())
		log.Printf("Greška prilikom kreiranja novog bloga: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(newBlog)
}

func (handler *BlogHandler) HandleLike(writer http.ResponseWriter, req *http.Request) {
	traceContext, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(req.Context(), "blog-handle-like")
	defer func() { span.End() }()

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

	liked, err := handler.BlogService.HandleLike(traceContext, blogId, request.AccountId)
	if err != nil {
		span.RecordError(err, trace.WithAttributes())
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
	log.Println("DEBUG BlogHandler: === Entering AddComment handler ===")
	log.Printf("DEBUG BlogHandler: Request Method: %s, URL: %s", req.Method, req.URL.Path)

	traceContext, span := (trace.TracerProvider)(handler.Tracer).Tracer(handler.ServiceName).Start(req.Context(), "blog-add-comment")
	defer func() {
		log.Println("DEBUG BlogHandler: === Exiting AddComment handler ===")
		span.End()
	}()

	vars := mux.Vars(req)
	blogIdStr := vars["blogId"]
	log.Printf("DEBUG BlogHandler: Extracted blogIdStr: %s", blogIdStr)

	authorId := req.Header.Get("X-Account-Id")
	log.Printf("DEBUG BlogHandler: Extracted X-Account-Id: %s", authorId)
	if authorId == "" {
		log.Println("ERROR BlogHandler: X-Account-Id header is empty. Authorization failed.")
		http.Error(writer, "Autorizacija neuspešna: User ID (X-Account-Id) nije pronađen.", http.StatusUnauthorized)
		return
	}

	blogId, err := primitive.ObjectIDFromHex(blogIdStr)
	if err != nil {
		log.Printf("ERROR BlogHandler: Invalid blog ID string '%s': %v", blogIdStr, err)
		http.Error(writer, "Invalid blog ID", http.StatusBadRequest)
		return
	}
	log.Printf("DEBUG BlogHandler: Converted blogId to ObjectID: %s", blogId.Hex())

	var dto dto.AddCommentDto
	// Pročitaj telo zahteva
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		log.Printf("ERROR BlogHandler: Failed to decode request body: %v", err)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("DEBUG BlogHandler: Decoded comment content: %s", dto.Content)

	comment := model.Comment{
		BlogId:   blogId,
		AuthorId: authorId,
		Content:  dto.Content,
	}
	log.Printf("DEBUG BlogHandler: Created comment model: %+v", comment)

	ctx, cancel := context.WithTimeout(traceContext, 5*time.Second)
	defer cancel()

	log.Println("DEBUG BlogHandler: Calling BlogService.AddComment...")
	err = handler.BlogService.AddComment(ctx, &comment)

	if err != nil {
		span.RecordError(err, trace.WithAttributes())
		log.Printf("ERROR BlogHandler: Error from BlogService.AddComment: %v", err)
		if err.Error() == "blog not found" || err.Error() == "author not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("DEBUG BlogHandler: Comment successfully added.")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(comment)
}
