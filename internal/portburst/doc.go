// Package portburst provides a sliding-window burst detector for port events.
//
// A Detector accumulates port-added events and compares the count within the
// configured window against a threshold. When the threshold is exceeded a
// Burst value is returned so callers can raise an alert or trigger a
// higher-frequency scan cycle.
//
// Usage:
//
//	det := portburst.New(10*time.Second, 5)
//	if b := det.Record(diff.Added); b != nil {
//		log.Printf("burst: %d new ports in window", b.Count)
//	}
package portburst
