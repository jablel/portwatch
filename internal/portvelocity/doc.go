// Package portvelocity measures the rate of change between consecutive port
// scans, expressed as a normalised score in the range [0, 1].
//
// A score of 0 means no ports changed between the two most recent scans.
// A score of 1 means every port observed across both scans was either newly
// added or removed — i.e. complete turnover.
//
// Typical usage:
//
//	tr := portvelocity.New()
//
//	for {
//		ports := scanner.Scan(...)
//		v := tr.Record(ports)
//		if v > 0.5 {
//			log.Printf("high port churn detected: velocity=%.2f", v)
//		}
//	}
package portvelocity
