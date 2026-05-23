// Package healthlog provides structured logging of health check events
// to a rotating log file or any io.Writer.
package healthlog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

// Logger writes health check results as structured log lines.
type Logger struct {
	mu sync.Mutex
	w  io.Writer
}

// New creates a Logger that writes to the given file path, creating it if
// necessary. The file is opened in append mode so existing entries are
// preserved across restarts.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("healthlog: open %q: %w", path, err)
	}
	return NewWithWriter(f), nil
}

// NewWithWriter creates a Logger that writes to w. Useful for testing.
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Log writes a single result as a tab-separated log line:
//
//	<RFC3339 timestamp>\t<url>\t<status>\t<latency_ms>ms\t[<error>]
func (l *Logger) Log(r checker.Result) {
	l.mu.Lock()
	defer l.mu.Unlock()

	status := "UP"
	if !r.Healthy {
		status = "DOWN"
	}

	errStr := ""
	if r.Err != nil {
		errStr = "\t" + r.Err.Error()
	}

	fmt.Fprintf(
		l.w,
		"%s\t%s\t%s\t%dms%s\n",
		r.Timestamp.UTC().Format(time.RFC3339),
		r.URL,
		status,
		r.Latency.Milliseconds(),
		errStr,
	)
}
