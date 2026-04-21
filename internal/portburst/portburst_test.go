package portburst_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portburst"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(n int, proto string) scanner.Port {
	return scanner.Port{Number: n, Protocol: proto}
}

func TestRecord_NoBurst_BelowThreshold(t *testing.T) {
	d := portburst.New(time.Second, 5)
	added := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	if b := d.Record(added); b != nil {
		t.Fatalf("expected nil burst, got count=%d", b.Count)
	}
}

func TestRecord_BurstWhenThresholdExceeded(t *testing.T) {
	d := portburst.New(time.Second, 2)
	added := []scanner.Port{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
		makePort(8080, "tcp"),
	}
	b := d.Record(added)
	if b == nil {
		t.Fatal("expected burst, got nil")
	}
	if b.Count != 3 {
		t.Fatalf("expected count=3, got %d", b.Count)
	}
}

func TestRecord_EmptyAdded_ReturnsNil(t *testing.T) {
	d := portburst.New(time.Second, 1)
	if b := d.Record(nil); b != nil {
		t.Fatal("expected nil for empty slice")
	}
}

func TestRecord_ZeroThreshold_NeverBursts(t *testing.T) {
	d := portburst.New(time.Second, 0)
	added := []scanner.Port{makePort(22, "tcp"), makePort(23, "tcp")}
	if b := d.Record(added); b != nil {
		t.Fatal("zero threshold should never trigger burst")
	}
}

func TestRecord_EventsExpireAfterWindow(t *testing.T) {
	d := portburst.New(50*time.Millisecond, 2)

	// First batch — below threshold on its own.
	d.Record([]scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")})

	// Wait for window to expire.
	time.Sleep(70 * time.Millisecond)

	// Second batch — should not combine with expired events.
	b := d.Record([]scanner.Port{makePort(8080, "tcp")})
	if b != nil {
		t.Fatalf("expected nil after window expiry, got count=%d", b.Count)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	d := portburst.New(time.Second, 1)
	d.Record([]scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")})
	d.Reset()
	b := d.Record([]scanner.Port{makePort(8080, "tcp")})
	if b != nil {
		t.Fatal("expected nil after reset")
	}
}
