package database

import (
	"context"
	"go-mongo-auth/configs"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectionManager() {
	client, err := mongo.NewClient(options.Client().ApplyURI(configs.Get("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), configs.GetChrono("mongo.timeout"))
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ping data base
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Mongo client created.")
	MongoClient = client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(configs.Get("mongo.database")).Collection(collectionName)
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
