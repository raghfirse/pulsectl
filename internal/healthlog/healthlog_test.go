package healthlog_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/healthlog"
)

func makeResult(url string, healthy bool, latencyMs int, err error) checker.Result {
	return checker.Result{
		URL:       url,
		Healthy:   healthy,
		Latency:   time.Duration(latencyMs) * time.Millisecond,
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Err:       err,
	}
}

func TestLogger_HealthyResult(t *testing.T) {
	var buf bytes.Buffer
	l := healthlog.NewWithWriter(&buf)
	l.Log(makeResult("https://example.com", true, 42, nil))

	line := buf.String()
	if !strings.Contains(line, "UP") {
		t.Errorf("expected UP in line, got: %q", line)
	}
	if !strings.Contains(line, "42ms") {
		t.Errorf("expected 42ms in line, got: %q", line)
	}
	if !strings.Contains(line, "https://example.com") {
		t.Errorf("expected URL in line, got: %q", line)
	}
}

func TestLogger_DownResult(t *testing.T) {
	var buf bytes.Buffer
	l := healthlog.NewWithWriter(&buf)
	l.Log(makeResult("https://fail.io", false, 200, errors.New("connection refused")))

	line := buf.String()
	if !strings.Contains(line, "DOWN") {
		t.Errorf("expected DOWN in line, got: %q", line)
	}
	if !strings.Contains(line, "connection refused") {
		t.Errorf("expected error message in line, got: %q", line)
	}
}

func TestLogger_TimestampFormat(t *testing.T) {
	var buf bytes.Buffer
	l := healthlog.NewWithWriter(&buf)
	l.Log(makeResult("https://example.com", true, 10, nil))

	if !strings.HasPrefix(buf.String(), "2024-06-01T12:00:00Z") {
		t.Errorf("unexpected timestamp format: %q", buf.String())
	}
}

func TestLogger_MultipleResults(t *testing.T) {
	var buf bytes.Buffer
	l := healthlog.NewWithWriter(&buf)
	l.Log(makeResult("https://example.com", true, 10, nil))
	l.Log(makeResult("https://fail.io", false, 300, errors.New("timeout")))

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "UP") {
		t.Errorf("expected first line to contain UP, got: %q", lines[0])
	}
	if !strings.Contains(lines[1], "DOWN") {
		t.Errorf("expected second line to contain DOWN, got: %q", lines[1])
	}
}

func TestNew_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "health.log")

	l, err := healthlog.New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	l.Log(makeResult("https://example.com", true, 5, nil))

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "UP") {
		t.Errorf("expected UP in file content, got: %q", string(data))
	}
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := healthlog.New("/nonexistent/dir/health.log")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
