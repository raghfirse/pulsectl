package reporter

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

// Summary holds aggregated stats for a single endpoint URL.
type Summary struct {
	URL       string
	Total     int
	Healthy   int
	Unhealthy int
	LastSeen  time.Time
}

// Reporter consumes checker results and prints status to an output writer.
type Reporter struct {
	out      io.Writer
	mu       sync.Mutex
	summaries map[string]*Summary
}

// New creates a Reporter that writes to stdout.
func New() *Reporter {
	return &Reporter{
		out:       os.Stdout,
		summaries: make(map[string]*Summary),
	}
}

// NewWithWriter creates a Reporter that writes to the given writer (useful for testing).
func NewWithWriter(w io.Writer) *Reporter {
	return &Reporter{
		out:       w,
		summaries: make(map[string]*Summary),
	}
}

// Consume reads from the results channel and logs each result until it is closed.
func (r *Reporter) Consume(results <-chan checker.Result) {
	for result := range results {
		r.record(result)
		r.print(result)
	}
}

// Summaries returns a copy of the current aggregated summaries.
func (r *Reporter) Summaries() map[string]Summary {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]Summary, len(r.summaries))
	for k, v := range r.summaries {
		out[k] = *v
	}
	return out
}

func (r *Reporter) record(res checker.Result) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.summaries[res.URL]
	if !ok {
		s = &Summary{URL: res.URL}
		r.summaries[res.URL] = s
	}
	s.Total++
	s.LastSeen = res.Timestamp
	if res.Healthy {
		s.Healthy++
	} else {
		s.Unhealthy++
	}
}

func (r *Reporter) print(res checker.Result) {
	status := "UP"
	if !res.Healthy {
		status = "DOWN"
	}
	fmt.Fprintf(r.out, "[%s] %-6s %s (%s)\n",
		res.Timestamp.Format(time.RFC3339),
		status,
		res.URL,
		res.Duration.Round(time.Millisecond),
	)
}
