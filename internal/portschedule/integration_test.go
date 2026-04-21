package portschedule_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portschedule"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestTracker_ConcurrentObserve(t *testing.T) {
	tr := portschedule.New(1)
	ports := []scanner.Port{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
		makePort(8080, "tcp"),
	}
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(h int) {
			defer wg.Done()
			tr.Observe(ports, now.Add(time.Duration(h)*time.Hour))
		}(i % 24)
	}
	wg.Wait() // must not race
}

func TestTracker_LearnThenViolate(t *testing.T) {
	tr := portschedule.New(2)
	p := makePort(5432, "tcp")
	day := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	// Learn at hours 8 and 9.
	tr.Observe([]scanner.Port{p}, day.Add(8*time.Hour))
	tr.Observe([]scanner.Port{p}, day.Add(9*time.Hour))

	// Seen at hour 2 — violation expected.
	v := tr.Observe([]scanner.Port{p}, day.Add(2*time.Hour))
	if len(v) == 0 {
		t.Fatal("expected a schedule violation, got none")
	}
	if v[0].Hour != 2 {
		t.Errorf("expected violation at hour 2, got %d", v[0].Hour)
	}
}
