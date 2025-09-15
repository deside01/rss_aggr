package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/go-chi/chi/v5"
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

func (apiCfc *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfc.DB.GetFeeds(r.Context())
	if err != nil {
		resWithErr(w, 500, fmt.Sprintf("hz: %v", err))
	}

	resWithJSON(w, 200, dbFeedsToFeeds(feeds))
}

func (apiCfc *apiConfig) handleGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfc.DB.GetUserFeeds(r.Context(), user.ID)
	if err != nil {
		resWithErr(w, 500, fmt.Sprintf("hz: %v", err))
	}

	resWithJSON(w, 200, dbFeedsToFeeds(feeds))
}

func (apiCfc *apiConfig) handleDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	paramID := chi.URLParam(r, "feedID")

	id, err := uuid.Parse(paramID)
	if err != nil {
		resWithErr(w, 400, fmt.Sprintf("hz: %v", err))
		return
	}

	_, err = apiCfc.DB.DeleteFeed(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resWithErr(w, 404, "not found")
			return
		}
		resWithErr(w, 400, fmt.Sprintf("err: %v", err))
		return
	}

	resWithJSON(w, 204, "")
}
