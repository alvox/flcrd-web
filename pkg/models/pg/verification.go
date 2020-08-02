package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VerificationModel struct {
	DB *pgxpool.Pool
}

func (m *VerificationModel) Create(code models.VerificationCode) (string, error) {
	stmt := `insert into flcrd.verification_code (user_id, code, code_exp) values ($1, $2, $3) returning code;`
	var c string
	err := m.DB.QueryRow(context.Background(), stmt, code.UserID, code.Code, code.CodeExp).Scan(&c)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if "auth_code_idx" == err.ConstraintName {
				return "", models.ErrUniqueViolation
			}
		}
		return "", err
	}
	return c, nil
}

func (m *VerificationModel) Get(code string) (*models.VerificationCode, error) {
	stmt := `select user_id, code, code_exp from flcrd.verification_code where code = $1;`
	c := &models.VerificationCode{}
	err := m.DB.QueryRow(context.Background(), stmt, code).Scan(&c.UserID, &c.Code, &c.CodeExp)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	c.CodeExp = c.CodeExp.UTC()
	return c, nil
}

func (m *VerificationModel) GetForUser(userID string) (*models.VerificationCode, error) {
	stmt := `select user_id, code, code_exp from flcrd.verification_code where user_id = $1;`
	c := &models.VerificationCode{}
	err := m.DB.QueryRow(context.Background(), stmt, userID).Scan(&c.UserID, &c.Code, &c.CodeExp)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	c.CodeExp = c.CodeExp.UTC()
	return c, nil
}

func (m *VerificationModel) Delete(code models.VerificationCode) error {
	stmt := `delete from flcrd.verification_code where user_id = $1;`
	r, err := m.DB.Exec(context.Background(), stmt, code.UserID)
	if err != nil {
		return err
	}
	if err := rowsCnt(r); err != nil {
		return err
	}
	return nil
}
