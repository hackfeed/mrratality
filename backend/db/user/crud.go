package userdb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUp(email, password string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := UserCol.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return User{}, errors.New("Failed to check email")
	}

	if count > 0 {
		return User{}, errors.New("This email is already taken")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return User{}, errors.New("Failed to hash password")
	}

	var user User

	user.Email = &email
	user.Password = &hashedPassword
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, _ := generateTokens(*user.Email, user.UserID)
	user.Token = &token
	user.RefreshToken = &refreshToken
	user.Files = []map[string]interface{}{}

	_, err = UserCol.InsertOne(ctx, user)
	if err != nil {
		return User{}, errors.New("Failed to create user")
	}

	return user, nil
}

func Login(email, password string) (User, error) {
	user, err := Read(bson.M{"email": email})
	if err != nil {
		return User{}, errors.New("User with given email doesn't exist")
	}

	err = verifyPassword(password, *user.Password)
	if err != nil {
		return User{}, errors.New("Password is incorrect")
	}

	token, refreshToken, err := generateTokens(*user.Email, user.UserID)
	if err != nil {
		return User{}, errors.New("Failed to generate tokens")
	}

	err = updateTokens(token, refreshToken, user.UserID)
	if err != nil {
		return User{}, errors.New("Failed to update tokens")
	}

	return user, nil
}

func Read(filter bson.M) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user User

	err := UserCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func Update(obj primitive.D, filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserCol.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", obj},
		},
		&opt,
	)

	if err != nil {
		return err
	}

	return nil
}
