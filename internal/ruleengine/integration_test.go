package ruleengine_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/ruleengine"
	"github.com/user/portwatch/internal/state"
)

// TestRuleEngine_ConcurrentEvaluate ensures no data races when multiple
// goroutines call Evaluate simultaneously.
func TestRuleEngine_ConcurrentEvaluate(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "all", Action: ruleengine.ActionAny,
		PortLow: 1, PortHigh: 65535,
	})

	diff := makeDiff(
		[]state.Port{makePort(8080, "tcp")},
		[]state.Port{makePort(9090, "tcp")},
	)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			matches := eng.Evaluate(diff)
			if len(matches) != 2 {
				t.Errorf("want 2 matches, got %d", len(matches))
			}
		}()
	}
	wg.Wait()
}

// TestRuleEngine_MultipleRulesCanMatchSamePort verifies that more than one
// rule may fire for the same port when their ranges overlap.
func TestRuleEngine_MultipleRulesCanMatchSamePort(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "broad", Action: ruleengine.ActionAdded,
		PortLow: 1, PortHigh: 65535,
	})
	_ = eng.Add(ruleengine.Rule{
		Name: "narrow", Action: ruleengine.ActionAdded,
		PortLow: 443, PortHigh: 443,
	})

	diff := makeDiff([]state.Port{makePort(443, "tcp")}, nil)
	matches := eng.Evaluate(diff)
	if len(matches) != 2 {
		t.Fatalf("want 2 matches (both rules), got %d", len(matches))
	}

	names := map[string]bool{}
	for _, m := range matches {
		names[m.Rule.Name] = true
	}
	if !names["broad"] || !names["narrow"] {
		t.Errorf("expected both 'broad' and 'narrow' to match, got %v", names)
	}
}
