// Package notifier provides webhook-based notification delivery
// for pulsectl health check alerts.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the JSON body sent to a webhook endpoint.
type Payload struct {
	Endpoint  string    `json:"endpoint"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Notifier sends alert payloads to a configured webhook URL.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a Notifier that posts to the given webhook URL.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

// NewWithClient creates a Notifier with a custom HTTP client (useful for testing).
func NewWithClient(webhookURL string, client *http.Client) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     client,
	}
}

// Notify sends a Payload to the configured webhook URL.
// It returns an error if marshalling or the HTTP request fails, or if
// the server responds with a non-2xx status code.
func (n *Notifier) Notify(p Payload) error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notifier: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("notifier: post to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
