package portlease

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Proto: proto}
}

func TestGrant_NoViolationBeforeExpiry(t *testing.T) {
	tr := New()
	p := makePort(8080, "tcp")
	tr.Grant(p, 10*time.Minute)

	violations := tr.Check([]scanner.Port{p})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_ViolationAfterExpiry(t *testing.T) {
	tr := New()
	now := time.Now()
	tr.now = func() time.Time { return now }

	p := makePort(9090, "tcp")
	tr.Grant(p, 5*time.Minute)

	// advance time past the TTL
	tr.now = func() time.Time { return now.Add(6 * time.Minute) }

	violations := tr.Check([]scanner.Port{p})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Port.Number != 9090 {
		t.Errorf("unexpected port in violation: %v", violations[0].Port)
	}
}

func TestCheck_NoViolationWhenPortClosed(t *testing.T) {
	tr := New()
	now := time.Now()
	tr.now = func() time.Time { return now }

	p := makePort(443, "tcp")
	tr.Grant(p, 1*time.Minute)
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }

	// port is not in current scan — no violation expected
	violations := tr.Check([]scanner.Port{})
	if len(violations) != 0 {
		t.Fatalf("expected no violations when port is closed, got %d", len(violations))
	}
}

func TestRevoke_RemovesLease(t *testing.T) {
	tr := New()
	now := time.Now()
	tr.now = func() time.Time { return now }

	p := makePort(22, "tcp")
	tr.Grant(p, 1*time.Minute)
	tr.Revoke(p)

	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	violations := tr.Check([]scanner.Port{p})
	if len(violations) != 0 {
		t.Fatalf("expected no violations after revoke, got %d", len(violations))
	}
}

func TestCheck_ZeroTTL_NeverExpires(t *testing.T) {
	tr := New()
	now := time.Now()
	tr.now = func() time.Time { return now }

	p := makePort(80, "tcp")
	tr.Grant(p, 0) // zero TTL

	tr.now = func() time.Time { return now.Add(365 * 24 * time.Hour) }
	violations := tr.Check([]scanner.Port{p})
	if len(violations) != 0 {
		t.Fatalf("expected no violations for zero-TTL lease, got %d", len(violations))
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Port: makePort(8080, "tcp"), Reason: "lease expired after 5m0s"}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty violation string")
	}
}
