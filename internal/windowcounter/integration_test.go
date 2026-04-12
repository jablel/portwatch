package windowcounter_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/windowcounter"
)

func TestWindowCounter_ConcurrentAdds(t *testing.T) {
	c := windowcounter.New(time.Second)
	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Add("tcp:80")
		}()
	}
	wg.Wait()
	if got := c.Count("tcp:80"); got != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, got)
	}
}

func TestWindowCounter_SlidingWindowDropsOldEvents(t *testing.T) {
	c := windowcounter.New(60 * time.Millisecond)

	// Add two events that will expire.
	c.Add("tcp:443")
	c.Add("tcp:443")
	time.Sleep(80 * time.Millisecond)

	// Add one fresh event; only it should be in the window.
	count := c.Add("tcp:443")
	if count != 1 {
		t.Fatalf("expected 1 after old events expired, got %d", count)
	}
}

func TestWindowCounter_ResetDoesNotAffectOtherKeys(t *testing.T) {
	c := windowcounter.New(time.Second)
	c.Add("tcp:80")
	c.Add("tcp:443")
	c.Reset("tcp:80")

	if got := c.Count("tcp:80"); got != 0 {
		t.Fatalf("tcp:80 should be 0 after reset, got %d", got)
	}
	if got := c.Count("tcp:443"); got != 1 {
		t.Fatalf("tcp:443 should still be 1, got %d", got)
	}
}
