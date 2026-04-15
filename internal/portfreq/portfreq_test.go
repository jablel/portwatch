package portfreq_test

import (
	"testing"

	"portwatch/internal/portfreq"
	"portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestFrequency_UnknownPortReturnsZero(t *testing.T) {
	tr := portfreq.New(5)
	if got := tr.Frequency(makePort(80, "tcp")); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFrequency_PresentInEveryScanReturnsOne(t *testing.T) {
	tr := portfreq.New(4)
	p := makePort(443, "tcp")
	for i := 0; i < 4; i++ {
		tr.Record([]scanner.Port{p})
	}
	if got := tr.Frequency(p); got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestFrequency_PresentInHalfScansReturnsHalf(t *testing.T) {
	tr := portfreq.New(4)
	p := makePort(8080, "tcp")
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{})
	tr.Record([]scanner.Port{})
	got := tr.Frequency(p)
	if got < 0.49 || got > 0.51 {
		t.Fatalf("expected ~0.5, got %v", got)
	}
}

func TestFrequency_DecaysWhenPortDisappears(t *testing.T) {
	tr := portfreq.New(3)
	p := makePort(22, "tcp")
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{p})
	tr.Record([]scanner.Port{}) // port absent – count decays

	got := tr.Frequency(p)
	if got >= 1.0 {
		t.Fatalf("expected decayed frequency, got %v", got)
	}
}

func TestAll_ReturnsAllTrackedPorts(t *testing.T) {
	tr := portfreq.New(2)
	tr.Record([]scanner.Port{makePort(80, "tcp"), makePort(53, "udp")})

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_EmptyBeforeAnyRecord(t *testing.T) {
	tr := portfreq.New(3)
	if got := tr.All(); len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestNew_WindowBelowOneClamped(t *testing.T) {
	tr := portfreq.New(0)
	p := makePort(80, "tcp")
	tr.Record([]scanner.Port{p})
	if got := tr.Frequency(p); got != 1.0 {
		t.Fatalf("expected 1.0 with window=1, got %v", got)
	}
}
