package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var notes = []Note{}
var nextID = 1

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for _, n := range notes {
		if n.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(n)
			return
		}
	}
	http.Error(w, "note not found", http.StatusNotFound)
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for i, n := range notes {
		if n.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "note not found", http.StatusNotFound)
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	case http.MethodPost:
		var n Note
		err := json.NewDecoder(r.Body).Decode(&n)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		n.ID = nextID
		nextID++
		notes = append(notes, n)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(n)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	}

}

func main() {
	http.HandleFunc("/notes", noteHandler)
	http.HandleFunc("GET /notes/{id}", getNoteHandler)
	http.HandleFunc("DELETE /notes/{id}", deleteNoteHandler)
	log.Println("server started on:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
