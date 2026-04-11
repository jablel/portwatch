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

// TestPipeline_DetectsNewPort opens a listener mid-run and verifies the diff
// captures the newly added port.
func TestPipeline_DetectsNewPort(t *testing.T) {
	// First pass — no open ports.
	sc := scanner.New("127.0.0.1", []int{}, "tcp")
	f := filter.New(nil)
	n, _ := notifier.New(notifier.BackendStdout, nil)
	p := pipeline.New(sc, f, n)

	first, err := p.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Open a real listener so the next scan can detect it.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	// Build a scanner scoped to that exact port.
	sc2 := scanner.New("127.0.0.1", []int{port}, "tcp")
	p2 := pipeline.New(sc2, f, n)

	second, err := p2.Run(context.Background(), first.Ports)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}

	if len(second.Diff.Added) == 0 {
		t.Error("expected at least one added port in diff")
	}
}

// TestPipeline_DetectsRemovedPort confirms that a closed listener is reported
// as removed in the diff.
func TestPipeline_DetectsRemovedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	sc := scanner.New("127.0.0.1", []int{port}, "tcp")
	f := filter.New(nil)
	n, _ := notifier.New(notifier.BackendStdout, nil)
	p := pipeline.New(sc, f, n)

	first, err := p.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Close the listener before the second scan.
	ln.Close()

	second, err := p.Run(context.Background(), first.Ports)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}

	if len(second.Diff.Removed) == 0 {
		t.Error("expected at least one removed port in diff")
	}
}
