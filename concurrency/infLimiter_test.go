package concurrency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InfLimiter_GetInfinityLimiter(t *testing.T) {
	assert.NotNil(t, GetInfinityLimiter(), "GetInfinityLimiter should return non-nil Limiter")
}

func Test_InfLimiter_AcquireAndRelease(t *testing.T) {
	t.Run("infLimiter acquire + release", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// infLimit doesn't really do anything, run some iterations anyway
		for i := 0; i < 100; i++ {
			if !assert.NoError(t, GetInfinityLimiter().Acquire(ctx), "infLimit should always acquire") {
				return
			}
		}

		// infLimit doesn't really do anything, run some iterations anyway
		for i := 0; i < 100; i++ {
			if !assert.NotPanics(t, func() { GetInfinityLimiter().Release() }, "infLimit release should not panic") {
				return
			}
		}
	})

	t.Run("infLimiter cancelled context", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if !assert.Error(t, GetInfinityLimiter().Acquire(ctx), "infLimit should fail on cancelled context") {
			return
		}
	})

}
