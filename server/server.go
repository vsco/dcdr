package server

import (
	"log"
	"syscall"

	"os"

	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/vsco/dcdr/server/middleware"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

type Server struct {
	Client     client.IFace
	middleware []interface{}
	config     *config.Config
	mux        *web.Mux
}

func New(cfg *config.Config, mux *web.Mux, cl client.IFace) (srv *Server) {
	srv = &Server{
		Client: cl,
		mux:    mux,
		config: cfg,
	}

	return
}

func NewDefault() (srv *Server) {
	cfg := config.DefaultConfig()
	client, err := client.New(cfg).Watch()

	if err != nil {
		printer.LogErrf("could not create client: %v", err)
	}

	srv = New(cfg, goji.DefaultMux, client)
	srv.Use(middleware.HTTPCachingHandler(client))

	return
}

func (srv *Server) Use(mdl ...web.MiddlewareType) *Server {
	for _, m := range mdl {
		srv.middleware = append(srv.middleware, m)
	}

	return srv
}

func (srv *Server) BindMux() *Server {
	for _, md := range srv.middleware {
		srv.mux.Use(md)
	}

	srv.mux.Get(srv.config.Server.Endpoint, handlers.FeaturesHandler(srv.Client))

	return srv
}

func (srv *Server) Serve() {
	srv.BindMux()

	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.PreHook(func() {
		printer.Logf("received kill signal, gracefully stopping...")
	})

	printer.Logf("pid: %d serving %s on %s", os.Getpid(), srv.config.Server.Endpoint, srv.config.Server.Host)
	err := graceful.ListenAndServe(srv.config.Server.Host, srv.mux)

	if err != nil {
		log.Println(err)
	}
}
