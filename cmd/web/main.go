package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/pg"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type application struct {
	decks interface {
		Create(string, string) (*string, error)
		Get(string) (*models.Deck, error)
		GetAll() ([]*models.Deck, error)
		Find(string) (*models.Deck, error)
		Update(*models.Deck) error
		Delete(string) error
	}
	flashcards interface {
		Create(*models.Flashcard) (*string, error)
		Get(string, string) (*models.Flashcard, error)
		GetAll(string) ([]*models.Flashcard, error)
		Update(*models.Flashcard) error
		Delete(string, string) error
	}
}

func main() {
	db, err := connectDB("postgres://flcrd:flcrd@localhost/flcrd?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	app := &application{
		decks:      &pg.DeckModel{DB: db},
		flashcards: &pg.FlashcardModel{DB: db},
	}

	srv := &http.Server{
		Addr:         ":5000",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server started")
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func connectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
