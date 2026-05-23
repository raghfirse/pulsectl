package pipeline_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/config"
	"github.com/user/pulsectl/internal/pipeline"
)

func makeConfig(webhookURL string, threshold int) *config.Config {
	return &config.Config{
		Alerting: config.AlertingConfig{
			Threshold: threshold,
			RateLimit: config.RateLimitConfig{Rate: 10, Burst: 5},
		},
		History:  config.HistoryConfig{MaxSize: 100},
		Webhook:  config.WebhookConfig{URL: webhookURL},
	}
}

func makeResult(url string, healthy bool) checker.Result {
	status := "UP"
	if !healthy {
		status = "DOWN"
	}
	return checker.Result{
		URL:       url,
		Status:    status,
		Healthy:   healthy,
		Timestamp: time.Now(),
	}
}

func TestPipeline_NoWebhookNoError(t *testing.T) {
	cfg := makeConfig("", 1)
	p := pipeline.New(cfg)
	// Should not panic or error even without a webhook configured.
	p.Process(makeResult("http://example.com", false))
}

func TestPipeline_WebhookFiredOnAlert(t *testing.T) {
	received := make(chan map[string]interface{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := makeConfig(server.URL, 1)
	p := pipeline.New(cfg)
	p.Process(makeResult("http://example.com", false))

	select {
	case payload := <-received:
		if payload["endpoint"] == nil {
			t.Error("expected endpoint in webhook payload")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("webhook was not called within timeout")
	}
}

func TestPipeline_HealthyResultNoWebhook(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := makeConfig(server.URL, 1)
	p := pipeline.New(cfg)
	p.Process(makeResult("http://example.com", true))

	time.Sleep(100 * time.Millisecond)
	if called {
		t.Error("webhook should not be called for a healthy result")
	}
}
