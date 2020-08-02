package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/jackc/pgconn"
)

func rowsCnt(r pgconn.CommandTag) error {
	count := r.RowsAffected()
	if count == 0 {
		return models.ErrNoRecord
	}
	return nil
}
