package db

import (
	"context"
	"fmt"
	"log"
	"time"
	"wetalk-academy/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(conf *config.Config) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(conf.Database.URI))
	if err != nil {
		log.Fatalf("[❌] Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("[❌] Failed to ping MongoDB: %v", err)
	}

	fmt.Println("[✅] Connected to MongoDB successfully")

	return &MongoDB{
		client:   client,
		Database: client.Database(conf.Database.Name),
	}
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}
