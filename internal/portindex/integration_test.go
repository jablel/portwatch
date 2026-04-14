package portindex_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portindex"
	"github.com/user/portwatch/internal/scanner"
)

func TestIndex_ConcurrentBuildAndQuery(t *testing.T) {
	idx := portindex.New()
	ports := makePorts()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idx.Build(ports)
			_ = idx.ByNumber(80)
			_ = idx.ByProtocol("tcp")
			_ = idx.ByTag("web")
		}()
	}
	wg.Wait()
}

func TestIndex_BuildThenQueryConsistency(t *testing.T) {
	idx := portindex.New()
	batch := []scanner.Port{
		{Number: 8080, Protocol: "tcp", Tags: []string{"proxy"}},
		{Number: 8443, Protocol: "tcp", Tags: []string{"proxy", "tls"}},
		{Number: 5353, Protocol: "udp", Tags: []string{"mdns"}},
	}
	idx.Build(batch)

	proxy := idx.ByTag("proxy")
	if len(proxy) != 2 {
		t.Fatalf("expected 2 proxy ports, got %d", len(proxy))
	}

	udp := idx.ByProtocol("udp")
	if len(udp) != 1 {
		t.Fatalf("expected 1 udp port, got %d", len(udp))
	}

	// Mutating returned slice must not affect index.
	proxy[0].Number = 0
	if got := idx.ByTag("proxy"); got[0].Number == 0 {
		t.Fatal("index returned a mutable reference")
	}
}
