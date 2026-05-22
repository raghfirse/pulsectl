package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/history"
	"github.com/user/pulsectl/internal/reporter"
)

func addResult(s *history.Store, url string, healthy bool) {
	s.Add(checker.Result{URL: url, Healthy: healthy, StatusCode: 200})
}

func TestPrintHistorySummary_Empty(t *testing.T) {
	store := history.New(10)
	var buf bytes.Buffer
	reporter.PrintHistorySummary(&buf, store)

	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected 'No history' message, got: %s", buf.String())
	}
}

func TestPrintHistorySummary_SingleEndpoint(t *testing.T) {
	store := history.New(10)
	addResult(store, "http://example.com", true)
	addResult(store, "http://example.com", true)
	addResult(store, "http://example.com", false)

	var buf bytes.Buffer
	reporter.PrintHistorySummary(&buf, store)
	out := buf.String()

	if !strings.Contains(out, "http://example.com") {
		t.Error("expected URL in output")
	}
	if !strings.Contains(out, "3") {
		t.Error("expected check count 3 in output")
	}
	if !strings.Contains(out, "66.7%") {
		t.Errorf("expected 66.7%% uptime, got: %s", out)
	}
}

func TestPrintHistorySummary_MultipleEndpoints(t *testing.T) {
	store := history.New(10)
	addResult(store, "http://alpha.com", true)
	addResult(store, "http://beta.com", false)

	var buf bytes.Buffer
	reporter.PrintHistorySummary(&buf, store)
	out := buf.String()

	if !strings.Contains(out, "http://alpha.com") {
		t.Error("expected alpha.com in output")
	}
	if !strings.Contains(out, "http://beta.com") {
		t.Error("expected beta.com in output")
	}
	if !strings.Contains(out, "100.0%") {
		t.Errorf("expected 100.0%% for alpha, output: %s", out)
	}
	if !strings.Contains(out, "0.0%") {
		t.Errorf("expected 0.0%% for beta, output: %s", out)
	}
}

func TestPrintHistorySummary_HeaderPresent(t *testing.T) {
	store := history.New(10)
	addResult(store, "http://x.com", true)

	var buf bytes.Buffer
	reporter.PrintHistorySummary(&buf, store)
	out := buf.String()

	if !strings.Contains(out, "ENDPOINT") || !strings.Contains(out, "UPTIME") {
		t.Errorf("expected table headers in output: %s", out)
	}
}
