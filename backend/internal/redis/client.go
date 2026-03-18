package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Client is our global Redis connection pool
var Client *redis.Client

// Ctx is a background context we use for Redis operations
var Ctx = context.Background()

// InitRedis connects to our local Docker Redis container
func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // The port we exposed in docker-compose
		Password: "",               // No password set in our simple Docker setup
		DB:       0,                // Use default DB
	})

	// Ping the database to make sure it's actually listening!
	pong, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Successfully connected to Redis! (Response: %s)", pong)
}