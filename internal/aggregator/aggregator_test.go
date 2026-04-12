package aggregator_test

import (
	"sort"
	"testing"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: uint16(n), Protocol: "tcp"}
	}
	return ports
}

func sortPorts(ports []scanner.Port) []scanner.Port {
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Number < ports[j].Number
	})
	return ports
}

func TestMerge_EmptyAggregator(t *testing.T) {
	a := aggregator.New()
	got := a.Merge()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestMerge_SingleSource(t *testing.T) {
	a := aggregator.New()
	a.Update("tcp", makePorts(80, 443))

	got := sortPorts(a.Merge())
	if len(got) != 2 || got[0].Number != 80 || got[1].Number != 443 {
		t.Fatalf("unexpected ports: %v", got)
	}
}

func TestMerge_DeduplicatesAcrossSources(t *testing.T) {
	a := aggregator.New()
	a.Update("src1", makePorts(80, 443))
	a.Update("src2", makePorts(443, 8080))

	got := sortPorts(a.Merge())
	if len(got) != 3 {
		t.Fatalf("expected 3 unique ports, got %d: %v", len(got), got)
	}
}

func TestUpdate_ReplacesExistingSource(t *testing.T) {
	a := aggregator.New()
	a.Update("src", makePorts(80))
	a.Update("src", makePorts(443))

	got := a.Merge()
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443 after update, got %v", got)
	}
}

func TestRemove_DeletesSource(t *testing.T) {
	a := aggregator.New()
	a.Update("src1", makePorts(80))
	a.Update("src2", makePorts(443))
	a.Remove("src1")

	got := a.Merge()
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443 after remove, got %v", got)
	}
}

func TestSources_ReturnsRegisteredNames(t *testing.T) {
	a := aggregator.New()
	a.Update("alpha", makePorts(80))
	a.Update("beta", makePorts(443))

	srcs := a.Sources()
	if len(srcs) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(srcs))
	}
}
