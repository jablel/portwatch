package portexpiry

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makePort(proto, addr string) scanner.Port {
	return scanner.Port{Proto: proto, Addr: addr}
}

func TestObserve_TracksNewPort(t *testing.T) {
	tr := New(time.Minute)
	p := makePort("tcp", "0.0.0.0:8080")
	tr.Observe([]scanner.Port{p})
	if tr.Age(p) == 0 {
		t.Fatal("expected non-zero age after observe")
	}
}

func TestObserve_RemovesAbsentPort(t *testing.T) {
	tr := New(time.Minute)
	p := makePort("tcp", "0.0.0.0:8080")
	tr.Observe([]scanner.Port{p})
	tr.Observe([]scanner.Port{})
	if tr.Age(p) != 0 {
		t.Fatal("expected zero age after port removed")
	}
}

func TestObserve_PreservesFirstSeen(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.now = func() time.Time { return now }
	p := makePort("tcp", "0.0.0.0:9090")
	tr.Observe([]scanner.Port{p})

	later := now.Add(5 * time.Second)
	tr.now = func() time.Time { return later }
	tr.Observe([]scanner.Port{p}) // second observe must not reset first-seen

	if tr.Age(p) != 5*time.Second {
		t.Fatalf("expected 5s age, got %v", tr.Age(p))
	}
}

func TestExpired_ReturnsPortsOverMaxAge(t *testing.T) {
	now := time.Now()
	tr := New(10 * time.Second)
	tr.now = func() time.Time { return now }

	old := makePort("tcp", "0.0.0.0:22")
	fresh := makePort("tcp", "0.0.0.0:80")
	tr.Observe([]scanner.Port{old, fresh})

	// Advance time so old port is expired, fresh is not.
	tr.now = func() time.Time { return now.Add(15 * time.Second) }
	// Re-observe to keep both ports tracked.
	tr.Observe([]scanner.Port{old, fresh})

	// Manually adjust fresh entry to be recent.
	tr.mu.Lock()
	tr.entries[portKey(fresh)] = Entry{FirstSeen: now.Add(14 * time.Second), Port: fresh}
	tr.mu.Unlock()

	expired := tr.Expired()
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired port, got %d", len(expired))
	}
	if expired[0].Port != old {
		t.Fatalf("expected expired port %v, got %v", old, expired[0].Port)
	}
}

func TestExpired_ZeroMaxAgeReturnsNil(t *testing.T) {
	tr := New(0)
	p := makePort("tcp", "0.0.0.0:443")
	tr.Observe([]scanner.Port{p})
	if tr.Expired() != nil {
		t.Fatal("expected nil for zero maxAge")
	}
}

func TestExpired_NegativeMaxAgeReturnsNil(t *testing.T) {
	tr := New(-time.Second)
	p := makePort("tcp", "0.0.0.0:443")
	tr.Observe([]scanner.Port{p})
	if tr.Expired() != nil {
		t.Fatal("expected nil for negative maxAge")
	}
}

func TestAge_UnknownPortReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	p := makePort("udp", "0.0.0.0:53")
	if tr.Age(p) != 0 {
		t.Fatal("expected zero age for unknown port")
	}
}
