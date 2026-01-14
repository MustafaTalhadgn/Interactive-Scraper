package main

import (
	"context"
	"fmt"
	"time"

	"InteractiveScraper/internal/transport"
)

func main() {
	limiter := transport.NewRateLimiter(1) // 1 req/sec

	fmt.Println("Testing rate limiter (1 req/sec + jitter)...")

	for i := 1; i <= 5; i++ {
		start := time.Now()

		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		elapsed := time.Since(start)
		fmt.Printf("Request %d allowed after %v\n", i, elapsed)
	}
}
