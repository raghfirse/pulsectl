// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently alerts or notifications are sent for a given endpoint.
//
// Tokens are consumed on each allowed action and replenished over time
// at a configurable rate. A burst capacity controls the maximum number
// of tokens that can accumulate.
package ratelimit
