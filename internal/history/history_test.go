package history_test

import (
	"testing"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/history"
)

func makeResult(url string, healthy bool) checker.Result {
	return checker.Result{URL: url, Healthy: healthy, StatusCode: 200}
}

func TestStore_AddAndGet(t *testing.T) {
	s := history.New(10)
	s.Add(makeResult("http://example.com", true))
	s.Add(makeResult("http://example.com", false))

	records := s.Get("http://example.com")
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestStore_MaxSizeEviction(t *testing.T) {
	s := history.New(3)
	for i := 0; i < 5; i++ {
		s.Add(makeResult("http://example.com", true))
	}

	records := s.Get("http://example.com")
	if len(records) != 3 {
		t.Fatalf("expected 3 records after eviction, got %d", len(records))
	}
}

func TestStore_UptimePercent_AllHealthy(t *testing.T) {
	s := history.New(10)
	for i := 0; i < 4; i++ {
		s.Add(makeResult("http://a.com", true))
	}

	pct := s.UptimePercent("http://a.com")
	if pct != 100.0 {
		t.Errorf("expected 100%%, got %.2f", pct)
	}
}

func TestStore_UptimePercent_Mixed(t *testing.T) {
	s := history.New(10)
	s.Add(makeResult("http://b.com", true))
	s.Add(makeResult("http://b.com", true))
	s.Add(makeResult("http://b.com", false))
	s.Add(makeResult("http://b.com", false))

	pct := s.UptimePercent("http://b.com")
	if pct != 50.0 {
		t.Errorf("expected 50%%, got %.2f", pct)
	}
}

func TestStore_UptimePercent_NoRecords(t *testing.T) {
	s := history.New(10)
	pct := s.UptimePercent("http://missing.com")
	if pct != 0 {
		t.Errorf("expected 0 for unknown URL, got %.2f", pct)
	}
}

func TestStore_URLs(t *testing.T) {
	s := history.New(10)
	s.Add(makeResult("http://one.com", true))
	s.Add(makeResult("http://two.com", false))

	urls := s.URLs()
	if len(urls) != 2 {
		t.Fatalf("expected 2 URLs, got %d", len(urls))
	}
}

func TestStore_GetReturnsCopy(t *testing.T) {
	s := history.New(10)
	s.Add(makeResult("http://example.com", true))

	records := s.Get("http://example.com")
	records[0].Result.Healthy = false

	original := s.Get("http://example.com")
	if !original[0].Result.Healthy {
		t.Error("Get should return a copy, not a reference")
	}
}
