package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Result holds the outcome of a single health-check poll.
type Result struct {
	Name       string
	URL        string
	StatusCode int
	Latency    time.Duration
	Err        error
	CheckedAt  time.Time
}

// Healthy returns true when the check succeeded with a 2xx status code.
func (r Result) Healthy() bool {
	return r.Err == nil && r.StatusCode >= 200 && r.StatusCode < 300
}

// String returns a human-readable summary of the result.
func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("[%s] %s — ERROR: %v", r.Name, r.URL, r.Err)
	}
	return fmt.Sprintf("[%s] %s — %d (%s)", r.Name, r.URL, r.StatusCode, r.Latency.Round(time.Millisecond))
}

// Checker performs HTTP health checks.
type Checker struct {
	client *http.Client
}

// New creates a Checker with the given request timeout.
func New(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

// Check sends a GET request to the given URL and returns a Result.
func (c *Checker) Check(ctx context.Context, name, url string) Result {
	result := Result{
		Name:      name,
		URL:       url,
		CheckedAt: time.Now(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Err = fmt.Errorf("building request: %w", err)
		return result
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Err = fmt.Errorf("executing request: %w", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	return result
}
