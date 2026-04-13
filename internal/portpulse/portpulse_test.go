package portpulse

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestPulse_UnknownPortReturnsZero(t *testing.T) {
	tr := New(5)
	if got := tr.Pulse(makePort(80, "tcp")); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestPulse_AfterSingleRecord(t *testing.T) {
	tr := New(4)
	tr.Record([]scanner.Port{makePort(443, "tcp")})
	got := tr.Pulse(makePort(443, "tcp"))
	// 1 observation out of window=4 → 0.25
	if got != 0.25 {
		t.Fatalf("expected 0.25, got %f", got)
	}
}

func TestPulse_SaturatesAtOne(t *testing.T) {
	tr := New(3)
	p := makePort(22, "tcp")
	for i := 0; i < 10; i++ {
		tr.Record([]scanner.Port{p})
	}
	if got := tr.Pulse(p); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestPulse_DecaysWhenPortDisappears(t *testing.T) {
	tr := New(4)
	p := makePort(8080, "tcp")
	// Fill to max
	for i := 0; i < 4; i++ {
		tr.Record([]scanner.Port{p})
	}
	if got := tr.Pulse(p); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
	// One tick without the port
	tr.Record(nil)
	got := tr.Pulse(p)
	if got != 0.75 {
		t.Fatalf("expected 0.75 after one miss, got %f", got)
	}
}

func TestPulse_PortRemovedWhenCountReachesZero(t *testing.T) {
	tr := New(2)
	p := makePort(3306, "tcp")
	tr.Record([]scanner.Port{p})
	tr.Record(nil) // count → 0, should be evicted
	tr.Record(nil)
	all := tr.All()
	if _, ok := all[portKey(p)]; ok {
		t.Fatal("expected port to be evicted from map")
	}
}

func TestAll_ReturnsAllTrackedPorts(t *testing.T) {
	tr := New(10)
	ports := []scanner.Port{
		makePort(80, "tcp"),
		makePort(53, "udp"),
	}
	tr.Record(ports)
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestNew_WindowClampedToOne(t *testing.T) {
	tr := New(0)
	if tr.window != 1 {
		t.Fatalf("expected window=1, got %d", tr.window)
	}
}
