package main

import "net/http"

func httpReadiness(w http.ResponseWriter, r *http.Request) {
	resWithJSON(w, 200, struct{}{})
}