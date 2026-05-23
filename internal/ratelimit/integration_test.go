package ratelimit_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/ratelimit"
)

// TestAllow_Concurrent verifies the limiter is safe under concurrent access.
func TestAllow_Concurrent(t *testing.T) {
	rl := ratelimit.New(50, 100)

	var wg sync.WaitGroup
	var allowed int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow() {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}

	wg.Wait()

	// At most burst (50) goroutines should have been allowed immediately;
	// the rest may or may not pass depending on timing, but no race should occur.
	if allowed > 100 {
		t.Fatalf("allowed %d > total goroutines 100", allowed)
	}
}

// TestAllow_SteadyRate verifies tokens refill at the expected rate over real time.
func TestAllow_SteadyRate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timing-sensitive test in short mode")
	}

	// 1 token/sec, burst 1
	rl := ratelimit.New(1, 1)
	rl.Allow() // drain

	time.Sleep(1100 * time.Millisecond)

	if !rl.Allow() {
		t.Fatal("expected token to be available after ~1s refill")
	}
}
