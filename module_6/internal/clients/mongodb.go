package clients

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, nil, err
	}

	log.Println("Successfully connected to MongoDB!")
	return client, client.Database(dbName), nil
}
