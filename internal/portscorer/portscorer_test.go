package portscorer_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscorer"
	"github.com/user/portwatch/internal/scanner"
)

// --- fakes ---

type fakeClassifier struct{ label string }

func (f *fakeClassifier) Classify(_ scanner.Port) string { return f.label }

type fakeTrencher struct{ label string }

func (f *fakeTrencher) Trend(_ scanner.Port) string { return f.label }

type fakeLifecycler struct{ label string }

func (f *fakeLifecycler) State(_ scanner.Port) string { return f.label }

func makePort(num uint16) scanner.Port {
	return scanner.Port{Number: num, Protocol: "tcp"}
}

func newScorer(cls, trend, life string) *portscorer.Scorer {
	return portscorer.New(
		&fakeClassifier{cls},
		&fakeTrencher{trend},
		&fakeLifecycler{life},
	)
}

// --- tests ---

func TestScore_NewDynamicRisingIsHighest(t *testing.T) {
	s := newScorer("dynamic", "rising", "new")
	got := s.Score(makePort(49152))
	// 10 + 15 + 20 = 45
	if got != 45 {
		t.Fatalf("expected 45, got %d", got)
	}
}

func TestScore_WellKnownStableActiveIsLow(t *testing.T) {
	s := newScorer("system", "stable", "active")
	got := s.Score(makePort(80))
	// -5 clamped to 0
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestScore_RegisteredFallingClosed(t *testing.T) {
	s := newScorer("registered", "falling", "closed")
	got := s.Score(makePort(8080))
	// 0 + 5 + 8 = 13
	if got != 13 {
		t.Fatalf("expected 13, got %d", got)
	}
}

func TestScoreAll_ReturnsMappedScores(t *testing.T) {
	s := newScorer("dynamic", "stable", "new")
	ports := []scanner.Port{makePort(1024), makePort(2048)}
	scores := s.ScoreAll(ports)
	if len(scores) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(scores))
	}
	for _, p := range ports {
		if _, ok := scores[p]; !ok {
			t.Errorf("missing score for port %v", p)
		}
	}
}

func TestScore_NeverNegative(t *testing.T) {
	s := newScorer("system", "falling", "active")
	got := s.Score(makePort(22))
	// -5 + 5 = 0
	if got < 0 {
		t.Fatalf("score must not be negative, got %d", got)
	}
}
