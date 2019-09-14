package mock

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"time"
)

var mockDeck = &models.Deck{
	ID:          "mock_deck_id_1",
	Name:        "Mock Deck",
	Description: "This is a mock deck",
	Public:      false,
	Created:     time.Now(),
	CreatedBy: models.DeckCreator{
		ID:   "test_user_1",
		Name: "Testuser1",
	},
}

type DeckModel struct{}

func (m *DeckModel) Create(name, description, createdBy string, public bool) (*string, error) {
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
	if userID == "test_user_1" {
		return []*models.Deck{mockDeck}, nil
	}
	return nil, nil
}

func (m *DeckModel) Delete(id string) error {
	return nil
}
