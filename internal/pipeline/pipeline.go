// Package pipeline wires together the scheduler, alerting, rate limiting,
// and notifier components into a single cohesive processing loop.
package pipeline

import (
	"log"

	"github.com/user/pulsectl/internal/alerting"
	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/config"
	"github.com/user/pulsectl/internal/history"
	"github.com/user/pulsectl/internal/notifier"
	"github.com/user/pulsectl/internal/ratelimit"
	"github.com/user/pulsectl/internal/reporter"
)

// Pipeline connects results from the scheduler to alerting, history,
// reporting, and optional webhook notification.
type Pipeline struct {
	cfg      *config.Config
	alerter  *alerting.Alerter
	store    *history.Store
	reporter *reporter.Reporter
	notifier *notifier.Notifier
	limiter  *ratelimit.Limiter
}

// New constructs a Pipeline from the provided config.
func New(cfg *config.Config) *Pipeline {
	p := &Pipeline{
		cfg:     cfg,
		alerter: alerting.New(cfg.Alerting.Threshold),
		store:   history.New(cfg.History.MaxSize),
		reporter: reporter.New(),
		limiter: ratelimit.New(cfg.Alerting.RateLimit.Rate, cfg.Alerting.RateLimit.Burst),
	}
	if cfg.Webhook.URL != "" {
		p.notifier = notifier.New(cfg.Webhook.URL)
	}
	return p
}

// Process handles a single checker result: stores it, reports it,
// evaluates alert conditions, and fires a webhook if warranted.
func (p *Pipeline) Process(result checker.Result) {
	p.store.Add(result)
	p.reporter.Print(result)

	alert, triggered := p.alerter.Evaluate(result)
	if !triggered {
		return
	}
	if p.notifier == nil {
		return
	}
	if !p.limiter.Allow() {
		log.Printf("[pipeline] rate limit reached, suppressing alert for %s", result.URL)
		return
	}
	if err := p.notifier.Send(notifier.Payload{
		Endpoint: alert.URL,
		Status:   alert.Status,
		Message:  alert.Message,
	}); err != nil {
		log.Printf("[pipeline] webhook error for %s: %v", result.URL, err)
	}
}
