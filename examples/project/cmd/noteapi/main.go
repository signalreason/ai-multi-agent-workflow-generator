package main

import (
    "encoding/json"
    "log"
    "net/http"

    "example-project/internal/notes"
)

type server struct {
    store *notes.Store
}

func (s *server) listNotes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(s.store.List()); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (s *server) addNote(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        Body string `json:"body"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    note, err := s.store.Add(payload.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(note); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    srv := &server{store: notes.NewStore()}

    mux := http.NewServeMux()
    mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            srv.listNotes(w, r)
        case http.MethodPost:
            srv.addNote(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    log.Println("listening on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
