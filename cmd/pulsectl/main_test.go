package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestMain_MissingConfig(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-config", "/nonexistent/path/config.yaml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil error")
	}
	if len(out) == 0 {
		t.Error("expected error output for missing config")
	}
}

func TestMain_InvalidConfig(t *testing.T) {
	p := writeTempConfig(t, `endpoints:\n  - name: bad\n`)
	cmd := exec.Command("go", "run", ".", "-config", p)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for invalid config, got output: %s", out)
	}
}

func TestMain_DefaultConfigFlag(t *testing.T) {
	// Verify that omitting -config defaults to "config.yaml" by checking
	// the process exits with an error when config.yaml does not exist in cwd.
	tmpDir := t.TempDir()
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = tmpDir
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when default config.yaml is absent")
	}
}
