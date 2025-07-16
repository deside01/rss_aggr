package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")
	if strings.TrimSpace(PORT) == "" {
		log.Fatal("Missing PORT env")
	}

	if strings.TrimSpace(HOST) == "" {
		log.Fatal("Missing HOST env")
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

	r.Mount("/v1", v1)

	addr := fmt.Sprintf("%v:%v", HOST, PORT)

	log.Printf("Server is running on %v port", PORT)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
