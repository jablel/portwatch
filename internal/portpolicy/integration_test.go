package portpolicy_test

import (
	"sync"
	"testing"

	"portwatch/internal/portpolicy"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

func TestEvaluator_ConcurrentAddAndEvaluate(t *testing.T) {
	e := portpolicy.New()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = e.Add(portpolicy.Policy{
				Name:     fmt.Sprintf("p%d", n),
				MinPort:  n * 100,
				MaxPort:  n*100 + 99,
				OnAdded:  true,
				Severity: portpolicy.SeverityWarn,
			})
		}(i)
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			diff := state.Diff{
				Added: []scanner.Port{{Number: 80, Protocol: "tcp"}},
			}
			_ = e.Evaluate(diff)
		}()
	}

	wg.Wait()
}

func TestEvaluator_AnyProtocolMatchesBoth(t *testing.T) {
	e := portpolicy.New()
	_ = e.Add(portpolicy.Policy{
		Name:     "any-proto",
		MinPort:  443,
		MaxPort:  443,
		Protocol: "",
		OnAdded:  true,
		Severity: portpolicy.SeverityCritical,
		Message:  "port 443 appeared",
	})

	diff := state.Diff{
		Added: []scanner.Port{
			{Number: 443, Protocol: "tcp"},
			{Number: 443, Protocol: "udp"},
		},
	}
	v := e.Evaluate(diff)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations for both protocols, got %d", len(v))
	}
}
