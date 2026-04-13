package portfence_test

import (
	"testing"

	"github.com/user/portwatch/internal/portfence"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestCheck_AllAllowed_NoViolations(t *testing.T) {
	f := portfence.New(false)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	if v := f.Check(ports); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_BlockedPort_ReturnsViolation(t *testing.T) {
	f := portfence.New(false)
	f.Block(makePort(8080, "tcp"))
	ports := []scanner.Port{makePort(80, "tcp"), makePort(8080, "tcp")}
	v := f.Check(ports)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Kind != portfence.ViolationBlocked {
		t.Errorf("expected ViolationBlocked, got %s", v[0].Kind)
	}
	if v[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", v[0].Port.Number)
	}
}

func TestCheck_StrictMode_UnknownPortViolates(t *testing.T) {
	f := portfence.New(true)
	f.Allow(makePort(443, "tcp"))
	ports := []scanner.Port{makePort(443, "tcp"), makePort(9000, "tcp")}
	v := f.Check(ports)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Kind != portfence.ViolationNotAllowed {
		t.Errorf("expected ViolationNotAllowed, got %s", v[0].Kind)
	}
}

func TestCheck_StrictMode_BlockedTakesPrecedence(t *testing.T) {
	f := portfence.New(true)
	f.Allow(makePort(80, "tcp"))
	f.Block(makePort(80, "tcp"))
	v := f.Check([]scanner.Port{makePort(80, "tcp")})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Kind != portfence.ViolationBlocked {
		t.Errorf("expected ViolationBlocked, got %s", v[0].Kind)
	}
}

func TestViolation_String(t *testing.T) {
	v := portfence.Violation{
		Port: makePort(22, "tcp"),
		Kind: portfence.ViolationBlocked,
	}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string from Violation.String()")
	}
}

func TestCheck_EmptyPorts_NoViolations(t *testing.T) {
	f := portfence.New(true)
	if v := f.Check(nil); len(v) != 0 {
		t.Fatalf("expected no violations for empty input, got %d", len(v))
	}
}
