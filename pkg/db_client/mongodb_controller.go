package db_client

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	mongoClient        *mongo.Client
	ctx                context.Context
	usersCollection    *mongo.Collection
	itemsCollection    *mongo.Collection
	sessionsCollection *mongo.Collection
}

func (c *Client) Close() error {
	return c.mongoClient.Disconnect(c.ctx)
}

func Get(timeout time.Duration, conn_string string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	clientOptions := options.Client().ApplyURI(conn_string)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, nil
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	itemsDatabase := client.Database("itemsData")
	usersCollection := itemsDatabase.Collection("users")
	itemsCollection := itemsDatabase.Collection("items")
	sessionsCollection := itemsDatabase.Collection("sessions")

	return &Client{mongoClient: client,
		itemsCollection:    itemsCollection,
		usersCollection:    usersCollection,
		sessionsCollection: sessionsCollection}, nil
}
