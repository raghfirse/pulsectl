// Package healthlog provides append-only structured logging of health check
// results produced by the checker package.
//
// Each result is written as a single newline-terminated, tab-separated line:
//
//	<timestamp>\t<url>\t<UP|DOWN>\t<latency>\t[<error>]
//
// Log files are safe for concurrent use and are opened in append mode so
// entries survive process restarts.
package healthlog
