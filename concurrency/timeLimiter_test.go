package concurrency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_TimeLimiter_Acquire(t *testing.T) {
	t.Run("timeLimiter.Acquire", func(t *testing.T) {
		t.Parallel()

		done := make(chan bool)

		limiter := NewTimeLimiter(time.Duration(50*time.Millisecond), done)
		if !assert.NotNil(t, limiter, "timeLimiter should not be nil") {
			return
		}

		acquired := 0
		start := time.Now()

		// we shouldn't be able to do 20 acquires in less than 1 second
		// 1000ms / 50ms = 20
		for i := 0; i < 20; i++ {
			if err := limiter.Acquire(context.Background()); err == nil {
				acquired++
			}
		}

		if !assert.Equal(t, 20, acquired, "timeLimiter acquired should = 20") {
			return
		}

		if !assert.GreaterOrEqual(t, time.Since(start).Microseconds(), int64(100000), "too many acquires in < 1 second") {
			return
		}
	})

	t.Run("rateLimiter.Acquire sampler", func(t *testing.T) {
		t.Parallel()

		acquired := 0

		done := make(chan bool)

		limiter := NewTimeLimiter(time.Duration(50*time.Millisecond), done)
		if !assert.NotNil(t, limiter, "timeLimiter should not be nil") {
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
		if !assert.Equal(t, 20, acquired, "acquire should be 20 for a 1 second period with 50ms interval") {
			return
		}
	})

	t.Run("timeLimiter cancelled context", func(t *testing.T) {
		t.Parallel()

		done := make(chan bool)

		limiter := NewTimeLimiter(time.Duration(50*time.Millisecond), done)
		if !assert.NotNil(t, limiter, "timeLimiter should not be nil") {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if !assert.Error(t, limiter.Acquire(ctx), "timeLimiter should fail on cancelled context") {
			return
		}
	})
}
