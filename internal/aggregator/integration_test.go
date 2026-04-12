package aggregator_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/scanner"
)

func TestAggregator_ConcurrentUpdates(t *testing.T) {
	a := aggregator.New()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			src := string(rune('a' + n))
			a.Update(src, []scanner.Port{
				{Number: uint16(8000 + n), Protocol: "tcp"},
			})
		}(i)
	}

	wg.Wait()

	got := a.Merge()
	if len(got) != 10 {
		t.Fatalf("expected 10 ports, got %d", len(got))
	}
}

func TestAggregator_UpdateThenRemoveLeavesCorrectPorts(t *testing.T) {
	a := aggregator.New()
	a.Update("tcp", makePorts(22, 80, 443))
	a.Update("udp", makePorts(53, 123))

	a.Remove("udp")
	got := sortPorts(a.Merge())

	expected := []uint16{22, 80, 443}
	if len(got) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(got))
	}
	for i, p := range got {
		if p.Number != expected[i] {
			t.Errorf("port[%d]: want %d, got %d", i, expected[i], p.Number)
		}
	}
}
