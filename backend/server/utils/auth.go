package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	userdb "backend/db/user"
)

type SignedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

func GenerateTokens(email, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:  email,
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(4)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateTokens(signedToken, signedRefreshToken, userId string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userdb.UserCol.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		return err
	}

	return nil
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

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(userPass, providedPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
	if err != nil {
		return err
	}

	return nil
}
