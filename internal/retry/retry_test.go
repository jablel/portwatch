package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func fastRetryer(cfg Config) *Retryer {
	r := New(cfg)
	r.sleep = func(time.Duration) {} // no-op so tests are instant
	return r
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	r := fastRetryer(Default())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	r := fastRetryer(Default())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retry, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsErrMaxAttempts(t *testing.T) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := fastRetryer(cfg)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_CancelledContextStopsRetries(t *testing.T) {
	r := fastRetryer(Default())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := r.Do(ctx, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls after cancelled context, got %d", calls)
	}
}

func TestDo_SingleAttemptNeverRetries(t *testing.T) {
	cfg := Config{MaxAttempts: 1, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := fastRetryer(cfg)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls)
	}
}
