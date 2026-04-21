package portсhadow_test

import (
	"sync"
	"testing"
	"time"

	portсhadow "portwatch/internal/portсhadow"
	"portwatch/internal/scanner"
)

func TestTracker_ConcurrentObserve(t *testing.T) {
	tr := portсhadow.New(5 * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			tr.Observe([]scanner.Port{{Number: n, Protocol: "tcp"}})
		}(i)
	}
	wg.Wait()

	// Drain shadows and active without panic
	_ = tr.Shadows()
	_ = tr.Active()
}

func TestTracker_RapidAppearDisappear(t *testing.T) {
	tr := portсhadow.New(30 * time.Second)

	ports := []scanner.Port{
		{Number: 5000, Protocol: "tcp"},
		{Number: 5001, Protocol: "tcp"},
	}

	tr.Observe(ports)
	tr.Observe(nil) // all gone immediately

	shadows := tr.Shadows()
	if len(shadows) != 2 {
		t.Fatalf("expected 2 shadow ports, got %d", len(shadows))
	}

	for _, s := range shadows {
		if s.Count != 1 {
			t.Errorf("port %d: expected count 1, got %d", s.Port.Number, s.Count)
		}
	}
}
