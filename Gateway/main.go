package main

import (
	"log"
	"net/http"

	"gateway.com/handler"
	"github.com/gorilla/mux"
)

func startServer(handler *handler.GatewayHandler) {
	router := mux.NewRouter()
	router.HandleFunc("/accounts/{anything:.*}", handler.HandleAccount).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/blogs/{anything:.*}", handler.HandleBlog).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/tours/{anything:.*}", handler.HandleTour).Methods("GET", "POST", "PUT", "DELETE")

	println("Gateway started...")
	log.Fatal(http.ListenAndServe(":8070", router))
}

func main() {
	// router := mux.NewRouter()
	// router.HandleFunc("/", main()).Methods("GET")
	handler := &handler.GatewayHandler{}
	startServer(handler)
}
