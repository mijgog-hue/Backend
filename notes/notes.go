package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Store interface {
	Create(title string) Note
	List() []Note
	Get(id int) (Note, bool)
	Delete(id int) bool
}

type MemoryNoteStore struct {
	notes  []Note
	nextID int
}

func NewMemoryStore() *MemoryNoteStore {
	return &MemoryNoteStore{notes: []Note{}, nextID: 1}
}

func (s *MemoryNoteStore) Create(text string) Note {
	n := Note{ID: s.nextID, Text: text}
	s.nextID++
	s.notes = append(s.notes, n)
	return n
}
func (s *MemoryNoteStore) List() []Note {
	return s.notes
}

func (s *MemoryNoteStore) Get(id int) (Note, bool) {
	for _, n := range s.notes {
		if n.ID == id {
			return n, true
		}
	}
	return Note{}, false
}

func (s *MemoryNoteStore) Delete(id int) bool {
	for i, n := range s.notes {
		if n.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return true
		}
	}
	return false
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

	s := NewMemoryStore()
	s.Create("first")
	s.Create("second")
	fmt.Println(s.Delete(1))  // true
	fmt.Println(s.Delete(99)) // false
	fmt.Println(s.List())
}
