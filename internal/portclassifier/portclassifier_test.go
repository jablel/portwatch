package portclassifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Proto: proto}
}

func TestClassify_SystemPort(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(80, "tcp")); got != portclassifier.TierSystem {
		t.Fatalf("expected system, got %s", got)
	}
}

func TestClassify_RegisteredPort(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(8080, "tcp")); got != portclassifier.TierRegistered {
		t.Fatalf("expected registered, got %s", got)
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(55000, "udp")); got != portclassifier.TierDynamic {
		t.Fatalf("expected dynamic, got %s", got)
	}
}

func TestClassify_PortZero_IsSystem(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(0, "tcp")); got != portclassifier.TierSystem {
		t.Fatalf("expected system, got %s", got)
	}
}

func TestClassify_Boundary1023_IsSystem(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(1023, "tcp")); got != portclassifier.TierSystem {
		t.Fatalf("expected system, got %s", got)
	}
}

func TestClassify_Boundary49151_IsRegistered(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(49151, "tcp")); got != portclassifier.TierRegistered {
		t.Fatalf("expected registered, got %s", got)
	}
}

func TestClassify_Boundary49152_IsDynamic(t *testing.T) {
	c := portclassifier.New()
	if got := c.Classify(makePort(49152, "tcp")); got != portclassifier.TierDynamic {
		t.Fatalf("expected dynamic, got %s", got)
	}
}

func TestClassify_CustomOverridesBuiltin(t *testing.T) {
	c := portclassifier.New()
	p := makePort(80, "tcp")
	c.Override(p, portclassifier.TierDynamic)
	if got := c.Classify(p); got != portclassifier.TierDynamic {
		t.Fatalf("expected dynamic override, got %s", got)
	}
}

func TestClassifyAll_AllPortsLabelled(t *testing.T) {
	c := portclassifier.New()
	ports := []scanner.Port{
		makePort(22, "tcp"),
		makePort(3000, "tcp"),
		makePort(60000, "udp"),
	}
	res := c.ClassifyAll(ports)
	if len(res) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res))
	}
	for _, p := range ports {
		if _, ok := res[p.String()]; !ok {
			t.Errorf("missing entry for %s", p.String())
		}
	}
}
