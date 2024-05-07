package concurrency

import (
	"context"
	"fmt"
	"sync/atomic"
)

// balanceChecker is a Limiter that is a wrapper around another Limiter
// and maintains the current balance between acquire / release
type balanceChecker struct {
	limiter Limiter
	balance int32
}

func NewBalanceChecker(limiter Limiter) Limiter {
	return &balanceChecker{limiter: limiter}
}

// Acquire the right to work
func (b *balanceChecker) Acquire(ctx context.Context) (err error) {
	// only update the balance if we acquire
	if err = b.limiter.Acquire(ctx); err == nil {
		atomic.AddInt32(&b.balance, 1)
	}

	return
}

// Release an acquired right to work so another worker can acquire
//
//	Notes
//		If the balance is (or goes) negative, this func panics
//
func (b *balanceChecker) Release() {
	b.limiter.Release()

	// if there is an over-release, panic
	if bal := atomic.AddInt32(&b.balance, -1); bal < 0 {
		panic(fmt.Errorf("semaphoreLimiter balance < 0"))
	}
}

// Balance returns the current balance for the Limiter
func (b *balanceChecker) Balance() int {
	return int(atomic.LoadInt32(&b.balance))
}
