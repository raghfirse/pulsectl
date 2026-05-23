package statuspage_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/history"
	"github.com/user/pulsectl/internal/statuspage"
)

// fakeStore implements statuspage.Store for testing.
type fakeStore struct {
	urls    []string
	results map[string][]history.Result
	uptime  map[string]float64
}

func (f *fakeStore) URLs() []string { return f.urls }

func (f *fakeStore) Get(url string) []history.Result { return f.results[url] }

func (f *fakeStore) UptimePercent(url string) float64 { return f.uptime[url] }

func newFakeStore() *fakeStore {
	return &fakeStore{
		results: make(map[string][]history.Result),
		uptime:  make(map[string]float64),
	}
}

func makeResult(url string, healthy bool) history.Result {
	return history.Result{
		URL:       url,
		Healthy:   healthy,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestHandler_Returns200(t *testing.T) {
	store := newFakeStore()
	h := statuspage.Handler(store)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_ContentTypeHTML(t *testing.T) {
	store := newFakeStore()
	h := statuspage.Handler(store)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html content-type, got %q", ct)
	}
}

func TestHandler_ShowsUpEndpoint(t *testing.T) {
	store := newFakeStore()
	store.urls = []string{"https://example.com"}
	store.results["https://example.com"] = []history.Result{makeResult("https://example.com", true)}
	store.uptime["https://example.com"] = 100.0

	h := statuspage.Handler(store)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	body := rec.Body.String()
	if !strings.Contains(body, "https://example.com") {
		t.Error("expected URL in body")
	}
	if !strings.Contains(body, "100.0%") {
		t.Error("expected uptime percentage in body")
	}
	if !strings.Contains(body, "class=\"up\"") {
		t.Error("expected up class in body")
	}
}

func TestHandler_ShowsDownEndpoint(t *testing.T) {
	store := newFakeStore()
	store.urls = []string{"https://broken.io"}
	store.results["https://broken.io"] = []history.Result{makeResult("https://broken.io", false)}
	store.uptime["https://broken.io"] = 42.5

	h := statuspage.Handler(store)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	body := rec.Body.String()
	if !strings.Contains(body, "class=\"down\"") {
		t.Error("expected down class in body")
	}
	if !strings.Contains(body, "42.5%") {
		t.Error("expected uptime percentage in body")
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	store := newFakeStore()
	h := statuspage.Handler(store)
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for empty store, got %d", rec.Code)
	}
}
