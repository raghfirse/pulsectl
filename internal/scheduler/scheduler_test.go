package scheduler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/config"
	"github.com/user/pulsectl/internal/scheduler"
)

func newTestConfig(urls []string, intervalSec int) *config.Config {
	endpoints := make([]config.Endpoint, len(urls))
	for i, u := range urls {
		endpoints[i] = config.Endpoint{URL: u, IntervalSeconds: intervalSec}
	}
	return &config.Config{
		DefaultIntervalSeconds: intervalSec,
		DefaultTimeoutSeconds:  2,
		Endpoints:              endpoints,
	}
}

func TestScheduler_ReceivesResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := newTestConfig([]string{server.URL}, 1)
	chk := checker.New(time.Duration(cfg.DefaultTimeoutSeconds) * time.Second)
	sched := scheduler.New(cfg, chk)

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	go sched.Start(ctx)

	var received int
	for r := range sched.Results() {
		if r.URL != server.URL {
			t.Errorf("unexpected URL: got %s", r.URL)
		}
		if !r.Healthy {
			t.Errorf("expected healthy result")
		}
		received++
	}

	if received < 1 {
		t.Errorf("expected at least 1 result, got %d", received)
	}
}

func TestScheduler_MultipleEndpoints(t *testing.T) {
	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverA.Close()

	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer serverB.Close()

	cfg := newTestConfig([]string{serverA.URL, serverB.URL}, 1)
	chk := checker.New(time.Duration(cfg.DefaultTimeoutSeconds) * time.Second)
	sched := scheduler.New(cfg, chk)

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	go sched.Start(ctx)

	seen := map[string]int{}
	for r := range sched.Results() {
		seen[r.URL]++
	}

	if seen[serverA.URL] < 1 {
		t.Errorf("expected results for serverA")
	}
	if seen[serverB.URL] < 1 {
		t.Errorf("expected results for serverB")
	}
}
