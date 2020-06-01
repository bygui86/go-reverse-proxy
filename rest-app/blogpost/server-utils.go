package blogpost

import (
	"fmt"
	"net/http"

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
)

func initBlogPosts() []*blogPost {
	return []*blogPost{
		buildBlogPost(post1_id, post1_title, post1_body),
		buildBlogPost(post2_id, post2_title, post2_body),
	}
}

func initRoutes() []*serverRoute {
	return []*serverRoute{}
}

// setupRouter - Create new Gorilla mux router
func (s *Server) setupRouter() {
	logging.Log.Debug("Create new router")

	s.router = mux.NewRouter().StrictSlash(true)

	// blogPosts
	s.addRoute(routerPostsRootUrl, s.getBlogPostByQuery, http.MethodGet,
		true, false, map[string]string{idKey: idValue})
	s.addRoute(routerPostsRootUrl, s.getBlogPosts, http.MethodGet,
		true, false, nil)
	s.addRoute(routerPostsRootUrl, s.createBlogPost, http.MethodPost,
		true, true, nil)
	s.addRoute(routerPostsRootUrl+routerIdUrlPath, s.getBlogPostByPath, http.MethodGet,
		true, false, nil)
	s.addRoute(routerPostsRootUrl+routerIdUrlPath, s.updateBlogPost, http.MethodPut,
		true, true, nil)
	s.addRoute(routerPostsRootUrl+routerIdUrlPath, s.deleteBlogPost, http.MethodDelete,
		true, false, nil)

	// routes
	s.addRoute(routerRoutesRootUrl, s.getRoutes, http.MethodGet,
		true, false, nil)

	// root
	s.addRoute(routerRootUrl, s.getRoot, http.MethodGet,
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
func (s *Server) addRoute(url string, handler func(http.ResponseWriter, *http.Request),
	method string, acceptHeader, contentTypeHeader bool, queries map[string]string) {

	route := s.router.HandleFunc(url, handler)
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

// setupHTTPServer - Create new HTTP server
func (s *Server) setupHTTPServer() {
	logging.SugaredLog.Debugf("Create new HTTP server on port %d", s.config.RestPort)

	if s.config != nil {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(httpServerHostFormat, s.config.RestHost, s.config.RestPort),
			Handler: s.router,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: httpServerWriteTimeoutDefault,
			ReadTimeout:  httpServerReadTimeoutDefault,
			IdleTimeout:  httpServerIdelTimeoutDefault,
		}
		return
	}

	logging.Log.Error("HTTP server creation failed: REST server configurations not initialized")
}
