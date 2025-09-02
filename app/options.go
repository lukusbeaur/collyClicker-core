// internal/app/options.go
package app

import (
	"time"

	"github.com/lukusbeaur/collyclicker-core/app/config"
)

type Option func(*config.Config)

func WithInputDir(dir string) Option {
	return func(c *config.Config) {
		if dir != "" { // default-to-sane if user passes ""
			c.InputDir = dir
		}
	}
}

func WithOutputDir(dir string) Option {
	return func(c *config.Config) {
		if dir != "" {
			c.OutputDir = dir
		}
	}
}

func WithAllowedDomains(domains ...string) Option {
	return func(c *config.Config) {
		if len(domains) > 0 {
			c.AllowedDomains = domains
		}
	}
}

func WithUserAgent(ua string) Option {
	return func(c *config.Config) {
		if ua != "" {
			c.UserAgent = ua
		}
	}
}

func WithParallelism(n int) Option {
	return func(c *config.Config) {
		if n > 0 {
			c.Parallelism = n
		}
	}
}

func WithDelays(delay, random time.Duration) Option {
	return func(c *config.Config) {
		if delay > 0 {
			c.Delay = delay
		}
		if random >= 0 {
			c.RandomDelay = random
		}
	}
}
