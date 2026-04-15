// Package portpolicy provides a rule-based policy engine for port change events.
//
// Policies are named rules that match ports by number range and protocol.
// They can be configured to fire on port additions, removals, or both.
// Each policy carries a severity (warn or critical) and a human-readable
// message that is included in the resulting Violation.
//
// Usage:
//
//	e := portpolicy.New()
//	_ = e.Add(portpolicy.Policy{
//		Name:     "no-telnet",
//		MinPort:  23,
//		MaxPort:  23,
//		Protocol: "tcp",
//		OnAdded:  true,
//		Severity: portpolicy.SeverityCritical,
//		Message:  "telnet must not be exposed",
//	})
//	violations := e.Evaluate(diff)
package portpolicy
