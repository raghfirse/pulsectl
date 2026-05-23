package metrics

import (
	"testing"
	"time"
)

func TestCounter_InitialSnapshot(t *testing.T) {
	c := New()
	s := c.Snapshot()
	if s.ChecksTotal != 0 || s.ChecksHealthy != 0 || s.ChecksDown != 0 || s.AlertsTotal != 0 {
		t.Errorf("expected zero counters, got %+v", s)
	}
}

func TestCounter_RecordHealthy(t *testing.T) {
	c := New()
	c.RecordCheck(true)
	c.RecordCheck(true)
	s := c.Snapshot()
	if s.ChecksTotal != 2 || s.ChecksHealthy != 2 || s.ChecksDown != 0 {
		t.Errorf("unexpected snapshot: %+v", s)
	}
}

func TestCounter_RecordDown(t *testing.T) {
	c := New()
	c.RecordCheck(false)
	s := c.Snapshot()
	if s.ChecksTotal != 1 || s.ChecksDown != 1 || s.ChecksHealthy != 0 {
		t.Errorf("unexpected snapshot: %+v", s)
	}
}

func TestCounter_Mixed(t *testing.T) {
	c := New()
	for i := 0; i < 3; i++ {
		c.RecordCheck(true)
	}
	for i := 0; i < 2; i++ {
		c.RecordCheck(false)
	}
	s := c.Snapshot()
	if s.ChecksTotal != 5 || s.ChecksHealthy != 3 || s.ChecksDown != 2 {
		t.Errorf("unexpected snapshot: %+v", s)
	}
}

func TestCounter_RecordAlert(t *testing.T) {
	c := New()
	c.RecordAlert()
	c.RecordAlert()
	if s := c.Snapshot(); s.AlertsTotal != 2 {
		t.Errorf("expected 2 alerts, got %d", s.AlertsTotal)
	}
}

func TestCounter_Uptime(t *testing.T) {
	now := time.Now()
	fakeClock := func() time.Time { return now }
	c := NewWithClock(fakeClock)

	// advance clock by 5 seconds
	now = now.Add(5 * time.Second)
	s := c.Snapshot()
	if s.Uptime < 5*time.Second {
		t.Errorf("expected uptime >= 5s, got %s", s.Uptime)
	}
}

func TestCounter_StartedAt(t *testing.T) {
	start := time.Now()
	c := NewWithClock(func() time.Time { return start })
	if !c.Snapshot().StartedAt.Equal(start) {
		t.Error("StartedAt mismatch")
	}
}
