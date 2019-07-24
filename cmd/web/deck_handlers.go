package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createDeck(w http.ResponseWriter, r *http.Request) {
	deck := readDeck(w, r)
	if deck == nil {
		return
	}
	id, err := app.decks.Create(deck.Name, deck.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deck, err = app.decks.Get(*id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, deck)
}

func (app *application) getDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, deck)
}

func (app *application) updateDeck(w http.ResponseWriter, r *http.Request) {
	deck := readDeck(w, r)
	if deck == nil {
		return
	}
	err := app.decks.Update(deck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deck, err = app.decks.Get(deck.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, deck)
}

func (app *application) deleteDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	err := app.decks.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}) {
	out, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readDeck(w http.ResponseWriter, r *http.Request) *models.Deck {
	if r.Body == nil {
		http.Error(w, "Please, send a request body", http.StatusBadRequest)
	}
	var deck models.Deck
	err := json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &deck
}
