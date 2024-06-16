package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/rohan031/identity/database"
	"github.com/rohan031/identity/router"
	"github.com/rohan031/identity/services"
)

func initServer() (*chi.Mux, *pgxpool.Pool) {
	// getting db connection pool
	pool, err := database.CreatePool()
	if err != nil {
		log.Fatal("Error connecting to database\n", err)
	}

	// redis connection
	client, err := database.CreateRedisClient()
	if err != nil {
		log.Fatal("Error connecting to redis\n", err)
	}
	services.SetConnections(pool, client)

	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.Logger)                               // logging middleware
	r.Use(middleware.AllowContentType("application/json")) // to only allow req body with json

	r.Mount("/", router.Router())

	return r, pool
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

	router, pool := initServer()
	defer pool.Close()

	log.Printf("Server is listening on PORT: %s", PORT)
	http.ListenAndServe(":"+PORT, router)
}
