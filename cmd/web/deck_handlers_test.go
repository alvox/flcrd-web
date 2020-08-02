package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetDeckForUser(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	token, err := generateAccessToken("test_user_1")
	require.Nil(t, err)

	status, _, resp := ts.get(t, "/v0/decks", *token)
	decks, valid := parseDecks(string(resp))
	require.True(t, valid)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, 1, len(*decks))
}
