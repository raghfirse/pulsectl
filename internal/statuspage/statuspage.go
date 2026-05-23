// Package statuspage provides a simple HTTP handler that renders a
// human-readable status page summarising the current health of all
// monitored endpoints.
package statuspage

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/user/pulsectl/internal/history"
)

const pageTmpl = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Pulsectl Status</title>
<style>
body{font-family:sans-serif;margin:2rem;}
table{border-collapse:collapse;width:100%;}
th,td{border:1px solid #ccc;padding:0.5rem 1rem;text-align:left;}
th{background:#f4f4f4;}
.up{color:green;font-weight:bold;}
.down{color:red;font-weight:bold;}
</style>
</head>
<body>
<h1>Pulsectl Status</h1>
<p>Generated: {{.Generated}}</p>
<table>
<tr><th>Endpoint</th><th>Status</th><th>Uptime</th><th>Last Checked</th></tr>
{{range .Rows}}
<tr>
  <td>{{.URL}}</td>
  <td class="{{if .Healthy}}up{{else}}down{{end}}">{{if .Healthy}}UP{{else}}DOWN{{end}}</td>
  <td>{{.Uptime}}</td>
  <td>{{.LastChecked}}</td>
</tr>
{{end}}
</table>
</body>
</html>`

type row struct {
	URL         string
	Healthy     bool
	Uptime      string
	LastChecked string
}

type pageData struct {
	Generated string
	Rows      []row
}

// Store is the subset of history.Store used by the status page handler.
type Store interface {
	URLs() []string
	Get(url string) []history.Result
	UptimePercent(url string) float64
}

// Handler returns an http.HandlerFunc that renders the status page.
func Handler(store Store) http.HandlerFunc {
	tmpl := template.Must(template.New("status").Parse(pageTmpl))

	return func(w http.ResponseWriter, r *http.Request) {
		urls := store.URLs()
		rows := make([]row, 0, len(urls))

		for _, u := range urls {
			results := store.Get(u)
			var healthy bool
			var lastChecked string
			if len(results) > 0 {
				last := results[len(results)-1]
				healthy = last.Healthy
				lastChecked = last.Timestamp.Format(time.RFC3339)
			} else {
				lastChecked = "never"
			}
			rows = append(rows, row{
				URL:         u,
				Healthy:     healthy,
				Uptime:      fmt.Sprintf("%.1f%%", store.UptimePercent(u)),
				LastChecked: lastChecked,
			})
		}

		data := pageData{
			Generated: time.Now().UTC().Format(time.RFC1123),
			Rows:      rows,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tmpl.Execute(w, data)
	}
}
