package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/google/uuid"
)

type CreateFeedParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (apiCfg *apiConfig) handlerFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)

	params := CreateFeedParams{}

	err := decoder.Decode(&params)
	if err != nil {
		resWithErr(w, 400, fmt.Sprint("Cannont parse JSON:", err))
		return
	}

	newFeed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Url:       params.Url,
	})
	if err != nil {
		resWithErr(w, 500, fmt.Sprint("Cannot create user:", err))
		return
	}

	resWithJSON(w, 201, dbFeedToFeed(newFeed))
}
