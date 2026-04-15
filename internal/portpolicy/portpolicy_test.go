package portpolicy

import (
	"testing"

	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func makeDiff(added, removed []scanner.Port) state.Diff {
	return state.Diff{Added: added, Removed: removed}
}

func TestAdd_InvalidRange(t *testing.T) {
	e := New()
	err := e.Add(Policy{Name: "bad", MinPort: 100, MaxPort: 50, OnAdded: true, Severity: SeverityWarn})
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	e := New()
	err := e.Add(Policy{Name: "", MinPort: 0, MaxPort: 1024, OnAdded: true, Severity: SeverityWarn})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestEvaluate_NoViolations_WhenNoPolicies(t *testing.T) {
	e := New()
	diff := makeDiff([]scanner.Port{makePort(80, "tcp")}, nil)
	if v := e.Evaluate(diff); len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestEvaluate_MatchesAddedPort(t *testing.T) {
	e := New()
	_ = e.Add(Policy{
		Name: "no-http", MinPort: 80, MaxPort: 80,
		Protocol: "tcp", OnAdded: true,
		Severity: SeverityCritical, Message: "HTTP opened",
	})
	diff := makeDiff([]scanner.Port{makePort(80, "tcp")}, nil)
	v := e.Evaluate(diff)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Policy != "no-http" {
		t.Errorf("unexpected policy name %q", v[0].Policy)
	}
	if v[0].Severity != SeverityCritical {
		t.Errorf("unexpected severity %q", v[0].Severity)
	}
}

func TestEvaluate_MatchesRemovedPort(t *testing.T) {
	e := New()
	_ = e.Add(Policy{
		Name: "ssh-must-stay", MinPort: 22, MaxPort: 22,
		Protocol: "tcp", OnRemoved: true,
		Severity: SeverityWarn, Message: "SSH closed",
	})
	diff := makeDiff(nil, []scanner.Port{makePort(22, "tcp")})
	v := e.Evaluate(diff)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestEvaluate_ProtocolMismatch_NoViolation(t *testing.T) {
	e := New()
	_ = e.Add(Policy{
		Name: "tcp-only", MinPort: 443, MaxPort: 443,
		Protocol: "tcp", OnAdded: true,
		Severity: SeverityWarn,
	})
	diff := makeDiff([]scanner.Port{makePort(443, "udp")}, nil)
	if v := e.Evaluate(diff); len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestEvaluate_MultipleMatchingPolicies(t *testing.T) {
	e := New()
	_ = e.Add(Policy{Name: "p1", MinPort: 0, MaxPort: 1023, OnAdded: true, Severity: SeverityWarn})
	_ = e.Add(Policy{Name: "p2", MinPort: 80, MaxPort: 80, OnAdded: true, Severity: SeverityCritical})
	diff := makeDiff([]scanner.Port{makePort(80, "tcp")}, nil)
	v := e.Evaluate(diff)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{
		Policy: "test", Port: makePort(8080, "tcp"),
		Severity: SeverityWarn, Message: "unexpected port",
	}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
