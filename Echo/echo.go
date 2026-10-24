package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type EchoMessage struct {
	Text string `json:"text"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EchoMessage
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func main() {
	http.HandleFunc("/echo", echoHandler)
	log.Println("Server Started on: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
