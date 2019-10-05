package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"database/sql"
)

func rowsCnt(r sql.Result) error {
	count, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return models.ErrNoRecord
	}
	return nil
}
