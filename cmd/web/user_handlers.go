package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
//	user := models.ParseUser(r)
//	if user == nil {
//		app.badRequest(w)
//		return
//	}
//	if errs := user.Validate(true); errs.Present() {
//		app.validationError(w, errs)
//		return
//	}
//	existingUser, err := app.users.GetByEmail(user.Email)
//	if err != nil && err != models.ErrNoRecord {
//		app.serverError(w, err)
//		return
//	}
//	if existingUser != nil {
//		app.emailNotUnique(w)
//		return
//	}
//	pwdHash, err := hashAndSalt(user.Password)
//	if err != nil {
//		app.serverError(w, err)
//	}
//	user.Password = pwdHash
//	user.Token.RefreshToken, user.Token.RefreshTokenExp = generateRefreshToken()
//	user.Status = "PENDING"
//	userId, err := app.users.Create(user)
//	if err != nil {
//		app.serverError(w, err)
//	}
//	user, err = app.users.Get(*userId)
//	if err != nil {
//		app.serverError(w, err)
//	}
//	accessToken, err := generateAccessToken(*userId)
//	if err != nil {
//		app.serverError(w, err)
//	}
//	user.Token.AccessToken = *accessToken
//	user.Password = ""
//	go app.sendConfirmation(user.ID, user.Name, user.Email)
//	reply(w, http.StatusCreated, user)
//}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("UserID")
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
	userID := r.Header.Get("UserID")
	existingUser, err := app.users.Get(userID)
	if modelError(app, err, w, "user") {
		return
	}
	if existingUser.Email != user.Email {
		err := app.deleteExistingCode(userID, w)
		if err != nil {
			return
		}
		//app.sendConfirmation(userID, user.Name, user.Email)
	}
	existingUser.Email = user.Email
	err = app.users.Update(existingUser)
	if err != nil {
		app.serverError(w, err)
	}
	reply(w, http.StatusOK, existingUser)
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
	//userID := r.Header.Get("UserID")
	//user, err := app.users.Get(userID)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//err = app.deleteExistingCode(userID, w)
	//if err != nil {
	//	return
	//}
	//go app.sendConfirmation(user.ID, user.Name, user.Email)
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
		app.errorLog.Println(err.Error())
		return
	}
	result, err := app.mailSender.SendConfirmation(email, userName, c)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println(result)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	err := app.users.Delete(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	reply(w, http.StatusOK, nil)
}

func generateVerificationCode(userID string) models.VerificationCode {
	return models.VerificationCode{
		UserID:  userID,
		Code:    uniuri.NewLen(40),
		CodeExp: time.Now().Add(24 * time.Hour * 5),
	}
}
