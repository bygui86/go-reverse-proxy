package reverse_proxy

import (
	"net/http"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
)

const (
	hostHeaderKey       = "Host"
	additionalHeaderKey = "X-Forwarded-Host"

	rootEndpoint = "/"
)

func NewReverseProxy(targetUrl string, port int) (*ReverseProxy, error) {
	logging.SugaredLog.Debugf("Create new reverse proxy for target %s", targetUrl)

	logging.Log.Debug("Create reverse proxy")
	proxy, err := createSingleHostReverseProxy(targetUrl)
	if err != nil {
		return nil, err
	}

	return &ReverseProxy{
		targetUrl:  targetUrl,
		port:       port,
		proxy:      proxy,
		errChannel: make(chan error, 1),
		running:    false,
	}, nil
}

func (r *ReverseProxy) Start() {
	logging.Log.Info("Start reverse proxy")

	http.HandleFunc(rootEndpoint, r.proxy.ServeHTTP)

	go r.startHttpServerController()

	go r.listenAndServe()

	logging.SugaredLog.Infof("Reverse proxy listening on port %d", r.port)
}

// TODO wrap http server with gorilla mux router
func (r *ReverseProxy) Shutdown(timeout int) {
	logging.Log.Warn("Shutdown reverse proxy")

	// http
	// if r.proxy != nil {
	// 	r.running = false
	// }
}
