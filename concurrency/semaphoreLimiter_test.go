package concurrency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SemaphoreLimiter_Acquire(t *testing.T) {
	t.Run("semaphoreLimiter acquire", func(t *testing.T) {
		t.Parallel()

		maxConcurrency := 5
		lim := NewSemaphoreLimiter(maxConcurrency)
		if !assert.NotNil(t, lim, "NewSemaphoreLimiter should not return nil") {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		for i := 0; i < maxConcurrency; i++ {
			if !assert.NoError(t, lim.Acquire(ctx), "expecting semaphoreLimiter.acquire to succeed") {
				return
			}
		}

		if !assert.Error(t, lim.Acquire(ctx), "expecting semaphoreLimiter.acquire to timeout") {
			return
		}
	})

	t.Run("semaphoreLimiter cancelled context", func(t *testing.T) {
		t.Parallel()

		limiter := NewSemaphoreLimiter(5)
		if !assert.NotNil(t, limiter, "semaphoreLimiter should not be nil") {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if !assert.Error(t, limiter.Acquire(ctx), "semaphoreLimiter should fail on cancelled context") {
			return
		}
	})

	t.Run("semaphoreLimiter async cancelled context", func(t *testing.T) {
		t.Parallel()

		limiter := NewSemaphoreLimiter(0)
		if !assert.NotNil(t, limiter, "semaphoreLimiter should not be nil") {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

		if !assert.Error(t, limiter.Acquire(ctx), "semaphoreLimiter should fail on cancelled context") {
			return
		}
	})

}
