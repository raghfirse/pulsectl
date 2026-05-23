package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/pulsectl/internal/ratelimit"
)

func TestAllow_FullBucket(t *testing.T) {
	rl := ratelimit.New(3, 1)
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow()=true on token %d", i)
		}
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	rl := ratelimit.New(2, 1)
	rl.Allow()
	rl.Allow()
	if rl.Allow() {
		t.Fatal("expected Allow()=false when bucket is empty")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }

	rl := ratelimit.NewWithClock(1, 2, clock) // 2 tokens/sec, burst 1
	rl.Allow()                                 // drain

	if rl.Allow() {
		t.Fatal("expected bucket empty after drain")
	}

	// Advance clock by 1 second — should refill 2 tokens, capped at max=1
	now = now.Add(1 * time.Second)
	if !rl.Allow() {
		t.Fatal("expected Allow()=true after refill")
	}
}

func TestAllow_PartialRefill(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }

	rl := ratelimit.NewWithClock(2, 1, clock) // 1 token/sec, burst 2
	rl.Allow()
	rl.Allow() // drain fully

	// Advance only 0.4 seconds — not enough to add a full token
	now = now.Add(400 * time.Millisecond)
	if rl.Allow() {
		t.Fatal("expected Allow()=false with partial refill")
	}

	// Advance another 0.7 seconds — total 1.1s, enough for 1 token
	now = now.Add(700 * time.Millisecond)
	if !rl.Allow() {
		t.Fatal("expected Allow()=true after full token refill")
	}
}

func TestAvailable_ReturnsTokenCount(t *testing.T) {
	rl := ratelimit.New(5, 1)
	if got := rl.Available(); got != 5 {
		t.Fatalf("expected 5 tokens, got %f", got)
	}
	rl.Allow()
	rl.Allow()
	if got := rl.Available(); got != 3 {
		t.Fatalf("expected 3 tokens after 2 consumed, got %f", got)
	}
}
