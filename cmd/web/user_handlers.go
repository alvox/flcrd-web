package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"net/http"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	user := readUser(w, r)
	if user == nil {
		return
	}
	err := validate(user)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	existingUser, err := app.users.GetByEmail(user.Email)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(w, err)
		return
	}
	if existingUser != nil {
		// log user already exists
		app.clientError(w, http.StatusBadRequest)
		return
	}
	pwdHash, err := hashAndSalt(user.Password)
	if err != nil {
		app.serverError(w, err)
	}
	user.Password = pwdHash
	userId, err := app.users.Create(user.Name, user.Email, user.Password)
	if err != nil {
		app.serverError(w, err)
	}
	user, err = app.users.GetByEmail(user.Email)
	if err != nil {
		app.serverError(w, err)
	}
	token, err := generateToken(*userId)
	if err != nil {
		app.serverError(w, err)
	}
	tokens := models.Token{
		AuthToken: *token,
	}
	user.Token = tokens
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, user)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user := readUser(w, r)
	if user == nil {
		return
	}
	err := validate(user)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	existingUser, err := app.users.GetByEmail(user.Email)
	if err != nil {
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if !checkPassword(existingUser.Password, user.Password) {
		app.notAuthorized(w)
		return
	}
	token, err := generateToken(existingUser.ID)
	if err != nil {
		app.serverError(w, err)
	}
	tokens := models.Token{
		AuthToken: *token,
	}
	existingUser.Token = tokens
	existingUser.Password = ""
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, existingUser)
}

func readUser(w http.ResponseWriter, r *http.Request) *models.User {
	if r.Body == nil {
		http.Error(w, "Please, send a request body", http.StatusBadRequest)
		return nil
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &user
}

func validate(user *models.User) error {
	err := validateEmailFormat(user.Email)
	if err != nil {
		return err
	}
	err = validateNotEmpty(user)
	if err != nil {
		return err
	}
	return nil
}

func validateNotEmpty(user *models.User) error {
	return nil // todo: implement
}

func validateEmailFormat(email string) error {
	return nil // todo: implement
}
