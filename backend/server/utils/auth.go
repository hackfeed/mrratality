package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

func GetExpirationTime(signedToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return 0, errors.New("token is invalid")
	}

	return claims.ExpiresAt, nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return &SignedDetails{}, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return &SignedDetails{}, errors.New("token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return &SignedDetails{}, errors.New("token is expired")
	}

	return claims, nil
}
