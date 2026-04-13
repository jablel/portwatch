package portclassifier_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/scanner"
)

func TestClassifier_ConcurrentClassify(t *testing.T) {
	c := portclassifier.New()
	ports := []scanner.Port{
		{Number: 80, Proto: "tcp"},
		{Number: 8080, Proto: "tcp"},
		{Number: 55000, Proto: "udp"},
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
	c := portclassifier.New()
	var ports []scanner.Port
	for i := uint16(0); i < 100; i++ {
		ports = append(ports, scanner.Port{Number: i * 655, Proto: "tcp"})
	}
	res := c.ClassifyAll(ports)
	if len(res) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(res))
	}
	for _, p := range ports {
		if _, ok := res[p.String()]; !ok {
			t.Errorf("port %s missing from ClassifyAll result", p.String())
		}
	}
}

func TestClassifier_CustomOverridesBuiltin(t *testing.T) {
	c := portclassifier.New()
	p := scanner.Port{Number: 443, Proto: "tcp"}

	if got := c.Classify(p); got != portclassifier.TierSystem {
		t.Fatalf("pre-override: expected system, got %s", got)
	}

	c.Override(p, portclassifier.TierRegistered)

	if got := c.Classify(p); got != portclassifier.TierRegistered {
		t.Fatalf("post-override: expected registered, got %s", got)
	}
}
