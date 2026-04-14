// Package portreport assembles a human-readable summary of the active port
// landscape by combining classification, trend, and lifecycle information.
//
// Usage:
//
//	r := portreport.New(classifier, trencher, lifecycler)
//	entries := r.Build(ports)
//	r.Write(os.Stdout, entries)
//
// Each Entry in the report contains the port number, protocol, class label
// (system / registered / dynamic), trend direction (rising / falling / stable),
// and lifecycle state (new / active / closed).
package portreport
