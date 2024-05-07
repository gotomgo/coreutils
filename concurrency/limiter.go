package concurrency

import "context"

// Limiter implements an acquire / release paradigm  to control the flow of
// work that can happen concurrently in a system
type Limiter interface {
	// Acquire the right to work
	Acquire(ctx context.Context) error

	// Release an acquired right to work so another worker can acquire
	Release()
}
