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
<<<<<<< HEAD
	router.HandleFunc("/follow/{anything:.*}", handler.HandleFollow).Methods("GET", "POST", "PUT", "DELETE")
=======
	router.HandleFunc("/tours/{anything:.*}", handler.HandleTour).Methods("GET", "POST", "PUT", "DELETE")

>>>>>>> 2725621ea2ff74fa594bf65ac24278dabd1ef6ff
	println("Gateway started...")
	log.Fatal(http.ListenAndServe(":8070", router))
}

func main() {
	// router := mux.NewRouter()
	// router.HandleFunc("/", main()).Methods("GET")
	handler := &handler.GatewayHandler{}
	startServer(handler)
}
