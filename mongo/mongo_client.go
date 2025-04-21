package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(ctx context.Context, uri string, databaseName string) *mongo.Client {
	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ Error connecting to MongoDB: %v\n", err)

		return nil
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Error pinging MongoDB: %v\n", err)

		return nil
	}

	log.Println("✅ Connected to MongoDB")

	return mongoClient
}
