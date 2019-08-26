package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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

func generateToken(userID string) (*string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
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

func validateToken(token string) error {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return appKey, nil
	})
	if err != nil {
		fmt.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return ErrNotAuthorized
		}
		return err
	}
	if !t.Valid {
		return ErrNotAuthorized
	}
	//if claims.UserID != userID {
	//	return ErrBadRequest
	//}
	return nil
}
