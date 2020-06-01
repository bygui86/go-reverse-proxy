package reverse_proxy

import (
	"net/http/httputil"
)

type ReverseProxy struct {
	targetUrl  string
	port       int
	proxy      *httputil.ReverseProxy
	errChannel chan error
	running    bool
}
