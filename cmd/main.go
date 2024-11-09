package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/CovetingAphid3/rate-limiter/internal/config"
	"github.com/CovetingAphid3/rate-limiter/internal/limiter"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
    config.Test()
    // Initialize Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis server address
        Password: "",           
        DB: 0,                   // Use default DB
    })

    // Initialize the Redis-based rate limiter (e.g., 5 requests per second, burst of 10)
    // rateLimiter := limiter.NewRedisLimiter(rdb, "my-rate-limiter", 5, 10)
	rateLimiter := limiter.NewRedisLimiter(rdb, "my-rate-limiter", 2, 5)

    // Define a handler that uses the rate limiter
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if rateLimiter.Allow() {
            fmt.Fprintf(w, "Request successful!")
        } else {
            http.Error(w, "Rate limit exceeded, try again later.", http.StatusTooManyRequests)
        }
    })

    // Start the server
    port := ":8080"
    log.Printf("Server is running on http://localhost%s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

