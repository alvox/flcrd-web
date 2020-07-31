package models

import (
	"encoding/json"
	"net/http"
)

func ParseUser(r *http.Request) *User {
	var user User
	if ok := parse(r, &user); !ok {
		return nil
	}
	return &user
}

func ParseTokens(r *http.Request) *Token {
	var token Token
	if ok := parse(r, &token); !ok {
		return nil
	}
	return &token
}

func ParseDeck(r *http.Request) *Deck {
	var deck Deck
	if ok := parse(r, &deck); !ok {
		return nil
	}
	return &deck
}

func ParseFlashcard(r *http.Request) *Flashcard {
	var flashcard Flashcard
	if ok := parse(r, &flashcard); !ok {
		return nil
	}
	return &flashcard
}

func ParseAuthRequest(r *http.Request) *AuthRequest {
	var authRequest AuthRequest
	if ok := parse(r, &authRequest); !ok {
		return nil
	}
	return &authRequest
}

func parse(r *http.Request, i interface{}) bool {
	if r.Body == nil {
		return false
	}
	err := json.NewDecoder(r.Body).Decode(i)
	if err != nil {
		return false
	}
	return true
}
