// app/config/config.go
package config

import "time"

type Config struct {
	InputDir  string
	OutputDir string

	AllowedDomains []string
	Async          bool
	UserAgent      string
	Parallelism    int
	Delay          time.Duration
	RandomDelay    time.Duration
	IgnoreRobots   bool
	AllowRevisit   bool

	DefaultHeaders map[string]string

	// Optional: proxy/debug if you want users to control them too
	UseProxy  bool
	ProxyList []string
	Debug     bool
}

func DefaultConfig() *Config {
	return &Config{
		InputDir:       "Input_Links/",
		OutputDir:      "Output_Data/",
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
