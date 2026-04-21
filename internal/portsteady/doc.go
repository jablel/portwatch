// Package portsteady measures the long-term stability of open ports.
//
// A Tracker observes successive port scan snapshots and computes a stability
// score in the range [0, 1] for each port:
//
//	1.0 — the port appeared in every scan within the configured window
//	0.0 — the port was never seen, or has disappeared
//
// Ports that vanish from a snapshot are immediately evicted so that the
// tracker never holds stale data.
//
// Typical usage:
//
//	tr := portsteady.New(10)          // 10-scan rolling window
//	tr.Observe(ports, time.Now())
//	score := tr.Stability(port)       // 0.0 – 1.0
//	uptime := tr.Uptime(port, time.Now())
package portsteady
