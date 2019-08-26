package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
)

type DeckModel struct {
	DB *sql.DB
}

func (m *DeckModel) Create(name, description string, private bool) (*string, error) {
	stmt := `insert into flcrd.deck (name, description, private) values ($1, $2, $3) returning id;`
	var id string
	err := m.DB.QueryRow(stmt, name, description, private).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *DeckModel) Update(deck *models.Deck) error {
	_, err := m.Get(deck.ID)
	if err != nil {
		return err
	}
	stmt := `update flcrd.deck set name = $1, description = $2, private = $3 where id = $4;`
	_, err = m.DB.Exec(stmt, deck.Name, deck.Description, deck.Private, deck.ID)
	return err
}

func (m *DeckModel) Get(id string) (*models.Deck, error) {
	stmt := `select id, name, description, created, private,
        (select count(*) from flcrd.flashcard where deck_id = deck.id) as cards_count
         from flcrd.deck where id = $1;`
	d := &models.Deck{}
	err := m.DB.QueryRow(stmt, id).Scan(&d.ID, &d.Name, &d.Description, &d.Created, &d.Private, &d.CardsCount)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	d.Created = d.Created.UTC()
	return d, nil
}

func (m *DeckModel) GetAll() ([]*models.Deck, error) {
	stmt := `select id, name, description, created, private,
       (select count(*) from flcrd.flashcard where deck_id = deck.id) as cards_count 
       from flcrd.deck;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	decks := []*models.Deck{}
	for rows.Next() {
		d := &models.Deck{}
		err = rows.Scan(&d.ID, &d.Name, &d.Description, &d.Created, &d.Private, &d.CardsCount)
		if err != nil {
			return nil, err
		}
		d.Created = d.Created.UTC()
		decks = append(decks, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return decks, nil
}

func (m *DeckModel) Delete(id string) error {
	_, err := m.Get(id)
	if err != nil {
		return err
	}
	stmt := `delete from flcrd.deck where id = $1;`
	_, err = m.DB.Exec(stmt, id)
	return err
}
