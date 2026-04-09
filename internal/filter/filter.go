package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines criteria for filtering ports.
type Rule struct {
	MinPort uint16
	MaxPort uint16
	Protocols []string
}

// Filter applies rules to a list of ports, returning only those that match.
type Filter struct {
	rules []Rule
}

// New creates a Filter with the given rules.
// If no rules are provided, all ports are accepted.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns the subset of ports that satisfy at least one rule.
// If the filter has no rules, all ports are returned unchanged.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	if len(f.rules) == 0 {
		return ports
	}

	var result []scanner.Port
	for _, p := range ports {
		if f.matches(p) {
			result = append(result, p)
		}
	}
	return result
}

// Exclude returns ports that do NOT match any rule (inverse of Apply).
func (f *Filter) Exclude(ports []scanner.Port) []scanner.Port {
	if len(f.rules) == 0 {
		return ports
	}

	var result []scanner.Port
	for _, p := range ports {
		if !f.matches(p) {
			result = append(result, p)
		}
	}
	return result
}

func (f *Filter) matches(p scanner.Port) bool {
	for _, r := range f.rules {
		if p.Number >= r.MinPort && p.Number <= r.MaxPort {
			if len(r.Protocols) == 0 {
				return true
			}
			for _, proto := range r.Protocols {
				if proto == p.Protocol {
					return true
				}
			}
		}
	}
	return false
}
