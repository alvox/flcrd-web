package mock

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"time"
)

var mockDeck = &models.Deck{
	ID:          "mock_deck_id_1",
	Name:        "Mock Deck",
	Description: "This is a mock deck",
	Private:     true,
	Created:     time.Now(),
}

type DeckModel struct{}

func (m *DeckModel) Create(name, description, createdBy string, private bool) (*string, error) {
	return &mockDeck.ID, nil
}

func (m *DeckModel) Update(deck *models.Deck) error {
	switch deck.ID {
	case "mock_deck_id_1":
		return nil
	default:
		return models.ErrNoRecord
	}
}

func (m *DeckModel) Get(id string) (*models.Deck, error) {
	switch id {
	case "mock_deck_id_1":
		return mockDeck, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *DeckModel) GetPublic() ([]*models.Deck, error) {
	return []*models.Deck{mockDeck}, nil
}

func (m *DeckModel) GetForUser(userID string) ([]*models.Deck, error) {
	return []*models.Deck{mockDeck}, nil
}

func (m *DeckModel) Delete(id string) error {
	return nil
}
