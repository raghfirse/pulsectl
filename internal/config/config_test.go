package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "pulsectl-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
endpoints:
  - name: example
    url: https://example.com/health
    interval: 10s
    timeout: 3s
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(cfg.Endpoints))
	}
	ep := cfg.Endpoints[0]
	if ep.Name != "example" {
		t.Errorf("expected name 'example', got %q", ep.Name)
	}
	if ep.Interval != 10*time.Second {
		t.Errorf("expected interval 10s, got %v", ep.Interval)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTempConfig(t, `
endpoints:
  - name: svc
    url: http://localhost:8080/ping
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ep := cfg.Endpoints[0]
	if ep.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", ep.Interval)
	}
	if ep.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", ep.Timeout)
	}
}

func TestLoad_MissingURL(t *testing.T) {
	path := writeTempConfig(t, `
endpoints:
  - name: broken
`)
	if _, err := Load(path); err == nil {
		t.Fatal("expected error for missing URL, got nil")
	}
}

func TestLoad_EmptyEndpoints(t *testing.T) {
	path := writeTempConfig(t, `endpoints: []`)
	if _, err := Load(path); err == nil {
		t.Fatal("expected error for empty endpoints, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	if _, err := Load("/nonexistent/path/config.yaml"); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
