package portvelocity

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

func TestRecord_FirstCallReturnsZero(t *testing.T) {
	tr := New()
	v := tr.Record([]scanner.Port{makePort("tcp", 80)})
	if v != 0 {
		t.Fatalf("expected 0 on first call, got %f", v)
	}
}

func TestRecord_NoChangeReturnsZero(t *testing.T) {
	tr := New()
	ports := []scanner.Port{makePort("tcp", 80), makePort("tcp", 443)}
	tr.Record(ports)
	v := tr.Record(ports)
	if v != 0 {
		t.Fatalf("expected 0 for identical scans, got %f", v)
	}
}

func TestRecord_AllPortsReplaced_ReturnsOne(t *testing.T) {
	tr := New()
	tr.Record([]scanner.Port{makePort("tcp", 80)})
	v := tr.Record([]scanner.Port{makePort("tcp", 9000)})
	if v != 1.0 {
		t.Fatalf("expected 1.0 for complete replacement, got %f", v)
	}
}

func TestRecord_HalfChanged_ReturnsHalf(t *testing.T) {
	tr := New()
	tr.Record([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443)})
	// Remove 443, add 8080 — 1 removed + 1 added out of denom 2
	v := tr.Record([]scanner.Port{makePort("tcp", 80), makePort("tcp", 8080)})
	if v != 1.0 {
		// denom = max(2, 2) = 2, changes = 2, velocity = 1.0
		t.Fatalf("expected 1.0, got %f", v)
	}
}

func TestRecord_OneAddedOutOfThree(t *testing.T) {
	tr := New()
	tr.Record([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443)})
	v := tr.Record([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443), makePort("tcp", 8080)})
	// added=1, removed=0, denom=max(2,3)=3 → 1/3 ≈ 0.333
	expected := 1.0 / 3.0
	if v < expected-0.001 || v > expected+0.001 {
		t.Fatalf("expected ~%f, got %f", expected, v)
	}
}

func TestLast_ReturnsLatestVelocity(t *testing.T) {
	tr := New()
	tr.Record([]scanner.Port{makePort("tcp", 80)})
	tr.Record([]scanner.Port{makePort("tcp", 9000)})
	if tr.Last() != 1.0 {
		t.Fatalf("expected Last() == 1.0, got %f", tr.Last())
	}
}

func TestReset_ClearsState(t *testing.T) {
	tr := New()
	tr.Record([]scanner.Port{makePort("tcp", 80)})
	tr.Record([]scanner.Port{makePort("tcp", 9000)})
	tr.Reset()
	if tr.Last() != 0 {
		t.Fatalf("expected Last() == 0 after Reset, got %f", tr.Last())
	}
	// First call after reset should again return 0
	v := tr.Record([]scanner.Port{makePort("tcp", 80)})
	if v != 0 {
		t.Fatalf("expected 0 after reset, got %f", v)
	}
}

func TestRecord_EmptyScans_ReturnsZero(t *testing.T) {
	tr := New()
	tr.Record(nil)
	v := tr.Record(nil)
	if v != 0 {
		t.Fatalf("expected 0 for two empty scans, got %f", v)
	}
}
