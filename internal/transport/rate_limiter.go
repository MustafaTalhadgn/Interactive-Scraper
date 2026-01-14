package transport

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter       *rate.Limiter
	jitterEnabled bool
	maxJitter     time.Duration
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(
			rate.Limit(rps),
			1,
		),
		jitterEnabled: true,
		maxJitter:     500 * time.Millisecond,
	}
}

func (r *RateLimiter) Wait(ctx context.Context) error {
	if err := r.limiter.Wait(ctx); err != nil {
		return err
	}

	if r.jitterEnabled {
		jitter := time.Duration(rand.Intn(int(r.maxJitter.Milliseconds()))) * time.Millisecond
		select {
		case <-time.After(jitter):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
