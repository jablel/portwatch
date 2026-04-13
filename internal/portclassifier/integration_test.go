package portclassifier_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/scanner"
)

func TestClassifier_ConcurrentClassify(t *testing.T) {
	c := portclassifier.New(nil)
	ports := []scanner.Port{
		{Number: 22, Protocol: "tcp"},
		{Number: 80, Protocol: "tcp"},
		{Number: 8080, Protocol: "tcp"},
		{Number: 60000, Protocol: "udp"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, p := range ports {
				_ = c.Classify(p)
			}
		}()
	}
	wg.Wait()
}

func TestClassifier_ClassifyAll_AllPortsLabelled(t *testing.T) {
	c := portclassifier.New(nil)
	ports := []scanner.Port{
		{Number: 443, Protocol: "tcp"},
		{Number: 3306, Protocol: "tcp"},
		{Number: 55000, Protocol: "udp"},
	}

	results := c.ClassifyAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(results))
	}
	for _, r := range results {
		if r.Tier == "" {
			t.Errorf("port %d has empty tier", r.Port.Number)
		}
	}
}

func TestClassifier_CustomOverridesBuiltin(t *testing.T) {
	custom := map[uint16]string{80: "my-web"}
	c := portclassifier.New(custom)

	p := scanner.Port{Number: 80, Protocol: "tcp"}
	r := c.Classify(p)
	if r.Label != "my-web" {
		t.Errorf("expected label 'my-web', got %q", r.Label)
	}
}
