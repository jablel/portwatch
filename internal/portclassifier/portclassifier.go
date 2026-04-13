// Package portclassifier assigns a tier label to each scanned port
// based on its number: system (0–1023), registered (1024–49151),
// or dynamic (49152–65535). Custom overrides take precedence.
package portclassifier

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Tier represents the classification bucket for a port.
type Tier string

const (
	TierSystem     Tier = "system"
	TierRegistered Tier = "registered"
	TierDynamic    Tier = "dynamic"
)

// Classifier assigns tiers to ports.
type Classifier struct {
	mu        sync.RWMutex
	overrides map[string]Tier // key: "<port>/<proto>"
}

// New returns a ready-to-use Classifier.
func New() *Classifier {
	return &Classifier{overrides: make(map[string]Tier)}
}

// Override sets a custom tier for a specific port+protocol pair.
func (c *Classifier) Override(p scanner.Port, t Tier) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.overrides[portKey(p)] = t
}

// Classify returns the tier for p, respecting any registered override.
func (c *Classifier) Classify(p scanner.Port) Tier {
	c.mu.RLock()
	if t, ok := c.overrides[portKey(p)]; ok {
		c.mu.RUnlock()
		return t
	}
	c.mu.RUnlock()
	return tierForNumber(p.Number)
}

// ClassifyAll returns a map from each port's String() representation to its tier.
func (c *Classifier) ClassifyAll(ports []scanner.Port) map[string]Tier {
	out := make(map[string]Tier, len(ports))
	for _, p := range ports {
		out[p.String()] = c.Classify(p)
	}
	return out
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Proto)
}

func tierForNumber(n uint16) Tier {
	switch {
	case n <= 1023:
		return TierSystem
	case n <= 49151:
		return TierRegistered
	default:
		return TierDynamic
	}
}
