package metrics

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Handler returns an http.HandlerFunc that writes a plain-text metrics
// snapshot to the response. It is intended for lightweight status pages.
func Handler(c *Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteTo(c, w)
	}
}

// WriteTo formats the current snapshot as plain text and writes it to w.
func WriteTo(c *Counter, w io.Writer) {
	s := c.Snapshot()
	fmt.Fprintf(w, "started_at:      %s\n", s.StartedAt.UTC().Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(w, "uptime_seconds:  %.0f\n", s.Uptime.Seconds())
	fmt.Fprintf(w, "checks_total:    %d\n", s.ChecksTotal)
	fmt.Fprintf(w, "checks_healthy:  %d\n", s.ChecksHealthy)
	fmt.Fprintf(w, "checks_down:     %d\n", s.ChecksDown)
	fmt.Fprintf(w, "alerts_total:    %d\n", s.AlertsTotal)
}

// Print writes the current snapshot to stdout.
func Print(c *Counter) {
	WriteTo(c, os.Stdout)
}
