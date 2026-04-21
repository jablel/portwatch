// Package portсhadow identifies "shadow" ports — ports that appear
// briefly within a configurable observation window and then disappear
// without persisting across multiple scans.
//
// Shadow ports may indicate:
//   - Short-lived services (e.g. ephemeral RPC endpoints)
//   - Port-scanning probes that briefly bind a socket
//   - Misconfigured services that crash on startup
//
// Usage:
//
//	tr := portсhadow.New(5 * time.Second)
//	tr.Observe(currentPorts)
//	for _, s := range tr.Shadows() {
//		fmt.Printf("shadow port detected: %v (seen %d time(s))\n", s.Port, s.Count)
//	}
package portсhadow
