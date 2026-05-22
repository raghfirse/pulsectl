package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/notifier"
)

func makePayload(endpoint, status, message string) notifier.Payload {
	return notifier.Payload{
		Endpoint:  endpoint,
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func TestNotifier_SendsCorrectPayload(t *testing.T) {
	var received notifier.Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.NewWithClient(ts.URL, ts.Client())
	p := makePayload("https://example.com", "DOWN", "threshold reached")

	if err := n.Notify(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Endpoint != p.Endpoint {
		t.Errorf("endpoint: got %q, want %q", received.Endpoint, p.Endpoint)
	}
	if received.Status != p.Status {
		t.Errorf("status: got %q, want %q", received.Status, p.Status)
	}
}

func TestNotifier_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notifier.NewWithClient(ts.URL, ts.Client())
	err := n.Notify(makePayload("https://example.com", "DOWN", "err"))
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestNotifier_UnreachableServerReturnsError(t *testing.T) {
	n := notifier.New("http://127.0.0.1:0/webhook")
	err := n.Notify(makePayload("https://example.com", "DOWN", "err"))
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}

func TestNotifier_ContentTypeIsJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type: got %q, want %q", ct, "application/json")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := notifier.NewWithClient(ts.URL, ts.Client())
	if err := n.Notify(makePayload("https://example.com", "UP", "ok")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
