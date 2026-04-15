// Package portсummary provides a concise summary of the current port landscape,
// including counts by protocol, class, and any recent changes.
package portсummary

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Summary holds aggregated statistics about a set of ports.
type Summary struct {
	Total    int
	ByProto  map[string]int
	Added    int
	Removed  int
}

// Builder computes a Summary from a port snapshot and an optional diff.
type Builder struct{}

// New returns a new Builder.
func New() *Builder { return &Builder{} }

// Build computes a Summary from the given ports and diff.
func (b *Builder) Build(ports []scanner.Port, diff state.Diff) Summary {
	s := Summary{
		Total:   len(ports),
		ByProto: make(map[string]int),
		Added:   len(diff.Added),
		Removed: len(diff.Removed),
	}
	for _, p := range ports {
		s.ByProto[p.Proto]++
	}
	return s
}

// Write renders the Summary as a human-readable table to w.
func (b *Builder) Write(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Total ports:\t%d\n", s.Total)

	protos := make([]string, 0, len(s.ByProto))
	for p := range s.ByProto {
		protos = append(protos, p)
	}
	sort.Strings(protos)
	for _, p := range protos {
		fmt.Fprintf(tw, "  %s:\t%d\n", p, s.ByProto[p])
	}

	fmt.Fprintf(tw, "Added:\t%d\n", s.Added)
	fmt.Fprintf(tw, "Removed:\t%d\n", s.Removed)
	tw.Flush()
}
