package portschedule

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func at(hour int) time.Time {
	return time.Date(2024, 1, 1, hour, 0, 0, 0, time.UTC)
}

func TestObserve_LearningPhase_NoViolations(t *testing.T) {
	tr := New(3)
	p := makePort(8080, "tcp")

	// First 3 observations at hour 10 — still learning, no violations.
	for i := 0; i < 3; i++ {
		v := tr.Observe([]scanner.Port{p}, at(10))
		if len(v) != 0 {
			t.Fatalf("expected no violations during learning, got %v", v)
		}
	}
}

func TestObserve_ViolationAfterLearning(t *testing.T) {
	tr := New(2)
	p := makePort(443, "tcp")

	// Learn at hour 9.
	for i := 0; i < 2; i++ {
		tr.Observe([]scanner.Port{p}, at(9))
	}

	// Now observe at hour 23 — should be a violation.
	v := tr.Observe([]scanner.Port{p}, at(23))
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Port.Number != 443 {
		t.Errorf("expected port 443, got %d", v[0].Port.Number)
	}
	if v[0].Hour != 23 {
		t.Errorf("expected hour 23, got %d", v[0].Hour)
	}
}

func TestObserve_NoViolationWithinLearnedHour(t *testing.T) {
	tr := New(1)
	p := makePort(22, "tcp")

	tr.Observe([]scanner.Port{p}, at(14))

	v := tr.Observe([]scanner.Port{p}, at(14))
	if len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
}

func TestObserve_MultiplePortsOnlyViolatingFlagged(t *testing.T) {
	tr := New(1)
	p1 := makePort(80, "tcp")
	p2 := makePort(9000, "tcp")

	tr.Observe([]scanner.Port{p1}, at(8))
	tr.Observe([]scanner.Port{p2}, at(8))

	v := tr.Observe([]scanner.Port{p1, p2}, at(20))
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestReset_ClearsSchedules(t *testing.T) {
	tr := New(1)
	p := makePort(3306, "tcp")

	tr.Observe([]scanner.Port{p}, at(5))
	tr.Reset()

	// After reset, port is unknown again — learning restarts, no violation.
	v := tr.Observe([]scanner.Port{p}, at(22))
	if len(v) != 0 {
		t.Errorf("expected no violations after reset, got %v", v)
	}
}

func TestViolation_ErrorMessage(t *testing.T) {
	v := Violation{Port: makePort(8080, "tcp"), Hour: 3, Expected: "hours [9 10]"}
	msg := v.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
