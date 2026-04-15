// Package portquota tracks per-protocol port count quotas and reports
// violations when the number of observed open ports exceeds a configured limit.
package portquota

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a quota breach for a specific protocol.
type Violation struct {
	Protocol string
	Limit    int
	Actual   int
}

func (v Violation) Error() string {
	return fmt.Sprintf("quota exceeded for %s: limit %d, actual %d", v.Protocol, v.Limit, v.Actual)
}

// Quota holds per-protocol limits.
type Quota struct {
	mu     sync.RWMutex
	limits map[string]int // protocol -> max allowed open ports
}

// New returns a Quota with no limits set.
func New() *Quota {
	return &Quota{limits: make(map[string]int)}
}

// Set defines the maximum number of open ports allowed for the given protocol.
// A limit of zero means unlimited.
func (q *Quota) Set(protocol string, limit int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.limits[protocol] = limit
}

// Check evaluates ports against all configured quotas and returns a slice of
// Violations — one per protocol that exceeds its limit. An empty slice means
// no violations.
func (q *Quota) Check(ports []scanner.Port) []Violation {
	q.mu.RLock()
	defer q.mu.RUnlock()

	counts := make(map[string]int)
	for _, p := range ports {
		counts[p.Protocol]++
	}

	var violations []Violation
	for proto, limit := range q.limits {
		if limit <= 0 {
			continue
		}
		actual := counts[proto]
		if actual > limit {
			violations = append(violations, Violation{
				Protocol: proto,
				Limit:    limit,
				Actual:   actual,
			})
		}
	}
	return violations
}

// Limits returns a copy of the current quota limits.
func (q *Quota) Limits() map[string]int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	out := make(map[string]int, len(q.limits))
	for k, v := range q.limits {
		out[k] = v
	}
	return out
}
