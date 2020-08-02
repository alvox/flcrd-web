package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type FlashcardModel struct {
	DB *pgxpool.Pool
}

func (m *FlashcardModel) Create(flashcard *models.Flashcard) (*string, error) {
	stmt := `insert into flcrd.flashcard (deck_id, front, rear) values ($1, $2, $3) returning id;`
	var id string
	err := m.DB.QueryRow(context.Background(), stmt, flashcard.DeckID, flashcard.Front, flashcard.Rear).Scan(&id)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if "flashcard_deck_id_fkey" == err.ConstraintName {
				return nil, models.ErrDeckNotFound
			}
		}
		return nil, err
	}
	return &id, nil
}

func (m *FlashcardModel) Get(deckID, flashcardID string) (*models.Flashcard, error) {
	stmt := `select id, deck_id, front, rear, created from flcrd.flashcard where id = $1 and deck_id = $2;`
	c := &models.Flashcard{}
	err := m.DB.QueryRow(context.Background(), stmt, flashcardID, deckID).Scan(&c.ID, &c.DeckID, &c.Front, &c.Rear, &c.Created)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	c.Created = c.Created.UTC()
	return c, nil
}

func (m *FlashcardModel) GetForUser(deckID, userID string) ([]*models.Flashcard, error) {
	stmt := `select f.id, f.deck_id, f.front, f.rear, f.created from flcrd.flashcard f
			join flcrd.deck d on d.id = f.deck_id
			where f.deck_id = $1 and d.created_by = $2;`
	rows, err := m.DB.Query(context.Background(), stmt, deckID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	flashcards := []*models.Flashcard{}
	for rows.Next() {
		c := &models.Flashcard{}
		err = rows.Scan(&c.ID, &c.DeckID, &c.Front, &c.Rear, &c.Created)
		if err != nil {
			return nil, err
		}
		c.Created = c.Created.UTC()
		flashcards = append(flashcards, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return flashcards, nil
}

func (m *FlashcardModel) GetPublic(deckID string) ([]*models.Flashcard, error) {
	stmt := `select f.id, f.deck_id, f.front, f.rear, f.created from flcrd.flashcard f 
        join flcrd.deck d on d.id = f.deck_id where d.id = $1 and d.public = true;`
	rows, err := m.DB.Query(context.Background(), stmt, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	flashcards := []*models.Flashcard{}
	for rows.Next() {
		c := &models.Flashcard{}
		err = rows.Scan(&c.ID, &c.DeckID, &c.Front, &c.Rear, &c.Created)
		if err != nil {
			return nil, err
		}
		c.Created = c.Created.UTC()
		flashcards = append(flashcards, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return flashcards, nil
}

func (m *FlashcardModel) Update(flashcard *models.Flashcard) error {
	stmt := `update flcrd.flashcard set deck_id = $1, front = $2, rear = $3 where id = $4;`
	r, err := m.DB.Exec(context.Background(), stmt, flashcard.DeckID, flashcard.Front, flashcard.Rear, flashcard.ID)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if "flashcard_deck_id_fkey" == err.ConstraintName {
				return models.ErrDeckNotFound
			}
		}
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}

func (m *FlashcardModel) Delete(deckID, flashcardID string) error {
	stmt := `delete from flcrd.flashcard where id = $1 and deck_id = $2;`
	r, err := m.DB.Exec(context.Background(), stmt, flashcardID, deckID)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}
