package concurrency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RetryConfig_Default(t *testing.T) {
	config := NewDefaultRetryConfig()

	assert.Equal(t, DefaultMinWait, config.MinWait)
	assert.Equal(t, DefaultMaxWait, config.MaxWait)
	assert.Equal(t, DefaultFactor, config.Factor)
	assert.Equal(t, DefaultJitter, config.Jitter)
}
