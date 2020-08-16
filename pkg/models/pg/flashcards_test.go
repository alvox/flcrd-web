package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testFlashcards = [2]*models.Flashcard{
	{
		ID:        "9f814806-e2df-4598-a323-1380d47b9c35",
		DeckID:    "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
		Front:     "Test Front 1 1",
		FrontType: "TEXT",
		Rear:      "Test Rear 1 1",
		RearType:  "TEXT",
		Created:   time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}, {
		ID:        "a5f65e8d-ca03-4d3e-978e-f5a612881231",
		DeckID:    "2601da50-56a6-41a1-a92e-5624598a7d19",
		Front:     "Test Front 2 2",
		FrontType: "TEXT",
		Rear:      "https://s3/testuser/testdeck/testimg.jpeg",
		RearType:  "IMAGE_URL",
		Created:   time.Date(2019, 5, 5, 15, 55, 0, 0, time.UTC),
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
		DeckID:    "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
		Front:     "Test Front",
		FrontType: "TEXT",
		Rear:      "Test Rear",
		RearType:  "https://testurl",
	}

	id, err := model.Create(c)
	require.Nil(t, err)
	flashcard, err := model.Get("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", *id)
	require.Nil(t, err)
	require.Equal(t, "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", flashcard.DeckID)
	require.Equal(t, "Test Front", flashcard.Front)
	require.Equal(t, "Test Rear", flashcard.Rear)
	require.Equal(t, "TEXT", flashcard.FrontType)
	require.Equal(t, "https://testurl", flashcard.RearType)
}

func TestFlashcardModel_Create_InvalidDeck(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	c := &models.Flashcard{
		DeckID:    "7af65126-d46c-4797-a329-09d283acc664",
		Front:     "Test Front",
		FrontType: "TEXT",
		Rear:      "Test Rear",
		RearType:  "TEXT",
	}
	_, err := model.Create(c)
	require.Equal(t, models.ErrDeckNotFound, err)
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
			deckId:        "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
			flashcardId:   "9f814806-e2df-4598-a323-1380d47b9c35",
			wantFlashcard: testFlashcards[0],
			wantError:     nil,
		},
		{
			name:          "Flashcard 5",
			deckId:        "2601da50-56a6-41a1-a92e-5624598a7d19",
			flashcardId:   "a5f65e8d-ca03-4d3e-978e-f5a612881231",
			wantFlashcard: testFlashcards[1],
			wantError:     nil,
		},
		{
			name:          "Non-existent Deck ID",
			deckId:        "7af65126-d46c-4797-a329-09d283acc664",
			flashcardId:   "9f814806-e2df-4598-a323-1380d47b9c35",
			wantFlashcard: nil,
			wantError:     models.ErrNoRecord,
		},
		{
			name:          "Non-existent Flashcard ID",
			deckId:        "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
			flashcardId:   "7af65126-d46c-4797-a329-09d283acc664",
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
			require.Equal(t, tt.wantError, err)
			require.Equal(t, tt.wantFlashcard, flashcard)
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
			deckId:    "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
			wantCount: 0,
		},
		{
			name:      "Deck 2",
			deckId:    "2601da50-56a6-41a1-a92e-5624598a7d19",
			wantCount: 2,
		},
		{
			name:      "Non-existent Deck",
			deckId:    "7af65126-d46c-4797-a329-09d283acc664",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := FlashcardModel{db}
			flashcards, err := model.GetPublic(tt.deckId)
			require.Nil(t, err)
			require.Equal(t, tt.wantCount, len(flashcards))
		})
	}
}

func TestFlashcardModel_GetForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcards, err := model.GetForUser("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", "40afbc9a-27e3-4b38-97f9-2930b8790a9f")
	require.Nil(t, err)
	require.Equal(t, 3, len(flashcards))
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
	flashcard.FrontType = "TEXT"
	flashcard.Rear = "https://s3/updatedurl"
	flashcard.RearType = "IMAGE_URL"
	err := model.Update(flashcard)
	require.Nil(t, err)
	updated, err := model.Get("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", "9f814806-e2df-4598-a323-1380d47b9c35")
	require.Nil(t, err)
	require.Equal(t, flashcard, updated)
}

func TestFlashcardModel_UpdateNonExistentDeck(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcard := &models.Flashcard{
		ID:      "9f814806-e2df-4598-a323-1380d47b9c35",
		DeckID:  "7af65126-d46c-4797-a329-09d283acc664",
		Front:   "Updated 1",
		Rear:    "Updated Description 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(flashcard)
	require.Equal(t, models.ErrDeckNotFound, err)
}

func TestFlashcardModel_UpdateNonExistentFlashcard(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	flashcard := &models.Flashcard{
		ID:      "7af65126-d46c-4797-a329-09d283acc664",
		DeckID:  "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
		Front:   "Updated 1",
		Rear:    "Updated Description 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(flashcard)
	require.Equal(t, models.ErrNoRecord, err)
}

func TestFlashcardModel_Delete_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	err := model.Delete("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", "9f814806-e2df-4598-a323-1380d47b9c35")
	require.Nil(t, err)
	_, err = model.Get("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e", "9f814806-e2df-4598-a323-1380d47b9c35")
	require.Equal(t, models.ErrNoRecord, err)
}

func TestFlashcardModel_Delete_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := FlashcardModel{db}
	err := model.Delete("7af65126-d46c-4797-a329-09d283acc664", "9f814806-e2df-4598-a323-1380d47b9c35")
	require.Equal(t, models.ErrNoRecord, err)
}
