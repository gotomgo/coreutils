package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jpillora/backoff"
	"github.com/stretchr/testify/assert"
)

func Test_Backoff_Params(t *testing.T) {
	b := backoff.Backoff{
		Min:    50 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: DefaultFactor,
		Jitter: false,
	}

	var sum time.Duration

	for i := 0; i < 10; i++ {
		wait := b.ForAttempt(float64(i))
		sum += wait
		fmt.Printf("retry %d => %v (total wait=%v)\n", i+1, wait, sum)
	}
}

func Test_RetryConfig_NewDefault(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		t.Parallel()

		rc := NewDefaultRetryConfig()

		if !assert.Equal(t, 50, rc.MinWait, "cfg.MinWait should be 50ms") {
			return
		}
		if !assert.Equal(t, 5000, rc.MaxWait, "cfg.MaxWait should be 5000s") {
			return
		}
		if !assert.Equal(t, DefaultFactor, rc.Factor, "cfg.Factor should be %d", DefaultFactor) {
			return
		}
		if !assert.True(t, rc.Jitter, "cfg.Jitter should be true") {
			return
		}
	})
}

func Test_RetryHandler_FromConfig(t *testing.T) {
	t.Run("from config", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandlerFromConfig(RetryConfig{
			MinWait: 25,
			MaxWait: 8000,
			Factor:  2.0,
			Jitter:  true,
		})

		if !assert.Equal(t, (25 * time.Millisecond).Seconds(), r.backoff.Min.Seconds(), "min should be 25ms") {
			return
		}
		if !assert.Equal(t, (8 * time.Second).Seconds(), r.backoff.Max.Seconds(), "max should be 8s") {
			return
		}
		if !assert.Equal(t, 2.0, r.backoff.Factor, "factor should be 2.0") {
			return
		}
		if !assert.True(t, r.backoff.Jitter, "jitter should be true") {
			return
		}
	})

	t.Run("skewed config", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandlerFromConfig(RetryConfig{
			MinWait: 25,
			MaxWait: 3,
			Factor:  0.8,
			Jitter:  true,
		})

		if !assert.Equal(t, MinWaitDuration.Seconds(), r.backoff.Min.Seconds(), "min should be %v", MinWaitDuration) {
			return
		}
		if !assert.Equal(t, (25 * time.Millisecond).Seconds(), r.backoff.Max.Seconds(), "max should be 25s") {
			return
		}
		if !assert.Equal(t, MinFactor, r.backoff.Factor, "factor should be %d", MinFactor) {
			return
		}
		if !assert.True(t, r.backoff.Jitter, "jitter should be true") {
			return
		}
	})
}

func Test_RetryHandler_WithMinMax(t *testing.T) {
	t.Run("normal min < max", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(25*time.Millisecond, 8*time.Second, 2.0, false)

		if !assert.Equal(t, (25 * time.Millisecond).Seconds(), r.backoff.Min.Seconds(), "min should be 25ms") {
			return
		}
		if !assert.Equal(t, (8 * time.Second).Seconds(), r.backoff.Max.Seconds(), "max should be 8s") {
			return
		}
	})

	t.Run("too low: min < 5ms < max", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(3*time.Millisecond, 8*time.Second, 2.0, false)

		if !assert.Equal(t, MinWaitDuration.Seconds(), r.backoff.Min.Seconds(), "min should be %v", MinWaitDuration) {
			return
		}
		if !assert.Equal(t, (8 * time.Second).Seconds(), r.backoff.Max.Seconds(), "max should be 8s") {
			return
		}
	})

	t.Run("inverted: min > max", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(50*time.Millisecond, 25*time.Millisecond, 2.0, false)

		if !assert.Equal(t, (25 * time.Millisecond).Seconds(), r.backoff.Min.Seconds(), "min should be 25ms") {
			return
		}
		if !assert.Equal(t, (50 * time.Millisecond).Seconds(), r.backoff.Max.Seconds(), "max should be 50ms") {
			return
		}
	})

	t.Run("inverted: min > max < 5ms", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(50*time.Millisecond, 3*time.Millisecond, 2.0, false)

		if !assert.Equal(t, MinWaitDuration.Seconds(), r.backoff.Min.Seconds(), "min should be %v", MinWaitDuration) {
			return
		}
		if !assert.Equal(t, (50 * time.Millisecond).Seconds(), r.backoff.Max.Seconds(), "max should be 50ms") {
			return
		}
	})
}

func Test_RetryHandler_WithFactor(t *testing.T) {
	t.Run("factor > 1", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(MinWaitDuration, 5*time.Second, 1.5, false)
		if !assert.Equal(t, 1.5, r.backoff.Factor, "factor should == 1.5") {
			return
		}
		r = NewRetryHandler(MinWaitDuration, 5*time.Second, 2.0, false)
		if !assert.Equal(t, 2.0, r.backoff.Factor, "factor should == 2.0") {
			return
		}
		r = NewRetryHandler(MinWaitDuration, 5*time.Second, 1.0, false)
		if !assert.Equal(t, 1.0, r.backoff.Factor, "factor should == 1.0") {
			return
		}
	})

	t.Run("factor < 1", func(t *testing.T) {
		t.Parallel()

		r := NewRetryHandler(MinWaitDuration, 5*time.Second, 0.5, false)

		if !assert.Equal(t, MinFactor, r.backoff.Factor, "factor should == %d", MinFactor) {
			return
		}
		r = NewRetryHandler(MinWaitDuration, 5*time.Second, 0.98, false)
		if !assert.Equal(t, MinFactor, r.backoff.Factor, "factor should == %d", MinFactor) {
			return
		}
	})
}

func Test_RetryHandler_ProcessWithRetry(t *testing.T) {
	// Note: because RetryHandler is immutable, it is safe for concurrent tests
	r := NewRetryHandler(MinWaitDuration, 5*time.Second, 2.0, false)

	t.Run("Just won't work", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()

		result, err := r.Execute(ctx, func(context.Context) (interface{}, error) {
			// fmt.Println("in target func")
			time.Sleep(200 * time.Millisecond)
			return nil, fmt.Errorf("this just isn't going to work")
		}, nil)

		assert.Error(t, err)
		assert.Nil(t, result)

		fmt.Println(err)
	})

	t.Run("Just works", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := r.Execute(ctx, func(context.Context) (interface{}, error) {
			// fmt.Println("in target func")
			time.Sleep(200 * time.Millisecond)
			return "success", nil
		}, nil)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("Eventually works", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		attempts := 0

		result, err := r.Execute(ctx, func(context.Context) (interface{}, error) {
			attempts++

			// fmt.Println("in target func")
			time.Sleep(200 * time.Millisecond)

			if attempts < 3 {
				return nil, fmt.Errorf("temporarily failing")
			}

			return "success", nil
		}, nil)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("Cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()

		go func(cancel context.CancelFunc) {
			time.Sleep(time.Second)
			cancel()
		}(cancel)

		result, err := r.Execute(ctx, func(context.Context) (interface{}, error) {
			// fmt.Println("in target func")
			time.Sleep(200 * time.Millisecond)
			return nil, fmt.Errorf("waiting to be cancelled")
		}, nil)

		assert.Error(t, err)
		assert.Nil(t, result)

		fmt.Println(err)
	})
}
