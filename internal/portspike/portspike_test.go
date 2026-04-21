package portspike

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Number: uint16(1024 + i), Protocol: "tcp"}
	}
	return ports
}

func TestRecord_FirstCallReturnsNil(t *testing.T) {
	d := New(0.5)
	if got := d.Record(makePorts(10)); got != nil {
		t.Fatalf("expected nil on first call, got %v", got)
	}
}

func TestRecord_NoSpikeBelowThreshold(t *testing.T) {
	d := New(0.5)
	d.Record(makePorts(10))
	if got := d.Record(makePorts(14)); got != nil { // 40 % < 50 %
		t.Fatalf("expected nil below threshold, got %v", got)
	}
}

func TestRecord_SpikeAtThreshold(t *testing.T) {
	d := New(0.5)
	d.Record(makePorts(10))
	spike := d.Record(makePorts(15)) // exactly 50 %
	if spike == nil {
		t.Fatal("expected spike at threshold, got nil")
	}
	if spike.Previous != 10 || spike.Current != 15 || spike.Delta != 5 {
		t.Fatalf("unexpected spike values: %+v", spike)
	}
}

func TestRecord_SpikeAboveThreshold(t *testing.T) {
	d := New(0.25)
	d.Record(makePorts(4))
	spike := d.Record(makePorts(8)) // 100 % > 25 %
	if spike == nil {
		t.Fatal("expected spike, got nil")
	}
	if spike.Ratio < 0.99 {
		t.Fatalf("expected ratio ~1.0, got %.4f", spike.Ratio)
	}
}

func TestRecord_DecreaseNeverSpikes(t *testing.T) {
	d := New(0.1)
	d.Record(makePorts(20))
	if got := d.Record(makePorts(5)); got != nil {
		t.Fatalf("decrease should never spike, got %v", got)
	}
}

func TestRecord_ZeroThresholdNeverSpikes(t *testing.T) {
	d := New(0)
	d.Record(makePorts(1))
	if got := d.Record(makePorts(1000)); got != nil {
		t.Fatalf("zero threshold should never spike, got %v", got)
	}
}

func TestRecord_PreviousZeroNeverSpikes(t *testing.T) {
	d := New(0.5)
	d.Record(makePorts(0))
	if got := d.Record(makePorts(10)); got != nil {
		t.Fatalf("previous=0 should not spike, got %v", got)
	}
}

func TestReset_ClearsBaseline(t *testing.T) {
	d := New(0.5)
	d.Record(makePorts(10))
	d.Reset()
	// After reset the next call is treated as first — no spike.
	if got := d.Record(makePorts(100)); got != nil {
		t.Fatalf("expected nil after reset, got %v", got)
	}
	// Now a spike can be detected again.
	spike := d.Record(makePorts(200)) // 100 % > 50 %
	if spike == nil {
		t.Fatal("expected spike after reset baseline, got nil")
	}
}

func TestSpike_String(t *testing.T) {
	s := Spike{Previous: 10, Current: 15, Delta: 5, Ratio: 0.5}
	got := s.String()
	if got == "" {
		t.Fatal("expected non-empty string")
	}
}
