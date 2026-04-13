// Package portscorer assigns a numeric risk score to a port based on
// classification, trend, and lifecycle state. Higher scores indicate
// ports that warrant closer attention.
package portscorer

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Weights applied to each contributing signal.
const (
	weightDynamic   = 10 // dynamic/ephemeral range ports are noisier
	weightRising    = 15 // port appearance frequency is increasing
	weightFalling   = 5  // port appearance frequency is decreasing
	weightNew       = 20 // port seen for the first time
	weightClosed    = 8  // port recently disappeared
	weightWellKnown = -5 // well-known ports are usually expected
)

// Classifier labels a port's range category.
type Classifier interface {
	Classify(p scanner.Port) string // "system", "registered", "dynamic"
}

// Trencher reports the trend direction for a port.
type Trencher interface {
	Trend(p scanner.Port) string // "rising", "falling", "stable"
}

// Lifecycler reports the lifecycle state of a port.
type Lifecycler interface {
	State(p scanner.Port) string // "new", "active", "closed", "flapping"
}

// Scorer computes risk scores for ports.
type Scorer struct {
	mu         sync.Mutex
	classifier Classifier
	trencher   Trencher
	lifecycler Lifecycler
}

// New returns a Scorer wired to the provided signal sources.
func New(c Classifier, t Trencher, l Lifecycler) *Scorer {
	return &Scorer{classifier: c, trencher: t, lifecycler: l}
}

// Score returns a non-negative integer risk score for p.
// A score of 0 means no elevated risk was detected.
func (s *Scorer) Score(p scanner.Port) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	score := 0

	switch s.classifier.Classify(p) {
	case "dynamic":
		score += weightDynamic
	case "system":
		score += weightWellKnown
	}

	switch s.trencher.Trend(p) {
	case "rising":
		score += weightRising
	case "falling":
		score += weightFalling
	}

	switch s.lifecycler.State(p) {
	case "new":
		score += weightNew
	case "closed":
		score += weightClosed
	}

	if score < 0 {
		return 0
	}
	return score
}

// ScoreAll returns a map of port → score for every port in ps.
func (s *Scorer) ScoreAll(ps []scanner.Port) map[scanner.Port]int {
	out := make(map[scanner.Port]int, len(ps))
	for _, p := range ps {
		out[p] = s.Score(p)
	}
	return out
}
