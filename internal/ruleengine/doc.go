// Package ruleengine provides a lightweight, rule-based evaluation engine for
// port change events.
//
// Rules are defined with a name, an action trigger (added, removed, or any),
// a port range, and an optional protocol filter. When Evaluate is called with
// a state.Diff the engine returns every (Rule, Port) pair whose conditions are
// satisfied, allowing callers to drive alerting or automation logic without
// hard-coding port-specific behaviour.
//
// Usage:
//
//	eng := ruleengine.New()
//	_ = eng.Add(ruleengine.Rule{
//		Name:     "high-port-opened",
//		Action:   ruleengine.ActionAdded,
//		PortLow:  1024,
//		PortHigh: 65535,
//	})
//	matches := eng.Evaluate(diff)
package ruleengine
