package main

import (
	"alexanderpopov.me/flcrd/pkg/models"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var appKey []byte

func initJwt(key string) {
	appKey = []byte(key)
}

type Claims struct {
	UserID string
	jwt.StandardClaims
}

func hashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hash string, pwd string) bool {
	byteHash := []byte(hash)
	bytePwd := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}

func generateAuthToken(userID string) (*string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(appKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func validateAuthToken(token string, couldBeExpired bool) (string, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return appKey, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		if err == jwt.ErrSignatureInvalid {
			return "", ErrNotAuthorized
		}
		if couldBeExpired && strings.Contains(err.Error(), "token is expired") {
			return claims.UserID, nil
		}
		return "", err
	}
	if !t.Valid {
		return "", ErrNotAuthorized
	}
	return claims.UserID, nil
}

func generateRefreshToken() (string, time.Time) {
	t := uniuri.NewLen(40)
	exp := time.Now().Add(24 * time.Hour * 30)
	return t, exp
}

func validateRefreshToken(token string, user *models.User) bool {
	if token != user.Token.RefreshToken {
		return false
	}
	if time.Now().After(user.Token.RefreshTokenExp) {
		return false
	}
	return true
}
