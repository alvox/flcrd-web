package pg

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"testing"
	"time"
)

func newTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	db, err := pgxpool.Connect(context.Background(), "postgres://test_flcrd:pass@localhost/test_flcrd?sslmode=disable")
	//db, err := sql.Open("postgres", "postgres://test_flcrd:pass@localhost/test_flcrd?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	setupScript, err := ioutil.ReadFile("../../../db/test/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(context.Background(), string(setupScript))
	if err != nil {
		t.Fatal(err)
	}
	return db, func() {
		teardownScript, err := ioutil.ReadFile("../../../db/test/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(teardownScript))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}

func parseTime(s string, t *testing.T) time.Time {
	tm, e := time.Parse(time.RFC3339, s)
	if e != nil {
		t.Errorf("unexpected error while preparing test data: %s", e.Error())
	}
	return tm
}
