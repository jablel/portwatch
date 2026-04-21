package portburst_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portburst"
	"github.com/user/portwatch/internal/scanner"
)

func TestDetector_ConcurrentRecord(t *testing.T) {
	d := portburst.New(time.Second, 100)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			d.Record([]scanner.Port{makePort(n+1024, "tcp")})
		}(i)
	}
	wg.Wait()
}

func TestDetector_BurstAcrossMultipleRecordCalls(t *testing.T) {
	d := portburst.New(time.Second, 3)

	// Three separate calls, each adding one port — total 4 within window.
	d.Record([]scanner.Port{makePort(80, "tcp")})
	d.Record([]scanner.Port{makePort(443, "tcp")})
	d.Record([]scanner.Port{makePort(22, "tcp")})
	b := d.Record([]scanner.Port{makePort(8080, "tcp")})

	if b == nil {
		t.Fatal("expected burst across multiple calls")
	}
	if b.Count != 4 {
		t.Fatalf("expected count=4, got %d", b.Count)
	}
}

func TestDetector_ResetUnblocksSubsequentCalls(t *testing.T) {
	d := portburst.New(time.Second, 2)

	// Trigger a burst.
	d.Record([]scanner.Port{makePort(80, "tcp"), makePort(443, "tcp"), makePort(22, "tcp")})

	// Reset and verify a small batch no longer triggers.
	d.Reset()
	b := d.Record([]scanner.Port{makePort(9090, "tcp")})
	if b != nil {
		t.Fatal("expected nil after reset")
	}
}
