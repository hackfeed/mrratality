package userdb

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

type signedDetails struct {
	Email  string
	UserID string
	jwt.StandardClaims
}

var DB *mongo.Client
var UserCol *mongo.Collection

func ConnectDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		panic("Failed to create MongoDB client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic("Failed to connect to MongoDB")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("Failed to ping MongoDB")
	}

	UserCol = client.Database("mrr").Collection("user")
	DB = client
}

func generateTokens(email, userID string) (string, string, error) {
	var token, refreshToken string

	claims := &signedDetails{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	refreshClaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(4)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return token, refreshToken, err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))

	return token, refreshToken, err
}

func updateTokens(signedToken, signedRefreshToken, userID string) error {
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	return Update(updateObj, bson.M{"user_id": userID})
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	return string(bytes), err
}

func verifyPassword(userPass, providedPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
}
