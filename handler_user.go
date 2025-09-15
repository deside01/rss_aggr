package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deside01/rss_aggr/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Name string `json:"name"`
}

func (apiCfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		resWithErr(w, 400, fmt.Sprint("Cannont parse JSON:", err))
		return
	}

	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		resWithErr(w, 500, fmt.Sprint("Cannot create user:", err))
		return
	}

	resWithJSON(w, 201, dbUserToUser(newUser))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	resWithJSON(w, 200, dbUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		resWithErr(w, 500, fmt.Sprintf("hz: %v", err))
	}

	resWithJSON(w, 200, dbPostsToPosts(posts))
}
