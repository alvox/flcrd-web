package pg

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("postgres", "postgres://test:pass@localhost/test_flcrd?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	setupScript, err := ioutil.ReadFile("../../../db/test/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(setupScript))
	if err != nil {
		t.Fatal(err)
	}
	return db, func() {
		teardownScript, err := ioutil.ReadFile("../../../db/test/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(teardownScript))
		if err != nil {
			t.Fatal(err)
		}
		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}
