package main

import (
	"fmt"
	"net/http"
	"time"
)

type RateLimiter struct {
	rate     int           // requests per second
	interval time.Duration // time interval between requests
	ticker   *time.Ticker
	limiter  chan struct{}
}

// NewRateLimiter creates a new rate limiter with the given rate (requests per second).
func NewRateLimiter(rate int) *RateLimiter {
	interval := time.Second / time.Duration(rate)
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		ticker:   time.NewTicker(interval),
		limiter:  make(chan struct{}, rate),
	}
}

// Start starts the rate limiter. It should be called before making requests.
func (rl *RateLimiter) Start() {
	for range rl.ticker.C {
		rl.limiter <- struct{}{}
	}
}

// Wait blocks until the rate limiter allows the next request.
func (rl *RateLimiter) Wait() {
	<-rl.limiter
}

// RateLimiterPlugin implements the Plugin interface for rate limiting.
type RateLimiterPlugin struct {
	limiter    *RateLimiter
	Settings   map[string]interface{}
	ActiveMode bool
}

// Handle handles the incoming request and applies rate limiting.
func (r *RateLimiterPlugin) Handle(req *http.Request) error {
	r.limiter.Wait()
	return nil
}

// GetType returns the type of the plugin.
func (r *RateLimiterPlugin) GetType() string {
	return "ratelimiter"
}

// GetSettings returns the settings of the plugin.
func (r *RateLimiterPlugin) GetSettings() map[string]interface{} {
	return r.Settings
}

// IsActiveMode returns whether the plugin is in active mode.
func (r *RateLimiterPlugin) IsActiveMode() bool {
	return r.ActiveMode
}

// Register the RateLimiter plugin in the plugin registry.
func NewRateLimiterPlugin(settings map[string]interface{}, activeMode bool) plugin {
	rate := int(settings["rate"].(float64))
	limiter := NewRateLimiter(rate)
	go limiter.Start()
	return &RateLimiterPlugin{limiter: limiter, Settings: settings, ActiveMode: activeMode}
}
