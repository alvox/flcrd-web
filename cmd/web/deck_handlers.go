package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
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
	id, err := app.decks.Create(deck.Name, deck.Description, r.Header.Get("UserID"), deck.Private)
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
	decks, err := app.decks.GetPublic()
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
