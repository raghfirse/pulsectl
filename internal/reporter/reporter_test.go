package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/reporter"
)

func makeResult(url string, healthy bool) checker.Result {
	return checker.Result{
		URL:       url,
		Healthy:   healthy,
		Status:    200,
		Duration:  42 * time.Millisecond,
		Timestamp: time.Now(),
	}
}

func TestReporter_PrintsUP(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(&buf)

	ch := make(chan checker.Result, 1)
	ch <- makeResult("http://example.com", true)
	close(ch)
	r.Consume(ch)

	out := buf.String()
	if !strings.Contains(out, "UP") {
		t.Errorf("expected UP in output, got: %s", out)
	}
	if !strings.Contains(out, "http://example.com") {
		t.Errorf("expected URL in output, got: %s", out)
	}
}

func TestReporter_PrintsDOWN(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(&buf)

	ch := make(chan checker.Result, 1)
	ch <- makeResult("http://down.example.com", false)
	close(ch)
	r.Consume(ch)

	out := buf.String()
	if !strings.Contains(out, "DOWN") {
		t.Errorf("expected DOWN in output, got: %s", out)
	}
}

func TestReporter_Summaries(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(&buf)

	ch := make(chan checker.Result, 4)
	ch <- makeResult("http://a.com", true)
	ch <- makeResult("http://a.com", true)
	ch <- makeResult("http://a.com", false)
	ch <- makeResult("http://b.com", true)
	close(ch)
	r.Consume(ch)

	summaries := r.Summaries()

	a := summaries["http://a.com"]
	if a.Total != 3 {
		t.Errorf("expected Total=3, got %d", a.Total)
	}
	if a.Healthy != 2 {
		t.Errorf("expected Healthy=2, got %d", a.Healthy)
	}
	if a.Unhealthy != 1 {
		t.Errorf("expected Unhealthy=1, got %d", a.Unhealthy)
	}

	b := summaries["http://b.com"]
	if b.Total != 1 {
		t.Errorf("expected Total=1 for b.com, got %d", b.Total)
	}
}
