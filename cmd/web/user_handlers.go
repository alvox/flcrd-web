package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	auth := models.ParseAuthRequest(r)
	if auth == nil {
		app.badRequest(w)
		return
	}
	if errs := auth.Validate(true); errs.Present() {
		app.validationError(w, errs)
		return
	}
	existingUser, err := app.users.GetByEmail(auth.Email)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(w, err)
		return
	}
	if existingUser != nil {
		app.emailNotUnique(w)
		return
	}
	pwdHash, err := hashAndSalt(auth.Password)
	if err != nil {
		app.serverError(w, err)
	}
	credentials := NewCredentials(pwdHash)
	user := NewUser(*auth)
	userId, err := app.users.Create(&user, &credentials)
	if err != nil {
		app.serverError(w, err)
	}
	createdUser, err := app.users.Get(*userId)
	if err != nil {
		app.serverError(w, err)
	}
	accessToken, err := generateAccessToken(*userId)
	if err != nil {
		app.serverError(w, err)
	}
	createdUser.Token.AccessToken = *accessToken
	go app.sendConfirmation(createdUser.ID, createdUser.Name, createdUser.Email)
	reply(w, http.StatusCreated, createdUser)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("UserID").(string)
	user, err := app.users.GetProfile(id)
	if modelError(app, err, w, "user") {
		return
	}
	reply(w, http.StatusOK, user)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	user := models.ParseUser(r)
	if user == nil {
		app.badRequest(w)
		return
	}
	if errs := user.ValidateForUpdate(); errs.Present() {
		app.validationError(w, errs)
		return
	}
	userID := r.Context().Value("UserID").(string)
	existingUser, err := app.users.Get(userID)
	if modelError(app, err, w, "user") {
		return
	}
	if existingUser.Email != user.Email {
		err := app.deleteExistingCode(userID, w)
		if err != nil {
			return
		}
		app.sendConfirmation(userID, user.Name, user.Email)
		existingUser.Status = "PENDING"
	}
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	err = app.users.Update(existingUser)
	if err != nil {
		app.serverError(w, err)
	}
	reply(w, http.StatusOK, existingUser)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	auth := models.ParseAuthRequest(r)
	if auth == nil {
		app.badRequest(w)
		return
	}
	if errs := auth.Validate(false); errs.Present() {
		app.validationError(w, errs)
		return
	}
	existingUser, err := app.users.GetByEmail(auth.Email)
	if err != nil {
		if err == models.ErrNoRecord {
			app.emailOrPasswordIncorrect(w)
			return
		}
		app.serverError(w, err)
		return
	}
	credentials, err := app.users.GetCredentials(existingUser.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !checkPassword(credentials.Password, auth.Password) {
		app.emailOrPasswordIncorrect(w)
		return
	}
	if credentials.Token.RefreshTokenExpired() {
		credentials.Token.RefreshToken, credentials.Token.RefreshTokenExp = generateRefreshToken()
		err = app.users.UpdateRefreshToken(credentials)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	accessToken, err := generateAccessToken(existingUser.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	existingUser.Token = credentials.Token
	existingUser.Token.AccessToken = *accessToken
	reply(w, http.StatusOK, existingUser)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	token := models.ParseTokens(r)
	if token == nil {
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
	credentials, err := app.users.GetCredentials(userID)
	if err != nil {
		if err == models.ErrNoRecord {
			app.accessTokenInvalid(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if !validateRefreshToken(token.RefreshToken, credentials) {
		app.refreshTokenInvalid(w)
		return
	}
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		app.serverError(w, err)
	}
	token.AccessToken = *accessToken
	reply(w, http.StatusOK, token)
}

func (app *application) activate(w http.ResponseWriter, r *http.Request) {
	c := mux.Vars(r)["code"]
	code, err := app.verification.Get(c)
	if err == models.ErrNoRecord {
		app.verificationCodeInvalid(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	if code.Expired() {
		app.verificationCodeInvalid(w)
		return
	}
	u, err := app.users.Get(code.UserID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if u.Status != "PENDING" {
		app.verificationCodeInvalid(w)
		return
	}
	u.Status = "ACTIVE"
	err = app.users.Update(u)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.verification.Delete(*code)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, nil)
}

func (app *application) resendConfirmation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("UserID").(string)
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if user.Status != "PENDING" {
		app.invalidUserStatus(w)
		return
	}
	err = app.deleteExistingCode(userID, w)
	if err != nil {
		return
	}
	go app.sendConfirmation(user.ID, user.Name, user.Email)
	reply(w, http.StatusOK, nil)
}

func (app *application) deleteExistingCode(userID string, w http.ResponseWriter) error {
	code, err := app.verification.GetForUser(userID)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(w, err)
		return err
	}
	if code != nil {
		err = app.verification.Delete(*code)
		if err != nil {
			app.serverError(w, err)
			return err
		}
	}
	return nil
}

func (app *application) sendConfirmation(userID, userName, email string) {
	code := generateVerificationCode(userID)
	c, err := app.verification.Create(code)
	if err != nil {
		log.Error().Err(err)
		return
	}
	result, err := app.mailSender.SendConfirmation(email, userName, c)
	if err != nil {
		log.Error().Err(err)
		return
	}
	log.Info().Str("ID", result.Id).Str("Message", result.Message).Msg("Email sent")
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("UserID").(string)
	err := app.users.Delete(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, nil)
}

func NewCredentials(pwdHash string) models.Credentials {
	return models.Credentials{
		Password: pwdHash,
		Token:    NewTokens(),
	}
}

func NewUser(auth models.AuthRequest) models.User {
	return models.User{
		Name:   auth.Name,
		Email:  auth.Email,
		Status: "PENDING",
	}
}
