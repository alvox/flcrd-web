package pg

import (
	_ "github.com/lib/pq"
	"testing"
)

func TestDeckModel_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}

	id, err := model.Create("Test Deck", "Deck, created from test")
	if err != nil {
		t.Error("Failed to create new deck")
	}

	deck, err := model.Get(*id)
	if err != nil {
		t.Error("failed to read created test deck")
	}
	if deck.Name != "Test Deck" {
		t.Errorf("invalid deck name: %s", deck.Name)
	}
	if deck.Description != "Deck, created from test" {
		t.Errorf("invalid deck description: %s", deck.Description)
	}
}
