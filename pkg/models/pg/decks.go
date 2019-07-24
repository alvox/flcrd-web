package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
)

type DeckModel struct {
	DB *sql.DB
}

func (m *DeckModel) Create(name, description string) (*string, error) {
	stmt := `insert into flcrd.deck (name, description) values ($1, $2) returning id;`
	var id string
	err := m.DB.QueryRow(stmt, name, description).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *DeckModel) Find(name string) (*models.Deck, error) {
	stmt := `select id, name, description, created from flcrd.deck 
    where name = $1;`
	d := &models.Deck{}
	err := m.DB.QueryRow(stmt, name).Scan(&d.ID, &d.Name, &d.Description, &d.Created)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (m *DeckModel) Update(deck *models.Deck) error {
	stmt := `update flcrd.deck set name = $1, description = $2 where id = $3;`
	_, err := m.DB.Exec(stmt, deck.Name, deck.Description, deck.ID)
	return err
}
