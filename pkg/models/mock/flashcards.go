package mock

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"time"
)

var mockFlashcards = []*models.Flashcard{
	{
		ID:      "test_flashcard_id_1",
		DeckID:  "test_deck_id_1",
		Front:   "Test Front 1",
		Rear:    "Test Rear 1",
		Created: time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}, {
		ID:      "test_flashcard_id_2",
		DeckID:  "test_deck_id_1",
		Front:   "Test Front 2",
		Rear:    "Test Rear 2",
		Created: time.Date(2019, 5, 5, 15, 55, 0, 0, time.UTC),
	},
}

type FlashcardModel struct{}

func (m *FlashcardModel) Create(flashcard *models.Flashcard) (*string, error) {
	return &mockFlashcards[0].ID, nil
}

func (m *FlashcardModel) Get(deckID, flashcardID string) (*models.Flashcard, error) {
	return mockFlashcards[0], nil
}

func (m *FlashcardModel) GetAll(deckID string) ([]*models.Flashcard, error) {
	return mockFlashcards, nil
}

func (m *FlashcardModel) Update(flashcard *models.Flashcard) error {
	return nil
}

func (m *FlashcardModel) Delete(deckID, flashcardID string) error {
	return nil
}
