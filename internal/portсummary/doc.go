// Package portсummary aggregates port scan results into a concise
// human-readable summary showing totals by protocol and recent diff counts.
//
// Usage:
//
//	b := portсummary.New()
//	s := b.Build(ports, diff)
//	b.Write(os.Stdout, s)
package portсummary
