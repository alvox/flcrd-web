package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var authDomain string

func initAuth(domain string) {
	authDomain = domain
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	Id        string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
}

func (c CustomClaims) Valid() error {
	return nil
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) supplyUserId(next http.Handler) http.Handler {
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
		fmt.Println("TOKEN: ", token)
		externalID, err := getExternalUserIdFromToken(token)
		if err != nil {
			app.accessTokenInvalid(w)
			return
		}
		userID, err := app.resolveUserID(externalID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		fmt.Println("USER ID: ", userID)
		r.Header.Add("UserID", userID)
		next.ServeHTTP(w, r)
	})
}

func (app *application) validateToken(next http.Handler) http.Handler {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			aud := "https://flashcards.rocks"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				fmt.Println("invalid audience")
				return token, errors.New("invalid audience")
			}
			iss := fmt.Sprintf("https://%s/", authDomain)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				fmt.Println("invalid issuer")
				return token, errors.New("invalid issuer")
			}
			cert, err := getPemCert(token)
			if err != nil {
				fmt.Println("invalid cert")
				return token, errors.New("invalid cert")
			}
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	}).Handler(next)
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", authDomain))
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()
	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}
	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}
	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}
	return cert, nil
}

func getExternalUserIdFromToken(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", ErrNotAuthorized
	}
	claims, ok := t.Claims.(*CustomClaims)
	if !ok {
		fmt.Println("can't parse claims")
		return "", ErrNotAuthorized
	}
	return claims.Subject, nil
}

func (app *application) resolveUserID(externalID string) (string, error) {
	user, err := app.users.Get(externalID)
	if err != nil {
		if err == models.ErrNoRecord {
			// todo: get user with email from auth0
			userID, err := app.users.Create(&models.User{
				ExternalID: externalID,
			})
			if err != nil {
				return "", err
			}
			return *userID, nil
		}
		return "", err
	}
	return user.ID, nil
}
