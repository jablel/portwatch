// Package portranker ranks observed ports by a composite priority score
// derived from classification, trend, and lifecycle state.
package portranker

import (
	"sort"
	"sync"

	"github.com/example/portwatch/internal/scanner"
)

// RankEntry holds a port and its computed rank score.
type RankEntry struct {
	Port  scanner.Port
	Score float64
}

// Ranker computes and returns a ranked list of ports.
type Ranker struct {
	mu      sync.Mutex
	weights Weights
}

// Weights controls how much each signal contributes to the final score.
type Weights struct {
	// ClassBonus is added for system ports (0-1023).
	ClassBonus float64
	// DynamicPenalty is subtracted for ephemeral/dynamic ports (49152+).
	DynamicPenalty float64
	// TCPBonus is added for TCP ports.
	TCPBonus float64
}

// DefaultWeights returns sensible default weights.
func DefaultWeights() Weights {
	return Weights{
		ClassBonus:     10.0,
		DynamicPenalty: 5.0,
		TCPBonus:       2.0,
	}
}

// New creates a Ranker with the given weights.
func New(w Weights) *Ranker {
	return &Ranker{weights: w}
}

// Rank returns ports sorted by descending priority score.
func (r *Ranker) Rank(ports []scanner.Port) []RankEntry {
	r.mu.Lock()
	w := r.weights
	r.mu.Unlock()

	entries := make([]RankEntry, len(ports))
	for i, p := range ports {
		entries[i] = RankEntry{Port: p, Score: r.score(p, w)}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Port.Number < entries[j].Port.Number
	})

	return entries
}

// SetWeights updates the ranking weights used for future calls.
func (r *Ranker) SetWeights(w Weights) {
	r.mu.Lock()
	r.weights = w
	r.mu.Unlock()
}

func (r *Ranker) score(p scanner.Port, w Weights) float64 {
	var s float64

	switch {
	case p.Number <= 1023:
		s += w.ClassBonus
	case p.Number >= 49152:
		s -= w.DynamicPenalty
	}

	if p.Protocol == "tcp" {
		s += w.TCPBonus
	}

	return s
}
