package ratelimiter

import (
	"log"
	"net/http"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/receiver"
	"time"
)

// Register this plugin in the plugin package
func init() {
	plugin.RegisterPlugin("ratelimiter", NewRateLimiterPlugin)
}

type RateLimiter struct {
	rate        int           // requests per second
	interval    time.Duration // time interval between requests
	ticker      *time.Ticker
	requestChan chan struct{}
}

// NewRateLimiter creates a new rate limiter with the given rate (requests per second).
func NewRateLimiter(rate int) *RateLimiter {
	interval := time.Second / time.Duration(rate)
	return &RateLimiter{
		rate:        rate,
		interval:    interval,
		ticker:      time.NewTicker(interval),
		requestChan: make(chan struct{}, rate),
	}
}

// Start starts the rate limiter. It should be called before making requests.
func (rl *RateLimiter) Start() {
	for range rl.ticker.C {
		rl.requestChan <- struct{}{}
	}
}

// Wait blocks until the rate limiter allows the next request.
func (rl *RateLimiter) Wait() {
	<-rl.requestChan
}

// RateLimiterPlugin implements the Plugin interface for rate limiting.
type RateLimiterPlugin struct {
	limiter    *RateLimiter
	Settings   map[string]interface{}
	ActiveMode bool
	Receivers  []receiver.NotificationReceiver
}

// Handle handles the incoming request and applies rate limiting.
func (r *RateLimiterPlugin) Handle(req *http.Request) error {
	// if the request channel is zero, there are no tokens in the bucket e.g. rate is being limited
	if len(r.limiter.requestChan) == 0 {
		r.Notify("request is being rate limited")
	}

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
	return false
}

func (r *RateLimiterPlugin) Notify(message string) {
	for _, receiver := range r.Receivers {
		err := receiver.Notify(message)
		if err != nil {
			log.Printf("unable to notify receiver. message %s - error: %v", message, err)
		}
	}
}

// Register the RateLimiter plugin in the plugin registry.
func NewRateLimiterPlugin(pluginSettings map[string]interface{}) plugin.Plugin {
	rate := int(pluginSettings["rate"].(float64))

	limiter := NewRateLimiter(rate)

	receivers, err := receiver.CreateReceivers(pluginSettings)
	if err != nil {
		panic(err)
	}

	limiterPlugin := &RateLimiterPlugin{
		limiter:   limiter,
		Settings:  pluginSettings,
		Receivers: receivers,
	}

	go limiter.Start()

	return limiterPlugin
}
