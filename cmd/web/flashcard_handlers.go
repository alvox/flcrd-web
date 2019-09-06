package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcard, valid := readFlashcard(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := flashcard.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcard.DeckID = deckID
	flashcardID, err := app.flashcards.Create(flashcard)
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcard, err = app.flashcards.Get(deckID, *flashcardID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, flashcard)
}

func (app *application) getFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	flashcardID := mux.Vars(r)["flashcardID"]
	flashcard, err := app.flashcards.Get(deckID, flashcardID)
	if err == models.ErrNoRecord {
		app.flashcardNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcard)
}

func (app *application) getPublicFlashcards(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcards, err := app.flashcards.GetPublic(deckID)
	if err != nil {
		app.serverError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcards)
}

func (app *application) getFlashcardsForUser(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcards, err := app.flashcards.GetForUser(deckID, r.Header.Get("UserID"))
	if err != nil {
		app.serverError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcards)
}

func (app *application) updateFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcard, valid := readFlashcard(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := flashcard.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcardID := mux.Vars(r)["flashcardID"]
	flashcard.ID = flashcardID
	err = app.flashcards.Update(flashcard)
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcard, err = app.flashcards.Get(flashcard.DeckID, flashcard.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcard)
}

func (app *application) deleteFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	flashcardID := mux.Vars(r)["flashcardID"]
	err = app.flashcards.Delete(deckID, flashcardID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readFlashcard(r *http.Request) (*models.Flashcard, bool) {
	if r.Body == nil {
		return nil, false
	}
	var flashcard models.Flashcard
	err := json.NewDecoder(r.Body).Decode(&flashcard)
	if err != nil {
		return nil, false
	}
	return &flashcard, true
}
