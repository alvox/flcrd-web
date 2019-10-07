package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) createDeck(w http.ResponseWriter, r *http.Request) {
	deck := models.ParseDeck(r)
	if deck == nil {
		app.badRequest(w)
		return
	}
	if errs := deck.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	id, err := app.decks.Create(deck.Name, deck.Description, r.Header.Get("UserID"), deck.Public)
	if err == models.ErrUniqueViolation {
		app.deckNameNotUnique(w)
		return
	}
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
	page, offset, limit, err := extractPaging(r)
	if err != nil {
		app.badRequest(w)
		return
	}
	var decks []*models.Deck
	var total int
	if terms == nil {
		decks, total, err = app.decks.GetPublic(offset, limit)
		addLinkHeader(w, page, limit, total)
	} else {
		decks, err = app.decks.Search(terms)
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, decks)
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
	if modelError(app, err, w, "deck") {
		return
	}
	reply(w, http.StatusOK, deck)
}

func (app *application) updateDeck(w http.ResponseWriter, r *http.Request) {
	deck := models.ParseDeck(r)
	if deck == nil {
		app.badRequest(w)
		return
	}
	deckID := mux.Vars(r)["deckID"]
	if errs := deck.ValidateWithID(deckID); errs.Present() {
		app.validationError(w, errs)
		return
	}
	err := app.decks.Update(deck)
	if err == models.ErrUniqueViolation {
		app.deckNameNotUnique(w)
		return
	}
	if modelError(app, err, w, "deck") {
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
	err := app.decks.Delete(id)
	if modelError(app, err, w, "deck") {
		return
	}
	reply(w, http.StatusNoContent, nil)
}
