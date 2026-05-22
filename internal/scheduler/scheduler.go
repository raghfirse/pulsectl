package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/config"
)

// Scheduler polls configured endpoints at their defined intervals.
type Scheduler struct {
	cfg     *config.Config
	checker *checker.Checker
	results chan checker.Result
	wg      sync.WaitGroup
}

// New creates a new Scheduler with the given config and checker.
func New(cfg *config.Config, chk *checker.Checker) *Scheduler {
	return &Scheduler{
		cfg:     cfg,
		checker: chk,
		results: make(chan checker.Result, len(cfg.Endpoints)*2),
	}
}

// Results returns the read-only channel of check results.
func (s *Scheduler) Results() <-chan checker.Result {
	return s.results
}

// Start launches a polling goroutine for each configured endpoint.
// It blocks until the context is cancelled, then waits for all workers to finish.
func (s *Scheduler) Start(ctx context.Context) {
	for _, ep := range s.cfg.Endpoints {
		s.wg.Add(1)
		go s.poll(ctx, ep)
	}

	<-ctx.Done()
	s.wg.Wait()
	close(s.results)
	log.Println("scheduler: all workers stopped")
}

// poll continuously checks a single endpoint until ctx is cancelled.
func (s *Scheduler) poll(ctx context.Context, ep config.Endpoint) {
	defer s.wg.Done()

	interval := time.Duration(ep.IntervalSeconds) * time.Second
	if interval <= 0 {
		interval = time.Duration(s.cfg.DefaultIntervalSeconds) * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("scheduler: starting poll for %s every %s", ep.URL, interval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result := s.checker.Check(ctx, ep.URL)
			select {
			case s.results <- result:
			default:
				log.Printf("scheduler: result channel full, dropping result for %s", ep.URL)
			}
		}
	}
}
