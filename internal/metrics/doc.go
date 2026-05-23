// Package metrics provides a lightweight, thread-safe counter for tracking
// pulsectl runtime statistics such as the total number of health checks
// performed, how many were healthy or down, how many alerts were fired,
// and how long the process has been running.
//
// Usage:
//
//	c := metrics.New()
//	c.RecordCheck(true)   // healthy
//	c.RecordCheck(false)  // down
//	c.RecordAlert()
//	snap := c.Snapshot()  // thread-safe point-in-time copy
package metrics
