package handler

import (
	"blog_project/model"
	"blog_project/proto/blog"
	"blog_project/service"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogGrpcHandler struct {
	blog.UnimplementedBlogServiceServer
	BlogService *service.BlogService
	Tracer      *trace.TracerProvider
	ServiceName string
}

func (h *BlogGrpcHandler) GetBlog(ctx context.Context, req *blog.GetBlogRequest) (*blog.GetBlogResponse, error) {
	_, span := (oteltrace.TracerProvider)(h.Tracer).Tracer(h.ServiceName).Start(ctx, "grpc-get-blog")
	defer span.End()

	log.Printf("GetBlog gRPC called for ID: %s", req.Id)

	blogData, err := h.BlogService.FindBlogById(ctx, req.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.RecordError(err, oteltrace.WithAttributes())
			return nil, status.Errorf(codes.NotFound, "Blog not found")
		}
		span.RecordError(err, oteltrace.WithAttributes())
		return nil, status.Errorf(codes.Internal, "Error fetching blog: %v", err)
	}

	var protoComments []*blog.Comment
	for _, c := range blogData.Comments {
		protoComments = append(protoComments, &blog.Comment{
			Id:         c.ID.Hex(),
			BlogId:     c.BlogId.Hex(),
			AuthorId:   c.AuthorId,
			AuthorName: c.AuthorName,
			Content:    c.Content,
			CreatedAt:  c.CreatedAt.Format(time.RFC3339),
		})
	}

	protoBlog := &blog.Blog{
		Id:        blogData.ID.Hex(),
		Title:     blogData.Title,
		Content:   blogData.Content,
		UserId:    blogData.UserId,
		ImageUrl:  blogData.ImageUrl,
		CreatedAt: blogData.CreatedAt.Format(time.RFC3339),
		Likes:     blogData.Likes,
		Comments:  protoComments,
	}

	return &blog.GetBlogResponse{Blog: protoBlog}, nil
}

func (h *BlogGrpcHandler) CreateBlog(ctx context.Context, req *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	_, span := (oteltrace.TracerProvider)(h.Tracer).Tracer(h.ServiceName).Start(ctx, "grpc-create-blog")
	defer span.End()

	log.Printf("CreateBlog gRPC called for blog with title: %s by user: %s", req.Title, req.UserId)

	if req.UserId == "" {
		span.RecordError(fmt.Errorf("User ID is empty"), oteltrace.WithAttributes())
		return nil, status.Errorf(codes.Unauthenticated, "Autorizacija neuspešna: User ID nije pronađen u zahtevu.")
	}

	if req.Title == "" || req.Content == "" {
		span.RecordError(fmt.Errorf("Title or content is empty"), oteltrace.WithAttributes())
		return nil, status.Errorf(codes.InvalidArgument, "Naslov i sadržaj su obavezna polja.")
	}

	newBlog := &model.Blog{
		Title:     req.Title,
		Content:   req.Content,
		UserId:    req.UserId,
		CreatedAt: time.Now(),
	}

	if len(req.GetImage()) > 0 {
		log.Printf("Received image data length: %d bytes", len(req.GetImage()))
		if len(req.GetImage()) > 0 {
			log.Printf("Received image data first byte (hex): %x", req.GetImage()[0])
			loggedImagePreview := string(req.GetImage())
			log.Printf("Received image data as string (first 100 chars): %s", loggedImagePreview[:min(len(loggedImagePreview), 100)])
			log.Printf("Received image data as string (last 100 chars): %s", loggedImagePreview[max(0, len(loggedImagePreview)-100):])
			log.Printf("Received image data as string length: %d", len(loggedImagePreview))
		}

		imageBytes := req.GetImage()

		ext := filepath.Ext(req.GetImageFilename())
		fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "blog_img", ext)
		filePath := filepath.Join(UploadDir, fileName)

		if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
			os.Mkdir(UploadDir, os.ModePerm)
		}

		err := os.WriteFile(filePath, imageBytes, 0644)
		if err != nil {
			log.Printf("Greška prilikom čuvanja fajla: %v", err)
			return nil, fmt.Errorf("greška prilikom čuvanja fajla: %w", err)
		}

		newBlog.ImageUrl = fileName
	}

	err := h.BlogService.Create(ctx, newBlog)
	if err != nil {
		span.RecordError(err, oteltrace.WithAttributes())
		log.Printf("Greška prilikom kreiranja bloga u servisu: %v", err)
		return nil, status.Errorf(codes.Internal, "Error creating blog: %v", err)
	}

	protoBlog := &blog.Blog{
		Id:        newBlog.ID.Hex(),
		Title:     newBlog.Title,
		Content:   newBlog.Content,
		UserId:    newBlog.UserId,
		ImageUrl:  newBlog.ImageUrl,
		CreatedAt: newBlog.CreatedAt.Format(time.RFC3339),
		Likes:     []string{},
		Comments:  []*blog.Comment{},
	}

	log.Printf("Blog '%s' uspešno kreiran putem gRPC-a.", newBlog.Title)
	return &blog.CreateBlogResponse{Blog: protoBlog}, nil
}
