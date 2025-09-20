package main

import (
	"log"
	"net/http"

	"gateway.com/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startServer(handler *handler.GatewayHandler) {
	router := mux.NewRouter()
	router.HandleFunc("/accounts/{anything:.*}", handler.HandleAccount).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/blogs/{anything:.*}", handler.HandleBlog).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.HandleFunc("/tours/{anything:.*}", handler.HandleTour).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)
	println("Gateway started...")
	log.Fatal(http.ListenAndServe(":8070", corsHandler))
}

func main() {
	handler := &handler.GatewayHandler{}
	handler.ConnectToGRPCServer()
	startServer(handler)
}
