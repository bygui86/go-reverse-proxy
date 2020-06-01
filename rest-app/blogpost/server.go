package blogpost

import (
	"context"
	"time"

	"github.com/bygui86/go-reverse-proxy/rest-app/logging"
)

const (
	idKey   = "id"
	idValue = "{" + idKey + ":[0-9]+}"

	urlSeparator        = "/"
	routerRootUrl       = urlSeparator
	routerPostsUrl      = "blogPosts"
	routerPostsRootUrl  = routerRootUrl + routerPostsUrl
	routerIdUrlPath     = urlSeparator + idValue
	routerRoutesUrl     = "routes"
	routerRoutesRootUrl = routerRootUrl + routerRoutesUrl

	httpServerHostFormat          = "%s:%d"
	httpServerWriteTimeoutDefault = time.Second * 15
	httpServerReadTimeoutDefault  = time.Second * 15
	httpServerIdelTimeoutDefault  = time.Second * 60
)

// NewRestServer - Create new REST server
func NewRestServer() *Server {
	logging.Log.Debug("Create new REST server")

	cfg := loadConfig()

	server := &Server{
		config:    cfg,
		blogPosts: initBlogPosts(),
		routes:    initRoutes(),
	}

	server.setupRouter()
	server.setupHTTPServer()
	return server
}

// Start - Start REST server
func (s *Server) Start() {
	logging.Log.Info("Start REST server")

	if s.httpServer != nil && !s.running {
		go func() {
			err := s.httpServer.ListenAndServe()
			if err != nil {
				logging.SugaredLog.Errorf("Error starting REST server: %s", err.Error())
			}
		}()
		s.running = true
		logging.SugaredLog.Infof("REST server listening on port %d", s.config.RestPort)
		return
	}

	logging.Log.Error("REST server start failed: HTTP server not initialized or HTTP server already running")
}

// Shutdown - Shutdown REST server
func (s *Server) Shutdown(timeout int) {
	logging.SugaredLog.Warnf("Shutdown REST server, timeout %d", timeout)

	if s.httpServer != nil && s.running {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logging.SugaredLog.Errorf("Error shutting down REST server: %s", err.Error())
		}
		s.running = false
		return
	}

	logging.Log.Error("REST server shutdown failed: HTTP server not initialized or HTTP server not running")
}
