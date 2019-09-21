package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (app *application) createDeck(w http.ResponseWriter, r *http.Request) {
	deck, valid := models.ParseDeck(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := deck.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	id, err := app.decks.Create(deck.Name, deck.Description, r.Header.Get("UserID"), deck.Public)
	if err != nil {
		app.serverError(w, err)
		return
	}
	deck, err = app.decks.Get(*id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusCreated, deck)
}

func (app *application) getPublicDecks(w http.ResponseWriter, r *http.Request) {
	terms := extractSearchTerms(r)
	var decks []*models.Deck
	var err error
	if terms == nil {
		decks, err = app.decks.GetPublic()
	} else {
		decks, err = app.decks.Search(terms)
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, decks)
}

func extractSearchTerms(r *http.Request) []string {
	query := r.URL.Query()
	if len(query) == 0 {
		return nil
	}
	q := query.Get("q")
	if len(q) == 0 {
		return nil
	}
	return strings.Split(q, ",")
}

func (app *application) getDecksForUser(w http.ResponseWriter, r *http.Request) {
	decks, err := app.decks.GetForUser(r.Header.Get("UserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, decks)
}

func (app *application) getDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	deck, err := app.decks.Get(id)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, deck)
}

func (app *application) updateDeck(w http.ResponseWriter, r *http.Request) {
	deck, valid := models.ParseDeck(r)
	if !valid {
		app.badRequest(w)
		return
	}
	deckID := mux.Vars(r)["deckID"]
	if errs := deck.ValidateWithID(deckID); errs.Present() {
		app.validationError(w, errs)
		return
	}
	err := app.decks.Update(deck)
	if err != nil {
		if err == models.ErrNoRecord {
			app.deckNotFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	deck, err = app.decks.Get(deck.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, deck)
}

func (app *application) deleteDeck(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(id)
	if err == models.ErrNoRecord {
		app.deckNotFound(w)
		return
	}
	err = app.decks.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusNoContent, nil)
}
