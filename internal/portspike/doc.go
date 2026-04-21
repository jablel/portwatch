// Package portspike provides spike detection for open-port counts.
//
// A Detector compares consecutive scan results and fires a Spike event
// when the number of open ports increases by more than a configured
// ratio in a single interval. This is useful for catching sudden
// bursts of newly opened ports that may indicate a misconfiguration or
// an intrusion.
//
// Example:
//
//	det := portspike.New(0.5) // alert on ≥50 % growth
//	for _, ports := range scanResults {
//		if spike := det.Record(ports); spike != nil {
//			log.Println(spike)
//		}
//	}
package portspike
