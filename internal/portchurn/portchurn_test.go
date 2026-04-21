package portchurn

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestChurn_FewerThanTwoObservations_ReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
	tr.Record([]scanner.Port{makePort(80, "tcp")}, t0)
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0 after single observation, got %v", got)
	}
}

func TestChurn_StablePorts_ReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	tr.Record(ports, t0)
	tr.Record(ports, t0.Add(10*time.Second))
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0 for stable ports, got %v", got)
	}
}

func TestChurn_AllPortsReplaced_ReturnsOne(t *testing.T) {
	tr := New(time.Minute)
	tr.Record([]scanner.Port{makePort(80, "tcp")}, t0)
	tr.Record([]scanner.Port{makePort(9000, "tcp")}, t0.Add(5*time.Second))
	if got := tr.Churn(); got != 1.0 {
		t.Fatalf("expected 1.0 for fully replaced ports, got %v", got)
	}
}

func TestChurn_HalfReplaced_ReturnsHalf(t *testing.T) {
	tr := New(time.Minute)
	tr.Record([]scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}, t0)
	tr.Record([]scanner.Port{makePort(80, "tcp"), makePort(9000, "tcp")}, t0.Add(5*time.Second))
	got := tr.Churn()
	// 3 unique ports total; 443 and 9000 each absent from one scan → 2 churned
	want := 2.0 / 3.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected ~%.4f, got %.4f", want, got)
	}
}

func TestChurn_WindowEvictsOldObservations(t *testing.T) {
	tr := New(10 * time.Second)
	// Old scan with a port that will later disappear.
	tr.Record([]scanner.Port{makePort(22, "tcp")}, t0)
	// Two recent scans with a stable port – old entry should be evicted.
	new1 := t0.Add(20 * time.Second)
	new2 := t0.Add(25 * time.Second)
	tr.Record([]scanner.Port{makePort(80, "tcp")}, new1)
	tr.Record([]scanner.Port{makePort(80, "tcp")}, new2)
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0 after old entry evicted, got %v", got)
	}
}

func TestChurn_EmptyScans_ReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(nil, t0)
	tr.Record(nil, t0.Add(5*time.Second))
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0 for empty scans, got %v", got)
	}
}

func TestReset_ClearsObservations(t *testing.T) {
	tr := New(time.Minute)
	tr.Record([]scanner.Port{makePort(80, "tcp")}, t0)
	tr.Record([]scanner.Port{makePort(9000, "tcp")}, t0.Add(5*time.Second))
	tr.Reset()
	if got := tr.Churn(); got != 0 {
		t.Fatalf("expected 0 after reset, got %v", got)
	}
}
