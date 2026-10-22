package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "unknown"
	}
	resp := Response{Message: "Hello, " + name}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
