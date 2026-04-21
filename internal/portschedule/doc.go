// Package portschedule provides time-of-day schedule learning and enforcement
// for monitored ports.
//
// A Tracker observes which hours of the day each port is normally active.
// After a configurable number of observations the schedule is considered
// learned, and any subsequent observation of a port outside its known active
// hours produces a Violation.
//
// Typical usage:
//
//	tr := portschedule.New(5) // enforce after 5 observations
//	violations := tr.Observe(ports, time.Now())
//	for _, v := range violations {
//		log.Println(v)
//	}
package portschedule
