package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDeckNotFound = errors.New("models: deck does not exist")
var ErrNonUniqueEmail = errors.New("models: user with this email already registered")

type Flashcard struct {
	ID      string    `json:"id"`
	DeckID  string    `json:"deck_id"`
	Front   string    `json:"front"`
	Rear    string    `json:"rear"`
	Created time.Time `json:"created"`
}

func ParseFlashcard(r *http.Request) (*Flashcard, bool) {
	if r.Body == nil {
		return nil, false
	}
	var flashcard Flashcard
	err := json.NewDecoder(r.Body).Decode(&flashcard)
	if err != nil {
		return nil, false
	}
	return &flashcard, true
}

type Deck struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CardsCount  int         `json:"cards_count"`
	Created     time.Time   `json:"created"`
	CreatedBy   DeckCreator `json:"created_by"`
	Public      bool        `json:"public"`
}

type DeckCreator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ParseDeck(r *http.Request) (*Deck, bool) {
	if r.Body == nil {
		return nil, false
	}
	var deck Deck
	err := json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		return nil, false
	}
	return &deck, true
}

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Created  time.Time `json:"created"`
	Token    Token     `json:"token,omitempty"`
}

func ParseUser(r *http.Request) (*User, bool) {
	if r.Body == nil {
		return nil, false
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, false
	}
	return &user, true
}

type Token struct {
	AccessToken     string    `json:"access_token,omitempty"`
	RefreshToken    string    `json:"refresh_token,omitempty"`
	RefreshTokenExp time.Time `json:"refresh_token_exp"`
}

func ParseTokens(r *http.Request) (*Token, bool) {
	if r.Body == nil {
		return nil, false
	}
	var token Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		return nil, false
	}
	return &token, true
}
