package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
	"github.com/lib/pq"
)

type VerificationModel struct {
	DB *sql.DB
}

func (m *VerificationModel) Create(code models.VerificationCode) (string, error) {
	stmt := `insert into flcrd.verification_code (user_id, code, code_exp) values ($1, $2, $3) returning code;`
	var c string
	err := m.DB.QueryRow(stmt, code.UserID, code.Code, code.CodeExp).Scan(&c)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if "unique_violation" == err.Code.Name() {
				return "", models.ErrNonUniqueCode
			}
		}
		return "", err
	}
	return c, nil
}

func (m *VerificationModel) Get(code string) (*models.VerificationCode, error) {
	stmt := `select user_id, code, code_exp from flcrd.verification_code where code = $1;`
	c := &models.VerificationCode{}
	err := m.DB.QueryRow(stmt, code).Scan(&c.UserID, &c.Code, &c.CodeExp)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}
	if err != nil {
		return nil, err
	}
	c.CodeExp = c.CodeExp.UTC()
	return c, nil
}

func (m *VerificationModel) Delete(code models.VerificationCode) error {
	_, err := m.Get(code.Code)
	if err != nil {
		return err
	}
	stmt := `delete from flcrd.verification_code where user_id = $1;`
	_, err = m.DB.Exec(stmt, code.UserID)
	return err
}
