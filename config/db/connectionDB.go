package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var LoanCollection *mongo.Collection
var LogCollection *mongo.Collection

func ConnectDB(connectionString string) {

    clientOptions := options.Client().ApplyURI(connectionString)

    client, err := mongo.NewClient(clientOptions)

    if err != nil {
        log.Fatalf("Error creating MongoDB client: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }

    Client = client
    LoanCollection = client.Database("loan").Collection("loans")
    UserCollection = client.Database("loan").Collection("users")
    LogCollection = client.Database("loan").Collection("logs")

}

func CreateTextIndex(collection *mongo.Collection) {
    indexModel := mongo.IndexModel{
        Keys: bson.D{
            {Key: "tags", Value: "text"},
            {Key: "autorname", Value: "text"},
        },
        Options: options.Index().SetDefaultLanguage("english"),
    }

    _, err := collection.Indexes().CreateOne(context.Background(), indexModel)
    if err != nil {
        log.Fatal(err)
    }
}