package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	// log.SetFlags(0)
	godotenv.Load()

	PORT := strings.TrimSpace(os.Getenv("PORT"))
	HOST := strings.TrimSpace(os.Getenv("HOST"))
	DB_URL := strings.TrimSpace(os.Getenv("DB_URL"))

	if PORT == "" {
		log.Fatal("Missing PORT env")
	}

	if HOST == "" {
		log.Fatal("Missing HOST env")
	}

	if DB_URL == "" {
		log.Fatal("Missing HOST env")
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Unable connect to DB:", err)
	}

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	go startScraper(db, 3, time.Minute)

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
	v1.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeed))
	v1.Delete("/feeds/{feedID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeed))
	v1.Get("/user/feeds", apiCfg.middlewareAuth(apiCfg.handleGetUserFeeds))
	v1.Get("/feeds", apiCfg.handleGetFeeds)
	v1.Post("/feeds/follow", apiCfg.middlewareAuth(apiCfg.handlerFeedFollow))
	v1.Get("/user/followed", apiCfg.middlewareAuth(apiCfg.handleGetUserFeedFollows))
	v1.Get("/user/feeds", apiCfg.middlewareAuth(apiCfg.handlerGetUserPosts))
	v1.Delete("/user/followed/{feedID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFollowedFeed))

	r.Mount("/v1", v1)

	addr := fmt.Sprintf("%v:%v", HOST, PORT)

	log.Printf("Server is running on %v port", PORT)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
