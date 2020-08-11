package models

import (
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ParseUser(r *http.Request, sanitizer *bluemonday.Policy) *User {
	var user User
	if ok := parse(r, &user); !ok {
		return nil
	}
	user.Name = sanitizer.Sanitize(user.Name)
	return &user
}

func ParseTokens(r *http.Request) *Token {
	var token Token
	if ok := parse(r, &token); !ok {
		return nil
	}
	return &token
}

func ParseDeck(r *http.Request, sanitizer *bluemonday.Policy) *Deck {
	var deck Deck
	if ok := parse(r, &deck); !ok {
		return nil
	}
	deck.Name = sanitizer.Sanitize(deck.Name)
	deck.Description = sanitizer.Sanitize(deck.Description)
	return &deck
}

func ParseFlashcard(r *http.Request, sanitizer *bluemonday.Policy) *Flashcard {
	var flashcard Flashcard
	if ok := parse(r, &flashcard); !ok {
		return nil
	}
	flashcard.Front = sanitizer.Sanitize(flashcard.Front)
	flashcard.Rear = sanitizer.Sanitize(flashcard.Rear)
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
		log.Error().Msg("request body is empty")
		return false
	}
	err := json.NewDecoder(r.Body).Decode(i)
	if err != nil {
		log.Error().Err(err).Msg("can't parse request body")
		return false
	}
	return true
}
