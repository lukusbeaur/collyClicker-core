package scraper

import (
	"log/slog"
	"time"

	colly "github.com/gocolly/colly/v2"
	"github.com/lukusbeaur/collyclicker-core/internal/Util"
)

// Users define these in app/handlers/*.go and return []SelectorHandler.
type SelectorHandler struct {
	Name     string
	Selector string
	Handler  func(e *colly.HTMLElement)
}

// ScraperConfig holds either a ready collector or knobs to build one.
type ScraperConfig struct {
	// Provide a collector OR leave nil and use settings below.
	Collector *colly.Collector

	// Builder settings (used only if Collector is nil)
	AllowedDomains []string
	Async          bool
	UserAgent      string
	Parallelism    int
	Delay          time.Duration
	RandomDelay    time.Duration
	IgnoreRobots   bool
	AllowRevisit   bool
	DefaultHeaders map[string]string

	// Proxies (optional)
	UseProxy  bool
	ProxyList []string // e.g., http://user:pass@host:port; rotate yourself if needed

	// Handlers for this scraping task
	LinkSelectors []SelectorHandler

	// Debug logging
	Debug bool
}

func defaultConfig() *ScraperConfig {
	return &ScraperConfig{
		AllowedDomains: []string{"scrapethissite.com", "www.scrapethissite.com"},
		Async:          false,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		Parallelism:    1,
		Delay:          2 * time.Second,
		RandomDelay:    4 * time.Second,
		IgnoreRobots:   false,
		AllowRevisit:   true,
		DefaultHeaders: map[string]string{
			"Referer":         "https://scrapethissite.com",
			"Accept-Language": "en-US,en;q=0.9",
		},
	}
}

func buildCollector(cfg *ScraperConfig) *colly.Collector {
	// If the user provided a collector, use it as-is.
	if cfg.Collector != nil {
		return cfg.Collector
	}

	c := colly.NewCollector(
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

	// Default headers
	c.OnRequest(func(r *colly.Request) {
		for k, v := range cfg.DefaultHeaders {
			r.Headers.Set(k, v)
		}
	})

	// Optional proxy: simple “first proxy” example. (Rotate in OnRequest if needed.)
	if cfg.UseProxy && len(cfg.ProxyList) > 0 {
		if err := c.SetProxy(cfg.ProxyList[0]); err != nil && cfg.Debug {
			Util.Logger.Warn("Proxy set failed",
				slog.String("proxy", cfg.ProxyList[0]),
				slog.Any("err", err))
		}
	}

	// Lightweight error logging
	c.OnError(func(r *colly.Response, err error) {
		if cfg.Debug {
			Util.Logger.Error("Request failed",
				slog.String("url", r.Request.URL.String()),
				slog.Int("status", r.StatusCode),
				slog.Any("err", err))
		}
	})

	return c
}

// CollyScraper encapsulates a configured collector + registered handlers.
type CollyScraper struct {
	Config    *ScraperConfig
	Collector *colly.Collector
}

// NewCollyScraper wires handlers into a (built or provided) collector.
func NewCollyScraper(cfg *ScraperConfig) *CollyScraper {
	if cfg == nil {
		cfg = defaultConfig()
	} else {
		// Merge zero-values with defaults (simple fill).
		def := defaultConfig()
		if cfg.AllowedDomains == nil || len(cfg.AllowedDomains) == 0 {
			cfg.AllowedDomains = def.AllowedDomains
		}
		if cfg.UserAgent == "" {
			cfg.UserAgent = def.UserAgent
		}
		if cfg.Parallelism <= 0 {
			cfg.Parallelism = def.Parallelism
		}
		if cfg.Delay <= 0 {
			cfg.Delay = def.Delay
		}
		if cfg.RandomDelay < 0 {
			cfg.RandomDelay = def.RandomDelay
		}
		if cfg.DefaultHeaders == nil {
			cfg.DefaultHeaders = def.DefaultHeaders
		}
	}

	c := buildCollector(cfg)

	for _, sh := range cfg.LinkSelectors {
		handler := sh // capture
		c.OnHTML(handler.Selector, handler.Handler)
	}

	return &CollyScraper{
		Config:    cfg,
		Collector: c,
	}
}

// Scrape visits the URL and waits for all callbacks to complete.
func (s *CollyScraper) Scrape(url string) error {
	err := s.Collector.Visit(url)
	s.Collector.Wait()
	return err
}
