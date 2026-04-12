package portclassifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto}
}

func TestClassify_SystemPort(t *testing.T) {
	c := portclassifier.New()
	r := c.Classify(makePort(80, "tcp"))
	if r.Tier != portclassifier.TierSystem {
		t.Fatalf("expected system, got %s", r.Tier)
	}
}

func TestClassify_RegisteredPort(t *testing.T) {
	c := portclassifier.New()
	r := c.Classify(makePort(8080, "tcp"))
	if r.Tier != portclassifier.TierRegistered {
		t.Fatalf("expected registered, got %s", r.Tier)
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	c := portclassifier.New()
	r := c.Classify(makePort(55000, "tcp"))
	if r.Tier != portclassifier.TierDynamic {
		t.Fatalf("expected dynamic, got %s", r.Tier)
	}
}

func TestClassify_PortZero_IsSystem(t *testing.T) {
	c := portclassifier.New()
	r := c.Classify(makePort(0, "tcp"))
	if r.Tier != portclassifier.TierSystem {
		t.Fatalf("expected system for port 0, got %s", r.Tier)
	}
}

func TestOverride_TakesPrecedence(t *testing.T) {
	c := portclassifier.New()
	c.Override(80, portclassifier.TierDynamic)
	r := c.Classify(makePort(80, "tcp"))
	if r.Tier != portclassifier.TierDynamic {
		t.Fatalf("expected override dynamic, got %s", r.Tier)
	}
}

func TestClassifyAll_ReturnsResultPerPort(t *testing.T) {
	c := portclassifier.New()
	ports := []scanner.Port{
		makePort(22, "tcp"),
		makePort(3000, "tcp"),
		makePort(60000, "udp"),
	}
	results := c.ClassifyAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("expected %d results, got %d", len(ports), len(results))
	}
	expected := []portclassifier.Tier{
		portclassifier.TierSystem,
		portclassifier.TierRegistered,
		portclassifier.TierDynamic,
	}
	for i, r := range results {
		if r.Tier != expected[i] {
			t.Errorf("port %d: expected %s, got %s", ports[i].Port, expected[i], r.Tier)
		}
	}
}

func TestClassifyAll_EmptySlice(t *testing.T) {
	c := portclassifier.New()
	results := c.ClassifyAll(nil)
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}
