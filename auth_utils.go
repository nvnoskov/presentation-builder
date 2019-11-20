package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gosimple/slug"
	"github.com/savsgio/go-logger"
)

var jwtSignKey = []byte("Pres-Builder-JWT")

type userCredential struct {
	Email []byte `json:"email"`
	Slug  string `json:"slug"`
	jwt.StandardClaims
}

func generateToken(email []byte) (userCredential, string, time.Time) {
	logger.Debugf("Create new token for user %s", email)

	expireAt := time.Now().Add(30 * time.Minute)

	// Embed User information to `token`
	creds := userCredential{
		Email: email,
		Slug:  slug.Make(string(email)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, creds)

	// token -> string. Only server knows the secret.
	tokenString, err := newToken.SignedString(jwtSignKey)
	if err != nil {
		logger.Error(err)
	}

	return creds, tokenString, expireAt
}

func validateToken(requestToken string) (*jwt.Token, *userCredential, error) {
	logger.Debug("Validating token...")

	user := &userCredential{}
	token, err := jwt.ParseWithClaims(requestToken, user, func(token *jwt.Token) (interface{}, error) {
		return jwtSignKey, nil
	})

	return token, user, err
}
