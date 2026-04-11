package tagger_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

// TestTagger_ConcurrentDefineAndTag ensures no data races when Define and Tag
// are called from multiple goroutines simultaneously.
func TestTagger_ConcurrentDefineAndTag(t *testing.T) {
	tg := tagger.New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n uint16) {
			defer wg.Done()
			tg.Define(n, "svc")
		}(uint16(10000 + i))
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n uint16) {
			defer wg.Done()
			_ = tg.Tag(scanner.Port{Number: n, Protocol: "tcp"})
		}(uint16(10000 + i))
	}

	wg.Wait()
}

// TestTagger_CustomTakesPrecedenceOverWellKnown verifies that a custom mapping
// set after construction overrides the built-in table.
func TestTagger_CustomTakesPrecedenceOverWellKnown(t *testing.T) {
	tg := tagger.New()
	tg.Define(443, "internal-proxy")

	got := tg.Tag(scanner.Port{Number: 443, Protocol: "tcp"})
	if got != "internal-proxy" {
		t.Fatalf("expected internal-proxy, got %q", got)
	}
}

// TestTagger_TagAllPreservesAllPorts checks that TagAll returns one entry per
// input port, including duplicates with different protocols.
func TestTagger_TagAllPreservesAllPorts(t *testing.T) {
	tg := tagger.New()
	ports := []scanner.Port{
		{Number: 80, Protocol: "tcp"},
		{Number: 80, Protocol: "udp"},
		{Number: 443, Protocol: "tcp"},
	}
	labels := tg.TagAll(ports)
	if len(labels) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(labels))
	}
}
