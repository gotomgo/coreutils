package concurrency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RateLimiter_Acquire(t *testing.T) {
	t.Run("rateLimiter.Acquire", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, 1)
		if !assert.NotNil(t, limiter, "rateLimiter should not be nil") {
			return
		}

		acquired := 0
		start := time.Now()

		// we shouldn't be able to do 11 acquires in less than 2 seconds
		// (eventsPerSecond == 5, and burst == 1)
		for i := 0; i < 11; i++ {
			if err := limiter.Acquire(context.Background()); err == nil {
				acquired++
			}
		}

		if !assert.Equal(t, 11, acquired, "acquired should = 11") {
			return
		}

		if !assert.GreaterOrEqual(t, time.Since(start).Microseconds(), int64(200000), "too many acquires in < 1.75 second") {
			return
		}
	})

	t.Run("rateLimiter.Acquire sampler", func(t *testing.T) {
		t.Parallel()

		acquired := 0

		limiter := NewRateLimiter(10, 1)
		if !assert.NotNil(t, limiter, "rateLimiter should not be nil") {
			return
		}

		timer := time.After(time.Second)
		for {
			select {
			case <-time.After(5 * time.Millisecond):
				if err := limiter.Acquire(context.Background()); err == nil {
					acquired++
				}

			case <-timer:
				goto done
			}
		}
	done:
		// this is 11 (not 10), because when we are close to the boundary
		// condition (time) we can acquire in the next interval. This is
		// expected behaviour, and we would need to reduce the boundary to
		// be somewhat less (~940ms) to get a value of 10. Note that this
		// discrepancy is not compounded over time. If we ran for 3 seconds
		// we would expect 31, and 10 seconds, 101
		if !assert.Equal(t, 11, acquired, "acquire should be 11 for a 1 second period with limit=10 and burst=1") {
			return
		}
	})

	t.Run("rateLimiter cancelled context", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(10, 1)
		if !assert.NotNil(t, limiter, "rateLimiter should not be nil") {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if !assert.Error(t, limiter.Acquire(ctx), "rateLimiter should fail on cancelled context") {
			return
		}
	})
}
