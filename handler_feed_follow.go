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

type FollowFeedParams struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (apiCfg *apiConfig) handlerFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)

	params := FollowFeedParams{}

	err := decoder.Decode(&params)
	if err != nil {
		resWithErr(w, 400, fmt.Sprint("Cannont parse JSON:", err))
		return
	}

	followFeed, err := apiCfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		resWithErr(w, 500, fmt.Sprint("Cannot follow:", err))
		return
	}

	resWithJSON(w, 201, dbFeedFollowToFeedFollow(followFeed))
}

func (apiCfc *apiConfig) handleGetUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfc.DB.GetUserFollowedFeeds(r.Context(), user.ID)
	if err != nil {
		resWithErr(w, 400, fmt.Sprintf("hz: %v", err))
	}

	resWithJSON(w, 200, dbFeedsFollowsToFeedsFollows(feeds))
}

func (apiCfg *apiConfig) handleDeleteFollowedFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	paramID := chi.URLParam(r, "feedID")
	id, err := uuid.Parse(paramID)
	if err != nil {
		resWithErr(w, 400, fmt.Sprintf("hz: %v", err))
		return
	}

	_, err = apiCfg.DB.DeleteFollowedFeed(r.Context(), database.DeleteFollowedFeedParams{
		ID:     id,
		UserID: user.ID,
	})
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

// func (apiCfc *apiConfig) handleGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
// 	feeds, err := apiCfc.DB.GetUserFeeds(r.Context(), user.ID)
// 	if err != nil {
// 		resWithErr(w, 500, fmt.Sprintf("hz: %v", err))
// 	}

// 	resWithJSON(w, 200, dbFeedsToFeeds(feeds))
// }
