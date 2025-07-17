package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	log.SetFlags(0)
	godotenv.Load()
	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")
	DB_URL := os.Getenv("DB_URL")

	if strings.TrimSpace(PORT) == "" {
		log.Fatal("Missing PORT env")
	}

	if strings.TrimSpace(HOST) == "" {
		log.Fatal("Missing HOST env")
	}

	if strings.TrimSpace(DB_URL) == "" {
		log.Fatal("Missing HOST env")
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Unable connect to DB:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	v1.Get("/", httpReadiness)
	v1.Get("/err", httpErr)
	v1.Post("/users", apiCfg.handlerUser)
	v1.Get("/users", apiCfg.handleGetUser)

	r.Mount("/v1", v1)

	addr := fmt.Sprintf("%v:%v", HOST, PORT)

	log.Printf("Server is running on %v port", PORT)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
