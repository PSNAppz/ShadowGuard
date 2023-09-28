package requestfilter

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"shadowguard/pkg/plugin"
	"shadowguard/pkg/publisher"

	"github.com/oschwald/geoip2-golang"
)

var dbPath = "plugins/ipfilter/GeoLite2-Country.mmdb"
var Type string = "requestfilter"

func init() {
	plugin.RegisterPlugin(Type, NewRequestFilterPlugin)
}

type RequestFilterPlugin struct {
	Settings        map[string]interface{}
	activeMode      bool
	ipBlacklist     []interface{}
	ipWhitelist     []interface{}
	regionWhitelist []interface{}
	regionBlacklist []interface{}
	publishers      []publisher.Publisher
}

func NewRequestFilterPlugin(settings map[string]interface{}) plugin.Plugin {
	publishers, err := publisher.CreatePublishers(settings)
	if err != nil {
		panic(err)
	}

	ipBlacklist, _ := settings["ip-blacklist"].([]interface{})
	ipWhitelist, _ := settings["ip-whitelist"].([]interface{})
	regionWhitelist, _ := settings["region-whitelist"].([]interface{})
	regionBlacklist, _ := settings["region-blacklist"].([]interface{})

	return &RequestFilterPlugin{
		Settings:        settings,
		activeMode:      settings["active_mode"].(bool),
		ipBlacklist:     ipBlacklist,
		ipWhitelist:     ipWhitelist,
		regionWhitelist: regionWhitelist,
		regionBlacklist: regionBlacklist,
		publishers:      publishers,
	}
}

func (p *RequestFilterPlugin) Type() string {
	return Type
}

func (p *RequestFilterPlugin) IsActiveMode() bool {
	return p.activeMode
}

func (p *RequestFilterPlugin) Notify(message string) {
	for _, pub := range p.publishers {
		err := pub.Publish(message)
		if err != nil {
			log.Printf("unable to notify publisher. message %s - error: %v", message, err)
		}
	}
}

func (p *RequestFilterPlugin) Handle(r *http.Request) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	// Handle IP-based restrictions
	if len(p.ipBlacklist) > 0 {
		for _, blacklistedIP := range p.ipBlacklist {
			if ip == blacklistedIP {
				return errors.New("IP address is blacklisted")
			}
		}
	}

	if len(p.ipWhitelist) > 0 {
		isWhitelisted := false
		for _, whitelistedIP := range p.ipWhitelist {
			if ip == whitelistedIP {
				isWhitelisted = true
				break
			}
		}
		if !isWhitelisted {
			return errors.New("IP address not whitelisted")
		}
	}

	// Handle region-based restrictions
	region := getRegionForIP(ip)
	if len(p.regionBlacklist) > 0 {

		for _, restrictedRegion := range p.regionBlacklist {
			if region == restrictedRegion {
				log.Printf("region %s is restricted by policy", region)
				return errors.New("access from this region is restricted by policy")
			}
		}
	}

	if len(p.regionWhitelist) > 0 {
		isRegionWhitelisted := false
		for _, allowedRegion := range p.regionWhitelist {
			if region == allowedRegion {
				isRegionWhitelisted = true
				break
			}
		}
		if !isRegionWhitelisted {
			return errors.New("access from this region is not whitelisted")
		}
	}

	return nil
}

func getRegionForIP(ip string) string {
	// Handle PATH for testing
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, b, _, _ := runtime.Caller(0)

		// Root directory of this source file
		root := filepath.Dir(b)
		dbPath = filepath.Join(root, "GeoLite2-Country.mmdb")
	}
	// Open the GeoIP2 database.
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Parse the IP
	parsedIP := net.ParseIP(ip)

	// Lookup the IP in the database
	record, err := db.Country(parsedIP)
	if err != nil {
		log.Fatal(err)
	}

	// Return the ISO country code (e.g., "US", "IN")
	return record.Country.IsoCode
}
