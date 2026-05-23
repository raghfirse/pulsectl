// Package metrics tracks runtime counters for pulsectl operations.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected counters.
type Snapshot struct {
	ChecksTotal   int64
	ChecksHealthy int64
	ChecksDown    int64
	AlertsTotal   int64
	Uptime        time.Duration
	StartedAt     time.Time
}

// Counter is a thread-safe metrics accumulator.
type Counter struct {
	mu            sync.Mutex
	checksTotal   int64
	checksHealthy int64
	checksDown    int64
	alertsTotal   int64
	startedAt     time.Time
	now           func() time.Time
}

// New returns a Counter initialised with the current wall clock.
func New() *Counter {
	return NewWithClock(time.Now)
}

// NewWithClock returns a Counter that uses the supplied clock function.
func NewWithClock(now func() time.Time) *Counter {
	return &Counter{startedAt: now(), now: now}
}

// RecordCheck increments the check counters based on whether the check was healthy.
func (c *Counter) RecordCheck(healthy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checksTotal++
	if healthy {
		c.checksHealthy++
	} else {
		c.checksDown++
	}
}

// RecordAlert increments the alert counter.
func (c *Counter) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns an immutable copy of the current counters.
func (c *Counter) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		ChecksTotal:   c.checksTotal,
		ChecksHealthy: c.checksHealthy,
		ChecksDown:    c.checksDown,
		AlertsTotal:   c.alertsTotal,
		Uptime:        c.now().Sub(c.startedAt),
		StartedAt:     c.startedAt,
	}
}
