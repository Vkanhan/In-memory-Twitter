package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(service *TwitterService) {
	handler := NewHandler(service)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	service := NewTwitterService()
	StartServer(service)
}
