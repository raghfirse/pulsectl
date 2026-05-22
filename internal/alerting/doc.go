// Package alerting implements threshold-based alerting for pulsectl.
//
// An Alerter tracks consecutive failures for each monitored endpoint.
// When the number of consecutive failures reaches a configured threshold,
// an alert is emitted to the configured writer (default: stdout).
//
// The failure counter for an endpoint is reset to zero whenever a healthy
// result is received, ensuring alerts only fire on sustained outages.
//
// Example usage:
//
//	a := alerting.New(3) // alert after 3 consecutive failures
//	for result := range results {
//		a.Evaluate(result)
//	}
package alerting
