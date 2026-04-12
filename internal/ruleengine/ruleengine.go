// Package ruleengine evaluates a set of named rules against a port diff
// and returns the subset of rules that matched.
package ruleengine

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/state"
)

// Action describes what triggered a rule match.
type Action string

const (
	ActionAdded   Action = "added"
	ActionRemoved Action = "removed"
	ActionAny     Action = "any"
)

// Rule defines a condition to evaluate against a port event.
type Rule struct {
	Name     string
	Action   Action
	PortLow  int
	PortHigh int
	Protocol string // empty means any
}

// Match is a rule that fired together with the port that triggered it.
type Match struct {
	Rule Rule
	Port state.Port
}

// Engine holds a collection of rules and evaluates them against diffs.
type Engine struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns an empty Engine.
func New() *Engine {
	return &Engine{}
}

// Add registers a rule with the engine. Returns an error if the rule name is
// empty or the port range is invalid.
func (e *Engine) Add(r Rule) error {
	if r.Name == "" {
		return fmt.Errorf("ruleengine: rule name must not be empty")
	}
	if r.PortHigh < r.PortLow {
		return fmt.Errorf("ruleengine: PortHigh (%d) < PortLow (%d)", r.PortHigh, r.PortLow)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = append(e.rules, r)
	return nil
}

// Evaluate returns all matches for the given diff.
func (e *Engine) Evaluate(diff state.Diff) []Match {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var matches []Match
	for _, r := range e.rules {
		if r.Action == ActionAdded || r.Action == ActionAny {
			for _, p := range diff.Added {
				if e.portMatches(r, p) {
					matches = append(matches, Match{Rule: r, Port: p})
				}
			}
		}
		if r.Action == ActionRemoved || r.Action == ActionAny {
			for _, p := range diff.Removed {
				if e.portMatches(r, p) {
					matches = append(matches, Match{Rule: r, Port: p})
				}
			}
		}
	}
	return matches
}

func (e *Engine) portMatches(r Rule, p state.Port) bool {
	if p.Number < r.PortLow || p.Number > r.PortHigh {
		return false
	}
	if r.Protocol != "" && r.Protocol != p.Protocol {
		return false
	}
	return true
}
