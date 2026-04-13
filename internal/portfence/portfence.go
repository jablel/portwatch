// Package portfence enforces port access policies by comparing observed
// ports against a configured allowlist and blocklist, emitting policy
// violations for any port that breaks the rules.
package portfence

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// ViolationKind describes why a port violated policy.
type ViolationKind string

const (
	ViolationBlocked    ViolationKind = "blocked"     // port is explicitly blocked
	ViolationNotAllowed ViolationKind = "not_allowed" // port is not in the allowlist
)

// Violation records a single policy breach.
type Violation struct {
	Port scanner.Port
	Kind ViolationKind
}

func (v Violation) String() string {
	return fmt.Sprintf("%s violates policy: %s", v.Port, v.Kind)
}

// Fence evaluates ports against allow and block lists.
type Fence struct {
	mu        sync.RWMutex
	allowlist map[string]struct{}
	blocklist map[string]struct{}
	strictMode bool // if true, ports not in allowlist are violations
}

// New creates a Fence. When strictMode is true every port must appear in the
// allowlist; when false only explicitly blocked ports produce violations.
func New(strictMode bool) *Fence {
	return &Fence{
		allowlist:  make(map[string]struct{}),
		blocklist:  make(map[string]struct{}),
		strictMode: strictMode,
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// Allow adds a port to the allowlist.
func (f *Fence) Allow(p scanner.Port) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowlist[portKey(p)] = struct{}{}
}

// Block adds a port to the blocklist.
func (f *Fence) Block(p scanner.Port) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blocklist[portKey(p)] = struct{}{}
}

// Check evaluates a slice of ports and returns any policy violations.
func (f *Fence) Check(ports []scanner.Port) []Violation {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var violations []Violation
	for _, p := range ports {
		k := portKey(p)
		if _, blocked := f.blocklist[k]; blocked {
			violations = append(violations, Violation{Port: p, Kind: ViolationBlocked})
			continue
		}
		if f.strictMode {
			if _, allowed := f.allowlist[k]; !allowed {
				violations = append(violations, Violation{Port: p, Kind: ViolationNotAllowed})
			}
		}
	}
	return violations
}
