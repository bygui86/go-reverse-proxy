package reverse_proxy

import (
	"context"
	"time"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
)

func NewReverseProxy(targetUrl string, host string, port int) (*ReverseProxy, error) {
	logging.SugaredLog.Debugf("Create new reverse proxy for target %s", targetUrl)

	logging.Log.Debug("Create reverse proxy")
	proxy, err := setupSingleHostReverseProxy(targetUrl)
	if err != nil {
		return nil, err
	}

	logging.Log.Debug("Create router")
	router := setupRouter(proxy)

	logging.Log.Debug("Create HTTP server")
	httpServer := setupHttpServer(router, host, port)

	return &ReverseProxy{
		targetUrl:  targetUrl,
		proxyPort:  port,
		proxy:      proxy,
		router:     router,
		httpServer: httpServer,
		errChannel: make(chan error, 1),
		running:    false,
	}, nil
}

func (r *ReverseProxy) Start() {
	logging.Log.Info("Start reverse proxy")

	if r.proxy != nil && r.router != nil && r.httpServer != nil && !r.running {
		go r.startHttpServerController()

		go r.listenAndServe()

		r.running = true
		logging.SugaredLog.Infof("Reverse proxy listening on proxyPort %d", r.proxyPort)
		return
	}

	logging.Log.Error("Reverse proxy start failed: proxy, router, HTTP server not initialized or HTTP server already running")
}

func (r *ReverseProxy) Shutdown(timeout int) {
	logging.SugaredLog.Warnf("Shutdown reverse proxy, timeout %d", timeout)

	if r.proxy != nil && r.router != nil && r.httpServer != nil && !r.running {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		err := r.httpServer.Shutdown(ctx)
		if err != nil {
			logging.SugaredLog.Errorf("Reverse proxy shutdown failed: %s", err.Error())
		}
		r.running = false
		return
	}

	logging.Log.Error("Reverse proxy shutdown failed: proxy, router, HTTP server not initialized or HTTP server not running")
}
