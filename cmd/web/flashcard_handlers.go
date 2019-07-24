package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcard := readFlashcard(w, r)
	if flashcard == nil {
		return
	}
	deckID := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(deckID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deck == nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	flashcard.DeckID = deckID
	flashcardID, err := app.flashcards.Create(flashcard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	flashcard, err = app.flashcards.Get(deckID, *flashcardID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, flashcard)
}

func (app *application) getFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	flashcardID := mux.Vars(r)["flashcardID"]
	flashcard, err := app.flashcards.Get(deckID, flashcardID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcard)
}

func (app *application) updateFlashcard(w http.ResponseWriter, r *http.Request) {
	flashcard := readFlashcard(w, r)
	if flashcard == nil {
		return
	}
	deckID := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(deckID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deck == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	flashcardID := mux.Vars(r)["flashcardID"]
	flashcard.ID = flashcardID
	err = app.flashcards.Update(flashcard)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	flashcard, err = app.flashcards.Get(flashcard.DeckID, flashcard.ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(flashcard)
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, flashcard)
}

func (app *application) deleteFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(deckID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deck == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	flashcardID := mux.Vars(r)["flashcardID"]
	err = app.flashcards.Delete(deckID, flashcardID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readFlashcard(w http.ResponseWriter, r *http.Request) *models.Flashcard {
	if r.Body == nil {
		http.Error(w, "Please, send a request body", http.StatusBadRequest)
		return nil
	}
	var flashcard models.Flashcard
	err := json.NewDecoder(r.Body).Decode(&flashcard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &flashcard
}
