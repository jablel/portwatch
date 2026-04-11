// Package circuitbreaker provides a thread-safe circuit breaker for portwatch.
//
// A Breaker transitions from closed (allowing calls) to open (rejecting calls)
// after a configurable number of consecutive failures. Once open, it remains
// open until a cooldown period has elapsed, after which the next call to Allow
// will succeed and the failure counter is reset.
//
// Typical usage:
//
//	br := circuitbreaker.New(5, 30*time.Second)
//
//	if err := br.Allow(); err != nil {
//	    // circuit is open — skip the operation
//	    return err
//	}
//	if err := doWork(); err != nil {
//	    br.RecordFailure()
//	    return err
//	}
//	br.RecordSuccess()
package circuitbreaker
