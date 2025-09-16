// scrape/selectors.go
package scraper

import (
	"log/slog"

	"github.com/gocolly/colly/v2"
	appcfg "github.com/lukusbeaur/collyclicker-core/app/config"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
)

type SelectorHandler struct {
	Name     string
	Selector string
	Handler  func(e *colly.HTMLElement)
}

// Adapter: build a colly.Collector from the *app* config.
// - handlers: the []SelectorHandler you got from app/handlers
// - collector: pass nil to build from cfg; pass a prebuilt collector to reuse it
func NewCollectorFromAppConfig(cfg *appcfg.Config, handlers []SelectorHandler, collector *colly.Collector) *colly.Collector {
	if cfg == nil {
		cfg = appcfg.DefaultConfig()
	}

	c := collector
	if c == nil {
		c = colly.NewCollector(
			colly.AllowedDomains(cfg.AllowedDomains...),
			colly.Async(cfg.Async),
			colly.UserAgent(cfg.UserAgent),
		)

		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: cfg.Parallelism,
			Delay:       cfg.Delay,
			RandomDelay: cfg.RandomDelay,
		})
		c.IgnoreRobotsTxt = cfg.IgnoreRobots
		c.AllowURLRevisit = cfg.AllowRevisit

		c.OnRequest(func(r *colly.Request) {
			for k, v := range cfg.DefaultHeaders {
				r.Headers.Set(k, v)
			}
		})

		if cfg.UseProxy && len(cfg.ProxyList) > 0 {
			if err := c.SetProxy(cfg.ProxyList[0]); err != nil && cfg.Debug {
				Util.Logger.Warn("Proxy set failed",
					slog.String("proxy", cfg.ProxyList[0]),
					slog.Any("err", err))
			}
		}

		c.OnError(func(r *colly.Response, err error) {
			if cfg.Debug {
				Util.Logger.Error("Request failed",
					slog.String("url", r.Request.URL.String()),
					slog.Int("status", r.StatusCode),
					slog.Any("err", err))
			}
		})
	}

	// Register handlers
	for _, h := range handlers {
		hLocal := h
		c.OnHTML(hLocal.Selector, func(e *colly.HTMLElement) { hLocal.Handler(e) })
	}

	return c
}
