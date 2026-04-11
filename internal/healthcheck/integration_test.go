package healthcheck_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

// TestHealthcheck_ConcurrentAccess verifies that RecordScan, RecordError, and
// Status can be called concurrently without data races.
func TestHealthcheck_ConcurrentAccess(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); c.RecordScan() }()
		go func() { defer wg.Done(); c.RecordError() }()
		go func() { defer wg.Done(); _ = c.Status() }()
	}
	wg.Wait()

	s := c.Status()
	if s.ScanCount == 0 {
		t.Fatal("expected at least one scan recorded")
	}
	if s.ErrorCount == 0 {
		t.Fatal("expected at least one error recorded")
	}
}

// TestHealthcheck_TransitionsHealthState verifies healthy → unhealthy transition.
func TestHealthcheck_TransitionsHealthState(t *testing.T) {
	staleness := 20 * time.Millisecond
	c := healthcheck.New(staleness)

	c.RecordScan()
	if !c.Status().Healthy {
		t.Fatal("should be healthy immediately after scan")
	}

	time.Sleep(staleness + 10*time.Millisecond)
	if c.Status().Healthy {
		t.Fatal("should be unhealthy after staleness window expires")
	}

	// recover
	c.RecordScan()
	if !c.Status().Healthy {
		t.Fatal("should be healthy again after new scan")
	}
}
