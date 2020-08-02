package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Create(user *models.User, credentials *models.Credentials) (*string, error) {
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	stmt := `insert into flcrd.user (name, email, status) values ($1, $2, $3) returning id;`
	var id string
	err = tx.QueryRow(ctx, stmt, user.Name, user.Email, user.Status).Scan(&id)
	if err != nil {
		_ = tx.Rollback(ctx)
		if err, ok := err.(*pgconn.PgError); ok {
			if "user_email_idx" == err.ConstraintName {
				return nil, models.ErrUniqueViolation
			}
		}
		return nil, err
	}
	stmt = `insert into flcrd.credentials (user_id, password, refresh_token, refresh_token_exp) values ($1, $2, $3, $4);`
	_, err = tx.Exec(ctx, stmt, id, credentials.Password, credentials.Token.RefreshToken, credentials.Token.RefreshTokenExp)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *UserModel) Get(userID string) (*models.User, error) {
	stmt := `select u.id, u.name, u.email, u.status, u.created, 
                    c.refresh_token, c.refresh_token_exp 
             from flcrd.user u
             left join flcrd.credentials c on c.user_id = u.id
             where u.id = $1;`
	u := &models.User{}
	err := m.DB.QueryRow(context.Background(), stmt, userID).Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.Created,
		&u.Token.RefreshToken, &u.Token.RefreshTokenExp)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	u.Created = u.Created.UTC()
	u.Token.RefreshTokenExp = u.Token.RefreshTokenExp.UTC()
	return u, nil
}

func (m *UserModel) GetProfile(userID string) (*models.User, error) {
	stmt := `
    select u.id, u.name, u.email, u.status,
        count(d.id) as decks_count,
        (select count(id) from flcrd.flashcard where deck_id in
             (select id from flcrd.deck where created_by = u.id)) as cards_count
    from flcrd.user u
    left join flcrd.deck d on d.created_by = u.id
    where u.id = $1 group by u.id;`
	u := &models.User{}
	err := m.DB.QueryRow(context.Background(), stmt, userID).Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.Stats.DecksCount, &u.Stats.CardsCount)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	stmt := `select id, name, email, status, created from flcrd.user where email = $1;`
	u := &models.User{}
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.Created)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	u.Created = u.Created.UTC()
	return u, nil
}

func (m *UserModel) Update(user *models.User) error {
	stmt := `update flcrd.user set name = $1, email = $2, status = $3 where id = $4;`
	r, err := m.DB.Exec(context.Background(), stmt, user.Name, user.Email, user.Status, user.ID)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Delete(userID string) error {
	ctx := context.Background()
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	r, err := tx.Exec(ctx, `delete from flcrd.user where id = $1;`, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := rowsCnt(r); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, `delete from flcrd.credentials where user_id = $1;`, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, `delete from flcrd.deck where created_by = $1;`, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// Credentials

func (m *UserModel) GetCredentials(userID string) (*models.Credentials, error) {
	stmt := `select user_id, password, refresh_token, refresh_token_exp 
             from flcrd.credentials 
             where user_id = $1;`
	c := &models.Credentials{}
	err := m.DB.QueryRow(context.Background(), stmt, userID).Scan(&c.UserID, &c.Password, &c.Token.RefreshToken, &c.Token.RefreshTokenExp)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	c.Token.RefreshTokenExp = c.Token.RefreshTokenExp.UTC()
	return c, nil
}

func (m *UserModel) UpdateRefreshToken(c *models.Credentials) error {
	stmt := `update flcrd.credentials set refresh_token = $1, refresh_token_exp = $2 where user_id = $3;`
	r, err := m.DB.Exec(context.Background(), stmt, c.Token.RefreshToken, c.Token.RefreshTokenExp, c.UserID)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}
