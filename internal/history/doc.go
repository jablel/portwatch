// Package history provides persistent, bounded storage of port-scan snapshots
// over time. Each snapshot (Entry) captures a UTC timestamp alongside the list
// of open ports observed during a single scan cycle.
//
// Usage:
//
//	h, err := history.Load("/var/lib/portwatch/history.json", 200)
//	if err != nil {
//		log.Fatal(err)
//	}
//	h.Add(ports)
//	if err := h.Save("/var/lib/portwatch/history.json"); err != nil {
//		log.Println("warning: could not save history:", err)
//	}
//
// Query helpers (Since, PortSeen, UniquePortsInRange) allow callers to
// interrogate the retained entries without deserialising the file again.
package history
