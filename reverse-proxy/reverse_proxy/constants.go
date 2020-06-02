package reverse_proxy

import "time"

const (
	httpServerHostFormat          = "%s:%d"
	httpServerWriteTimeoutDefault = time.Second * 15
	httpServerReadTimeoutDefault  = time.Second * 15
	httpServerIdelTimeoutDefault  = time.Second * 60

	rootEndpoint    = "/"
	forwardEndpoint = "/{.*}"

	xForwardedHost = "X-Forwarded-Host" // reverse proxy host
	hostHeaderKey  = "Host"
)
