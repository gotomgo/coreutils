package concurrency

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// timeLimiter implements Limiter using a time-based rate
type timeLimiter struct {
	rate      int64
	semaphore chan struct{}
	done      chan bool
}

// NewTimeLimiter creates an instance of timeLimiter
//
//	Notes
//		timeLimiter uses a time.Timer to drive when a successful Acquire can
//		occur, and uses an unbuffered channel to control blocking and
//		unblocking.
//
//		The timer requires the use of a go routine, making this form
//		of limiter more expensive than rateLimiter, but it is more consistent
//		in that it is truly periodic, not approximately periodic.
//
//		In addition, the rate can be changed (SetRate) at anytime (based on
//		environmental conditions for example)
//
func NewTimeLimiter(rate time.Duration, done chan bool) Limiter {
	return (&timeLimiter{
		rate:      int64(rate),
		semaphore: make(chan struct{}),
		done:      done,
	}).start()
}

// Acquire acquires the right to perform work, or blocks until that right
// can be obtained (or the context times out / is cancelled)
func (l *timeLimiter) Acquire(ctx context.Context) (result error) {
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

// Release is a NOP as acquire is based on time elapsed
func (l *timeLimiter) Release() {}

// SetRate sets the rate for the timer atomically
//
//	Notes
//		SetRate does not re-program the timer immediately. The timer resets
//		itself after the next interval triggers
//
func (l *timeLimiter) SetRate(rate time.Duration) {
	atomic.StoreInt64(&l.rate, int64(rate))
}

// getRate gets the current timer rate atomically
func (l *timeLimiter) getRate() time.Duration {
	dur := atomic.LoadInt64(&l.rate)
	return time.Duration(dur)
}

// start runs the go routine that processes the timer intervals
func (l *timeLimiter) start() Limiter {
	go func() {
		timer := time.NewTimer(l.getRate())

		for {
			select {
			// our work is done
			case <-l.done:
				return
			case <-timer.C:
				timer.Stop()
				select {
				// read will release a waiter, if any
				case <-l.semaphore:
				default:
				}
				// start the cycle again
				timer.Reset(l.getRate())
			}
		}
	}()

	return l
}
