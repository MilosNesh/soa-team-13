package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gateway.com/handler"
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
	conn, err := grpc.DialContext(
		context.Background(),
		"stakeholders:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	client := stakeholders.NewStakeholdersServiceClient(conn)
	err = stakeholders.RegisterStakeholdersServiceHandlerClient(
		context.Background(),
		gwmux,
		client,
	)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
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
		if err := gwServer.ListenAndServe(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	if err = gwServer.Close(); err != nil {
		log.Fatalln("error while stopping server: ", err)
	}
	log.Println("Started GRPC...")
}

func main() {
	log.Println("Starting application...") // Log koji označava početak rada aplikacije
	handler := &handler.GatewayHandler{}
	go startServer(handler)
	connectToGRPCServer()
}
