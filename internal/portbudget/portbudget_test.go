package portbudget_test

import (
	"testing"

	"github.com/user/portwatch/internal/portbudget"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Number: uint16(8000 + i), Protocol: "tcp"}
	}
	return ports
}

func TestCheck_UnderLimit_ReturnsNil(t *testing.T) {
	b := portbudget.New(10)
	if v := b.Check(makePorts(5)); v != nil {
		t.Fatalf("expected nil violation, got %v", v)
	}
}

func TestCheck_AtLimit_ReturnsNil(t *testing.T) {
	b := portbudget.New(5)
	if v := b.Check(makePorts(5)); v != nil {
		t.Fatalf("expected nil violation at limit, got %v", v)
	}
}

func TestCheck_ExceedsLimit_ReturnsViolation(t *testing.T) {
	b := portbudget.New(3)
	ports := makePorts(5)
	v := b.Check(ports)
	if v == nil {
		t.Fatal("expected violation, got nil")
	}
	if v.Limit != 3 {
		t.Errorf("limit: want 3, got %d", v.Limit)
	}
	if v.Observed != 5 {
		t.Errorf("observed: want 5, got %d", v.Observed)
	}
	if len(v.Excess) != 2 {
		t.Errorf("excess len: want 2, got %d", len(v.Excess))
	}
}

func TestCheck_ZeroLimit_NeverViolates(t *testing.T) {
	b := portbudget.New(0)
	if v := b.Check(makePorts(100)); v != nil {
		t.Fatalf("expected nil for zero limit, got %v", v)
	}
}

func TestCheck_NegativeLimit_NeverViolates(t *testing.T) {
	b := portbudget.New(-1)
	if v := b.Check(makePorts(50)); v != nil {
		t.Fatalf("expected nil for negative limit, got %v", v)
	}
}

func TestSetLimit_UpdatesEnforcement(t *testing.T) {
	b := portbudget.New(10)
	ports := makePorts(5)
	if v := b.Check(ports); v != nil {
		t.Fatal("expected no violation before tightening limit")
	}
	b.SetLimit(3)
	if v := b.Check(ports); v == nil {
		t.Fatal("expected violation after tightening limit")
	}
}

func TestViolation_ErrorMessage(t *testing.T) {
	b := portbudget.New(2)
	v := b.Check(makePorts(4))
	if v == nil {
		t.Fatal("expected violation")
	}
	msg := v.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
