package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/pg"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	decks interface {
		Create(string, string) (*string, error)
		Get(string) (*models.Deck, error)
		GetAll() ([]*models.Deck, error)
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
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	port := flag.String("port", ":5000", "Application port")
	dsn := flag.String("dsn", "postgres://flcrd:flcrd@flcrd-test-db/flcrd?sslmode=disable", "Postgres data source")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := connectDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		decks:      &pg.DeckModel{DB: db},
		flashcards: &pg.FlashcardModel{DB: db},
		infoLog:    infoLog,
		errorLog:   errorLog,
	}

	srv := &http.Server{
		Addr:         *port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s port", *port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
