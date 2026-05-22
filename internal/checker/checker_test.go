package checker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

func TestCheck_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := checker.New(5 * time.Second)
	res := c.Check(context.Background(), "test-ok", ts.URL)

	if !res.Healthy() {
		t.Fatalf("expected healthy result, got: %s", res)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
	if res.Err != nil {
		t.Errorf("expected no error, got %v", res.Err)
	}
}

func TestCheck_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	c := checker.New(5 * time.Second)
	res := c.Check(context.Background(), "test-503", ts.URL)

	if res.Healthy() {
		t.Fatal("expected unhealthy result for 503")
	}
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", res.StatusCode)
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	c := checker.New(2 * time.Second)
	res := c.Check(context.Background(), "bad-url", "://not-a-url")

	if res.Err == nil {
		t.Fatal("expected an error for invalid URL, got nil")
	}
	if res.Healthy() {
		t.Fatal("expected unhealthy result for invalid URL")
	}
}

func TestCheck_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := checker.New(50 * time.Millisecond)
	res := c.Check(context.Background(), "timeout-test", ts.URL)

	if res.Err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if res.Healthy() {
		t.Fatal("expected unhealthy result on timeout")
	}
}

func TestResult_String(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := checker.New(5 * time.Second)
	res := c.Check(context.Background(), "str-test", ts.URL)

	s := res.String()
	if len(s) == 0 {
		t.Error("expected non-empty string representation")
	}
}
