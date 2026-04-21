package portvelocity_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portvelocity"
	"github.com/user/portwatch/internal/scanner"
)

func TestVelocity_ConcurrentRecord(t *testing.T) {
	tr := portvelocity.New()

	base := []scanner.Port{
		{Protocol: "tcp", Number: 80},
		{Protocol: "tcp", Number: 443},
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.Record(base)
		}()
	}
	wg.Wait()

	// No race or panic is the primary assertion; velocity must be in [0,1].
	v := tr.Last()
	if v < 0 || v > 1 {
		t.Fatalf("velocity out of range: %f", v)
	}
}

func TestVelocity_GradualChurnStaysBelow1(t *testing.T) {
	tr := portvelocity.New()

	// Seed with 10 ports.
	initial := make([]scanner.Port, 10)
	for i := range initial {
		initial[i] = scanner.Port{Protocol: "tcp", Number: 1000 + i}
	}
	tr.Record(initial)

	// Replace one port per scan — velocity should be low.
	current := make([]scanner.Port, len(initial))
	copy(current, initial)

	for step := 0; step < 5; step++ {
		current[step] = scanner.Port{Protocol: "tcp", Number: 2000 + step}
		v := tr.Record(current)
		if v > 0.5 {
			t.Fatalf("step %d: velocity %f unexpectedly high", step, v)
		}
	}
}
