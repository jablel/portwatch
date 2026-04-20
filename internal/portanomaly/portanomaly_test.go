package portanomaly

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func feed(d *Detector, ports []scanner.Port, times int) {
	for i := 0; i < times; i++ {
		d.Record(ports)
	}
}

func TestAnomalies_NoDataReturnsEmpty(t *testing.T) {
	d := New(0.3, 5)
	result := d.Anomalies([]scanner.Port{makePort(80, "tcp")})
	if len(result) != 0 {
		t.Fatalf("expected no anomalies, got %d", len(result))
	}
}

func TestAnomalies_StablePortNotFlagged(t *testing.T) {
	d := New(0.3, 5)
	p := makePort(443, "tcp")
	// Fill both windows with consistent presence.
	feed(d, []scanner.Port{p}, 15)
	result := d.Anomalies([]scanner.Port{p})
	if len(result) != 0 {
		t.Fatalf("stable port should not be flagged, got %v", result)
	}
}

func TestAnomalies_SuddenlyAbsentPortFlagged(t *testing.T) {
	d := New(0.3, 5)
	p := makePort(8080, "tcp")
	// Baseline: always present.
	feed(d, []scanner.Port{p}, 10)
	// Recent: always absent.
	feed(d, []scanner.Port{}, 10)

	result := d.Anomalies([]scanner.Port{p})
	if len(result) == 0 {
		t.Fatal("expected anomaly for suddenly absent port")
	}
	if result[0].Deviation >= 0 {
		t.Errorf("expected negative deviation, got %+.2f", result[0].Deviation)
	}
}

func TestAnomalies_SuddenlyPresentPortFlagged(t *testing.T) {
	d := New(0.3, 5)
	p := makePort(9999, "tcp")
	// Baseline: always absent.
	feed(d, []scanner.Port{}, 10)
	// Recent: always present.
	feed(d, []scanner.Port{p}, 10)

	result := d.Anomalies([]scanner.Port{p})
	if len(result) == 0 {
		t.Fatal("expected anomaly for suddenly present port")
	}
	if result[0].Deviation <= 0 {
		t.Errorf("expected positive deviation, got %+.2f", result[0].Deviation)
	}
}

func TestAnomalies_BelowThresholdNotFlagged(t *testing.T) {
	d := New(0.9, 5) // very high threshold
	p := makePort(22, "tcp")
	feed(d, []scanner.Port{p}, 10)
	feed(d, []scanner.Port{}, 10)

	// deviation is ~1.0 which equals threshold — use 0.9 threshold so just
	// below means we need deviation < 0.9 to not flag; here it will flag.
	// Re-test with threshold > max possible deviation.
	d2 := New(1.1, 5)
	feed(d2, []scanner.Port{p}, 10)
	feed(d2, []scanner.Port{}, 10)
	result := d2.Anomalies([]scanner.Port{p})
	if len(result) != 0 {
		t.Fatalf("expected no anomaly when deviation below threshold, got %v", result)
	}
}

func TestAnomalyString(t *testing.T) {
	a := Anomaly{
		Port:      makePort(80, "tcp"),
		Baseline:  0.9,
		Recent:    0.1,
		Deviation: -0.8,
	}
	s := a.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestNew_WindowClampedToOne(t *testing.T) {
	d := New(0.3, 0)
	if d.window != 1 {
		t.Errorf("expected window=1 after clamp, got %d", d.window)
	}
}
