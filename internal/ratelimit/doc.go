// Package ratelimit implements a token-bucket rate limiter for controlling
// the frequency of outbound health-check requests in pulsectl.
//
// Usage:
//
//	rl := ratelimit.New(10, 2) // burst of 10, refill 2/sec
//	if rl.Allow() {
//	    // perform check
//	}
//
// The limiter is safe for concurrent use.
package ratelimit
