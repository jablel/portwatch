// Package portguard enforces an allowlist of expected ports and flags
// any port that falls outside the approved set as a violation.
package portguard

import (
	"fmt"
	"sync"

	"portwatch/internal/scanner"
)

// Violation describes a port that was not in the allowlist.
type Violation struct {
	Port   scanner.Port
	Reason string
}

func (v Violation) String() string {
	return fmt.Sprintf("violation: %s – %s", v.Port, v.Reason)
}

// Guard holds the set of allowed ports and evaluates scanned results.
type Guard struct {
	mu      sync.RWMutex
	allowed map[string]struct{}
}

// New returns a Guard pre-loaded with the given allowlist.
func New(allowed []scanner.Port) *Guard {
	g := &Guard{allowed: make(map[string]struct{}, len(allowed))}
	for _, p := range allowed {
		g.allowed[key(p)] = struct{}{}
	}
	return g
}

// Allow adds a port to the allowlist.
func (g *Guard) Allow(p scanner.Port) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.allowed[key(p)] = struct{}{}
}

// Revoke removes a port from the allowlist.
func (g *Guard) Revoke(p scanner.Port) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.allowed, key(p))
}

// Check evaluates a slice of ports and returns any violations.
func (g *Guard) Check(ports []scanner.Port) []Violation {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var violations []Violation
	for _, p := range ports {
		if _, ok := g.allowed[key(p)]; !ok {
			violations = append(violations, Violation{
				Port:   p,
				Reason: "not in allowlist",
			})
		}
	}
	return violations
}

// IsAllowed reports whether a single port is in the allowlist.
func (g *Guard) IsAllowed(p scanner.Port) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.allowed[key(p)]
	return ok
}

func key(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}
