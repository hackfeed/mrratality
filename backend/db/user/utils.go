package userdb

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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
