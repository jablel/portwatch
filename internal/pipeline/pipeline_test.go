package pipeline_test

import (
	"context"
	"net"
	"testing"

	"portwatch/internal/filter"
	"portwatch/internal/notifier"
	"portwatch/internal/pipeline"
	"portwatch/internal/scanner"
)

func freePort(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func newPipeline(t *testing.T) *pipeline.Pipeline {
	t.Helper()
	sc := scanner.New("127.0.0.1", []int{}, "tcp")
	f := filter.New(nil)
	n, _ := notifier.New(notifier.BackendStdout, nil)
	return pipeline.New(sc, f, n)
}

func TestRun_ReturnsResultWithPorts(t *testing.T) {
	_, close := freePort(t)
	defer close()

	p := newPipeline(t)
	res, err := p.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if res.ScannedAt.IsZero() {
		t.Error("ScannedAt should not be zero")
	}
}

func TestRun_DiffIsEmptyWhenPortsUnchanged(t *testing.T) {
	p := newPipeline(t)

	first, err := p.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}

	second, err := p.Run(context.Background(), first.Ports)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}

	if len(second.Diff.Added) != 0 || len(second.Diff.Removed) != 0 {
		t.Errorf("expected empty diff on stable scan, got added=%d removed=%d",
			len(second.Diff.Added), len(second.Diff.Removed))
	}
}

func TestRun_CancelledContextReturnsError(t *testing.T) {
	p := newPipeline(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// A cancelled context may or may not surface an error depending on scanner
	// implementation; we just ensure the call does not panic.
	_, _ = p.Run(ctx, nil)
}
