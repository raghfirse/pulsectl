package reporter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/pulsectl/internal/history"
)

// PrintHistorySummary writes a formatted uptime summary table for all tracked
// endpoints to w, using data from the provided history Store.
func PrintHistorySummary(w io.Writer, store *history.Store) {
	urls := store.URLs()
	if len(urls) == 0 {
		fmt.Fprintln(w, "No history available.")
		return
	}
	sort.Strings(urls)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ENDPOINT\tCHECKS\tUPTIME")
	fmt.Fprintln(tw, "--------\t------\t------")

	for _, url := range urls {
		records := store.Get(url)
		pct := store.UptimePercent(url)
		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\n", url, len(records), pct)
	}

	tw.Flush()
}
