package main

import (
	"net/http"
	"strings"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			app.accessTokenInvalid(w)
			return
		}
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			app.accessTokenInvalid(w)
			return
		}
		token := authHeader[1]
		userID, err := validateAuthToken(token, false)
		if err != nil {
			app.accessTokenInvalid(w)
			return
		}
		r.Header.Add("UserID", userID)
		next.ServeHTTP(w, r)
	})
}
