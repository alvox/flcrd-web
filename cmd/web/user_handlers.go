package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
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
	//todo: create token
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, userId)
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

func hashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
