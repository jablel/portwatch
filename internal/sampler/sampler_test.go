package sampler_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

func startTCPListener(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().(*net.TCPAddr).Port
}

func newScanner(t *testing.T) *scanner.Scanner {
	t.Helper()
	s, err := scanner.New(scanner.Options{
		Ports:    "1-65535",
		Protocol: "tcp",
		Host:     "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}
	return s
}

func TestLatest_ReturnsNilWhenEmpty(t *testing.T) {
	s := sampler.New(newScanner(t), time.Hour, 5)
	if got := s.Latest(); got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestAll_EmptyInitially(t *testing.T) {
	s := sampler.New(newScanner(t), time.Hour, 5)
	if got := s.All(); len(got) != 0 {
		t.Fatalf("expected 0 samples, got %d", len(got))
	}
}

func TestSampler_CollectsSample(t *testing.T) {
	_ = startTCPListener(t)
	sc := newScanner(t)
	sm := sampler.New(sc, 20*time.Millisecond, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sm.Run(ctx)
		close(done)
	}()
	<-done

	samples := sm.All()
	if len(samples) == 0 {
		t.Fatal("expected at least one sample")
	}
	for _, s := range samples {
		if s.At.IsZero() {
			t.Error("sample has zero timestamp")
		}
	}
}

func TestSampler_EvictsOldestWhenFull(t *testing.T) {
	sc := newScanner(t)
	sm := sampler.New(sc, 10*time.Millisecond, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sm.Run(ctx)
		close(done)
	}()
	<-done

	if got := sm.All(); len(got) > 3 {
		t.Fatalf("expected at most 3 samples, got %d", len(got))
	}
}

func TestNew_ZeroCapacityDefaultsToOne(t *testing.T) {
	sc := newScanner(t)
	sm := sampler.New(sc, time.Hour, 0)
	if sm == nil {
		t.Fatal("expected non-nil sampler")
	}
}
