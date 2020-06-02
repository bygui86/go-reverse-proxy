package blogpost

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-reverse-proxy/rest-app/logging"
)

const (
	post1_id    = "1"
	post1_title = "My first BlogPost"
	post1_body  = "This is the content of my first BlogPost"

	post2_id    = "2"
	post2_title = "My second BlogPost"
	post2_body  = "This is the content of my second BlogPost"

	post3_id    = "3"
	post3_title = "My third BlogPost"
	post3_body  = "This is the content of my third BlogPost"
)

func initBlogPosts() []*blogPost {
	return []*blogPost{
		buildBlogPost(post1_id, post1_title, post1_body),
		buildBlogPost(post2_id, post2_title, post2_body),
		buildBlogPost(post3_id, post3_title, post3_body),
	}
}

func initRoutes() []*serverRoute {
	return []*serverRoute{}
}

// setupRouter - Setup new Gorilla mux router
func setupRouter() *mux.Router {
	logging.Log.Debug("Setup new router")

	return mux.NewRouter().StrictSlash(true)
}

func (s *Server) setupHandlers() {
	// blogPosts
	addRoute(s.router, blogPostsRootUrl, s.getBlogPostByQuery, http.MethodGet,
		true, false, map[string]string{idKey: idValue})
	addRoute(s.router, blogPostsRootUrl, s.getBlogPosts, http.MethodGet,
		true, false, nil)
	addRoute(s.router, blogPostsRootUrl, s.createBlogPost, http.MethodPost,
		true, true, nil)
	addRoute(s.router, blogPostsRootUrl+blogPostIdEndpointPath, s.getBlogPostByPath, http.MethodGet,
		true, false, nil)
	addRoute(s.router, blogPostsRootUrl+blogPostIdEndpointPath, s.updateBlogPost, http.MethodPut,
		true, true, nil)
	addRoute(s.router, blogPostsRootUrl+blogPostIdEndpointPath, s.deleteBlogPost, http.MethodDelete,
		true, false, nil)

	// routes
	addRoute(s.router, routesRootUrl, s.getRoutes, http.MethodGet,
		true, false, nil)

	// root
	addRoute(s.router, rootUrl, s.getRoot, http.MethodGet,
		false, false, nil)
}

/*
	If we need to store all routes for some reason, we can do it in a struct like following ...
		// Handler is responsible for defining a HTTP request serverRoute and corresponding handler.
		type Handler struct {
			Route func(r *mux.Route)   // Receives a serverRoute to modify, like adding path, methods, etc.
			Func  http.HandlerFunc     // HTTP Handler
		}
	... and use it as we wish.
*/
func addRoute(router *mux.Router, url string, handler func(http.ResponseWriter, *http.Request),
	method string, acceptHeader, contentTypeHeader bool, queries map[string]string) {

	route := router.HandleFunc(url, handler)
	route.Methods(method)
	if acceptHeader {
		route.Headers(acceptHeaderKey, applicationJsonValue)
	}
	if contentTypeHeader {
		route.Headers(contentTypeHeaderKey, applicationJsonValue)
	}
	if queries != nil && len(queries) > 0 {
		for key, value := range queries {
			route.Queries(key, value)
		}
	}
}

// setupHttpServer - Setup new HTTP server
func setupHttpServer(router *mux.Router, host string, port int) *http.Server {
	logging.SugaredLog.Debugf("Setup new HTTP server")

	return &http.Server{
		Addr:    fmt.Sprintf(httpServerHostFormat, host, port),
		Handler: router,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: httpServerWriteTimeoutDefault,
		ReadTimeout:  httpServerReadTimeoutDefault,
		IdleTimeout:  httpServerIdelTimeoutDefault,
	}
}

func (s *Server) startHttpServerController() {
	logging.Log.Debug("Start HTTP server controller")

	for err := range s.errChannel {
		logging.SugaredLog.Errorf("HTTP server failed and stopped working: %s", err.Error())
		s.running = false
		os.Exit(502)
	}
}

func (s *Server) listenAndServe() {
	logging.SugaredLog.Debugf("Listen and serve on port %d", s.config.RestPort)
	s.errChannel <- s.httpServer.ListenAndServe()
}
