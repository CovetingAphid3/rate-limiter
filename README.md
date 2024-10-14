
# Distributed Rate Limiter

A simple yet effective distributed rate limiter implemented in Go, utilizing Redis for token storage and management. This project allows you to control the rate of incoming requests to your services, ensuring fair usage and protection against abuse.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Example](#example)
- [Logging](#logging)
- [Contributing](#contributing)
- [License](#license)

## Features

- Distributed rate limiting using Redis.
- Configurable rate and burst limits.
- Detailed logging for monitoring and debugging.
- Simple integration with existing Go applications.

## Installation

To install the required dependencies and set up the project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/CovetingAphid3/rate-limiter.git
   cd rate-limiter
   ```

2. Install Go modules:

   ```bash
   go mod tidy
   ```

3. Ensure you have Redis running locally or access to a Redis server.

## Usage

1. Import the rate limiter package in your Go application:

   ```go
   import "github.com/CovetingAphid3/rate-limiter/internal/limiter"
   ```

2. Initialize the Redis client and the rate limiter in your main application file:

   ```go
   package main

   import (
       "github.com/go-redis/redis/v8"
       "github.com/CovetingAphid3/rate-limiter/internal/limiter"
       "context"
   )

   func main() {
       ctx := context.Background()

       // Create a Redis client
       client := redis.NewClient(&redis.Options{
           Addr: "localhost:6379", // Redis server address
       })

       // Initialize the rate limiter
       rateLimiter := limiter.NewRedisLimiter(client, "my-rate-limiter", 5, 10)

       // Your application logic...
   }
   ```

## Configuration

- **Rate**: The number of requests allowed per time unit (e.g., 5 requests per second).
- **Burst**: The maximum number of requests that can be processed in a single burst (e.g., 10 requests).

## Example

Here’s a simple example of a HTTP server that uses the rate limiter:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/go-redis/redis/v8"
    "github.com/CovetingAphid3/rate-limiter/internal/limiter"
    "context"
)

func main() {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    rateLimiter := limiter.NewRedisLimiter(client, "my-rate-limiter", 5, 10)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if rateLimiter.Allow() {
            fmt.Fprintf(w, "Request successful!")
        } else {
            http.Error(w, "Rate limit exceeded, try again later.", http.StatusTooManyRequests)
        }
    })

    port := ":8080"
    log.Printf("Server is running on http://localhost%s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
```

## Logging

The rate limiter includes detailed logging to monitor the request flow and the status of the tokens. Logs are written to the standard output and can help diagnose issues and understand system behavior.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

```

