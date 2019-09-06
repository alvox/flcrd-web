package models

import (
	"errors"
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

type Deck struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CardsCount  int         `json:"cards_count"`
	Cards       []Flashcard `json:"cards"`
	Created     time.Time   `json:"created"`
	CreatedBy   string      `json:"created_by"`
	Private     bool        `json:"private"`
}

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Created  time.Time `json:"created"`
	Token    Token     `json:"token,omitempty"`
}

type Token struct {
	AccessToken     string    `json:"access_token,omitempty"`
	RefreshToken    string    `json:"refresh_token,omitempty"`
	RefreshTokenExp time.Time `json:"refresh_token_exp"`
}
