package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway.com/handler"
	"gateway.com/proto/blog"
	"gateway.com/proto/stakeholders"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startServer(handler *handler.GatewayHandler) {
	router := mux.NewRouter()
	router.PathPrefix("/uploads/").Handler(http.HandlerFunc(handler.HandleStaticFiles))
	router.HandleFunc("/accounts/{anything:.*}", handler.HandleAccount).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/blogs/{anything:.*}", handler.HandleBlog).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/tours/{anything:.*}", handler.HandleTour).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/shopping/{anything:.*}", handler.HandleShopping).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)
	println("Gateway started...")
	log.Fatal(http.ListenAndServe(":8070", corsHandler))
}

func connectToGRPCServer() {
	gwmux := runtime.NewServeMux()

	connStakeholders, err := grpc.DialContext(
		context.Background(),
		"stakeholders:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial Stakeholders server:", err)
	}
	defer connStakeholders.Close()

	clientStakeholders := stakeholders.NewStakeholdersServiceClient(connStakeholders)
	err = stakeholders.RegisterStakeholdersServiceHandlerClient(
		context.Background(),
		gwmux,
		clientStakeholders,
	)
	if err != nil {
		log.Fatalln("Failed to register Stakeholders gateway:", err)
	}

	connBlog, err := grpc.DialContext(
		context.Background(),
		"blog:50052",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial Blog server:", err)
	}
	defer connBlog.Close()

	clientBlog := blog.NewBlogServiceClient(connBlog)
	err = blog.RegisterBlogServiceHandlerClient(
		context.Background(),
		gwmux,
		clientBlog,
	)
	if err != nil {
		log.Fatalln("Failed to register Blog gateway:", err)
	}

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(gwmux)

	gwServer := &http.Server{
		Addr:    ":8071",
		Handler: corsHandler,
	}

	go func() {
		if err := gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Gateway server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	log.Println("Shutting down Gateway server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = gwServer.Shutdown(ctx); err != nil {
		log.Fatalln("Error while stopping server: ", err)
	}
	log.Println("GRPC Gateway server stopped.")
}

func main() {
	log.Println("Starting application...") // Log koji označava početak rada aplikacije
	handler := &handler.GatewayHandler{}
	go startServer(handler)
	connectToGRPCServer()
}
