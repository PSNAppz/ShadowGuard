// plugins.go

package plugins

import (
	_ "shadowguard/plugins/ipfilter"
	_ "shadowguard/plugins/monitor"
	_ "shadowguard/plugins/ratelimiter"
)
