// Package portclassifier categorises ports into well-known service tiers
// (system, registered, dynamic) and exposes a simple Classify API.
package portclassifier

import "github.com/user/portwatch/internal/scanner"

// Tier represents the classification tier of a port number.
type Tier string

const (
	// TierSystem covers IANA system (well-known) ports 0-1023.
	TierSystem Tier = "system"
	// TierRegistered covers IANA registered ports 1024-49151.
	TierRegistered Tier = "registered"
	// TierDynamic covers ephemeral / dynamic ports 49152-65535.
	TierDynamic Tier = "dynamic"
)

// Result holds the classification outcome for a single port.
type Result struct {
	Port scanner.Port
	Tier Tier
}

// Classifier classifies ports into tiers.
type Classifier struct {
	overrides map[uint16]Tier
}

// New returns a Classifier with no custom overrides.
func New() *Classifier {
	return &Classifier{overrides: make(map[uint16]Tier)}
}

// Override registers a custom tier for a specific port number, taking
// precedence over the default range-based classification.
func (c *Classifier) Override(port uint16, tier Tier) {
	c.overrides[port] = tier
}

// Classify returns the Tier for a single scanner.Port.
func (c *Classifier) Classify(p scanner.Port) Result {
	if t, ok := c.overrides[uint16(p.Port)]; ok {
		return Result{Port: p, Tier: t}
	}
	return Result{Port: p, Tier: tierForNumber(p.Port)}
}

// ClassifyAll classifies a slice of ports and returns a Result per port.
func (c *Classifier) ClassifyAll(ports []scanner.Port) []Result {
	out := make([]Result, len(ports))
	for i, p := range ports {
		out[i] = c.Classify(p)
	}
	return out
}

// tierForNumber maps a raw port number to its default Tier.
func tierForNumber(n int) Tier {
	switch {
	case n <= 1023:
		return TierSystem
	case n <= 49151:
		return TierRegistered
	default:
		return TierDynamic
	}
}
