package main

import (
	"net/http"
	"testing"
)

func TestGetDeckForUser(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	token, err := generateAccessToken("test_user_1")
	if err != nil {
		t.Fatal(err)
	}

	status, _, resp := ts.get(t, "/v0/decks", *token)
	decks, valid := parseDecks(string(resp))
	if !valid {
		t.Error()
	}
	if status != http.StatusOK {
		t.Errorf("status: want 200; got %d", status)
	}
	if len(*decks) != 1 {
		t.Errorf("expected 1 deck; got %d", len(*decks))
	}
}
