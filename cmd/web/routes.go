package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v0/decks", app.createDeck).Methods("POST")
	return router
}

func (app *application) createDeck(w http.ResponseWriter, r *http.Request) {
	deck := readJSON(w, r)
	if deck == nil {
		return
	}
	id, err := app.decks.Create(deck.Name, deck.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	deck, err = app.decks.Get(*id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	writeJsonResponse(w, deck)
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}) {
	out, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func readJSON(w http.ResponseWriter, r *http.Request) *models.Deck {
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
