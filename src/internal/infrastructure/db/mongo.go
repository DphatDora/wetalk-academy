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

var (
	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
)

func InitMongoDB(conf *config.Config) {
	uri := conf.Database.URI
	dbName := conf.Database.Name

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %s", err)
	}

	fmt.Println("[✅] Connect to MongoDB successfully")

	mongoClient = client
	mongoDatabase = client.Database(dbName)
}

func GetDB() *mongo.Database {
	if mongoDatabase == nil {
		panic("[❌] Connection to MongoDB is not setup")
	}
	return mongoDatabase
}

func CloseMongoDB() error {
	if mongoClient == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongoClient.Disconnect(ctx)
}
