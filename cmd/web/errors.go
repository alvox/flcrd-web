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

func (e ApiError) str() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
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
	handleError(app, w, http.StatusUnauthorized, &ApiError{
		Code:    "002",
		Message: "access token invalid or expired",
	})
}

func (app *application) refreshTokenInvalid(w http.ResponseWriter) {
	handleError(app, w, http.StatusUnauthorized, &ApiError{
		Code:    "003",
		Message: "refresh token invalid or expired",
	})
}

func (app *application) badRequest(w http.ResponseWriter) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:    "004",
		Message: "can't read request",
	})
}

func (app *application) validationError(w http.ResponseWriter, errs *models.ValidationErrors) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:             "005",
		Message:          "request validation failed",
		ValidationErrors: errs.Errors,
	})
}

func (app *application) emailOrPasswordIncorrect(w http.ResponseWriter) {
	handleError(app, w, http.StatusUnauthorized, &ApiError{
		Code:    "006",
		Message: "email or password incorrect",
	})
}

func (app *application) emailNotUnique(w http.ResponseWriter) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:    "007",
		Message: "user with this email already registered",
	})
}

func (app *application) deckNameNotUnique(w http.ResponseWriter) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:    "008",
		Message: "deck with this name already exists",
	})
}

func (app *application) notFound(w http.ResponseWriter, model string) {
	handleError(app, w, http.StatusNotFound, &ApiError{
		Code:    "009",
		Message: fmt.Sprintf("%s not found", model),
	})
}

func (app *application) verificationCodeInvalid(w http.ResponseWriter) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:    "010",
		Message: "verification code invalid or expired",
	})
}

func (app *application) invalidUserStatus(w http.ResponseWriter) {
	handleError(app, w, http.StatusBadRequest, &ApiError{
		Code:    "011",
		Message: "user status is not suitable for this operation",
	})
}

func handleError(app *application, w http.ResponseWriter, status int, e *ApiError) {
	w.WriteHeader(status)
	app.errorLog.Println(e.str())
	writeJsonResponse(w, e)
}
