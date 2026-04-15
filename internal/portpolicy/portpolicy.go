// Package portpolicy evaluates a set of named policies against a port diff,
// producing a list of violations for ports that breach defined rules.
package portpolicy

import (
	"fmt"
	"sync"

	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

// Severity classifies how serious a policy violation is.
type Severity string

const (
	SeverityWarn     Severity = "warn"
	SeverityCritical Severity = "critical"
)

// Policy describes a rule applied to newly-opened or closed ports.
type Policy struct {
	Name     string
	MinPort  int
	MaxPort  int
	Protocol string // "tcp", "udp", or "" for any
	OnAdded  bool   // trigger when a port is added
	OnRemoved bool  // trigger when a port is removed
	Severity Severity
	Message  string
}

// Violation is produced when a port matches a policy.
type Violation struct {
	Policy   string
	Port     scanner.Port
	Severity Severity
	Message  string
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s – %s (%s)", v.Severity, v.Policy, v.Message, v.Port)
}

// Evaluator holds registered policies and evaluates them against diffs.
type Evaluator struct {
	mu       sync.RWMutex
	policies []Policy
}

// New returns an empty Evaluator.
func New() *Evaluator {
	return &Evaluator{}
}

// Add registers a policy. Returns an error if the port range is invalid.
func (e *Evaluator) Add(p Policy) error {
	if p.MinPort < 0 || p.MaxPort > 65535 || p.MinPort > p.MaxPort {
		return fmt.Errorf("portpolicy: invalid range %d-%d for policy %q", p.MinPort, p.MaxPort, p.Name)
	}
	if p.Name == "" {
		return fmt.Errorf("portpolicy: policy name must not be empty")
	}
	e.mu.Lock()
	e.policies = append(e.policies, p)
	e.mu.Unlock()
	return nil
}

// Evaluate checks the diff against all registered policies and returns violations.
func (e *Evaluator) Evaluate(diff state.Diff) []Violation {
	e.mu.RLock()
	policies := make([]Policy, len(e.policies))
	copy(policies, e.policies)
	e.mu.RUnlock()

	var violations []Violation
	for _, pol := range policies {
		if pol.OnAdded {
			for _, p := range diff.Added {
				if matches(pol, p) {
					violations = append(violations, Violation{
						Policy: pol.Name, Port: p,
						Severity: pol.Severity, Message: pol.Message,
					})
				}
			}
		}
		if pol.OnRemoved {
			for _, p := range diff.Removed {
				if matches(pol, p) {
					violations = append(violations, Violation{
						Policy: pol.Name, Port: p,
						Severity: pol.Severity, Message: pol.Message,
					})
				}
			}
		}
	}
	return violations
}

func matches(pol Policy, p scanner.Port) bool {
	if p.Number < pol.MinPort || p.Number > pol.MaxPort {
		return false
	}
	if pol.Protocol != "" && p.Protocol != pol.Protocol {
		return false
	}
	return true
}
