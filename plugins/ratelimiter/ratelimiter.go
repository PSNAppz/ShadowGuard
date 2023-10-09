package ratelimiter

import (
	"log"
	"net/http"
	"shadowguard/pkg/database"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"
	"time"
)

var Type string = "ratelimiter"

// Register this plugin in the plugin package
func init() {
	plugin.RegisterPlugin(Type, NewRateLimiterPlugin)
}

// RateLimiterPlugin implements the Plugin interface for rate limiting.
type RateLimiterPlugin struct {
	Settings   map[string]interface{}
	limiter    *RateLimiter
	publishers []publisher.Publisher
	db         *database.Database
}

// GetType returns the type of the plugin.
func (r *RateLimiterPlugin) Type() string {
	return Type
}

// IsActiveMode returns whether the plugin is in active mode.
func (r *RateLimiterPlugin) IsActiveMode() bool {
	return false
}

func (r *RateLimiterPlugin) Notify(message string) {
	for _, publisher := range r.publishers {
		err := publisher.Publish(message)
		if err != nil {
			log.Printf("unable to notify publisher. message %s - error: %v", message, err)
		}
	}
}

// Register the RateLimiter plugin in the plugin registry.
func NewRateLimiterPlugin(pluginSettings map[string]interface{}, db *database.Database) plugin.Plugin {
	rate := int(pluginSettings["rate"].(float64))

	limiter := NewRateLimiter(rate)

	publishers, err := publisher.CreatePublishers(pluginSettings)
	if err != nil {
		panic(err)
	}

	limiterPlugin := &RateLimiterPlugin{
		Settings:   pluginSettings,
		limiter:    limiter,
		publishers: publishers,
		db:         db,
	}

	go limiter.Start()

	return limiterPlugin
}

// Handle handles the incoming request and applies rate limiting.
func (r *RateLimiterPlugin) Handle(req *http.Request) error {
	// if the request channel is zero, there are no tokens in the bucket e.g. rate is being limited
	if len(r.limiter.requestChan) == 0 {
		r.Notify("request is being rate limited")
	}

	r.limiter.Wait(r.db, req)
	return nil
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
func (rl *RateLimiter) Wait(db *database.Database, r *http.Request) {
	requestModel, err := database.NewRequest(r, Type)
	if err != nil {
		log.Println("unable to database rate limited request. ")
	}
	db.Insert(requestModel)
	<-rl.requestChan
}
