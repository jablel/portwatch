// Package portreport provides a summary report of the current port landscape,
// combining classifier, trend, and lifecycle data into a single snapshot.
package portreport

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a single row of the port report.
type Entry struct {
	Port      scanner.Port
	Class     string
	Trend     string
	State     string
	FirstSeen time.Time
	LastSeen  time.Time
}

// Classifier labels a port with a class string (system/registered/dynamic).
type Classifier interface {
	Classify(p scanner.Port) string
}

// Trencher returns the trend direction for a port.
type Trencher interface {
	Trend(p scanner.Port) string
}

// Lifecycler returns the lifecycle state for a port.
type Lifecycler interface {
	State(p scanner.Port) string
}

// Reporter builds and writes port reports.
type Reporter struct {
	classifier Classifier
	trencher   Trencher
	lifecycler Lifecycler
}

// New creates a Reporter with the provided enrichment sources.
func New(c Classifier, t Trencher, l Lifecycler) *Reporter {
	return &Reporter{classifier: c, trencher: t, lifecycler: l}
}

// Build constructs a slice of Entry values from the given ports.
func (r *Reporter) Build(ports []scanner.Port) []Entry {
	entries := make([]Entry, 0, len(ports))
	for _, p := range ports {
		entries = append(entries, Entry{
			Port:  p,
			Class: r.classifier.Classify(p),
			Trend: r.trencher.Trend(p),
			State: r.lifecycler.State(p),
		})
	}
	return entries
}

// Write renders entries as a tab-aligned table to w.
func (r *Reporter) Write(w io.Writer, entries []Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTO\tCLASS\tTREND\tSTATE")
	for _, e := range entries {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
			e.Port.Number, e.Port.Proto, e.Class, e.Trend, e.State)
	}
	return tw.Flush()
}
