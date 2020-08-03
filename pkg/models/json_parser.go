package models

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
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
