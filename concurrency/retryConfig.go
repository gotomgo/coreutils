package concurrency

import "time"

const (
	// MinWaitDuration is the minimum configurable wait time for backoff policies
	MinWaitDuration = 5 * time.Millisecond

	// MinFactor is the minimum exponential factor for backoff policies
	MinFactor = 1.0

	// DefaultMinWait is the default duration for MinWait expressed as
	// Milliseconds
	DefaultMinWait = 50

	// DefaultMaxWait is the default duration for MaxWait expressed as
	// Milliseconds
	DefaultMaxWait = 5000

	// DefaultFactor is the default exponential factor for backoff policies
	DefaultFactor = 2.0

	// DefaultJitter indicates the default setting for jitter
	DefaultJitter = true
)

// RetryConfig is a configuration for retry backoff
type RetryConfig struct {
	MinWait int     `yaml:"minWait" json:"minWait"`
	MaxWait int     `yaml:"maxWait" json:"maxWait"`
	Factor  float64 `yaml:"factor" json:"factor"`
	Jitter  bool    `yaml:"jitter" json:"jitter"`
}

// NewDefaultRetryConfig creates an instance of RetryConfig with default values
//
//	Default Field Values
//		MinWait - DefaultMinWait	(50 milliseconds)
//		MaxWait - DefaultMaxWait	(5000 milliseconds)
//		Factor  - DefaultFactor		(2.0, squared exponential growth)
//		Jitter  - DefaultJitter		(jitter is enabled by default)
//
func NewDefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MinWait: DefaultMinWait,
		MaxWait: DefaultMaxWait,
		Factor:  DefaultFactor,
		Jitter:  DefaultJitter,
	}
}
