// Package retry provides exponential-backoff retry logic for use throughout
// portwatch.
//
// Basic usage:
//
//	r := retry.New(retry.Default())
//	err := r.Do(ctx, func() error {
//		return doSomethingFallible()
//	})
//	if errors.Is(err, retry.ErrMaxAttempts) {
//		// all attempts exhausted
//	}
//
// The Retryer backs off exponentially between attempts, doubling the delay
// each time up to MaxDelay. The context is checked before every attempt and
// during each sleep, so cancellations are acted on promptly.
package retry
