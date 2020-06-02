package reverse_proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

type ReverseProxy struct {
	targetUrl  string
	proxyPort  int
	proxy      *httputil.ReverseProxy
	router     *mux.Router
	httpServer *http.Server
	errChannel chan error
	running    bool
}
