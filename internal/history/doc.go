// Package history provides an in-memory, bounded store for health-check
// results. It tracks results per endpoint URL and exposes uptime statistics
// computed from the accumulated records.
//
// Usage:
//
//	store := history.New(100) // keep last 100 results per endpoint
//	store.Add(result)         // record a checker.Result
//	pct := store.UptimePercent("http://example.com") // 0–100
package history
