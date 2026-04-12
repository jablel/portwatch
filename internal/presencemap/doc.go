// Package presencemap tracks consecutive-scan presence streaks for open ports.
//
// After each scan cycle, call Observe with the current port list. The map
// increments a streak counter for every port that remains open and removes
// entries for ports that have disappeared.
//
// Use Stable to retrieve ports that have been continuously present for at
// least N consecutive scans — useful for suppressing alerts on transient
// ports and focusing attention on newly stable listeners.
//
// Example:
//
//	pm := presencemap.New()
//	pm.Observe(ports)
//	stable := pm.Stable(3) // ports seen in 3+ consecutive scans
package presencemap
