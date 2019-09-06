package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

var ErrNotAuthorized = errors.New("not authorized")

type ApiError struct {
	Code             string                    `json:"code"`
	Message          string                    `json:"message"`
	ValidationErrors []*models.ValidationError `json:"validation_errors,omitempty"`
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)
	w.WriteHeader(http.StatusInternalServerError)
	writeJsonResponse(w, &ApiError{
		Code:    "001",
		Message: "server error",
	})
}

func (app *application) accessTokenInvalid(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	writeJsonResponse(w, &ApiError{
		Code:    "002",
		Message: "access token invalid or expired",
	})
}

func (app *application) refreshTokenInvalid(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	writeJsonResponse(w, &ApiError{
		Code:    "003",
		Message: "refresh token invalid or expired",
	})
}

func (app *application) badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	writeJsonResponse(w, &ApiError{
		Code:    "004",
		Message: "can't read request body",
	})
}

func (app *application) validationError(w http.ResponseWriter, errs *models.ValidationErrors) {
	w.WriteHeader(http.StatusBadRequest)
	writeJsonResponse(w, &ApiError{
		Code:             "005",
		Message:          "request validation failed",
		ValidationErrors: errs.Errors,
	})
}

func (app *application) emailOrPasswordIncorrect(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	writeJsonResponse(w, &ApiError{
		Code:    "006",
		Message: "email or password incorrect",
	})
}

func (app *application) duplicatedEmail(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	writeJsonResponse(w, &ApiError{
		Code:    "007",
		Message: "user with this email already registered",
	})
}

func (app *application) deckNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	writeJsonResponse(w, &ApiError{
		Code:    "008",
		Message: "deck not found",
	})
}

func (app *application) flashcardNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	writeJsonResponse(w, &ApiError{
		Code:    "009",
		Message: "flashcard not found",
	})
}
