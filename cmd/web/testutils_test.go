package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/mock"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApp(t *testing.T) *application {
	return &application{
		decks:      &mock.DeckModel{},
		flashcards: &mock.FlashcardModel{},
		users:      &mock.UserModel{},
		infoLog:    log.New(ioutil.Discard, "", 0),
		errorLog:   log.New(ioutil.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) post(t *testing.T, urlPath, body string) (int, http.Header, []byte) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	response, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, response
}

func parseUser(resp string) (*models.User, bool) {
	var user models.User
	err := json.NewDecoder(strings.NewReader(resp)).Decode(&user)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func parseError(resp string) (*ApiError, bool) {
	var e ApiError
	err := json.NewDecoder(strings.NewReader(resp)).Decode(&e)
	if err != nil {
		return nil, false
	}
	return &e, true
}
