// Package ratelimit provides a simple token-bucket rate limiter
// to prevent overwhelming endpoints during health checks.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which health checks are dispatched.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to maxTokens burst and refills
// at ratePerSec tokens per second.
func New(maxTokens float64, ratePerSec float64) *Limiter {
	return NewWithClock(maxTokens, ratePerSec, time.Now)
}

// NewWithClock creates a Limiter with an injectable clock for testing.
func NewWithClock(maxTokens float64, ratePerSec float64, clock func() time.Time) *Limiter {
	return &Limiter{
		tokens:   maxTokens,
		max:      maxTokens,
		rate:     ratePerSec,
		lastTick: clock(),
		clock:    clock,
	}
}

// Allow reports whether one token is available and consumes it.
// Returns false if the rate limit is exceeded.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}

	l.tokens--
	return true
}

// Available returns the current token count (approximate, for observability).
func (l *Limiter) Available() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
