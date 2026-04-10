package debounce_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/debounce"
)

func TestTrigger_FiresAfterWait(t *testing.T) {
	d := debounce.New(30 * time.Millisecond)

	var called int32
	d.Trigger("port:8080", func(key string) {
		atomic.AddInt32(&called, 1)
	})

	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected callback to fire once, got %d", called)
	}
}

func TestTrigger_CoalescesRapidCalls(t *testing.T) {
	d := debounce.New(40 * time.Millisecond)

	var count int32
	for i := 0; i < 5; i++ {
		d.Trigger("port:9090", func(key string) {
			atomic.AddInt32(&count, 1)
		})
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(80 * time.Millisecond)
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected exactly 1 call after coalescing, got %d", count)
	}
}

func TestTrigger_DifferentKeysAreIndependent(t *testing.T) {
	d := debounce.New(30 * time.Millisecond)

	var mu sync.Mutex
	fired := make(map[string]int)

	cb := func(key string) {
		mu.Lock()
		fired[key]++
		mu.Unlock()
	}

	d.Trigger("port:80", cb)
	d.Trigger("port:443", cb)

	time.Sleep(70 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if fired["port:80"] != 1 || fired["port:443"] != 1 {
		t.Fatalf("expected each key fired once, got %v", fired)
	}
}

func TestTrigger_ZeroWaitCallsImmediately(t *testing.T) {
	d := debounce.New(0)

	var called int32
	d.Trigger("port:22", func(key string) {
		atomic.AddInt32(&called, 1)
	})

	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected immediate call for zero wait duration")
	}
}

func TestCancel_StopsPendingTimer(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)

	var called int32
	d.Trigger("port:3000", func(key string) {
		atomic.AddInt32(&called, 1)
	})
	d.Cancel("port:3000")

	time.Sleep(80 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected callback to be cancelled")
	}
}

func TestPending_CountsActiveTimers(t *testing.T) {
	d := debounce.New(100 * time.Millisecond)

	d.Trigger("port:8080", func(string) {})
	d.Trigger("port:8081", func(string) {})

	if got := d.Pending(); got != 2 {
		t.Fatalf("expected 2 pending, got %d", got)
	}

	time.Sleep(150 * time.Millisecond)
	if got := d.Pending(); got != 0 {
		t.Fatalf("expected 0 pending after timers fire, got %d", got)
	}
}
