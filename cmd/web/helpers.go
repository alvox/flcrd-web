package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

var ErrBadRequest = errors.New("bad request")
var ErrNotAuthorized = errors.New("not authorized")
var ErrEmailFormatInvalid = errors.New("invalid email format")
var ErrUserAlreadyExists = errors.New("user with this email already exists")

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) notAuthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
