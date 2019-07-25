package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
)

type FlashcardModel struct {
	DB *sql.DB
}

func (m *FlashcardModel) Create(flashcard *models.Flashcard) (*string, error) {
	stmt := `insert into flcrd.flashcard (deck_id, front, rear) values ($1, $2, $3) returning id;`
	var id string
	err := m.DB.QueryRow(stmt, flashcard.DeckID, flashcard.Front, flashcard.Rear).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *FlashcardModel) Get(deckID, flashcardID string) (*models.Flashcard, error) {
	stmt := `select id, deck_id, front, rear, created from flcrd.flashcard where id = $1 and deck_id = $2;`
	c := &models.Flashcard{}
	err := m.DB.QueryRow(stmt, flashcardID, deckID).Scan(&c.ID, &c.DeckID, &c.Front, &c.Rear, &c.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (m *FlashcardModel) Update(flashcard *models.Flashcard) error {
	stmt := `update flcrd.flashcard set deck_id = $1, front = $2, rear = $3 where id = $4;`
	_, err := m.DB.Exec(stmt, flashcard.DeckID, flashcard.Front, flashcard.Rear, flashcard.ID)
	return err
}

func (m *FlashcardModel) Delete(deckID, flashcardID string) error {
	stmt := `delete from flcrd.flashcard where id = $1 and deck_id = $2;`
	_, err := m.DB.Exec(stmt, flashcardID, deckID)
	return err
}
