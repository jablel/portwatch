package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

// TestRetry_EventualSuccess verifies that a flaky operation succeeds once
// transient errors clear, using real (short) delays.
func TestRetry_EventualSuccess(t *testing.T) {
	cfg := retry.Config{
		MaxAttempts: 5,
		BaseDelay:   5 * time.Millisecond,
		MaxDelay:    20 * time.Millisecond,
	}
	r := retry.New(cfg)

	var calls int32
	err := r.Do(context.Background(), func() error {
		n := atomic.AddInt32(&calls, 1)
		if n < 4 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

// TestRetry_ContextTimeoutAborts ensures a deadline cancels in-flight retries.
func TestRetry_ContextTimeoutAborts(t *testing.T) {
	cfg := retry.Config{
		MaxAttempts: 10,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    200 * time.Millisecond,
	}
	r := retry.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()

	err := r.Do(ctx, func() error {
		return errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected an error due to context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("unexpected error: %v", err)
	}
}
