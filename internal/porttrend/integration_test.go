package porttrend_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/porttrend"
	"github.com/user/portwatch/internal/scanner"
)

func TestTracker_ConcurrentRecord(t *testing.T) {
	tr := porttrend.New(time.Minute)
	p := scanner.Port{Protocol: "tcp", Number: 8080}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.Record([]scanner.Port{p})
		}()
	}
	wg.Wait()

	// Must not panic and must return a valid trend.
	trend := tr.Trend(p)
	if trend != porttrend.Stable && trend != porttrend.Rising && trend != porttrend.Falling {
		t.Fatalf("unexpected trend value: %d", trend)
	}
}

func TestTracker_WindowEvictsOldObservations(t *testing.T) {
	tr := porttrend.New(80 * time.Millisecond)
	p := scanner.Port{Protocol: "tcp", Number: 3306}

	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{p})

	// Wait for the window to expire.
	time.Sleep(120 * time.Millisecond)

	// After eviction the port should look stable (no data).
	if got := tr.Trend(p); got != porttrend.Stable {
		t.Fatalf("expected Stable after window expiry, got %s", got)
	}
}
