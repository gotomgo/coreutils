package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// WorkRequest represents 1 unit of work to be performed by a BatchRequest
//
//	Fields
//		CorrelationID - an identifier for the work. Unique to the batch
//		Target        - the target func that performs the work
//
type WorkRequest struct {
	CorrelationID string
	Target        TargetFunc
}

// WorkResult represents the result of a WorkRequest performed as part of a
// BatchRequest
//
//	Fields
//		CorrelationID - the identifier from the WorkRequest
//		Result        - the value returned by WorkRequest.Target
//		Err			  - the error returned by WorkRequest.Target
//
type WorkResult struct {
	CorrelationID string
	Result        interface{}
	Err           error
}

// BatchResponse is the cumulative result of a BatchRequest
//
//	Fields
//		CorrelationID - the unique id from the BatchRequest
//		Results       - the result of each WorkRequest by
//			WorkRequest.CorrelationID
//
type BatchResponse struct {
	CorrelationID string
	Results       map[string]*WorkResult
}

// BatchRequest represents a collection of work to be executed as a parallel
// batch, with an aggregated result (BatchResponse)
type BatchRequest struct {
	correlationID string
	limiter       Limiter
	workRequests  []WorkRequest
}

// WithWorkRequest adds a WorkRequest to the batch
func (br *BatchRequest) WithWorkRequest(req WorkRequest) *BatchRequest {
	br.workRequests = append(br.workRequests, req)
	return br
}

// requestResultChan is used to communicate WorkResult's to the batch processor
type requestResultChan chan *WorkResult

// Execute executes the work requests in a BatchRequest
func (br *BatchRequest) Execute(ctx context.Context) (*BatchResponse, error) {
	results := make(requestResultChan, len(br.workRequests))

	// setup a response struct
	response := BatchResponse{
		CorrelationID: br.correlationID,
		Results:       map[string]*WorkResult{},
	}

	// don't modify BatchRequest so we need a var so we can substitute infinity
	// semaphoreLimiter when not specified
	limiter := br.limiter

	if limiter == nil {
		limiter = GetInfinityLimiter()
	}

	wg := sync.WaitGroup{}

	for _, req := range br.workRequests {
		wg.Add(1)
		go func(req WorkRequest, resultChan requestResultChan) {
			defer wg.Done()

			// setup a request result
			result := &WorkResult{CorrelationID: req.CorrelationID}

			// block until we are under concurrency limits, timeout, or context
			// is cancelled
			if err := limiter.Acquire(ctx); err != nil {
				// timeout or context cancelled, so update result and send to
				// results chan
				result.Err = err
				results <- result
				return
			}

			// release our acquire
			defer limiter.Release()

			// handle target panic
			defer func() {
				if err := recover(); err != nil {
					// err does not have to be type error
					if _, ok := err.(error); !ok {
						// Force err to be an error
						err = fmt.Errorf("target panic: %v", err)
					}

					// update result and send to results chan
					result.Err = err.(error)
					results <- result
				}
			}()

			// call the target and update results
			result.Result, result.Err = req.Target(ctx)
			results <- result
		}(req, results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		response.Results[result.CorrelationID] = result
	}

	return &response, nil
}
