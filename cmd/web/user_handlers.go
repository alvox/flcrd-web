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
	if errs := user.Validate(true); len(errs) > 0 {
		err := map[string]interface{}{"validationError": errs}
		w.WriteHeader(http.StatusBadRequest)
		writeJsonResponse(w, err)
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
	user.Token.RefreshToken, user.Token.RefreshTokenExp = generateRefreshToken()
	userId, err := app.users.Create(user)
	if err != nil {
		app.serverError(w, err)
	}
	user, err = app.users.GetByEmail(user.Email)
	if err != nil {
		app.serverError(w, err)
	}
	authToken, err := generateAuthToken(*userId)
	if err != nil {
		app.serverError(w, err)
	}
	user.Token.AuthToken = *authToken
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, user)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user := readUser(w, r)
	if user == nil {
		return
	}
	if errs := user.Validate(false); len(errs) > 0 {
		err := map[string]interface{}{"validationError": errs}
		w.WriteHeader(http.StatusBadRequest)
		writeJsonResponse(w, err)
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
	existingUser.Token.RefreshToken, existingUser.Token.RefreshTokenExp = generateRefreshToken()
	err = app.users.UpdateRefreshToken(existingUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	authToken, err := generateAuthToken(existingUser.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	existingUser.Token.AuthToken = *authToken
	existingUser.Password = ""
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, existingUser)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	token := readTokens(w, r)
	if token == nil {
		return
	}
	if errs := token.Validate(); len(errs) > 0 {
		err := map[string]interface{}{"validationError": errs}
		w.WriteHeader(http.StatusBadRequest)
		writeJsonResponse(w, err)
		return
	}
	userID, err := validateAuthToken(token.AuthToken, true)
	if err != nil {
		app.notAuthorized(w)
		return
	}
	user, err := app.users.Get(userID)
	if err != nil {
		if err == models.ErrNoRecord {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if !validateRefreshToken(token.RefreshToken, user) {
		app.notAuthorized(w)
		return
	}
	authToken, err := generateAuthToken(userID)
	if err != nil {
		app.serverError(w, err)
	}
	token.AuthToken = *authToken
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, token)
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

func readTokens(w http.ResponseWriter, r *http.Request) *models.Token {
	if r.Body == nil {
		http.Error(w, "Please, send a request body", http.StatusBadRequest)
		return nil
	}
	var token models.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	return &token
}
