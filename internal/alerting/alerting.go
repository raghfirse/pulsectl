// Package alerting provides threshold-based alerting for endpoint health checks.
// It tracks consecutive failures and emits alerts when thresholds are crossed.
package alerting

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

// Alert represents a triggered alert for an endpoint.
type Alert struct {
	Endpoint         string
	ConsecutiveFails int
	LastFailure      time.Time
	Message          string
}

// Alerter tracks consecutive failures per endpoint and writes alerts.
type Alerter struct {
	mu               sync.Mutex
	threshold        int
	consecutiveFails map[string]int
	lastFailure      map[string]time.Time
	w                io.Writer
}

// New creates an Alerter that fires after threshold consecutive failures.
func New(threshold int) *Alerter {
	return NewWithWriter(threshold, os.Stdout)
}

// NewWithWriter creates an Alerter that writes alerts to w.
func NewWithWriter(threshold int, w io.Writer) *Alerter {
	return &Alerter{
		threshold:        threshold,
		consecutiveFails: make(map[string]int),
		lastFailure:      make(map[string]time.Time),
		w:                w,
	}
}

// Evaluate processes a checker result and emits an alert if the failure
// threshold has been reached. It resets the counter on a healthy result.
func (a *Alerter) Evaluate(result checker.Result) *Alert {
	a.mu.Lock()
	defer a.mu.Unlock()

	if result.Healthy {
		delete(a.consecutiveFails, result.Endpoint)
		delete(a.lastFailure, result.Endpoint)
		return nil
	}

	a.consecutiveFails[result.Endpoint]++
	a.lastFailure[result.Endpoint] = result.Timestamp

	fails := a.consecutiveFails[result.Endpoint]
	if fails >= a.threshold {
		alert := &Alert{
			Endpoint:         result.Endpoint,
			ConsecutiveFails: fails,
			LastFailure:      result.Timestamp,
			Message: fmt.Sprintf("ALERT: %s has been DOWN for %d consecutive checks",
				result.Endpoint, fails),
		}
		fmt.Fprintln(a.w, alert.Message)
		return alert
	}
	return nil
}

// ConsecutiveFails returns the current consecutive failure count for an endpoint.
func (a *Alerter) ConsecutiveFails(endpoint string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.consecutiveFails[endpoint]
}
