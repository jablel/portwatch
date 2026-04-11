package rollup_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
)

func TestRollup_CoalescesRapidBurst(t *testing.T) {
	var calls int32
	r := rollup.New(50*time.Millisecond, func(d rollup.Diff) {
		atomic.AddInt32(&calls, 1)
	})

	for i := uint16(1000); i < 1020; i++ {
		r.Add(makeDiff([]uint16{i}, nil))
	}

	time.Sleep(200 * time.Millisecond)

	if n := atomic.LoadInt32(&calls); n != 1 {
		t.Fatalf("expected exactly 1 handler call, got %d", n)
	}
}

func TestRollup_SeparateBurstsProduceSeparateCalls(t *testing.T) {
	var calls int32
	r := rollup.New(30*time.Millisecond, func(d rollup.Diff) {
		atomic.AddInt32(&calls, 1)
	})

	r.Add(makeDiff([]uint16{1111}, nil))
	time.Sleep(120 * time.Millisecond) // let first window close

	r.Add(makeDiff([]uint16{2222}, nil))
	time.Sleep(120 * time.Millisecond) // let second window close

	if n := atomic.LoadInt32(&calls); n != 2 {
		t.Fatalf("expected 2 handler calls, got %d", n)
	}
}
