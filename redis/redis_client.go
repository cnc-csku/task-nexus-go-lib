package redis

import (
	"github.com/redis/go-redis/v9"

	"context"
	"log"
)

func NewRedisClient(ctx context.Context, uri string) *redis.Client {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		log.Fatalln("Redis URI is invalid: ", err)
	}

	client := redis.NewClient(opts)

	// test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("🚫 Cannot connect to Redis | ", err)
	} else {
		log.Println("✅ Connected to Redis")
	}

	return client
}
