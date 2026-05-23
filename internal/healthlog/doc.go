// Package healthlog provides append-only structured logging of health check
// results produced by the checker package.
//
// Each result is written as a single newline-terminated, tab-separated line:
//
//	<timestamp>\t<url>\t<UP|DOWN>\t<latency>\t[<error>]
//
// The timestamp field uses RFC3339 format (e.g. 2006-01-02T15:04:05Z07:00).
// The latency field is expressed in milliseconds. The error field is omitted
// when the check succeeds and contains a human-readable message otherwise.
//
// Log files are safe for concurrent use and are opened in append mode so
// entries survive process restarts. Callers should close the log when done
// to ensure all buffered writes are flushed to disk.
package healthlog
