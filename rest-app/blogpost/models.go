package blogpost

import (
	"net/http"

	"github.com/gorilla/mux"
)

type config struct {
	RestHost string
	RestPort int
}

type Server struct {
	config     *config
	router     *mux.Router
	httpServer *http.Server
	errChannel chan error
	running    bool
	blogPosts  []*blogPost
	routes     []*serverRoute
}

type blogPost struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type serverRoute struct {
	PathTemplate    string   `json:"pathTemplate"`
	PathRegexp      string   `json:"pathRegexp"`
	QueriesTemplate []string `json:"queriesTemplate"`
	QueriesRegexp   []string `json:"queriesRegexp"`
	Methods         []string `json:"methods"`
}

type errorResponse struct {
	Request string `json:"request"`
	Message string `json:"message"`
}
