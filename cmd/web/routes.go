package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v0/decks", app.createDeck).Methods("POST")
	router.HandleFunc("/v0/decks/{deckID}", app.getDeck).Methods("GET")
	router.HandleFunc("/v0/decks", app.updateDeck).Methods("PUT")
	router.HandleFunc("/v0/decks/{deckID}", app.deleteDeck).Methods("DELETE")
	return router
}
