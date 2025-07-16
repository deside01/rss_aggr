package main

import "net/http"

func httpErr(w http.ResponseWriter, r *http.Request) {
	resWithErr(w, 400, "something went wrong")
}