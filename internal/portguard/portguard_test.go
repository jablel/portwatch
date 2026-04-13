package portguard_test

import (
	"testing"

	"portwatch/internal/portguard"
	"portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestCheck_AllAllowed_NoViolations(t *testing.T) {
	g := portguard.New([]scanner.Port{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
	})

	violations := g.Check([]scanner.Port{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
	})

	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_UnknownPort_ReturnsViolation(t *testing.T) {
	g := portguard.New([]scanner.Port{makePort(80, "tcp")})

	violations := g.Check([]scanner.Port{
		makePort(80, "tcp"),
		makePort(9999, "tcp"),
	})

	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Number != 9999 {
		t.Errorf("expected port 9999, got %d", violations[0].Port.Number)
	}
}

func TestAllow_AddsPortDynamically(t *testing.T) {
	g := portguard.New(nil)
	p := makePort(8080, "tcp")
	g.Allow(p)

	if !g.IsAllowed(p) {
		t.Error("expected port to be allowed after Allow()")
	}
}

func TestRevoke_RemovesPort(t *testing.T) {
	p := makePort(22, "tcp")
	g := portguard.New([]scanner.Port{p})
	g.Revoke(p)

	if g.IsAllowed(p) {
		t.Error("expected port to be revoked")
	}
}

func TestCheck_ProtocolDistinct(t *testing.T) {
	g := portguard.New([]scanner.Port{makePort(53, "tcp")})

	violations := g.Check([]scanner.Port{makePort(53, "udp")})

	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for udp/53, got %d", len(violations))
	}
}

func TestViolation_String(t *testing.T) {
	v := portguard.Violation{
		Port:   makePort(1234, "tcp"),
		Reason: "not in allowlist",
	}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty violation string")
	}
}
