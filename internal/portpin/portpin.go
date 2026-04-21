// Package portpin tracks which ports have been explicitly "pinned" as
// expected and flags any that appear or disappear outside that set.
package portpin

import (
	"fmt"
	"sync"

	"portwatch/internal/scanner"
)

// Violation describes a port that violated the pinned set.
type Violation struct {
	Port   scanner.Port
	Reason string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: %s", v.Port, v.Reason)
}

// Pinner holds the set of pinned (expected) ports and checks observed
// port lists against it.
type Pinner struct {
	mu     sync.RWMutex
	pinned map[string]scanner.Port
}

// New returns an empty Pinner.
func New() *Pinner {
	return &Pinner{pinned: make(map[string]scanner.Port)}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// Pin marks a port as expected.
func (p *Pinner) Pin(port scanner.Port) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pinned[portKey(port)] = port
}

// Unpin removes a port from the expected set.
func (p *Pinner) Unpin(port scanner.Port) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pinned, portKey(port))
}

// Pinned returns a snapshot of all currently pinned ports.
func (p *Pinner) Pinned() []scanner.Port {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]scanner.Port, 0, len(p.pinned))
	for _, port := range p.pinned {
		out = append(out, port)
	}
	return out
}

// Check compares observed ports against the pinned set and returns
// violations for unexpected ports and missing pinned ports.
func (p *Pinner) Check(observed []scanner.Port) []Violation {
	p.mu.RLock()
	defer p.mu.RUnlock()

	obsMap := make(map[string]scanner.Port, len(observed))
	for _, port := range observed {
		obsMap[portKey(port)] = port
	}

	var violations []Violation

	// unexpected ports: observed but not pinned
	for key, port := range obsMap {
		if _, ok := p.pinned[key]; !ok {
			violations = append(violations, Violation{Port: port, Reason: "unexpected port observed"})
		}
	}

	// missing ports: pinned but not observed
	for key, port := range p.pinned {
		if _, ok := obsMap[key]; !ok {
			violations = append(violations, Violation{Port: port, Reason: "pinned port missing"})
		}
	}

	return violations
}
