package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	// Middleware
	router.Use(app.logRequest)
	// Decks
	router.HandleFunc("/v0/decks", app.createDeck).Methods("POST")
	router.HandleFunc("/v0/decks", app.getDecks).Methods("GET")
	router.HandleFunc("/v0/decks/{deckID}", app.getDeck).Methods("GET")
	router.HandleFunc("/v0/decks/{deckID}", app.updateDeck).Methods("PUT")
	router.HandleFunc("/v0/decks/{deckID}", app.deleteDeck).Methods("DELETE")
	// Flashcards
	router.HandleFunc("/v0/decks/{deckID}/flashcards", app.createFlashcard).Methods("POST")
	router.HandleFunc("/v0/decks/{deckID}/flashcards", app.getFlashcards).Methods("GET")
	router.HandleFunc("/v0/decks/{deckID}/flashcards/{flashcardID}", app.getFlashcard).Methods("GET")
	router.HandleFunc("/v0/decks/{deckID}/flashcards/{flashcardID}", app.updateFlashcard).Methods("PUT")
	router.HandleFunc("/v0/decks/{deckID}/flashcards/{flashcardID}", app.deleteFlashcard).Methods("DELETE")
	// Users
	router.HandleFunc("/v0/users/register", app.registerUser).Methods("POST")
	return router
}
