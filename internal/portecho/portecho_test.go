package portecho_test

import (
	"testing"

	"github.com/user/portwatch/internal/portecho"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: uint16(number), Proto: proto}
}

func TestScore_UnknownPortReturnsZero(t *testing.T) {
	tr := portecho.New()
	if got := tr.Score(makePort(80, "tcp")); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestScore_PresentInEveryScansReturnsOne(t *testing.T) {
	tr := portecho.New()
	p := makePort(443, "tcp")
	for i := 0; i < 5; i++ {
		tr.Record([]scanner.Port{p})
	}
	if got := tr.Score(p); got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestScore_PresentInHalfScansReturnsHalf(t *testing.T) {
	tr := portecho.New()
	p := makePort(8080, "tcp")
	other := makePort(9090, "tcp")

	// scan 1: both present
	tr.Record([]scanner.Port{p, other})
	// scan 2: only other present
	tr.Record([]scanner.Port{other})

	got := tr.Score(p)
	// p seen 1 out of 2 scans → 0.5
	if got != 0.5 {
		t.Fatalf("expected 0.5, got %v", got)
	}
}

func TestAll_ReturnsAllTrackedPorts(t *testing.T) {
	tr := portecho.New()
	ports := []scanner.Port{
		makePort(22, "tcp"),
		makePort(53, "udp"),
	}
	tr.Record(ports)

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := portecho.New()
	p := makePort(80, "tcp")
	tr.Record([]scanner.Port{p})

	all := tr.All()
	// Mutate the copy – should not affect internal state.
	for k, e := range all {
		e.Seen = 999
		all[k] = e
	}

	if got := tr.Score(p); got != 1.0 {
		t.Fatalf("internal state was mutated; expected 1.0, got %v", got)
	}
}

func TestRecord_MultipleCallsAccumulate(t *testing.T) {
	tr := portecho.New()
	p := makePort(3306, "tcp")

	for i := 0; i < 3; i++ {
		tr.Record([]scanner.Port{p})
	}

	all := tr.All()
	e, ok := all["tcp:3306"]
	if !ok {
		t.Fatal("expected entry for tcp:3306")
	}
	if e.Seen != 3 || e.Total != 3 {
		t.Fatalf("expected Seen=3 Total=3, got Seen=%d Total=%d", e.Seen, e.Total)
	}
}
