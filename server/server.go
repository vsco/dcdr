package server

import (
	"net/http"

	"os"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/vsco/dcdr/server/middleware"
)

// Middleware helper type for handlers that receive a `Client`.
type Middleware func(client.IFace) func(http.Handler) http.Handler

// Server HTTP API for accessing `Features`
type Server struct {
	Client     client.IFace
	Router     *mux.Router
	middleware []Middleware
	config     *config.Config
}

// NewDefault creates a new `Server` using `config.hcl`.
func NewDefault() (srv *Server, err error) {
	cfg := config.DefaultConfig()
	client, err := client.New(cfg)

	if err != nil {
		return nil, err
	}

	srv = New(cfg, client)

	return
}

// New create a new `Server`
func New(cfg *config.Config, dcdr client.IFace) (srv *Server) {
	srv = &Server{
		Client: dcdr,
		Router: mux.NewRouter(),
		config: cfg,
	}

	return
}

// RegisterRoutes binds `Endpoint` to the `FeaturesHandler`.
func (srv *Server) RegisterRoutes() {
	srv.Router.Handle(srv.config.Server.Endpoint, srv.FeaturesHandler()).Methods("GET")
}

// FeaturesHandler delegates to `handlers.FeaturesHandler` and adds the
// middleware chain.
func (srv *Server) FeaturesHandler() http.Handler {
	fn := handlers.FeaturesHandler(srv.Client)

	return srv.WithMiddleware(http.HandlerFunc(fn))
}

// Use appends `Middleware` to the internal chain.
func (srv *Server) Use(h ...Middleware) {
	srv.middleware = append(srv.middleware, h...)
}

// ServeHTTP registers the `HTTPCachingHandler` and sets up the route
// handlers and logging.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Use(middleware.HTTPCachingHandler)
	srv.RegisterRoutes()

	logger := gh.CombinedLoggingHandler(os.Stdout, srv.Router)
	logger.ServeHTTP(w, r)
}

// Serve starts the server on the configured `Host`.
func (srv *Server) Serve() error {
	return http.ListenAndServe(srv.config.Server.Host, srv)
}

// WithMiddleware adds the middleware chain to `h` passing each the `Client`.
func (srv *Server) WithMiddleware(h http.Handler) http.Handler {
	for _, mw := range srv.middleware {
		h = mw(srv.Client)(h)
	}

	return h
}
