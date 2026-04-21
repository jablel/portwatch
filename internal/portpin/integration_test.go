package portpin_test

import (
	"sync"
	"testing"

	"portwatch/internal/portpin"
	"portwatch/internal/scanner"
)

func TestPinner_ConcurrentPinAndCheck(t *testing.T) {
	p := portpin.New()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			p.Pin(scanner.Port{Number: n, Protocol: "tcp"})
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = p.Check([]scanner.Port{{Number: 0, Protocol: "tcp"}})
		}()
	}

	wg.Wait()

	if len(p.Pinned()) == 0 {
		t.Error("expected pinned ports to be non-empty after concurrent pins")
	}
}

func TestPinner_PinUnpinCycle(t *testing.T) {
	p := portpin.New()
	port := scanner.Port{Number: 8080, Protocol: "tcp"}

	p.Pin(port)
	if v := p.Check([]scanner.Port{}); len(v) != 1 {
		t.Fatalf("expected 1 violation before unpin, got %d", len(v))
	}

	p.Unpin(port)
	if v := p.Check([]scanner.Port{}); len(v) != 0 {
		t.Fatalf("expected 0 violations after unpin, got %d", len(v))
	}
}
