package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_ClosedBelowThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordSuccess_ResetsFailures(t *testing.T) {
	b := circuitbreaker.New(2, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatal("expected closed after success")
	}
}

func TestAllow_ResetsAfterCooldown(t *testing.T) {
	now := time.Now()
	b := circuitbreaker.New(1, 50*time.Millisecond)
	b.(*struct{ *circuitbreaker.Breaker }) // type assertion not needed; use exported clock hook via New
	_ = b

	// Use the package-level test helper instead.
	b2 := circuitbreaker.NewWithClock(1, 50*time.Millisecond, func() time.Time { return now })
	b2.RecordFailure()
	if err := b2.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatal("expected open")
	}
	// advance clock past cooldown
	now = now.Add(100 * time.Millisecond)
	if err := b2.Allow(); err != nil {
		t.Fatalf("expected closed after cooldown, got %v", err)
	}
}

func TestState_String(t *testing.T) {
	if circuitbreaker.StateClosed.String() != "closed" {
		t.Fatal("wrong closed string")
	}
	if circuitbreaker.StateOpen.String() != "open" {
		t.Fatal("wrong open string")
	}
}

func TestAllow_ZeroCooldownNeverResets(t *testing.T) {
	b := circuitbreaker.New(1, 0)
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatal("expected open with zero cooldown")
	}
}
