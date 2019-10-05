package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDeckNotFound = errors.New("models: deck does not exist")
var ErrNonUniqueEmail = errors.New("models: user with this email already registered")
var ErrNonUniqueCode = errors.New("models: verification code for user already exists")

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
	Created     time.Time   `json:"created"`
	CreatedBy   DeckCreator `json:"created_by"`
	Public      bool        `json:"public"`
}

type DeckCreator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Token    Token     `json:"token,omitempty"`
}

type Token struct {
	AccessToken     string    `json:"access_token,omitempty"`
	RefreshToken    string    `json:"refresh_token,omitempty"`
	RefreshTokenExp time.Time `json:"refresh_token_exp"`
}

type VerificationCode struct {
	UserID  string
	Code    string
	CodeExp time.Time
}

func (c VerificationCode) Expired() bool {
	return time.Now().After(c.CodeExp)
}
