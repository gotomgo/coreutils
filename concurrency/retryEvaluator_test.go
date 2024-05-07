package concurrency

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NoRetryEvaluator(t *testing.T) {
	assert.False(t, NoRetryEval(nil), "NoRetryEval should always return false")
	assert.False(t, NoRetryEval(fmt.Errorf("an error")), "NoRetryEval should always return false")
}

func Test_AlwaysRetryEvaluator(t *testing.T) {
	assert.True(t, AlwaysRetryEval(nil), "AlwaysRetryEval should always return true")
	assert.True(t, AlwaysRetryEval(fmt.Errorf("an error")), "AlwaysRetryEval should always return true")
}

func Test_AllRetryEvaluator(t *testing.T) {
	assert.False(t, AllRetryEval(NoRetryEval, AlwaysRetryEval)(nil), "expecting AllRetryEval to return false")
	assert.True(t, AllRetryEval(AlwaysRetryEval, func(error) bool {
		return true
	})(nil), "expecting AllRetryEval to return true")
}

func Test_AnyRetryEvaluator(t *testing.T) {
	assert.True(t, AnyRetryEval(NoRetryEval, AlwaysRetryEval)(nil), "expecting AnyRetryEval to return true")
	assert.False(t, AnyRetryEval(NoRetryEval, func(error) bool {
		return false
	})(nil), "expecting AnyRetryEval to return false")
}

func Test_RetryCountEvaluator(t *testing.T) {
	n := 3
	rc := RetryCountEval(n)
	for i := 0; i < n; i++ {
		assert.True(t, rc(nil), "expecting RetryCountEval to return true")
	}
	assert.False(t, rc(nil), "expecting RetryCountEval to return false after N retries")
}
