package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	h := NewHandler()
	router := mux.NewRouter()
	h.InitRoutes(router)

	log.Printf("Yayy!, server's running on PORT: %s", "8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Oops!, Failed to start server: %v", err)
	}
	defer h.postgresClient.Close()
	defer h.redisClient.Close()
}
