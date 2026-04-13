package portevict_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portevict"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: uint16(number), Proto: proto}
}

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestObserve_NoEvictionWhilePortPresent(t *testing.T) {
	tr := portevict.New(64)
	ports := []scanner.Port{makePort(80, "tcp")}

	tr.Observe(ports, t0)
	tr.Observe(ports, t0.Add(time.Second))

	if got := tr.Evicted(); len(got) != 0 {
		t.Fatalf("expected 0 evictions, got %d", len(got))
	}
}

func TestObserve_EvictsDisappearedPort(t *testing.T) {
	tr := portevict.New(64)
	ports := []scanner.Port{makePort(443, "tcp")}

	tr.Observe(ports, t0)
	tr.Observe(nil, t0.Add(5*time.Second))

	evicted := tr.Evicted()
	if len(evicted) != 1 {
		t.Fatalf("expected 1 eviction, got %d", len(evicted))
	}
	rec := evicted[0]
	if rec.Port.Number != 443 {
		t.Errorf("expected port 443, got %d", rec.Port.Number)
	}
	if rec.Duration != 0 {
		t.Errorf("expected zero duration for single-scan port, got %v", rec.Duration)
	}
	if rec.EvictedAt != t0.Add(5*time.Second) {
		t.Errorf("unexpected EvictedAt: %v", rec.EvictedAt)
	}
}

func TestObserve_DurationSpansFirstToLastSeen(t *testing.T) {
	tr := portevict.New(64)
	port := []scanner.Port{makePort(8080, "tcp")}

	tr.Observe(port, t0)
	tr.Observe(port, t0.Add(10*time.Second))
	tr.Observe(nil, t0.Add(20*time.Second))

	evicted := tr.Evicted()
	if len(evicted) != 1 {
		t.Fatalf("expected 1 eviction, got %d", len(evicted))
	}
	if evicted[0].Duration != 10*time.Second {
		t.Errorf("expected 10s duration, got %v", evicted[0].Duration)
	}
}

func TestEvicted_EvictsOldestWhenFull(t *testing.T) {
	tr := portevict.New(2)

	for _, n := range []int{80, 443, 8080} {
		tr.Observe([]scanner.Port{makePort(n, "tcp")}, t0)
		tr.Observe(nil, t0.Add(time.Second))
	}

	evicted := tr.Evicted()
	if len(evicted) != 2 {
		t.Fatalf("expected 2 records (maxLen), got %d", len(evicted))
	}
	// oldest (port 80) should have been dropped
	if evicted[0].Port.Number == 80 {
		t.Errorf("port 80 should have been evicted from the ring")
	}
}

func TestActiveCount_TracksLivePorts(t *testing.T) {
	tr := portevict.New(64)
	ports := []scanner.Port{makePort(22, "tcp"), makePort(80, "tcp")}

	tr.Observe(ports, t0)
	if got := tr.ActiveCount(); got != 2 {
		t.Errorf("expected 2 active, got %d", got)
	}

	tr.Observe(nil, t0.Add(time.Second))
	if got := tr.ActiveCount(); got != 0 {
		t.Errorf("expected 0 active after eviction, got %d", got)
	}
}

func TestObserve_PortReturnsClearsEvictionHistory(t *testing.T) {
	tr := portevict.New(64)
	port := []scanner.Port{makePort(3306, "tcp")}

	tr.Observe(port, t0)
	tr.Observe(nil, t0.Add(time.Second)) // evicted
	tr.Observe(port, t0.Add(2*time.Second)) // returns
	tr.Observe(nil, t0.Add(3*time.Second)) // evicted again

	if got := tr.Evicted(); len(got) != 2 {
		t.Fatalf("expected 2 eviction records, got %d", len(got))
	}
}
