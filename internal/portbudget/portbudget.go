// Package portbudget enforces a maximum number of concurrently open ports
// and reports violations when the observed count exceeds the configured limit.
package portbudget

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a budget breach.
type Violation struct {
	Limit    int
	Observed int
	Excess   []scanner.Port
}

func (v Violation) Error() string {
	return fmt.Sprintf("port budget exceeded: limit %d, observed %d (%d excess)",
		v.Limit, v.Observed, len(v.Excess))
}

// Budget enforces an upper bound on open port count.
type Budget struct {
	mu    sync.Mutex
	limit int
}

// New creates a Budget with the given limit. A limit <= 0 is treated as
// unlimited and Check will never return a violation.
func New(limit int) *Budget {
	return &Budget{limit: limit}
}

// SetLimit updates the maximum allowed port count.
func (b *Budget) SetLimit(limit int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.limit = limit
}

// Limit returns the current configured limit.
func (b *Budget) Limit() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.limit
}

// Check evaluates ports against the budget. It returns a non-nil *Violation
// when the count exceeds the limit. The Excess slice contains the ports beyond
// the limit, preserving the original order.
func (b *Budget) Check(ports []scanner.Port) *Violation {
	b.mu.Lock()
	limit := b.limit
	b.mu.Unlock()

	if limit <= 0 || len(ports) <= limit {
		return nil
	}

	return &Violation{
		Limit:    limit,
		Observed: len(ports),
		Excess:   ports[limit:],
	}
}
