package redisclient

import (
    "github.com/go-redis/redis/v8"
    "log"
    "context"
)

var ctx = context.Background()

// NewRedisClient initializes and returns a Redis client
func NewRedisClient(addr, password string, db int) *redis.Client {
    // Initialize Redis client with given options
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,     // Redis server address
        Password: password, // Redis password, leave empty if no password is set
        DB:       db,       // Redis database number (default is 0)
    })

    // Test the connection
    err := rdb.Ping(ctx).Err()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }

    log.Println("Connected to Redis")
    return rdb
}

