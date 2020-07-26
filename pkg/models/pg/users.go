package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
	"github.com/lib/pq"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(user *models.User) (*string, error) {
	stmt := `insert into flcrd.user (email, external_id) 
             values ($1, $2) returning id;`
	var id string
	err := m.DB.QueryRow(stmt, user.Email, user.ExternalID).Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if "unique_violation" == err.Code.Name() {
				return nil, models.ErrUniqueViolation
			}
		}
		return nil, err
	}
	return &id, nil
}

func (m *UserModel) Get(userID string) (*models.User, error) {
	stmt := `select id, email, external_id, created from flcrd.user where external_id = $1;`
	u := &models.User{}
	err := m.DB.QueryRow(stmt, userID).Scan(&u.ID, &u.Email, &u.ExternalID, &u.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	u.Created = u.Created.UTC()
	return u, nil
}

func (m *UserModel) GetProfile(userID string) (*models.User, error) {
	stmt := `
    select u.id, u.email,
        count(d.id) as decks_count,
        (select count(id) from flcrd.flashcard where deck_id in
             (select id from flcrd.deck where created_by = u.id)) as cards_count
    from flcrd.user u
    left join flcrd.deck d on d.created_by = u.id
    where u.id = $1 group by u.id;`
	u := &models.User{}
	err := m.DB.QueryRow(stmt, userID).Scan(&u.ID, &u.Email, &u.Stats.DecksCount, &u.Stats.CardsCount)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	u.Created = u.Created.UTC()
	return u, nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	stmt := `select id, email, external_id, created from flcrd.user where email = $1;`
	u := &models.User{}
	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Email, &u.ExternalID, &u.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	u.Created = u.Created.UTC()
	return u, nil
}

func (m *UserModel) Update(user *models.User) error {
	stmt := `update flcrd.user set email = $1 where id = $2;`
	r, err := m.DB.Exec(stmt, user.Email, user.ID)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Delete(userID string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	r, err := tx.Exec(`delete from flcrd.user where id = $1;`, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := rowsCnt(r); err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.Exec(`delete from flcrd.deck where created_by = $1;`, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
