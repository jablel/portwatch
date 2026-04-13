package portlifecycle

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestObserve_NewPort(t *testing.T) {
	tr := New()
	p := makePort(80, "tcp")
	tr.Observe([]scanner.Port{p}, t0)

	e, ok := tr.Get(p)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.State != StateNew {
		t.Errorf("want StateNew, got %s", e.State)
	}
	if e.SeenCount != 1 {
		t.Errorf("want SeenCount=1, got %d", e.SeenCount)
	}
	if !e.FirstSeen.Equal(t0) {
		t.Errorf("unexpected FirstSeen: %v", e.FirstSeen)
	}
}

func TestObserve_ActiveAfterSecondScan(t *testing.T) {
	tr := New()
	p := makePort(443, "tcp")
	tr.Observe([]scanner.Port{p}, t0)
	tr.Observe([]scanner.Port{p}, t0.Add(time.Second))

	e, _ := tr.Get(p)
	if e.State != StateActive {
		t.Errorf("want StateActive, got %s", e.State)
	}
	if e.SeenCount != 2 {
		t.Errorf("want SeenCount=2, got %d", e.SeenCount)
	}
}

func TestObserve_ClosedWhenAbsent(t *testing.T) {
	tr := New()
	p := makePort(8080, "tcp")
	tr.Observe([]scanner.Port{p}, t0)
	tr.Observe([]scanner.Port{}, t0.Add(time.Second))

	e, _ := tr.Get(p)
	if e.State != StateClosed {
		t.Errorf("want StateClosed, got %s", e.State)
	}
	if e.MissCount != 1 {
		t.Errorf("want MissCount=1, got %d", e.MissCount)
	}
}

func TestObserve_MissCountResetsOnReturn(t *testing.T) {
	tr := New()
	p := makePort(22, "tcp")
	tr.Observe([]scanner.Port{p}, t0)
	tr.Observe([]scanner.Port{}, t0.Add(time.Second))
	tr.Observe([]scanner.Port{p}, t0.Add(2*time.Second))

	e, _ := tr.Get(p)
	if e.MissCount != 0 {
		t.Errorf("want MissCount=0 after reappearance, got %d", e.MissCount)
	}
	if e.State != StateActive {
		t.Errorf("want StateActive after reappearance, got %s", e.State)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := New()
	ports := []scanner.Port{
		makePort(80, "tcp"),
		makePort(53, "udp"),
		makePort(443, "tcp"),
	}
	tr.Observe(ports, t0)

	all := tr.All()
	if len(all) != 3 {
		t.Errorf("want 3 entries, got %d", len(all))
	}
}

func TestGet_MissingPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(makePort(9999, "tcp"))
	if ok {
		t.Error("expected ok=false for unknown port")
	}
}

func TestEntry_String(t *testing.T) {
	e := Entry{
		Port:      makePort(80, "tcp"),
		State:     StateActive,
		SeenCount: 5,
		MissCount: 0,
	}
	s := e.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
