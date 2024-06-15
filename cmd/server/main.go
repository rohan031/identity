package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rohan031/identity/router"
)

func initServer() *chi.Mux {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Logger)                               // logging middleware
	r.Use(middleware.AllowContentType("application/json")) // to only allow req body with json

	r.Mount("/", router.Router())

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env file: %v\n", err)
	}

	PORT := "8080"
	if port := os.Getenv("PORT"); port != "" {
		PORT = port
	}

	router := initServer()

	log.Printf("Server is listening on PORT: %s", PORT)
	http.ListenAndServe(":"+PORT, router)
}
