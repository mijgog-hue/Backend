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
type API struct {
	store Store
}

func NewMemoryStore() *MemoryNoteStore {
	return &MemoryNoteStore{notes: []Note{}, nextID: 1}
}

func (s *MemoryNoteStore) Create(title string) Note {
	n := Note{ID: s.nextID, Text: title}
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

func (a *API) createNote(w http.ResponseWriter, r *http.Request) {
	var n Note
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	created := a.store.Create(n.Text)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (a *API) listNote(w http.ResponseWriter, r *http.Request) {
	notes := a.store.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (a *API) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	n, ok := a.store.Get(id)
	if !ok {
		http.Error(w, "note not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

func (a *API) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ok := a.store.Delete(id)

	if !ok {
		http.Error(w, "note not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	store := NewMemoryStore()
	api := &API{store: store}

	http.HandleFunc("GET /notes", api.listNote)
	http.HandleFunc("/notes", api.createNote)
	http.HandleFunc("GET /notes/{id}", api.getNote)
	http.HandleFunc("DELETE /notes/{id}", api.deleteNote)
	log.Println("server started on:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
