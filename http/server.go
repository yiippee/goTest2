package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Test)

	http.ListenAndServe(":6001", mux)
}

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "hello,world.")
}
