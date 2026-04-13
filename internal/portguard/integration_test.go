package portguard_test

import (
	"sync"
	"testing"

	"portwatch/internal/portguard"
	"portwatch/internal/scanner"
)

func TestGuard_ConcurrentAllowAndCheck(t *testing.T) {
	g := portguard.New(nil)

	var wg sync.WaitGroup
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			g.Allow(scanner.Port{Number: n, Protocol: "tcp"})
		}(i)
	}
	wg.Wait()

	ports := make([]scanner.Port, 20)
	for i := range ports {
		ports[i] = scanner.Port{Number: i + 1, Protocol: "tcp"}
	}

	violations := g.Check(ports)
	if len(violations) != 0 {
		t.Fatalf("expected no violations after concurrent Allow, got %d", len(violations))
	}
}

func TestGuard_RevokeAndCheckConcurrently(t *testing.T) {
	ports := make([]scanner.Port, 10)
	for i := range ports {
		ports[i] = scanner.Port{Number: i + 1, Protocol: "tcp"}
	}

	g := portguard.New(ports)

	var wg sync.WaitGroup
	for _, p := range ports[:5] {
		wg.Add(1)
		go func(pp scanner.Port) {
			defer wg.Done()
			g.Revoke(pp)
		}(p)
	}
	wg.Wait()

	// ports 6-10 should still be allowed
	for _, p := range ports[5:] {
		if !g.IsAllowed(p) {
			t.Errorf("port %d should still be allowed", p.Number)
		}
	}
}
