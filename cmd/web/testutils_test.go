package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/mock"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApp(t *testing.T) *application {
	return &application{
		decks:        &mock.DeckModel{},
		flashcards:   &mock.FlashcardModel{},
		users:        &mock.UserModel{},
		verification: &mock.VerificationModel{},
		mailSender:   &MockMailSender{},
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath, token string) (int, http.Header, []byte) {
	req, err := http.NewRequest("GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	rs, err := ts.Client().Do(req)
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

func parseDecks(resp string) (*[]models.Deck, bool) {
	var decks []models.Deck
	err := json.NewDecoder(strings.NewReader(resp)).Decode(&decks)
	if err != nil {
		return nil, false
	}
	return &decks, true
}

type MockMailSender struct{}

func (s *MockMailSender) SendConfirmation(to, name, code string) (*SendMessageResponse, error) {
	return &SendMessageResponse{
		Message: "OK",
		Id:      "TEST_MAIL_ID",
	}, nil
}
