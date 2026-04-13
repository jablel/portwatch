package porttrend

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestTrend_StableWhenEvenlyDistributed(t *testing.T) {
	tr := New(2 * time.Second)
	p := makePort("tcp", 80)

	// record equal observations – trend must be stable
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{p})

	if got := tr.Trend(p); got != Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestTrend_UnknownPortIsStable(t *testing.T) {
	tr := New(time.Minute)
	p := makePort("tcp", 9999)
	if got := tr.Trend(p); got != Stable {
		t.Fatalf("expected Stable for unknown port, got %s", got)
	}
}

func TestTrend_RisingWhenMoreInRecentHalf(t *testing.T) {
	tr := New(4 * time.Second)
	p := makePort("tcp", 443)

	// Simulate old observations by injecting directly.
	old := time.Now().Add(-3 * time.Second)
	tr.mu.Lock()
	k := portKey(p)
	tr.buckets[k] = []bucket{{count: 1, at: old}}
	tr.mu.Unlock()

	// Two recent observations.
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{p})

	if got := tr.Trend(p); got != Rising {
		t.Fatalf("expected Rising, got %s", got)
	}
}

func TestTrend_FallingWhenMoreInOlderHalf(t *testing.T) {
	tr := New(4 * time.Second)
	p := makePort("udp", 53)

	old := time.Now().Add(-3 * time.Second)
	tr.mu.Lock()
	k := portKey(p)
	tr.buckets[k] = []bucket{
		{count: 1, at: old},
		{count: 1, at: old},
	}
	tr.mu.Unlock()

	// One recent observation.
	tr.Record([]scanner.Port{p})

	if got := tr.Trend(p); got != Falling {
		t.Fatalf("expected Falling, got %s", got)
	}
}

func TestEvict_RemovesStaleEntries(t *testing.T) {
	tr := New(100 * time.Millisecond)
	p := makePort("tcp", 22)
	tr.Record([]scanner.Port{p})

	time.Sleep(150 * time.Millisecond)

	tr.mu.Lock()
	tr.evict(time.Now())
	_, exists := tr.buckets[portKey(p)]
	tr.mu.Unlock()

	if exists {
		t.Fatal("expected stale bucket to be evicted")
	}
}

func TestTrend_String(t *testing.T) {
	cases := []struct {
		trend Trend
		want  string
	}{
		{Stable, "stable"},
		{Rising, "rising"},
		{Falling, "falling"},
	}
	for _, c := range cases {
		if got := c.trend.String(); got != c.want {
			t.Errorf("Trend(%d).String() = %q, want %q", c.trend, got, c.want)
		}
	}
}
