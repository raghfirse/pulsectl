package uptime

import (
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

func makeResult(url string, healthy bool) checker.Result {
	return checker.Result{
		URL:     url,
		Healthy: healthy,
	}
}

func TestUptimePercent_NoData(t *testing.T) {
	tr := New(time.Minute)
	if got := tr.UptimePercent("http://example.com"); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestUptimePercent_AllHealthy(t *testing.T) {
	tr := New(time.Minute)
	for i := 0; i < 5; i++ {
		tr.Record(makeResult("http://a.com", true))
	}
	if got := tr.UptimePercent("http://a.com"); got != 100 {
		t.Fatalf("expected 100, got %f", got)
	}
}

func TestUptimePercent_AllDown(t *testing.T) {
	tr := New(time.Minute)
	for i := 0; i < 4; i++ {
		tr.Record(makeResult("http://b.com", false))
	}
	if got := tr.UptimePercent("http://b.com"); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestUptimePercent_Mixed(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(makeResult("http://c.com", true))
	tr.Record(makeResult("http://c.com", true))
	tr.Record(makeResult("http://c.com", false))
	tr.Record(makeResult("http://c.com", false))
	if got := tr.UptimePercent("http://c.com"); got != 50 {
		t.Fatalf("expected 50, got %f", got)
	}
}

func TestUptimePercent_EvictsOldEntries(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }

	tr := NewWithClock(time.Minute, clock)

	// Record two old (unhealthy) entries 2 minutes ago.
	now = now.Add(-2 * time.Minute)
	tr.Record(makeResult("http://d.com", false))
	tr.Record(makeResult("http://d.com", false))

	// Advance clock to present; record one healthy entry.
	now = now.Add(2 * time.Minute)
	tr.Record(makeResult("http://d.com", true))

	// Only the recent healthy entry should remain.
	if got := tr.UptimePercent("http://d.com"); got != 100 {
		t.Fatalf("expected 100 after eviction, got %f", got)
	}
}

func TestUptimePercent_MultipleURLsIsolated(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(makeResult("http://x.com", true))
	tr.Record(makeResult("http://y.com", false))

	if got := tr.UptimePercent("http://x.com"); got != 100 {
		t.Fatalf("x: expected 100, got %f", got)
	}
	if got := tr.UptimePercent("http://y.com"); got != 0 {
		t.Fatalf("y: expected 0, got %f", got)
	}
}
