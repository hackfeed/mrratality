package userdb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	Email        *string            `bson:"email" validate:"email,required"`
	Password     *string            `bson:"password" validate:"required,min=6"`
	Token        *string            `bson:"token"`
	RefreshToken *string            `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}
