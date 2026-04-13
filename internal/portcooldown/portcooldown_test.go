package portcooldown

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

func TestAbsentFor_UnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.AbsentFor(makePort("tcp", 8080))
	if ok {
		t.Fatal("expected false for unknown port")
	}
}

func TestAbsentFor_KnownPortReturnsElapsed(t *testing.T) {
	now := time.Now()
	tr := New()
	tr.clock = func() time.Time { return now }

	p := makePort("tcp", 443)
	tr.Observe([]scanner.Port{p})

	// Advance clock by 5 seconds without observing the port again.
	tr.clock = func() time.Time { return now.Add(5 * time.Second) }

	dur, ok := tr.AbsentFor(p)
	if !ok {
		t.Fatal("expected true for known port")
	}
	if dur != 5*time.Second {
		t.Fatalf("expected 5s, got %v", dur)
	}
}

func TestInCooldown_WithinWindow(t *testing.T) {
	now := time.Now()
	tr := New()
	tr.clock = func() time.Time { return now }

	p := makePort("tcp", 22)
	tr.Observe([]scanner.Port{p})

	tr.clock = func() time.Time { return now.Add(2 * time.Second) }

	if !tr.InCooldown(p, 10*time.Second) {
		t.Fatal("expected port to be in cooldown")
	}
}

func TestInCooldown_OutsideWindow(t *testing.T) {
	now := time.Now()
	tr := New()
	tr.clock = func() time.Time { return now }

	p := makePort("tcp", 22)
	tr.Observe([]scanner.Port{p})

	tr.clock = func() time.Time { return now.Add(15 * time.Second) }

	if tr.InCooldown(p, 10*time.Second) {
		t.Fatal("expected port to be outside cooldown")
	}
}

func TestInCooldown_ZeroDurationNeverInCooldown(t *testing.T) {
	tr := New()
	p := makePort("udp", 53)
	tr.Observe([]scanner.Port{p})

	if tr.InCooldown(p, 0) {
		t.Fatal("zero duration should never be in cooldown")
	}
}

func TestReset_ClearsState(t *testing.T) {
	tr := New()
	p := makePort("tcp", 80)
	tr.Observe([]scanner.Port{p})
	tr.Reset()

	_, ok := tr.AbsentFor(p)
	if ok {
		t.Fatal("expected state to be cleared after Reset")
	}
}

func TestObserve_UpdatesLastSeen(t *testing.T) {
	now := time.Now()
	tr := New()
	tr.clock = func() time.Time { return now }

	p := makePort("tcp", 8443)
	tr.Observe([]scanner.Port{p})

	tr.clock = func() time.Time { return now.Add(30 * time.Second) }
	tr.Observe([]scanner.Port{p}) // port is present again

	tr.clock = func() time.Time { return now.Add(31 * time.Second) }

	dur, ok := tr.AbsentFor(p)
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if dur != 1*time.Second {
		t.Fatalf("expected 1s absence, got %v", dur)
	}
}
