package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type DeckModel struct {
	DB *pgxpool.Pool
}

func (m *DeckModel) Create(name, description, createdBy string, public bool) (*string, error) {
	stmt := `insert into flcrd.deck (name, description, created_by, public, search_tokens) 
             values ($1, $2, $3, $4, to_tsvector($5)) returning id;`
	var id string
	err := m.DB.QueryRow(context.Background(), stmt, name, description, createdBy, public, fmt.Sprint(name, " ", description)).Scan(&id)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if "deck_name_user_idx" == err.ConstraintName {
				return nil, models.ErrUniqueViolation
			}
		}
		return nil, err
	}
	return &id, nil
}

func (m *DeckModel) Update(deck *models.Deck) error {
	stmt := `update flcrd.deck set name = $1, description = $2, public = $3, 
                 search_tokens = to_tsvector($4) 
             where id = $5;`
	r, err := m.DB.Exec(context.Background(), stmt, deck.Name, deck.Description, deck.Public, fmt.Sprint(deck.Name, " ", deck.Description), deck.ID)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if "deck_name_user_idx" == err.ConstraintName {
				return models.ErrUniqueViolation
			}
		}
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}

func (m *DeckModel) Get(id string) (*models.Deck, error) {
	stmt := `select d.id, d.name, d.description, d.created, d.public,
                 (select count(*) from flcrd.flashcard where deck_id = d.id) as cards_count,
                 u.id, u.name
             from flcrd.deck d
             left join flcrd.user u on u.id = d.created_by
             where d.id = $1;`
	d := &models.Deck{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&d.ID, &d.Name, &d.Description, &d.Created, &d.Public,
		&d.CardsCount, &d.CreatedBy.ID, &d.CreatedBy.Name)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	d.Created = d.Created.UTC()
	return d, nil
}

func (m *DeckModel) GetPublic(offset, limit int) ([]*models.Deck, int, error) {
	stmt := `select count(*) from flcrd.deck where public = true;`
	var total int
	err := m.DB.QueryRow(context.Background(), stmt).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	stmt = `select d.id, d.name, d.description, d.created, d.public,
                 (select count(*) from flcrd.flashcard where deck_id = d.id) as cards_count,
                 u.id, u.name
             from flcrd.deck d
             left join flcrd.user u on u.id = d.created_by
             where d.public = true order by d.created offset $1 limit $2;`
	rows, err := m.DB.Query(context.Background(), stmt, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	decks, err := read(rows)
	return decks, total, err
}

func (m *DeckModel) GetForUser(userID string) ([]*models.Deck, error) {
	stmt := `select d.id, d.name, d.description, d.created, d.public,
                 (select count(*) from flcrd.flashcard where deck_id = d.id) as cards_count,
                 u.id, u.name
             from flcrd.deck d
             left join flcrd.user u on u.id = d.created_by
             where u.id = $1 order by d.created;`
	rows, err := m.DB.Query(context.Background(), stmt, userID)
	if err != nil {
		return nil, err
	}
	return read(rows)
}

func (m *DeckModel) Delete(id string) error {
	stmt := `delete from flcrd.deck where id = $1;`
	r, err := m.DB.Exec(context.Background(), stmt, id)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}

func (m *DeckModel) Search(terms []string) ([]*models.Deck, error) {
	t := strings.Join(terms[:], "<->")
	stmt := `select d.id, d.name, d.description, d.created, d.public,
                 (select count(*) from flcrd.flashcard where deck_id = d.id) as cards_count,
                 u.id, u.name
             from flcrd.deck d
             left join flcrd.user u on u.id = d.created_by
             where d.public = true and d.search_tokens @@ to_tsquery($1) order by d.created;`
	rows, err := m.DB.Query(context.Background(), stmt, t)
	if err != nil {
		return nil, err
	}
	return read(rows)
}

func read(rows pgx.Rows) ([]*models.Deck, error) {
	defer rows.Close()
	decks := []*models.Deck{}
	var err error
	for rows.Next() {
		d := &models.Deck{}
		err = rows.Scan(&d.ID, &d.Name, &d.Description, &d.Created, &d.Public,
			&d.CardsCount, &d.CreatedBy.ID, &d.CreatedBy.Name)
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
