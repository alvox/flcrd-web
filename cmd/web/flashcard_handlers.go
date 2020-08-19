package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"path/filepath"
)

// todo: check max image size
const maxFileSize = int64(10240000) // 10MB

func (app *application) createFlashcard(w http.ResponseWriter, r *http.Request) {
	f := &models.Flashcard{
		Front:     app.sanitizer.Sanitize(r.FormValue("front")),
		FrontType: app.sanitizer.Sanitize(r.FormValue("front_type")),
		Rear:      app.sanitizer.Sanitize(r.FormValue("rear")),
		RearType:  app.sanitizer.Sanitize(r.FormValue("rear_type")),
	}
	if errs := f.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if modelError(app, err, w, "deck") {
		return
	}
	f.DeckID = deckID
	f.ID = uuid.NewV4().String()

	if f.FrontType == "IMAGE" {
		fileName, err := uploadImageToS3("front", f.ID, r, app.awsCredentials)
		if err != nil {
			app.serverError(w, err)
			return
		}
		f.Front = fileName
	}
	if f.RearType == "IMAGE" {
		fileName, err := uploadImageToS3("rear", f.ID, r, app.awsCredentials)
		if err != nil {
			app.serverError(w, err)
			return
		}
		f.Rear = fileName
	}

	flashcardID, err := app.flashcards.Create(f)
	if err != nil {
		app.serverError(w, err)
		return
	}
	f, err = app.flashcards.Get(deckID, *flashcardID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusCreated, f)
}

func (app *application) getFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	flashcardID := mux.Vars(r)["flashcardID"]
	flashcard, err := app.flashcards.Get(deckID, flashcardID)
	if modelError(app, err, w, "flashcard") {
		return
	}
	reply(w, http.StatusOK, flashcard)
}

func (app *application) getPublicFlashcards(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if modelError(app, err, w, "deck") {
		return
	}
	flashcards, err := app.flashcards.GetPublic(deckID)
	if err != nil {
		app.serverError(w, err)
	}
	reply(w, http.StatusOK, flashcards)
}

func (app *application) getFlashcardsForUser(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if modelError(app, err, w, "deck") {
		return
	}
	flashcards, err := app.flashcards.GetForUser(deckID, r.Context().Value("UserID").(string))
	if err != nil {
		app.serverError(w, err)
	}
	reply(w, http.StatusOK, flashcards)
}

func (app *application) updateFlashcard(w http.ResponseWriter, r *http.Request) {
	f := &models.Flashcard{
		Front:     app.sanitizer.Sanitize(r.FormValue("front")),
		FrontType: app.sanitizer.Sanitize(r.FormValue("front_type")),
		Rear:      app.sanitizer.Sanitize(r.FormValue("rear")),
		RearType:  app.sanitizer.Sanitize(r.FormValue("rear_type")),
	}
	if errs := f.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	f.DeckID = mux.Vars(r)["deckID"]
	_, err := app.decks.Get(f.DeckID)
	if modelError(app, err, w, "deck") {
		return
	}
	f.ID = mux.Vars(r)["flashcardID"]

	if f.FrontType == "IMAGE" {
		fileName, err := uploadImageToS3("front", f.ID, r, app.awsCredentials)
		if err != nil {
			app.serverError(w, err)
			return
		}
		f.Front = fileName
	}
	if f.RearType == "IMAGE" {
		fileName, err := uploadImageToS3("rear", f.ID, r, app.awsCredentials)
		if err != nil {
			app.serverError(w, err)
			return
		}
		f.Rear = fileName
	}

	err = app.flashcards.Update(f)
	if err != nil {
		app.serverError(w, err)
		return
	}
	f, err = app.flashcards.Get(f.DeckID, f.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, f)
}

func (app *application) deleteFlashcard(w http.ResponseWriter, r *http.Request) {
	deckID := mux.Vars(r)["deckID"]
	_, err := app.decks.Get(deckID)
	if modelError(app, err, w, "deck") {
		return
	}
	flashcardID := mux.Vars(r)["flashcardID"]
	err = app.flashcards.Delete(deckID, flashcardID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusNoContent, nil)
}

// todo: resize images if they're bigger than 1MB
func uploadImageToS3(side, cardID string, r *http.Request, c *credentials.Credentials) (string, error) {
	file, fileHeader, err := r.FormFile(fmt.Sprintf("%s_image", side))
	if err != nil {
		return "", err
	}
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: c,
	})
	if err != nil {
		return "", err
	}

	size := fileHeader.Size
	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("images/%s/%s%s", cardID, side, filepath.Ext(fileHeader.Filename))

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String("flcrd-img-orig"),
		Key:    aws.String(fileName),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		return "", err
	}
	return fileName, nil
}
