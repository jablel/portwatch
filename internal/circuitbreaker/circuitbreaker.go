// Package circuitbreaker implements a simple circuit breaker that opens after
// a threshold of consecutive errors and resets after a cooldown period.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// Breaker is a circuit breaker that tracks consecutive failures.
type Breaker struct {
	mu          sync.Mutex
	threshold   int
	cooldown    time.Duration
	failures    int
	openedAt    time.Time
	now         func() time.Time
}

// New creates a Breaker that opens after threshold consecutive errors and
// resets after the given cooldown duration.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		now:       time.Now,
	}
}

// Allow returns nil if the call is permitted, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.failures >= b.threshold {
		if b.cooldown > 0 && b.now().Sub(b.openedAt) >= b.cooldown {
			b.failures = 0
		} else {
			return ErrOpen
		}
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
}

// RecordFailure increments the failure counter and opens the breaker if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures == b.threshold {
		b.openedAt = b.now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.failures >= b.threshold {
		if b.cooldown > 0 && b.now().Sub(b.openedAt) >= b.cooldown {
			return StateClosed
		}
		return StateOpen
	}
	return StateClosed
}
