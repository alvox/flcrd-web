package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDeckNotFound = errors.New("models: deck does not exist")
var ErrUniqueViolation = errors.New("models: unique violation")

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
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Email      string `json:"email"`
}

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Created    time.Time `json:"created"`
	ExternalID string    `json:"external_id"`
	Stats      Stats     `json:"stats,omitempty"`
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
