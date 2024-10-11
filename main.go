package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/CovetingAphid3/rate-limiter/limiter"
)

func main() {
    // Initialize the rate limiter (e.g., 5 requests per second, burst of 10)
    rateLimiter := limiter.NewLimiter(5, 10)

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
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

