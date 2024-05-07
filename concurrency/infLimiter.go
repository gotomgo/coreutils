package concurrency

import (
	"context"
)

// infLimiter implements a Limiter that allows for infinite concurrent work
type infLimiter struct{}

// Acquire always succeeds
func (l *infLimiter) Acquire(ctx context.Context) (err error) {
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
	}

	return
}

// Release always succeeds
func (l *infLimiter) Release() {}

// create a singleton that can be referenced (allowable because we have no
// state to protect)
var _infLimiter = &infLimiter{}

// GetInfinityLimiter returns an instance of infLimiter
func GetInfinityLimiter() Limiter {
	return _infLimiter
}
