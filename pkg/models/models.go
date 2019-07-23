package models

import "time"

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
	Cards       []Flashcard `json:"cards"`
	Created     time.Time   `json:"created"`
}
