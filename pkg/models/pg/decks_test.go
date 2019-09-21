package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	_ "github.com/lib/pq"
	"reflect"
	"testing"
	"time"
)

var testDecks = [2]*models.Deck{
	{
		ID:          "test_deck_id_1",
		Name:        "Test Name 1",
		Description: "Test Description 1",
		Created:     time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
		Public:      false,
		CardsCount:  3,
		CreatedBy: models.DeckCreator{
			ID:   "testuser_id_1",
			Name: "Testuser1",
		},
	}, {
		ID:          "test_deck_id_2",
		Name:        "Test Name 2",
		Description: "Test Description 2",
		Created:     time.Date(2019, 2, 2, 12, 22, 0, 0, time.UTC),
		Public:      true,
		CardsCount:  2,
		CreatedBy: models.DeckCreator{
			ID:   "testuser_id_2",
			Name: "Testuser2",
		},
	},
}

func TestDeckModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}

	id, err := model.Create("Test Deck", "Deck, created from test", "testuser_id_1", true)
	if err != nil {
		t.Errorf("failed to create new deck: %s", err.Error())
	}
	deck, err := model.Get(*id)
	if err != nil {
		t.Errorf("failed to read created test deck: %s", err.Error())
	}
	if deck.Name != "Test Deck" {
		t.Errorf("invalid deck name: %s", deck.Name)
	}
	if deck.Description != "Deck, created from test" {
		t.Errorf("invalid deck description: %s", deck.Description)
	}
}

func TestDeckModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name      string
		deckId    string
		wantDeck  *models.Deck
		wantError error
	}{
		{
			name:      "Deck 1",
			deckId:    "test_deck_id_1",
			wantDeck:  testDecks[0],
			wantError: nil,
		},
		{
			name:      "Deck 2",
			deckId:    "test_deck_id_2",
			wantDeck:  testDecks[1],
			wantError: nil,
		},
		{
			name:      "Non-existent ID",
			deckId:    "test_deck_id_5",
			wantDeck:  nil,
			wantError: models.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			deck, err := model.Get(tt.deckId)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if !reflect.DeepEqual(deck, tt.wantDeck) {
				t.Errorf("want %v; got %v", tt.wantDeck, deck)
			}
		})
	}
}

func TestDeckModel_GetForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name      string
		userId    string
		wantDecks []*models.Deck
	}{
		{
			name:      "User 1",
			userId:    "testuser_id_1",
			wantDecks: []*models.Deck{testDecks[0]},
		},
		{
			name:      "User 2",
			userId:    "testuser_id_2",
			wantDecks: []*models.Deck{testDecks[1]},
		},
		{
			name:      "Non-existent user",
			userId:    "non-existing-testuser",
			wantDecks: []*models.Deck{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			deck, err := model.GetForUser(tt.userId)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(deck, tt.wantDecks) {
				t.Errorf("want %v; got %v", tt.wantDecks, deck)
			}
		})
	}
}

func TestDeckModel_GetPublic_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	decks, err := model.GetPublic()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(decks) != 1 {
		t.Errorf("unexpected collection size: want %d, got %d", 1, len(decks))
	}
	if !reflect.DeepEqual(testDecks[1], decks[0]) {
		t.Errorf("want %v; got %v", testDecks[1], decks[0])
	}
}

func TestDeckModel_Update_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	deck := testDecks[0]
	deck.Name = "Updated name"
	deck.Description = "Updated description"
	err := model.Update(deck)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	updated, err := model.Get("test_deck_id_1")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(deck, updated) {
		t.Errorf("want %v; got %v", deck, updated)
	}
}

func TestDeckModel_Update_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	deck := &models.Deck{
		ID:          "test_deck_id_5",
		Name:        "Updated 1",
		Description: "Updated Description 1",
		Created:     time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(deck)
	if err != models.ErrNoRecord {
		t.Errorf("unexpected error: want %s; got %s", models.ErrNoRecord, err)
	}
}

func TestDeckModel_Delete_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	err := model.Delete("test_deck_id_1")
	if err != nil {
		t.Errorf("failed to delete deck: %s", err)
	}
	_, err = model.Get("test_deck_id_1")
	if err != models.ErrNoRecord {
		t.Error("deck hasn't been deleted in delete test")
	}
}

func TestDeckModel_Delete_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	err := model.Delete("test_deck_id_5")
	if err != models.ErrNoRecord {
		t.Errorf("unexpected error: want %s; got %s", models.ErrNoRecord, err)
	}
}

func TestDeckModel_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name    string
		terms   []string
		wantLen int
	}{
		{
			name:    "All",
			terms:   []string{"Test"},
			wantLen: 2,
		},
		{
			name:    "First",
			terms:   []string{"Description", "1"},
			wantLen: 1,
		},
		{
			name:    "All, two terms",
			terms:   []string{"Test", "Name"},
			wantLen: 2,
		},
		{
			name:    "None",
			terms:   []string{"Not", "Here"},
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			decks, err := model.Search(tt.terms)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if len(decks) != tt.wantLen {
				t.Errorf("unexpected collection size: want %d, got %d", tt.wantLen, len(decks))
			}
		})
	}
}
