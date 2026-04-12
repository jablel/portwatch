package presencemap_test

import (
	"sort"
	"testing"

	"github.com/user/portwatch/internal/presencemap"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(numbers ...int) []scanner.Port {
	ports := make([]scanner.Port, len(numbers))
	for i, n := range numbers {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestStreak_StartsAtOne(t *testing.T) {
	pm := presencemap.New()
	ports := makePorts(80)
	pm.Observe(ports)
	if got := pm.Streak(ports[0]); got != 1 {
		t.Fatalf("expected streak 1, got %d", got)
	}
}

func TestStreak_IncrementsOnConsecutiveObservations(t *testing.T) {
	pm := presencemap.New()
	p := makePorts(443)
	for i := 1; i <= 4; i++ {
		pm.Observe(p)
		if got := pm.Streak(p[0]); got != i {
			t.Fatalf("after %d observations expected streak %d, got %d", i, i, got)
		}
	}
}

func TestStreak_ResetsWhenPortDisappears(t *testing.T) {
	pm := presencemap.New()
	p := makePorts(22)
	pm.Observe(p)
	pm.Observe(p)
	// port disappears
	pm.Observe([]scanner.Port{})
	if got := pm.Streak(p[0]); got != 0 {
		t.Fatalf("expected streak 0 after disappearance, got %d", got)
	}
}

func TestStreak_UnknownPortReturnsZero(t *testing.T) {
	pm := presencemap.New()
	p := scanner.Port{Number: 9999, Protocol: "tcp"}
	if got := pm.Streak(p); got != 0 {
		t.Fatalf("expected 0 for untracked port, got %d", got)
	}
}

func TestStable_ReturnsPortsAboveThreshold(t *testing.T) {
	pm := presencemap.New()
	all := makePorts(80, 443, 8080)
	for i := 0; i < 3; i++ {
		pm.Observe(all)
	}
	// add a new port that has only 1 streak
	pm.Observe(append(all, makePorts(9090)...))

	stable := pm.Stable(3)
	if len(stable) != 3 {
		t.Fatalf("expected 3 stable ports, got %d", len(stable))
	}
	nums := make([]int, len(stable))
	for i, s := range stable {
		nums[i] = s.Number
	}
	sort.Ints(nums)
	expected := []int{80, 443, 8080}
	for i, v := range expected {
		if nums[i] != v {
			t.Fatalf("stable ports mismatch: got %v", nums)
		}
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	pm := presencemap.New()
	pm.Observe(makePorts(80, 443))
	pm.Reset()
	if got := pm.Stable(1); len(got) != 0 {
		t.Fatalf("expected empty after reset, got %d entries", len(got))
	}
}
