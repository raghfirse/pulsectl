package alerting_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/alerting"
	"github.com/user/pulsectl/internal/checker"
)

func makeResult(endpoint string, healthy bool) checker.Result {
	return checker.Result{
		Endpoint:  endpoint,
		Healthy:   healthy,
		Timestamp: time.Now(),
	}
}

func TestAlerter_NoAlertBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alerting.NewWithWriter(3, &buf)

	for i := 0; i < 2; i++ {
		alert := a.Evaluate(makeResult("http://example.com", false))
		if alert != nil {
			t.Errorf("expected no alert before threshold, got %v", alert)
		}
	}
	if buf.Len() > 0 {
		t.Errorf("expected no output before threshold, got: %s", buf.String())
	}
}

func TestAlerter_AlertAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alerting.NewWithWriter(3, &buf)
	endpoint := "http://example.com"

	for i := 0; i < 3; i++ {
		a.Evaluate(makeResult(endpoint, false))
	}

	if !strings.Contains(buf.String(), "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), endpoint) {
		t.Errorf("expected endpoint in output, got: %s", buf.String())
	}
}

func TestAlerter_AlertReturnsStruct(t *testing.T) {
	var buf bytes.Buffer
	a := alerting.NewWithWriter(2, &buf)
	endpoint := "http://api.example.com"

	a.Evaluate(makeResult(endpoint, false))
	alert := a.Evaluate(makeResult(endpoint, false))

	if alert == nil {
		t.Fatal("expected alert struct, got nil")
	}
	if alert.Endpoint != endpoint {
		t.Errorf("expected endpoint %s, got %s", endpoint, alert.Endpoint)
	}
	if alert.ConsecutiveFails != 2 {
		t.Errorf("expected 2 consecutive fails, got %d", alert.ConsecutiveFails)
	}
}

func TestAlerter_ResetsOnHealthy(t *testing.T) {
	var buf bytes.Buffer
	a := alerting.NewWithWriter(3, &buf)
	endpoint := "http://example.com"

	a.Evaluate(makeResult(endpoint, false))
	a.Evaluate(makeResult(endpoint, false))
	a.Evaluate(makeResult(endpoint, true)) // recovery

	if a.ConsecutiveFails(endpoint) != 0 {
		t.Errorf("expected 0 consecutive fails after recovery, got %d", a.ConsecutiveFails(endpoint))
	}
}

func TestAlerter_IndependentEndpoints(t *testing.T) {
	var buf bytes.Buffer
	a := alerting.NewWithWriter(2, &buf)

	a.Evaluate(makeResult("http://a.com", false))
	a.Evaluate(makeResult("http://b.com", false))

	if a.ConsecutiveFails("http://a.com") != 1 {
		t.Errorf("expected 1 fail for a.com")
	}
	if a.ConsecutiveFails("http://b.com") != 1 {
		t.Errorf("expected 1 fail for b.com")
	}
}
