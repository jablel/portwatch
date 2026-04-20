// Package portanomaly detects statistically anomalous port activity.
//
// A Detector maintains two rolling windows of per-port frequency samples:
// a baseline window (older observations) and a recent window (newer
// observations). After each call to Record the two averages are compared;
// ports whose deviation exceeds the configured threshold are surfaced via
// Anomalies.
//
// Typical usage:
//
//	det := portanomaly.New(0.4, 10) // 40 % deviation, 10-scan windows
//	for _, snapshot := range scans {
//		det.Record(snapshot)
//	}
//	for _, a := range det.Anomalies(currentPorts) {
//		log.Println(a)
//	}
package portanomaly
