package watcher_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

// startTCPListener binds an ephemeral TCP port and returns the listener and port number.
func startTCPListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	return ln, ln.Addr().(*net.TCPAddr).Port
}

func TestWatcher_DetectsAddedPort(t *testing.T) {
	s := scanner.New("127.0.0.1", 1, 65535)

	// Start watcher with an empty baseline.
	w := watcher.New(s, 50*time.Millisecond)
	w.Start([]scanner.Port{})
	defer w.Stop()

	// Open a port after the watcher has started.
	ln, _ := startTCPListener(t)
	defer ln.Close()

	select {
	case evt := <-w.Events():
		if len(evt.Added) == 0 {
			t.Errorf("expected at least one added port, got none")
		}
		if evt.DetectedAt.IsZero() {
			t.Error("DetectedAt should not be zero")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenNothingChanges(t *testing.T) {
	// Pre-open a port so baseline and subsequent scans agree.
	ln, _ := startTCPListener(t)
	defer ln.Close()

	s := scanner.New("127.0.0.1", 1, 65535)
	initial, err := s.Scan()
	if err != nil {
		t.Fatalf("initial scan failed: %v", err)
	}

	w := watcher.New(s, 50*time.Millisecond)
	w.Start(initial)
	defer w.Stop()

	select {
	case evt := <-w.Events():
		t.Errorf("unexpected change event: %+v", evt)
	case <-time.After(300 * time.Millisecond):
		// success — no spurious events
	}
}

func TestWatcher_StopClosesEventChannel(t *testing.T) {
	s := scanner.New("127.0.0.1", 1, 1024)
	w := watcher.New(s, 100*time.Millisecond)
	w.Start([]scanner.Port{})
	w.Stop()

	// After Stop the channel must be closed within a reasonable time.
	select {
	case _, ok := <-w.Events():
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for channel close")
	}
}
