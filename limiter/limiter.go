package limiter

import (
    "time"
)

// Limiter struct to hold rate limit details
type Limiter struct {
    rate  int           // requests per time unit
    burst int           // max burst of requests allowed
    tokens int          // current number of tokens
    last time.Time      // last time tokens were added
    interval time.Duration // time between token additions
}

// NewLimiter creates a new rate limiter
func NewLimiter(rate, burst int) *Limiter {
    return &Limiter{
        rate:     rate,
        burst:    burst,
        tokens:   burst,
        last:     time.Now(),
        interval: time.Second / time.Duration(rate),
    }
}

// Allow checks if a request can pass the rate limiter
func (l *Limiter) Allow() bool {
    now := time.Now()
    elapsed := now.Sub(l.last)

    // Add tokens based on elapsed time
    l.tokens += int(elapsed / l.interval)
    if l.tokens > l.burst {
        l.tokens = l.burst
    }

    l.last = now

    // Check if there are enough tokens for a request
    if l.tokens > 0 {
        l.tokens--
        return true
    }
    return false
}

