package concurrency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BalanceChecker_NewBalanceChecker(t *testing.T) {
	bc := NewBalanceChecker(GetInfinityLimiter())
	assert.NotNil(t, bc, "NewBalancerChecker should not return nil")
	// assumption here that GetInfinityLimiter() returns a singleton
	assert.Equal(t, bc.(*balanceChecker).limiter, GetInfinityLimiter(), "expecting semaphoreLimiter to be GetInfinityLimiter()")
}

func Test_BalanceChecker_Balance(t *testing.T) {
	bc := NewBalanceChecker(GetInfinityLimiter())
	assert.NotNil(t, bc, "NewBalancerChecker should not return nil")

	for i := 0; i < 3; i++ {
		assert.NoError(t, bc.Acquire(context.Background()), "expecting balanceChecker.Acquire to return nil")
		assert.Equal(t, i+1, bc.(*balanceChecker).Balance(), "expecting balance to be %d", i+1)
	}

	for i := 0; i < 3; i++ {
		bc.Release()
		assert.Equal(t, 2-i, bc.(*balanceChecker).Balance(), "expecting balance to be %d", 2-i)
	}
}

func Test_BalanceChecker_OverRelease(t *testing.T) {
	bc := NewBalanceChecker(GetInfinityLimiter())
	assert.NotNil(t, bc, "NewBalancerChecker should not return nil")

	assert.Panics(t, func() { bc.Release() }, "expecting panic after balanceChecker.Release()")

}
