// plugins.go

package plugins

import (
	_ "AegisGuard/plugins/ipfilter"
	_ "AegisGuard/plugins/monitor"
	_ "AegisGuard/plugins/ratelimiter"
)
