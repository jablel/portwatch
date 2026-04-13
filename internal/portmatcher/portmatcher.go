// Package portmatcher provides pattern-based matching for ports using
// glob-style rules such as "80", "8000-8999", or "*:tcp".
package portmatcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Rule represents a single match pattern.
type Rule struct {
	raw      string
	minPort  uint16
	maxPort  uint16
	protocol string // "tcp", "udp", or "" for any
}

// Matcher holds a set of rules and matches ports against them.
type Matcher struct {
	rules []Rule
}

// New returns an empty Matcher.
func New() *Matcher {
	return &Matcher{}
}

// Add parses and appends a rule string. Accepted formats:
//   - "80"         — exact port, any protocol
//   - "8000-8999"  — port range, any protocol
//   - "443:tcp"    — exact port, specific protocol
//   - "8000-8999:udp" — range with protocol
func (m *Matcher) Add(pattern string) error {
	r, err := parseRule(pattern)
	if err != nil {
		return fmt.Errorf("portmatcher: invalid pattern %q: %w", pattern, err)
	}
	m.rules = append(m.rules, r)
	return nil
}

// Match returns true if port satisfies at least one rule.
func (m *Matcher) Match(p scanner.Port) bool {
	for _, r := range m.rules {
		if r.matches(p) {
			return true
		}
	}
	return false
}

// MatchAny returns all ports from the slice that match at least one rule.
func (m *Matcher) MatchAny(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if m.Match(p) {
			out = append(out, p)
		}
	}
	return out
}

func (r Rule) matches(p scanner.Port) bool {
	if r.protocol != "" && !strings.EqualFold(r.protocol, p.Protocol) {
		return false
	}
	return p.Number >= r.minPort && p.Number <= r.maxPort
}

func parseRule(pattern string) (Rule, error) {
	r := Rule{raw: pattern}
	parts := strings.SplitN(pattern, ":", 2)
	portPart := parts[0]
	if len(parts) == 2 {
		proto := strings.ToLower(parts[1])
		if proto != "tcp" && proto != "udp" {
			return r, fmt.Errorf("unknown protocol %q", parts[1])
		}
		r.protocol = proto
	}
	if idx := strings.Index(portPart, "-"); idx != -1 {
		lo, err1 := parsePort(portPart[:idx])
		hi, err2 := parsePort(portPart[idx+1:])
		if err1 != nil || err2 != nil {
			return r, fmt.Errorf("invalid range %q", portPart)
		}
		if lo > hi {
			return r, fmt.Errorf("range start %d > end %d", lo, hi)
		}
		r.minPort, r.maxPort = lo, hi
	} else {
		v, err := parsePort(portPart)
		if err != nil {
			return r, err
		}
		r.minPort, r.maxPort = v, v
	}
	return r, nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q", s)
	}
	return uint16(n), nil
}
