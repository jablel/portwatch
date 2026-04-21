package portсhadow

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestObserve_ShadowDetectedWhenPortDisappearsWithinWindow(t *testing.T) {
	tr := New(5 * time.Second)

	tr.Observe([]scanner.Port{makePort(8080, "tcp")})
	tr.Observe([]scanner.Port{}) // port gone

	shadows := tr.Shadows()
	if len(shadows) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(shadows))
	}
	if shadows[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", shadows[0].Port.Number)
	}
}

func TestObserve_NoPersistentPortNotShadow(t *testing.T) {
	now := time.Now()
	tr := New(1 * time.Second)
	tr.now = func() time.Time { return now }

	tr.Observe([]scanner.Port{makePort(443, "tcp")})
	// advance time beyond window
	tr.now = func() time.Time { return now.Add(2 * time.Second) }
	tr.Observe([]scanner.Port{makePort(443, "tcp")})
	tr.now = func() time.Time { return now.Add(3 * time.Second) }
	tr.Observe([]scanner.Port{}) // disappears after window

	shadows := tr.Shadows()
	if len(shadows) != 0 {
		t.Errorf("expected no shadows, got %d", len(shadows))
	}
}

func TestObserve_ActiveContainsCurrentPorts(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Observe([]scanner.Port{makePort(22, "tcp"), makePort(80, "tcp")})

	active := tr.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active ports, got %d", len(active))
	}
}

func TestShadows_ResetsAfterRead(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Observe([]scanner.Port{makePort(9999, "tcp")})
	tr.Observe([]scanner.Port{})

	first := tr.Shadows()
	second := tr.Shadows()

	if len(first) != 1 {
		t.Fatalf("expected 1 shadow on first read, got %d", len(first))
	}
	if len(second) != 0 {
		t.Errorf("expected 0 shadows on second read, got %d", len(second))
	}
}

func TestObserve_ZeroWindow_NeverShadow(t *testing.T) {
	tr := New(0)
	tr.Observe([]scanner.Port{makePort(1234, "udp")})
	tr.Observe([]scanner.Port{})

	if got := tr.Shadows(); len(got) != 0 {
		t.Errorf("zero window should never produce shadows, got %d", len(got))
	}
}

func TestObserve_CountTracksObservations(t *testing.T) {
	tr := New(10 * time.Second)
	p := makePort(8080, "tcp")

	tr.Observe([]scanner.Port{p})
	tr.Observe([]scanner.Port{p})
	tr.Observe([]scanner.Port{p})

	active := tr.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active port, got %d", len(active))
	}
	if active[0].Count != 3 {
		t.Errorf("expected count 3, got %d", active[0].Count)
	}
}
