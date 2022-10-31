package database

import (
	"context"
	"go-mongo-auth/internal/config"
	"time"

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

	MongoClient struct {
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

	return &MongoClient{
		client:      client,
		mongoConfig: mongoConfig,
	}, nil
}

func (c *MongoClient) GetCollection(collectionName string) *mongo.Collection {
	collection := c.client.Database(c.mongoConfig.Database).Collection(collectionName)
	return collection
}

func (c *MongoClient) CreateOneDocument(collectionName string, doc any) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := c.GetCollection(collectionName)

	return collection.InsertOne(ctx, doc)
}

func (c *MongoClient) FindOneDocument(collectionName string, filter primitive.M) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := c.GetCollection(collectionName)

	return collection.FindOne(ctx, filter)
}
