package reverse_proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
)

func setupRouter(proxy *httputil.ReverseProxy) *mux.Router {
	logging.Log.Debug("Setup new router")

	router := mux.NewRouter().StrictSlash(true)

	// INFO part 1: to forward all requests including root url, we must specify a forwarding for the root url as well
	router.HandleFunc(rootEndpoint, proxy.ServeHTTP)
	// INFO part 2: otherwise we could manage the root url internally in the reverse proxy
	// http.Handle(rootEndpoint, router)

	router.HandleFunc(forwardEndpoint, proxy.ServeHTTP)
	return router
}

func setupHttpServer(router *mux.Router, host string, port int) *http.Server {
	logging.Log.Debug("Setup HTTP server")

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(httpServerHostFormat, host, port),
		Handler: router,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: httpServerWriteTimeoutDefault,
		ReadTimeout:  httpServerReadTimeoutDefault,
		IdleTimeout:  httpServerIdelTimeoutDefault,
	}
	return httpServer
}

func (r *ReverseProxy) listenAndServe() {
	logging.SugaredLog.Debugf("Listen and serve on port %d", r.proxyPort)

	// using http reverse proxy only
	// r.errChannel <- http.ListenAndServe(fmt.Sprintf(httpServerHostFormat, r.proxyPort), nil)

	// using gorilla mux router with reverse proxy
	// r.errChannel <- http.ListenAndServe(fmt.Sprintf(httpServerHostFormat, r.proxyPort), r.router)

	// using http server containing gorilla mux router with reverse proxy
	r.errChannel <- r.httpServer.ListenAndServe()
}

func (r *ReverseProxy) startHttpServerController() {
	logging.Log.Debug("Start HTTP server controller")

	for err := range r.errChannel {
		logging.SugaredLog.Errorf("HTTP server failed and stopped working: %s", err.Error())
		r.running = false
		os.Exit(502)
	}
}
