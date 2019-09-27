package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	// Middleware
	router.Use(app.logRequest)
	auth := alice.New(app.validateToken)
	// Decks
	router.Handle("/v0/decks", auth.ThenFunc(app.createDeck)).Methods("POST")
	router.Handle("/v0/decks", auth.ThenFunc(app.getDecksForUser)).Methods("GET")
	router.Handle("/v0/decks/{deckID}", auth.ThenFunc(app.getDeck)).Methods("GET")
	router.Handle("/v0/decks/{deckID}", auth.ThenFunc(app.updateDeck)).Methods("PUT")
	router.Handle("/v0/decks/{deckID}", auth.ThenFunc(app.deleteDeck)).Methods("DELETE")
	// Flashcards
	router.Handle("/v0/decks/{deckID}/flashcards", auth.ThenFunc(app.createFlashcard)).Methods("POST")
	router.Handle("/v0/decks/{deckID}/flashcards", auth.ThenFunc(app.getFlashcardsForUser)).Methods("GET")
	router.Handle("/v0/decks/{deckID}/flashcards/{flashcardID}", auth.ThenFunc(app.getFlashcard)).Methods("GET")
	router.Handle("/v0/decks/{deckID}/flashcards/{flashcardID}", auth.ThenFunc(app.updateFlashcard)).Methods("PUT")
	router.Handle("/v0/decks/{deckID}/flashcards/{flashcardID}", auth.ThenFunc(app.deleteFlashcard)).Methods("DELETE")
	// Users
	router.HandleFunc("/v0/users/register", app.registerUser).Methods("POST")
	router.HandleFunc("/v0/users/login", app.login).Methods("POST")
	router.HandleFunc("/v0/users/refresh", app.refreshToken).Methods("POST")
	router.HandleFunc("/v0/users/activate/{code}", app.activate).Methods("POST")
	router.Handle("/v0/users/code", auth.ThenFunc(app.resendConfirmation)).Methods("POST")
	// Public routes
	router.HandleFunc("/v0/public/decks", app.getPublicDecks).Methods("GET")
	router.HandleFunc("/v0/public/decks/{deckID}/flashcards", app.getPublicFlashcards).Methods("GET")

	return router
}
