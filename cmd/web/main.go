package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/pg"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type application struct {
	decks interface {
		Create(string, string) (*string, error)
		Find(string) (*models.Deck, error)
	}
}

func main() {
	db, err := connectDB("postgres://flcrd:flcrd@localhost/flcrd?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	app := &application{
		decks: &pg.DeckModel{DB: db},
	}

	id, err := app.decks.Create("My d", "This is my new deck")
	if err != nil {
		log.Fatal(err)
	}

	deck, err := app.decks.Find("My d")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*deck)
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
