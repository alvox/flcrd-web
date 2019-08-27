package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"reflect"
	"testing"
	"time"
)

var testFlashcards = [2]*models.Flashcard{
	{
		ID:      "test_flashcard_id_1",
		DeckID:  "test_deck_id_1",
		Front:   "Test Front 1 1",
		Rear:    "Test Rear 1 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}, {
		ID:      "test_flashcard_id_5",
		DeckID:  "test_deck_id_2",
		Front:   "Test Front 2 2",
		Rear:    "Test Rear 2 2",
		Created: time.Date(2019, 5, 5, 15, 55, 0, 0, time.UTC),
	},
}

func TestFlashcardModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}

	c := &models.Flashcard{
		DeckID: "test_deck_id_1",
		Front:  "Test Front",
		Rear:   "Test Rear",
	}

	id, err := model.Create(c)
	if err != nil {
		t.Errorf("unexpected error when creating flashcard: %s", err)
	}
	flashcard, err := model.Get("test_deck_id_1", *id)
	if err != nil {
		t.Errorf("unexpected error when reading test flashcard: %s", err)
	}
	checkField("test_deck_id_1", flashcard.DeckID, t)
	checkField("Test Front", flashcard.Front, t)
	checkField("Test Rear", flashcard.Rear, t)
}

func TestFlashcardModel_Create_InvalidDeck(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	c := &models.Flashcard{
		DeckID: "test_deck_id_5",
		Front:  "Test Front",
		Rear:   "Test Rear",
	}
	_, err := model.Create(c)
	if err != models.ErrDeckNotFound {
		t.Errorf("unexpected error: want %s; got %s", models.ErrDeckNotFound, err)
	}
}

func TestFlashcardModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name          string
		deckId        string
		flashcardId   string
		wantFlashcard *models.Flashcard
		wantError     error
	}{
		{
			name:          "Flashcard 1",
			deckId:        "test_deck_id_1",
			flashcardId:   "test_flashcard_id_1",
			wantFlashcard: testFlashcards[0],
			wantError:     nil,
		},
		{
			name:          "Flashcard 5",
			deckId:        "test_deck_id_2",
			flashcardId:   "test_flashcard_id_5",
			wantFlashcard: testFlashcards[1],
			wantError:     nil,
		},
		{
			name:          "Non-existent Deck ID",
			deckId:        "test_deck_id_5",
			flashcardId:   "test_flashcard_id_1",
			wantFlashcard: nil,
			wantError:     models.ErrNoRecord,
		},
		{
			name:          "Non-existent Flashcard ID",
			deckId:        "test_deck_id_1",
			flashcardId:   "test_flashcard_id_10",
			wantFlashcard: nil,
			wantError:     models.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := FlashcardModel{db}
			flashcard, err := model.Get(tt.deckId, tt.flashcardId)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if !reflect.DeepEqual(flashcard, tt.wantFlashcard) {
				t.Errorf("want %v; got %v", tt.wantFlashcard, flashcard)
			}
		})
	}
}

func TestFlashcardModel_GetPublic_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name      string
		deckId    string
		wantCount int
	}{
		{
			name:      "Deck 1",
			deckId:    "test_deck_id_1",
			wantCount: 0,
		},
		{
			name:      "Deck 2",
			deckId:    "test_deck_id_2",
			wantCount: 2,
		},
		{
			name:      "Non-existent Deck",
			deckId:    "test_deck_id_5",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := FlashcardModel{db}
			flashcards, err := model.GetPublic(tt.deckId)
			if err != nil {
				t.Errorf("unexpected error while reading flashcards: %s", err)
			}
			if tt.wantCount != len(flashcards) {
				t.Errorf("want %d; got %d", tt.wantCount, len(flashcards))
			}
		})
	}
}

func TestFlashcardModel_Update_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcard := testFlashcards[0]
	flashcard.Front = "Updated Front"
	flashcard.Rear = "Updated Rear"
	err := model.Update(flashcard)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	updated, err := model.Get("test_deck_id_1", "test_flashcard_id_1")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(flashcard, updated) {
		t.Errorf("want %v; got %v", flashcard, updated)
	}
}

func TestFlashcardModel_UpdateNonExistentDeck(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcard := &models.Flashcard{
		ID:      "test_flashcard_id_1",
		DeckID:  "test_deck_id_5",
		Front:   "Updated 1",
		Rear:    "Updated Description 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(flashcard)
	if err != models.ErrDeckNotFound {
		t.Errorf("unexpected error: want %s; got %s", models.ErrDeckNotFound, err)
	}
}

func TestFlashcardModel_UpdateNonExistentFlashcard(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcard := &models.Flashcard{
		ID:      "test_flashcard_id_10",
		DeckID:  "test_deck_id_1",
		Front:   "Updated 1",
		Rear:    "Updated Description 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(flashcard)
	if err != models.ErrNoRecord {
		t.Errorf("unexpected error: want %s; got %s", models.ErrNoRecord, err)
	}
}

func TestFlashcardModel_Delete_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	err := model.Delete("test_deck_id_1", "test_flashcard_id_1")
	if err != nil {
		t.Errorf("failed to delete flashcard: %s", err)
	}
	_, err = model.Get("test_deck_id_1", "test_flashcard_id_1")
	if err != models.ErrNoRecord {
		t.Error("flashcard hasn't been deleted in delete test")
	}
}

func TestFlashcardModel_Delete_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	err := model.Delete("test_deck_id_5", "test_flashcard_id_1")
	if err != models.ErrNoRecord {
		t.Errorf("unexpected error: want %s; got %s", models.ErrNoRecord, err)
	}
}

func checkField(expected, actual string, t *testing.T) {
	if expected != actual {
		t.Errorf("want %s; got %s", expected, actual)
	}
}
