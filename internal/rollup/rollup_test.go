package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/state"
)

func makeDiff(added, removed []uint16) rollup.Diff {
	d := rollup.Diff{}
	for _, p := range added {
		d.Added = append(d.Added, state.Port{Port: p, Proto: "tcp"})
	}
	for _, p := range removed {
		d.Removed = append(d.Removed, state.Port{Port: p, Proto: "tcp"})
	}
	return d
}

func TestAdd_FlushesAfterWindow(t *testing.T) {
	done := make(chan rollup.Diff, 1)
	r := rollup.New(20*time.Millisecond, func(d rollup.Diff) { done <- d })

	r.Add(makeDiff([]uint16{8080}, nil))

	select {
	case d := <-done:
		if len(d.Added) != 1 || d.Added[0].Port != 8080 {
			t.Fatalf("unexpected diff: %+v", d)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for flush")
	}
}

func TestAdd_MergesMultipleDiffs(t *testing.T) {
	done := make(chan rollup.Diff, 1)
	r := rollup.New(40*time.Millisecond, func(d rollup.Diff) { done <- d })

	r.Add(makeDiff([]uint16{8080}, nil))
	r.Add(makeDiff([]uint16{9090}, nil))
	r.Add(makeDiff(nil, []uint16{443}))

	select {
	case d := <-done:
		if len(d.Added) != 2 {
			t.Fatalf("expected 2 added ports, got %d", len(d.Added))
		}
		if len(d.Removed) != 1 {
			t.Fatalf("expected 1 removed port, got %d", len(d.Removed))
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out")
	}
}

func TestFlush_ImmediateFlush(t *testing.T) {
	done := make(chan rollup.Diff, 1)
	r := rollup.New(10*time.Second, func(d rollup.Diff) { done <- d })

	r.Add(makeDiff([]uint16{22}, nil))
	r.Flush()

	select {
	case d := <-done:
		if len(d.Added) != 1 {
			t.Fatalf("expected 1 added port, got %d", len(d.Added))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for immediate flush")
	}
}

func TestFlush_NoPendingIsNoop(t *testing.T) {
	called := false
	r := rollup.New(10*time.Millisecond, func(d rollup.Diff) { called = true })
	r.Flush()
	time.Sleep(30 * time.Millisecond)
	if called {
		t.Fatal("handler should not have been called with no pending diff")
	}
}

func TestAdd_ZeroWindowFlushesImmediately(t *testing.T) {
	done := make(chan rollup.Diff, 1)
	r := rollup.New(0, func(d rollup.Diff) { done <- d })

	r.Add(makeDiff([]uint16{3000}, nil))

	select {
	case d := <-done:
		if len(d.Added) != 1 {
			t.Fatalf("unexpected diff: %+v", d)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out")
	}
}
