// Package retry provides a simple exponential-backoff retry mechanism
// for transient errors encountered during port scanning and alerting.
package retry

import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds the parameters that control retry behaviour.
type Config struct {
	// MaxAttempts is the total number of tries (including the first).
	MaxAttempts int
	// BaseDelay is the wait time before the second attempt.
	BaseDelay time.Duration
	// MaxDelay caps the exponential growth of the delay.
	MaxDelay time.Duration
}

// Default returns a Config suitable for most in-process retries.
func Default() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
	}
}

// Retryer executes operations with automatic retries.
type Retryer struct {
	cfg   Config
	sleep func(time.Duration) // injectable for tests
}

// New creates a Retryer with the provided Config.
func New(cfg Config) *Retryer {
	return &Retryer{
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

// Do calls fn up to MaxAttempts times, backing off exponentially between
// attempts. It returns nil on the first success, ErrMaxAttempts when all
// attempts fail, or ctx.Err() if the context is cancelled.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.cfg.BaseDelay
	var lastErr error

	for attempt := 0; attempt < r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if lastErr = fn(); lastErr == nil {
			return nil
		}

		if attempt < r.cfg.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay *= 2
			if delay > r.cfg.MaxDelay {
				delay = r.cfg.MaxDelay
			}
		}
	}

	return ErrMaxAttempts
}
