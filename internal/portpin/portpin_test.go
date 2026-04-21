package portpin

import (
	"testing"

	"portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestCheck_AllPinned_NoViolations(t *testing.T) {
	p := New()
	p.Pin(makePort(80, "tcp"))
	p.Pin(makePort(443, "tcp"))

	violations := p.Check([]scanner.Port{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
	})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_UnexpectedPort_ReturnsViolation(t *testing.T) {
	p := New()
	p.Pin(makePort(80, "tcp"))

	violations := p.Check([]scanner.Port{
		makePort(80, "tcp"),
		makePort(9999, "tcp"),
	})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Number != 9999 {
		t.Errorf("expected violation for port 9999, got %d", violations[0].Port.Number)
	}
}

func TestCheck_MissingPinnedPort_ReturnsViolation(t *testing.T) {
	p := New()
	p.Pin(makePort(80, "tcp"))
	p.Pin(makePort(443, "tcp"))

	violations := p.Check([]scanner.Port{makePort(80, "tcp")})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Number != 443 {
		t.Errorf("expected violation for port 443, got %d", violations[0].Port.Number)
	}
}

func TestUnpin_RemovesPort(t *testing.T) {
	p := New()
	p.Pin(makePort(80, "tcp"))
	p.Unpin(makePort(80, "tcp"))

	violations := p.Check([]scanner.Port{})
	if len(violations) != 0 {
		t.Fatalf("expected no violations after unpin, got %d", len(violations))
	}
}

func TestPinned_ReturnsAllPinnedPorts(t *testing.T) {
	p := New()
	p.Pin(makePort(22, "tcp"))
	p.Pin(makePort(53, "udp"))

	pinned := p.Pinned()
	if len(pinned) != 2 {
		t.Fatalf("expected 2 pinned ports, got %d", len(pinned))
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Port: makePort(8080, "tcp"), Reason: "unexpected port observed"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty violation string")
	}
}
