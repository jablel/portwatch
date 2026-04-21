package portsteady

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func tick(base time.Time, n int) time.Time {
	return base.Add(time.Duration(n) * time.Second)
}

func TestStability_UnknownPortReturnsZero(t *testing.T) {
	tr := New(5)
	if got := tr.Stability(makePort("tcp", 80)); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestStability_PresentInAllScansReturnsOne(t *testing.T) {
	tr := New(4)
	p := makePort("tcp", 443)
	for i := 0; i < 4; i++ {
		tr.Observe([]scanner.Port{p}, tick(epoch, i))
	}
	if got := tr.Stability(p); got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestStability_PresentInHalfScansReturnsHalf(t *testing.T) {
	tr := New(4)
	p := makePort("tcp", 8080)
	// seed two observations so scans == 2
	tr.Observe([]scanner.Port{p}, tick(epoch, 0))
	tr.Observe([]scanner.Port{p}, tick(epoch, 1))
	// The score is scans/window = 2/4 = 0.5
	if got := tr.Stability(p); got != 0.5 {
		t.Fatalf("expected 0.5, got %v", got)
	}
}

func TestStability_ScansCapAtWindow(t *testing.T) {
	tr := New(3)
	p := makePort("udp", 53)
	for i := 0; i < 10; i++ {
		tr.Observe([]scanner.Port{p}, tick(epoch, i))
	}
	if got := tr.Stability(p); got != 1.0 {
		t.Fatalf("expected 1.0 (capped), got %v", got)
	}
}

func TestStability_PortRemovedWhenAbsent(t *testing.T) {
	tr := New(5)
	p := makePort("tcp", 22)
	tr.Observe([]scanner.Port{p}, tick(epoch, 0))
	// next scan without p — tracker should evict it
	tr.Observe([]scanner.Port{}, tick(epoch, 1))
	if got := tr.Stability(p); got != 0 {
		t.Fatalf("expected 0 after eviction, got %v", got)
	}
}

func TestUptime_UnknownPortReturnsZero(t *testing.T) {
	tr := New(5)
	if d := tr.Uptime(makePort("tcp", 80), epoch); d != 0 {
		t.Fatalf("expected zero duration, got %v", d)
	}
}

func TestUptime_GrowsOverTime(t *testing.T) {
	tr := New(5)
	p := makePort("tcp", 9090)
	tr.Observe([]scanner.Port{p}, epoch)
	later := epoch.Add(5 * time.Minute)
	if d := tr.Uptime(p, later); d != 5*time.Minute {
		t.Fatalf("expected 5m, got %v", d)
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	tr := New(0)
	if tr.window != defaultWindow {
		t.Fatalf("expected default window %d, got %d", defaultWindow, tr.window)
	}
}
