package concurrency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func Test_FailFastLimiter_New(t *testing.T) {
	t.Run("failFastLimiter.NewFailFastLimiter", func(t *testing.T) {
		t.Parallel()

		limiter := NewFailFastLimiter(20)
		if !assert.NotNil(t, limiter, "failFastLimiter should not be nil") {
			return
		}

		if !assert.Equal(t, rate.Limit(1), limiter.(*failFastLimiter).lim.Limit()) {
			return
		}
		if !assert.Equal(t, 20, limiter.(*failFastLimiter).lim.Burst()) {
			return
		}
	})

	t.Run("failFastLimiter.NewFailFastLimiter bad events/sec", func(t *testing.T) {
		t.Parallel()

		limiter := NewFailFastLimiter(-20)
		if !assert.NotNil(t, limiter, "failFastLimiter should not be nil") {
			return
		}

		if !assert.Equal(t, rate.Limit(0), limiter.(*failFastLimiter).lim.Limit()) {
			return
		}
		if !assert.Equal(t, 0, limiter.(*failFastLimiter).lim.Burst()) {
			return
		}
	})
}

func Test_FailFastLimiter_Acquire(t *testing.T) {
	t.Run("failFastLimiter.Acquire", func(t *testing.T) {
		t.Parallel()

		limiter := NewFailFastLimiter(5)
		if !assert.NotNil(t, limiter, "failFastLimiter should not be nil") {
			return
		}

		acquired := 0
		start := time.Now()

		for i := 0; i < int(2*5); i++ {
			if err := limiter.Acquire(context.Background()); err == nil {
				acquired++
			}
		}

		if !assert.Equal(t, 5, acquired, "acquired should = 5") {
			return
		}

		// make sure the 10 acquire attempts with 5 grants happened in under 10 milliseconds
		// 10 milliseconds is *very* generous
		if assert.Less(t, time.Since(start).Milliseconds(), (10 * time.Millisecond).Milliseconds(), "acquire appears to have blocked") {
			return
		}

		// run the acquires again. None should succeed
		for i := 0; i < int(2*5); i++ {
			if err := limiter.Acquire(context.Background()); err == nil {
				acquired++
			}
		}

		// acquired should still be 5 as not enough time has passed to allow more
		if !assert.Equal(t, 5, acquired, "acquired should = 5") {
			return
		}

		// wait a full second
		time.Sleep(time.Second)

		// acquire must succeed after full second has passed...
		if !assert.NoError(t, limiter.Acquire(context.Background()), "acquire should work due to passage of time") {
			return
		}
	})

	t.Run("failFastLimiter.Acquire sampler", func(t *testing.T) {
		t.Parallel()

		acquired := 0

		limiter := NewFailFastLimiter(10)
		if !assert.NotNil(t, limiter, "failFastLimiter should not be nil") {
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
		if !assert.Equal(t, 10, acquired, "acquire should be 10 for a 1 second period with limit=10 and burst=1") {
			return
		}
	})

	t.Run("failFastLimiter cancelled context", func(t *testing.T) {
		t.Parallel()

		limiter := NewFailFastLimiter(10)
		if !assert.NotNil(t, limiter, "failFastLimiter should not be nil") {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if !assert.Error(t, limiter.Acquire(ctx), "failFastLimiter should fail on cancelled context") {
			return
		}
	})

}
