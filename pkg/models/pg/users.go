package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(name, email, passwordHash string) (*string, error) {
	stmt := `insert into flcrd.user (name, email, password) values ($1, $2, $3) returning id;`
	var id string
	err := m.DB.QueryRow(stmt, name, email, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	stmt := `select id, name, email, password, created from flcrd.user where email = $1;`
	d := &models.User{}
	err := m.DB.QueryRow(stmt, email).Scan(&d.ID, &d.Name, &d.Email, &d.Password, &d.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	d.Created = d.Created.UTC()
	return d, nil
}
