package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	h := r
	http.ListenAndServe(":7777", h)
}
