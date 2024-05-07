package concurrency

// RetryEvaluator is used to evaluate an error and determine if a Retry is
// warranted (i.e., the error is transient)
type RetryEvaluator func(error) bool

// NoRetryEval is a RetryEvaluator that returns false (indicating not to retry)
func NoRetryEval(error) bool {
	return false
}

// AlwaysRetryEval is a RetryEvaluator that returns true (indicating to retry)
func AlwaysRetryEval(error) bool {
	return true
}

// RetryCountEval returns a RetryEvaluator that retries N times and then
// cancels further retries
//
//	Notes
//		A value of n=0 means no retries
//		Any initial value of n that is < 0 is considered 0
//
func RetryCountEval(n int) func(error) bool {
	// negative values are not allowed
	if n < 0 {
		n = 0
	}

	return func(error) bool {
		if n <= 0 {
			return false
		}

		n--
		return true
	}
}

// AllRetryEval allows one or more RetryEvaluators to be executed in series
// to determine if a retry should occur. All evaluations must return true
//
//	Notes
//		Each RetryEvaluator has the power to Veto a retry, and a single veto
//		prevents the retry
//
func AllRetryEval(retries ...RetryEvaluator) func(error) bool {
	return func(err error) bool {
		for _, rh := range retries {
			if !rh(err) {
				return false
			}
		}

		return true
	}
}

// AnyRetryEval allows one or more RetryEvaluators to be executed in
// series to determine if a retry should occur. A single evaluator must return
// true
//
//	Notes
//		Each RetryEvaluator has the power to Veto a retry, and a single veto
//		prevents the retry
//
func AnyRetryEval(retries ...RetryEvaluator) func(error) bool {
	return func(err error) bool {
		for _, rh := range retries {
			if rh(err) {
				return true
			}
		}

		return false
	}
}
