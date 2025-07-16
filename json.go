package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func resWithErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error:", msg)
	}

	resWithJSON(w, code, Error{
		Error: msg,
	})
}

func resWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Unable to marshal JSON: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

