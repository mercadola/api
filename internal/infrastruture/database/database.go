package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClientInterface interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

func NewClient(uri string, ctx context.Context) (*mongo.Client, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return mongoClient, nil
}
