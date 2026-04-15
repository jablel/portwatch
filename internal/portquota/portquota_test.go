package portquota

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs ...struct{ num int; proto string }) []scanner.Port {
	var ports []scanner.Port
	for _, s := range specs {
		ports = append(ports, scanner.Port{Number: s.num, Protocol: s.proto})
	}
	return ports
}

func TestCheck_NoLimits_NoViolations(t *testing.T) {
	q := New()
	ports := makePorts(
		struct{ num int; proto string }{80, "tcp"},
		struct{ num int; proto string }{443, "tcp"},
	)
	if v := q.Check(ports); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_UnderLimit_NoViolations(t *testing.T) {
	q := New()
	q.Set("tcp", 5)
	ports := makePorts(
		struct{ num int; proto string }{80, "tcp"},
		struct{ num int; proto string }{443, "tcp"},
	)
	if v := q.Check(ports); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestCheck_ExceedsLimit_ReturnsViolation(t *testing.T) {
	q := New()
	q.Set("tcp", 2)
	ports := makePorts(
		struct{ num int; proto string }{80, "tcp"},
		struct{ num int; proto string }{443, "tcp"},
		struct{ num int; proto string }{8080, "tcp"},
	)
	v := q.Check(ports)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Protocol != "tcp" || v[0].Limit != 2 || v[0].Actual != 3 {
		t.Errorf("unexpected violation: %+v", v[0])
	}
}

func TestCheck_ZeroLimit_NeverViolates(t *testing.T) {
	q := New()
	q.Set("udp", 0)
	ports := makePorts(
		struct{ num int; proto string }{53, "udp"},
		struct{ num int; proto string }{123, "udp"},
	)
	if v := q.Check(ports); len(v) != 0 {
		t.Fatalf("expected no violations with zero limit, got %v", v)
	}
}

func TestCheck_MultipleProtocols_IndependentViolations(t *testing.T) {
	q := New()
	q.Set("tcp", 1)
	q.Set("udp", 1)
	ports := makePorts(
		struct{ num int; proto string }{80, "tcp"},
		struct{ num int; proto string }{443, "tcp"},
		struct{ num int; proto string }{53, "udp"},
		struct{ num int; proto string }{123, "udp"},
	)
	v := q.Check(ports)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(v), v)
	}
}

func TestLimits_ReturnsCopy(t *testing.T) {
	q := New()
	q.Set("tcp", 10)
	l := q.Limits()
	l["tcp"] = 999
	if q.Limits()["tcp"] != 10 {
		t.Error("Limits should return a copy, not a reference")
	}
}

func TestViolation_Error(t *testing.T) {
	v := Violation{Protocol: "tcp", Limit: 3, Actual: 7}
	got := v.Error()
	if got == "" {
		t.Error("Error() should return a non-empty string")
	}
}
