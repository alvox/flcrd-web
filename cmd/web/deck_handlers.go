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
	id, err := app.decks.Create(deck.Name, deck.Description, deck.Private)
	if err != nil {
		app.serverError(w, err)
		return
	}
	deck, err = app.decks.Get(*id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, deck)
}

func (app *application) getDecks(w http.ResponseWriter, r *http.Request) {
	decks, err := app.decks.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, decks)
}

func (app *application) getDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
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
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	}
	err = app.decks.Update(deck)
	if err != nil {
		app.serverError(w, err)
		return
	}
	deck, err = app.decks.Get(deck.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, deck)
}

func (app *application) deleteDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	}
	err = app.decks.Delete(id)
	if err != nil {
		app.serverError(w, err)
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
