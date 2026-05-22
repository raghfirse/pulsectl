package history

import (
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
)

// Record holds a single check result with a timestamp.
type Record struct {
	Timestamp time.Time
	Result    checker.Result
}

// Store maintains a bounded in-memory history of check results per endpoint.
type Store struct {
	mu      sync.RWMutex
	records map[string][]Record
	maxSize int
}

// New creates a new Store that keeps at most maxSize records per endpoint.
func New(maxSize int) *Store {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Store{
		records: make(map[string][]Record),
		maxSize: maxSize,
	}
}

// Add appends a result to the history for its endpoint URL.
func (s *Store) Add(r checker.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := Record{Timestamp: time.Now(), Result: r}
	key := r.URL
	s.records[key] = append(s.records[key], rec)

	if len(s.records[key]) > s.maxSize {
		s.records[key] = s.records[key][len(s.records[key])-s.maxSize:]
	}
}

// Get returns a copy of all records for a given URL.
func (s *Store) Get(url string) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	src := s.records[url]
	out := make([]Record, len(src))
	copy(out, src)
	return out
}

// UptimePercent calculates the percentage of healthy checks for a URL.
// Returns 0 if no records exist.
func (s *Store) UptimePercent(url string) float64 {
	records := s.Get(url)
	if len(records) == 0 {
		return 0
	}
	var healthy int
	for _, r := range records {
		if r.Result.Healthy {
			healthy++
		}
	}
	return float64(healthy) / float64(len(records)) * 100
}

// URLs returns all tracked endpoint URLs.
func (s *Store) URLs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	urls := make([]string, 0, len(s.records))
	for u := range s.records {
		urls = append(urls, u)
	}
	return urls
}
