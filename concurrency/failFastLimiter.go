package concurrency

import (
	"context"
	"fmt"
	"math"

	"golang.org/x/time/rate"
)

// failFastLimiter implements a Limiter based on rate.Limiter
type failFastLimiter struct {
	lim *rate.Limiter
}

// NewFailFastLimiter creates an instance of rateLimiter based on the number
// of acquires / sec we want to achieve.
//
//	Notes
//		FailFastLimiter allocates everything to burst with an eventsPerSecond
//		== 1 (except when burst is 0)
//
//		This arrangement allows for all the acquires to happen instantaneously
//		rather than be subject to sub-intervals of time, and the formula
//
//		acquiresPerSecond := (eventsPerSecond+burst)-1
//
//		is satisfied because eventsPerSecond =1 and the -1 cancel leaving
//		burst == acquiresPerSecond
//
//
func NewFailFastLimiter(acquiresPerSecond int) Limiter {
	// rate.Limit allows any value. Just clip to 0
	if acquiresPerSecond < 0 {
		acquiresPerSecond = 0
	}

	burst := acquiresPerSecond
	eventsPerSecond := math.Min(float64(burst), 1.0)

	return &failFastLimiter{
		lim: rate.NewLimiter(rate.Limit(eventsPerSecond), burst),
	}
}

// Acquire the right to perform work
//
//	Notes
//		The rate.Limiter.Allow() func is used to determine if the acquire can
//		be granted
//
func (l *failFastLimiter) Acquire(ctx context.Context) (err error) {
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		if !l.lim.Allow() {
			err = fmt.Errorf("rate limit exceeded")
		}
	}

	return
}

// Release releases a previous acquire
//
//	Notes
//		This func is a NOP as the rate.Limiter replenishes itself based on
//		the passage of time
//
func (l *failFastLimiter) Release() {}
