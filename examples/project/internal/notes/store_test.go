package notes

import "testing"

func TestStoreAddAndList(t *testing.T) {
    store := NewStore()

    if _, err := store.Add(""); err == nil {
        t.Fatal("expected error for empty body")
    }

    note, err := store.Add("hello")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    notes := store.List()
    if len(notes) != 1 {
        t.Fatalf("expected 1 note, got %d", len(notes))
    }
    if notes[0].ID != note.ID {
        t.Fatalf("expected note id %d, got %d", note.ID, notes[0].ID)
    }
}
