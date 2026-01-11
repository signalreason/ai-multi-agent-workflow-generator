package notes

import (
    "errors"
    "sync"
    "time"
)

// Note captures a single text note.
type Note struct {
    ID        int
    Body      string
    CreatedAt time.Time
}

// Store keeps notes in memory for demo purposes.
type Store struct {
    mu    sync.Mutex
    next  int
    notes []Note
}

// NewStore creates a new Store.
func NewStore() *Store {
    return &Store{next: 1}
}

// Add stores a new note.
func (s *Store) Add(body string) (Note, error) {
    if body == "" {
        return Note{}, errors.New("body is required")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    note := Note{
        ID:        s.next,
        Body:      body,
        CreatedAt: time.Now().UTC(),
    }
    s.next++
    s.notes = append(s.notes, note)
    return note, nil
}

// List returns all notes.
func (s *Store) List() []Note {
    s.mu.Lock()
    defer s.mu.Unlock()

    out := make([]Note, len(s.notes))
    copy(out, s.notes)
    return out
}
