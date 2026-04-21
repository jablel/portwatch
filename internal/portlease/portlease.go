// Package portlease tracks temporary port reservations with TTL-based expiry.
// A lease represents a port that is expected to be open for a bounded duration;
// violations are raised when a port outlives its lease or disappears before it expires.
package portlease

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a lease breach.
type Violation struct {
	Port   scanner.Port
	Reason string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: %s", v.Port, v.Reason)
}

type lease struct {
	port      scanner.Port
	grantedAt time.Time
	ttl       time.Duration
}

func (l lease) expired(now time.Time) bool {
	return l.ttl > 0 && now.After(l.grantedAt.Add(l.ttl))
}

// Tracker manages port leases.
type Tracker struct {
	mu     sync.Mutex
	leases map[string]lease
	now    func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		leases: make(map[string]lease),
		now:    time.Now,
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s/%d", p.Proto, p.Number)
}

// Grant registers a lease for p that expires after ttl.
// A zero ttl means the lease never expires on its own.
func (t *Tracker) Grant(p scanner.Port, ttl time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.leases[portKey(p)] = lease{port: p, grantedAt: t.now(), ttl: ttl}
}

// Revoke removes an existing lease.
func (t *Tracker) Revoke(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.leases, portKey(p))
}

// Check evaluates current against active leases and returns any violations.
// A violation is raised when a leased port has exceeded its TTL while still open.
func (t *Tracker) Check(current []scanner.Port) []Violation {
	now := t.now()
	present := make(map[string]struct{}, len(current))
	for _, p := range current {
		present[portKey(p)] = struct{}{}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var violations []Violation
	for k, l := range t.leases {
		_, open := present[k]
		if open && l.expired(now) {
			violations = append(violations, Violation{
				Port:   l.port,
				Reason: fmt.Sprintf("lease expired after %s", l.ttl),
			})
		}
	}
	return violations
}
