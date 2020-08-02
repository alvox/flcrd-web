package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"alexanderpopov.me/flcrd/pkg/models/pg"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	decks interface {
		Create(string, string, string, bool) (*string, error)
		Get(string) (*models.Deck, error)
		GetForUser(string) ([]*models.Deck, error)
		GetPublic(int, int) ([]*models.Deck, int, error)
		Update(*models.Deck) error
		Delete(string) error
		Search([]string) ([]*models.Deck, error)
	}
	flashcards interface {
		Create(*models.Flashcard) (*string, error)
		Get(string, string) (*models.Flashcard, error)
		GetForUser(string, string) ([]*models.Flashcard, error)
		GetPublic(string) ([]*models.Flashcard, error)
		Update(*models.Flashcard) error
		Delete(string, string) error
	}
	users interface {
		Create(user *models.User, credentials *models.Credentials) (*string, error)
		Get(string) (*models.User, error)
		GetProfile(string) (*models.User, error)
		GetByEmail(string) (*models.User, error)
		Update(user *models.User) error
		Delete(string) error
		GetCredentials(string) (*models.Credentials, error)
		UpdateRefreshToken(credentials *models.Credentials) error
	}
	verification interface {
		Create(code models.VerificationCode) (string, error)
		Get(string) (*models.VerificationCode, error)
		GetForUser(string) (*models.VerificationCode, error)
		Delete(code models.VerificationCode) error
	}
	mailSender interface {
		SendConfirmation(string, string, string) (*SendMessageResponse, error)
	}
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	port := flag.String("port", ":5000", "Application port")
	dsn := flag.String("dsn", "postgres://flcrd:flcrd@localhost/flcrd?sslmode=disable", "Postgres data source")
	key := flag.String("appkey", "test-key", "Application key")
	mailUrl := flag.String("mail_api_url", "", "URL to the email service")
	mailKey := flag.String("mail_api_key", "", "API key to the email service")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := connectDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		decks:        &pg.DeckModel{DB: db},
		flashcards:   &pg.FlashcardModel{DB: db},
		users:        &pg.UserModel{DB: db},
		verification: &pg.VerificationModel{DB: db},
		mailSender: &MailSender{
			baseUrl: *mailUrl,
			apiKey:  *mailKey,
		},
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:         *port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	initJwt(*key)
	infoLog.Printf("Starting server on %s port", *port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func connectDB(dsn string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
