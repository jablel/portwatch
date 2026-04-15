package portwatch_test

import (
	"context"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/portwatch"
	"github.com/user/portwatch/internal/scanner"
)

// stubScanner returns a fixed port list on every call.
type stubScanner struct {
	ports []scanner.Port
}

func (s *stubScanner) Scan(_ context.Context) ([]scanner.Port, error) {
	return s.ports, nil
}

func newWatcher(ports []scanner.Port) (*portwatch.Watcher, *strings.Builder) {
	buf := &strings.Builder{}
	n, _ := notifier.New(notifier.BackendStdout, buf)
	sc := &scanner.Scanner{} // replaced via stub in practice; kept for type
	_ = sc
	// Build watcher using the stub indirectly through the public constructor
	// by wrapping with a real scanner that we override below.
	w := portwatch.NewWithScanner(&stubScanner{ports: ports}, n)
	return w, buf
}

func TestTick_ReturnsScannedPorts(t *testing.T) {
	ports := []scanner.Port{
		{Number: 80, Protocol: "tcp"},
		{Number: 443, Protocol: "tcp"},
	}
	w, _ := newWatcher(ports)

	res, err := w.Tick(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(res.Ports))
	}
	if res.ScannedAt.IsZero() {
		t.Error("ScannedAt should not be zero")
	}
}

func TestTick_DiffIsEmptyOnSecondIdenticalScan(t *testing.T) {
	ports := []scanner.Port{{Number: 22, Protocol: "tcp"}}
	w, _ := newWatcher(ports)

	if _, err := w.Tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	res, err := w.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Diff.Added) != 0 || len(res.Diff.Removed) != 0 {
		t.Errorf("expected empty diff, got added=%v removed=%v", res.Diff.Added, res.Diff.Removed)
	}
}

func TestReset_ClearsPreviousState(t *testing.T) {
	ports := []scanner.Port{{Number: 8080, Protocol: "tcp"}}
	w, _ := newWatcher(ports)

	if _, err := w.Tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	w.Reset()

	res, err := w.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Diff.Added) != 1 {
		t.Errorf("expected 1 added port after reset, got %d", len(res.Diff.Added))
	}
}
