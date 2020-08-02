package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("remote_address", r.RemoteAddr).
			Str("proto", r.Proto).
			Str("method", r.Method).
			Str("uri", r.URL.RequestURI()).
			Msg("REQUEST")
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
		userID, err := validateAccessToken(token, false)
		if err != nil {
			app.accessTokenInvalid(w)
			return
		}
		ctx := context.WithValue(r.Context(), "UserID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
