package main

import (
	"fmt"
	"net/http"

	"github.com/deside01/rss_aggr/internal/auth"
	"github.com/deside01/rss_aggr/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			resWithErr(w, 403, fmt.Sprint("Auth error: ", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			resWithErr(w, 404, fmt.Sprint("Couldn't get user: ", err))
			return
		}

		handler(w, r, user)
	}
}
