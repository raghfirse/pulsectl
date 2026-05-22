// Package notifier delivers alert notifications to external systems via
// HTTP webhooks.
//
// When pulsectl detects that a monitored endpoint has exceeded its
// configured failure threshold, the Notifier marshals a structured
// Payload and POSTs it as JSON to the configured webhook URL.
//
// Usage:
//
//	n := notifier.New("https://hooks.example.com/alert")
//	err := n.Notify(notifier.Payload{
//		Endpoint:  "https://api.example.com/health",
//		Status:    "DOWN",
//		Message:   "3 consecutive failures",
//		Timestamp: time.Now(),
//	})
package notifier
