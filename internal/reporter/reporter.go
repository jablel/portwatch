package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Format controls the output format of the reporter.
type Format string

const (
	FormatText Format = "text"
	FormatCSV  Format = "csv"
)

// Reporter writes port change reports to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to out in the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{out: out, format: format}
}

// Report writes a summary of the diff to the reporter's output.
func (r *Reporter) Report(diff state.Diff) error {
	switch r.format {
	case FormatCSV:
		return r.writeCSV(diff)
	default:
		return r.writeText(diff)
	}
}

func (r *Reporter) writeText(diff state.Diff) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	ts := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(w, "Port Change Report\t%s\n", ts)
	fmt.Fprintln(w, "---")
	for _, p := range diff.Added {
		fmt.Fprintf(w, "ADDED\t%s\n", p)
	}
	for _, p := range diff.Removed {
		fmt.Fprintf(w, "REMOVED\t%s\n", p)
	}
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		fmt.Fprintln(w, "No changes detected.")
	}
	return w.Flush()
}

func (r *Reporter) writeCSV(diff state.Diff) error {
	ts := time.Now().UTC().Format(time.RFC3339)
	for _, p := range diff.Added {
		fmt.Fprintf(r.out, "%s,ADDED,%s\n", ts, p)
	}
	for _, p := range diff.Removed {
		fmt.Fprintf(r.out, "%s,REMOVED,%s\n", ts, p)
	}
	return nil
}
