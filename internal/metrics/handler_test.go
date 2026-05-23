package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWriteTo_ContainsAllFields(t *testing.T) {
	c := New()
	c.RecordCheck(true)
	c.RecordCheck(false)
	c.RecordAlert()

	var sb strings.Builder
	WriteTo(c, &sb)
	out := sb.String()

	for _, want := range []string{
		"started_at:",
		"uptime_seconds:",
		"checks_total:    2",
		"checks_healthy:  1",
		"checks_down:     1",
		"alerts_total:    1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, out)
		}
	}
}

func TestHandler_Returns200(t *testing.T) {
	c := New()
	h := Handler(c)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_BodyNotEmpty(t *testing.T) {
	c := NewWithClock(func() time.Time { return time.Now() })
	h := Handler(c)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}
