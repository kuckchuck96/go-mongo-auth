package database

import (
	"context"
	"go-mongo-auth/internal/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IMongoClient interface {
		GetCollection(string) *mongo.Collection
		CreateOneDocument(string, any) (*mongo.InsertOneResult, error)
		FindOneDocument(string, primitive.M) *mongo.SingleResult
	}

	mongoClient struct {
		client      *mongo.Client
		mongoConfig config.Mongo
	}
)

func NewMongoClient(mongoConfig config.Mongo) (IMongoClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConfig.Uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoConfig.Timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// ping data base
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mongoClient{
		client:      client,
		mongoConfig: mongoConfig,
	}, nil
}

func getContextWithTimeout(client *mongoClient) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), client.mongoConfig.ContextTimeout)
}

func (client *mongoClient) GetCollection(collectionName string) *mongo.Collection {
	collection := client.client.Database(client.mongoConfig.Database).Collection(collectionName)
	return collection
}

func (client *mongoClient) CreateOneDocument(collectionName string, doc any) (*mongo.InsertOneResult, error) {
	ctx, cancel := getContextWithTimeout(client)
	defer cancel()

	collection := client.GetCollection(collectionName)

	return collection.InsertOne(ctx, doc)
}

func (client *mongoClient) FindOneDocument(collectionName string, filter primitive.M) *mongo.SingleResult {
	ctx, cancel := getContextWithTimeout(client)
	defer cancel()

	collection := client.GetCollection(collectionName)

	return collection.FindOne(ctx, filter)
}
