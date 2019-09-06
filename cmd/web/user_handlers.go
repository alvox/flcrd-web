package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"net/http"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	user, valid := readUser(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := user.Validate(true); errs.Present() {
		app.validationError(w, errs)
		return
	}
	existingUser, err := app.users.GetByEmail(user.Email)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(w, err)
		return
	}
	if existingUser != nil {
		// log user already exists
		app.duplicatedEmail(w)
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
	accessToken, err := generateAccessToken(*userId)
	if err != nil {
		app.serverError(w, err)
	}
	user.Token.AccessToken = *accessToken
	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, user)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user, valid := readUser(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := user.Validate(false); errs.Present() {
		app.validationError(w, errs)
		return
	}
	existingUser, err := app.users.GetByEmail(user.Email)
	if err != nil {
		if err == models.ErrNoRecord {
			app.emailOrPasswordIncorrect(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if !checkPassword(existingUser.Password, user.Password) {
		app.emailOrPasswordIncorrect(w)
		return
	}
	existingUser.Token.RefreshToken, existingUser.Token.RefreshTokenExp = generateRefreshToken()
	err = app.users.UpdateRefreshToken(existingUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	accessToken, err := generateAccessToken(existingUser.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	existingUser.Token.AccessToken = *accessToken
	existingUser.Password = ""
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, existingUser)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	token, valid := readTokens(r)
	if !valid {
		app.badRequest(w)
		return
	}
	if errs := token.Validate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	userID, err := validateAccessToken(token.AccessToken, true)
	if err != nil {
		app.accessTokenInvalid(w)
		return
	}
	user, err := app.users.Get(userID)
	if err != nil {
		if err == models.ErrNoRecord {
			app.accessTokenInvalid(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if !validateRefreshToken(token.RefreshToken, user) {
		app.refreshTokenInvalid(w)
		return
	}
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		app.serverError(w, err)
	}
	token.AccessToken = *accessToken
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, token)
}

func readUser(r *http.Request) (*models.User, bool) {
	if r.Body == nil {
		return nil, false
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func readTokens(r *http.Request) (*models.Token, bool) {
	if r.Body == nil {
		return nil, false
	}
	var token models.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		return nil, false
	}
	return &token, true
}
