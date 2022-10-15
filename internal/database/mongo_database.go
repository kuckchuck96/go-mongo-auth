package database

import (
	"context"
	"go-mongo-auth/internal/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectionManager() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Get("mongo.uri")))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.GetChrono("mongo.timeout"))
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// ping data base
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Mongo client created.")
	MongoClient = client

	return nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(config.Get("mongo.database")).Collection(collectionName)
	return collection
}

func CreateOneDocument(collectionName string, doc any) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetCollection(MongoClient, collectionName)

	return collection.InsertOne(ctx, doc)
}

func FindOneDocument(collectionName string, filter primitive.M) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetCollection(MongoClient, collectionName)

	return collection.FindOne(ctx, filter)
}
