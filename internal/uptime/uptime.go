// Package uptime provides a simple uptime tracker that computes
// rolling uptime percentages for monitored endpoints over a sliding
// time window.
package uptime

import (
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

// Clock allows time to be injected for testing.
type Clock func() time.Time

// entry records a single check outcome and when it occurred.
type entry struct {
	at      time.Time
	health  bool
}

// Tracker maintains a sliding-window uptime record per endpoint.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	clock   Clock
	records map[string][]entry
}

// New creates a Tracker with the given sliding window duration.
func New(window time.Duration) *Tracker {
	return NewWithClock(window, time.Now)
}

// NewWithClock creates a Tracker with a custom clock (useful for tests).
func NewWithClock(window time.Duration, clock Clock) *Tracker {
	return &Tracker{
		window:  window,
		clock:   clock,
		records: make(map[string][]entry),
	}
}

// Record adds a checker.Result to the tracker.
func (t *Tracker) Record(r checker.Result) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	t.records[r.URL] = append(t.records[r.URL], entry{
		at:     now,
		health: r.Healthy,
	})
	t.evict(r.URL, now)
}

// UptimePercent returns the percentage of healthy checks within the
// window for the given URL. Returns 0 if no data is available.
func (t *Tracker) UptimePercent(url string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	t.evict(url, now)

	entries := t.records[url]
	if len(entries) == 0 {
		return 0
	}

	healthy := 0
	for _, e := range entries {
		if e.health {
			healthy++
		}
	}
	return float64(healthy) / float64(len(entries)) * 100
}

// evict removes entries older than the window. Must be called with mu held.
func (t *Tracker) evict(url string, now time.Time) {
	cutoff := now.Add(-t.window)
	entries := t.records[url]
	i := 0
	for i < len(entries) && entries[i].at.Before(cutoff) {
		i++
	}
	t.records[url] = entries[i:]
}
