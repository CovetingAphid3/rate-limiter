package limiter

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
    client *redis.Client
    key    string
    burst  int
    interval time.Duration
    script string // Add a field to hold the Lua script
}

// NewRedisLimiter creates a new Redis-based rate limiter

func NewRedisLimiter(client *redis.Client, key string, rate int, burst int) *RedisLimiter {
    limiter := &RedisLimiter{
        client:   client,
        key:      key,
        burst:    burst,
        interval: time.Duration(1000/rate) * time.Millisecond,
    }

    // Initialize Redis keys
    ctx := context.Background()
    client.Set(ctx, key+":tokens", burst, 0) // Set initial tokens to burst limit
    client.Set(ctx, key+":last", time.Now().UnixNano()/1e6, 0) // Set last timestamp to now
    log.Printf("INFO: Initialized rate limiter with key: '%s', burst: %d", key, burst)

    return limiter
}


// Allow checks if a request can pass the rate limiter
func (r *RedisLimiter) Allow() bool {
    ctx := context.Background()

    // Retrieve current tokens and last timestamp from Redis
    tokensStr, err := r.client.Get(ctx, r.key+":tokens").Result()
    if err != nil {
        log.Printf("Error retrieving tokens: %v", err)
        return false
    }
    lastStr, err := r.client.Get(ctx, r.key+":last").Result()
    if err != nil {
        log.Printf("Error retrieving last time: %v", err)
        return false
    }

    // Convert retrieved values
    tokens := 0
    if tokensStr != "" {
        tokens, _ = strconv.Atoi(tokensStr)
    }
    lastTime, _ := strconv.ParseInt(lastStr, 10, 64)

    // Get the current time in milliseconds
    currentTime := time.Now().UnixNano() / 1e6
    log.Printf("Current Time: %d, Last Time: %d, Current Tokens: %d", currentTime, lastTime, tokens)

    // Calculate the elapsed time and add tokens
    elapsed := currentTime - lastTime
    newTokens := tokens + int(elapsed/int64(r.interval.Milliseconds()))

    // Cap tokens to the burst limit
    if newTokens > r.burst {
        newTokens = r.burst
    }

    log.Printf("New Tokens After Replenishment: %d (Max Burst: %d)", newTokens, r.burst)

    // Check if we can allow the request
    if newTokens > 0 {
        newTokens-- // Decrement a token for the allowed request
        r.client.Set(ctx, r.key+":tokens", newTokens, 0) // Update token count
        r.client.Set(ctx, r.key+":last", currentTime, 0) // Update last time
        log.Printf("Request allowed. Tokens remaining: %d", newTokens)
        return true
    }

    // Update last time if we are at limit
    r.client.Set(ctx, r.key+":last", currentTime, 0)
    log.Printf("Request denied. Rate limit exceeded. Tokens remaining: %d", newTokens)
    return false
}

