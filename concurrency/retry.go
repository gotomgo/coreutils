package concurrency

import (
	"context"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

const (
	// "Infinite" Duration
	InfDuration = time.Duration(1<<63 - 1)
)

// TargetFunc is the func to be executed, and is generally a closure that calls
// the actual func (with parameters)
type TargetFunc func(ctx context.Context) (interface{}, error)

// RetryHandler implements retries using parameters that control backoff
//
//	Notes
//		RetryHandler has immutable state, so it is safe for concurrent use
//
type RetryHandler struct {
	// backoff is only used for parameters. state is never changed.
	// This allows for concurrent use of the RetryHandler
	backoff backoff.Backoff
}

// NewRetryHandler creates an instance of RetryHandler with parameters
// describing the backoff policy
//
//	Notes
//		If min > max they will be swapped to create a valid range
//		If either min or max is < MinWaitDuration (5ms), the param is set to
//		  MinWaitDuration
//		If factor is < MinFactor (1.0) it is set to MinFactor
//
func NewRetryHandler(min, max time.Duration, factor float64, jitter bool) *RetryHandler {
	return (&RetryHandler{
		backoff: backoff.Backoff{
			Jitter: jitter,
		}}).withMinMax(min, max).withFactor(factor)
}

// NewRetryHandlerFromConfig creates an instance of RetryHandler using backoff
// policy specified by a RetryConfig
func NewRetryHandlerFromConfig(cfg RetryConfig) *RetryHandler {
	return (&RetryHandler{
		backoff: backoff.Backoff{
			Jitter: cfg.Jitter,
		}}).withMinMax(
		time.Duration(cfg.MinWait)*time.Millisecond,
		time.Duration(cfg.MaxWait)*time.Millisecond).
		withFactor(cfg.Factor)
}

// withMinMax sets the min/max bounds for the backoff policy
//
//	Notes
//		If min > max they will be swapped to create a valid range
//		If either min or max is < MinWaitDuration (5ms), the param is set to
//	  	  MinWaitDuration
//
func (r *RetryHandler) withMinMax(min, max time.Duration) *RetryHandler {
	if min < MinWaitDuration {
		min = MinWaitDuration
	}

	if max < MinWaitDuration {
		max = MinWaitDuration
	}

	if min > max {
		temp := max
		max = min
		min = temp
	}

	r.backoff.Min = min
	r.backoff.Max = max

	return r
}

// withFactor sets the exponential factor for the backoff policy of a
// RetryHandler
//
//	Notes
//		If factor is < MinFactor (1.0) it is set to MinFactor
//
func (r *RetryHandler) withFactor(factor float64) *RetryHandler {
	if factor < MinFactor {
		factor = MinFactor
	}

	r.backoff.Factor = factor

	return r
}

// Execute executes a target function with appropriate retry logic
//
//	Notes
//		if the target returns an error, and retryEval is nil, or is non-nil and
//		indicates the current error is transient (by returning true), then the
//		target func is retried until successful, the context deadline expires,
//		or retryEval returns false for a subsequent error
//
//		If the context does not have a deadline the target func will be retried
//		until it is successful.
//
//		If you are looking for N retries rather than a time based approach than
//		you should build retry count state into a RetryEvaluator func and
//		return false on a non-transient error, or if the maximum retires has
//		been exceeded. Retries based on retry count is non-preferred. Timeouts
//		are a better mechanism for retries as they express how long you are
//		willing to wait, a non-arbitrary value, whereas a retry count is rather
//		arbitrary in the scheme of things.
//
//		See the RetryCountEval func for a RetryEvaluator that performs N
//		retries
//
func (r *RetryHandler) Execute(
	ctx context.Context,
	target TargetFunc,
	retryEval RetryEvaluator) (result interface{}, err error) {
	// track the number of attempts for backoff calculations
	attempts := 0

	for {
		// Check if ctx is already cancelled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// execute the target func returning the result if successful
		if result, err = target(ctx); err == nil {
			return
		}

		// force result to nil JIC
		result = nil

		// evaluate error for retry?
		if retryEval != nil {
			// don't retry? return original error
			if !retryEval(err) {
				return
			}
		}

		// assume infinite wait unless context has a deadline
		ctxDeadline := InfDuration
		if deadline, ok := ctx.Deadline(); ok {
			ctxDeadline = time.Until(deadline)

			// check the deadline now, before waiting
			if ctxDeadline <= 0 {
				err = errors.Wrap(err, "context deadline exceeded")
				return
			}
		}

		var waitTime time.Duration

		// get the wait specified for this retry attempt
		// (1st retry is attempt 0)
		delay := r.backoff.ForAttempt(float64(attempts))

		// use the shorter waitTime
		if delay < ctxDeadline {
			waitTime = delay
		} else {
			waitTime = ctxDeadline
		}

		// wait for the backoff period to expire or the context to be cancelled
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), err.Error())

			// Note: this is an alloc/dealloc per iteration. We can create a
			// timer outside of the for loop and reuse that timer on each
			// iteration
		case <-time.After(waitTime):
			// run target again
			attempts++
		}
	}
}
