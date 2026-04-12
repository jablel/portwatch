package ruleengine_test

import (
	"testing"

	"github.com/user/portwatch/internal/ruleengine"
	"github.com/user/portwatch/internal/state"
)

func makeDiff(added, removed []state.Port) state.Diff {
	return state.Diff{Added: added, Removed: removed}
}

func makePort(number int, proto string) state.Port {
	return state.Port{Number: number, Protocol: proto}
}

func TestAdd_InvalidName(t *testing.T) {
	eng := ruleengine.New()
	if err := eng.Add(ruleengine.Rule{PortLow: 80, PortHigh: 80}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_InvalidRange(t *testing.T) {
	eng := ruleengine.New()
	err := eng.Add(ruleengine.Rule{Name: "bad", PortLow: 100, PortHigh: 80})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestEvaluate_MatchesAddedPort(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "http", Action: ruleengine.ActionAdded,
		PortLow: 80, PortHigh: 80, Protocol: "tcp",
	})
	diff := makeDiff([]state.Port{makePort(80, "tcp")}, nil)
	matches := eng.Evaluate(diff)
	if len(matches) != 1 {
		t.Fatalf("want 1 match, got %d", len(matches))
	}
	if matches[0].Rule.Name != "http" {
		t.Errorf("unexpected rule name: %s", matches[0].Rule.Name)
	}
}

func TestEvaluate_NoMatchOnWrongAction(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "only-removed", Action: ruleengine.ActionRemoved,
		PortLow: 80, PortHigh: 80,
	})
	diff := makeDiff([]state.Port{makePort(80, "tcp")}, nil)
	if matches := eng.Evaluate(diff); len(matches) != 0 {
		t.Fatalf("want 0 matches, got %d", len(matches))
	}
}

func TestEvaluate_ActionAnyMatchesBoth(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "any-443", Action: ruleengine.ActionAny,
		PortLow: 443, PortHigh: 443,
	})
	diff := makeDiff(
		[]state.Port{makePort(443, "tcp")},
		[]state.Port{makePort(443, "tcp")},
	)
	matches := eng.Evaluate(diff)
	if len(matches) != 2 {
		t.Fatalf("want 2 matches, got %d", len(matches))
	}
}

func TestEvaluate_ProtocolFilterExcludes(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "tcp-only", Action: ruleengine.ActionAdded,
		PortLow: 53, PortHigh: 53, Protocol: "tcp",
	})
	diff := makeDiff([]state.Port{makePort(53, "udp")}, nil)
	if matches := eng.Evaluate(diff); len(matches) != 0 {
		t.Fatalf("want 0 matches, got %d", len(matches))
	}
}

func TestEvaluate_RangeMatch(t *testing.T) {
	eng := ruleengine.New()
	_ = eng.Add(ruleengine.Rule{
		Name: "high", Action: ruleengine.ActionAdded,
		PortLow: 1024, PortHigh: 65535,
	})
	ports := []state.Port{
		makePort(1024, "tcp"),
		makePort(32000, "tcp"),
		makePort(65535, "tcp"),
		makePort(80, "tcp"), // out of range
	}
	diff := makeDiff(ports, nil)
	matches := eng.Evaluate(diff)
	if len(matches) != 3 {
		t.Fatalf("want 3 matches, got %d", len(matches))
	}
}
