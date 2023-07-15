package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Internal server reached. Called from IP: ", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/internal", handler)
	fmt.Println("Starting server on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
