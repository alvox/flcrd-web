package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDeckNotFound = errors.New("models: deck does not exist")
var ErrUniqueViolation = errors.New("models: unique violation")

type Flashcard struct {
	ID        string    `json:"id"`
	DeckID    string    `json:"deck_id"`
	Front     string    `json:"front"`
	FrontType string    `json:"front_type"`
	Rear      string    `json:"rear"`
	RearType  string    `json:"rear_type"`
	Created   time.Time `json:"created"`
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
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Status  string    `json:"status"`
	Created time.Time `json:"created,omitempty"`
	Token   Token     `json:"token,omitempty"`
	Stats   Stats     `json:"stats,omitempty"`
}

type Token struct {
	AccessToken     string    `json:"access_token,omitempty"`
	RefreshToken    string    `json:"refresh_token,omitempty"`
	RefreshTokenExp time.Time `json:"refresh_token_exp,omitempty"`
}

func (t Token) RefreshTokenExpired() bool {
	return time.Now().After(t.RefreshTokenExp)
}

type Stats struct {
	DecksCount int `json:"decks_count"`
	CardsCount int `json:"cards_count"`
}

type VerificationCode struct {
	UserID  string
	Code    string
	CodeExp time.Time
}

func (c VerificationCode) Expired() bool {
	return time.Now().After(c.CodeExp)
}

type Credentials struct {
	UserID   string
	Password string
	Token    Token
}

type AuthRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
