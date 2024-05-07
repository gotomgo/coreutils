package concurrency

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// semaphoreLimiter implements a Limiter that uses a semaphore construct
type semaphoreLimiter struct {
	semaphore chan struct{}
}

// NewSemaphoreLimiter creates an instance of semaphoreLimiter with a max
// semaphore size
func NewSemaphoreLimiter(maxConcurrency int) Limiter {
	return &semaphoreLimiter{
		semaphore: make(chan struct{}, maxConcurrency),
	}
}

// Acquire acquires the right to perform work, or blocks until that right
// can be obtained (or the context times out / is cancelled)
func (l *semaphoreLimiter) Acquire(ctx context.Context) (result error) {
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// use infinite wait unless context has a deadline
	waitLimit := rate.InfDuration
	if deadline, ok := ctx.Deadline(); ok {
		waitLimit = time.Until(deadline)
	}

	select {
	// try to acquire by writing to semaphore channel
	case l.semaphore <- struct{}{}:
		// semaphore acquired so we are good to go

	// context cancelled
	case <-ctx.Done():
		result = ctx.Err()

	// wait time expired
	case <-time.After(waitLimit):
		select {
		// check the context just in case
		case <-ctx.Done():
			result = ctx.Err()
		default:
			result = fmt.Errorf("timeout in acquire (context deadline exceeded)")
		}
	}

	return
}

// Release releases a semaphore allowing another worker to acquire
func (l *semaphoreLimiter) Release() {
	// read from the semaphore to free up a slot
	<-l.semaphore
}
