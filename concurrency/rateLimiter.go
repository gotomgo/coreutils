package concurrency

import (
	"context"

	"golang.org/x/time/rate"
)

// rateLimiter implements a Limiter that uses rate.Limiter
type rateLimiter struct {
	lim *rate.Limiter
}

// NewRateLimiter creates an instance of rateLimiter based on events per
// second and a burst size
func NewRateLimiter(eventsPerSecond rate.Limit, burst int) Limiter {
	// rate.Limit allows any value. Just clip to 0
	if eventsPerSecond < 0.0 {
		eventsPerSecond = 0.0
	}

	return &rateLimiter{
		lim: rate.NewLimiter(eventsPerSecond, burst),
	}
}

// Acquire acquires the right to work using rate.Limiter.Wait
func (l *rateLimiter) Acquire(ctx context.Context) error {
	return l.lim.Wait(ctx)
}

// Release releases a previous acquire
//
//	Notes
//		This func is a NOP as the rate.Limiter replenishes itself based on
//		the passage of time
//
func (l *rateLimiter) Release() {}
