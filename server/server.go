package server

import (
	"log"
	"syscall"

	"os"

	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

type server struct {
	Client     client.ClientIFace
	middleware []interface{}
	config     *config.Config
	mux        *web.Mux
}

func New(cfg *config.Config, mux *web.Mux, cl client.ClientIFace) (srv *server) {
	srv = &server{
		Client: cl,
		mux:    mux,
		config: cfg,
	}

	return
}

func NewDefault() (srv *server) {
	cfg := config.DefaultConfig()
	client, err := client.New(cfg).Watch()

	if err != nil {
		printer.LogErr("could not create client: %v", err)
	}

	srv = New(cfg, goji.DefaultMux, client)

	return
}

func (srv *server) Use(mdl ...web.MiddlewareType) *server {
	for _, m := range mdl {
		srv.middleware = append(srv.middleware, m)
	}

	return srv
}

func (srv *server) Serve() {
	for _, md := range srv.middleware {
		srv.mux.Use(md)
	}

	srv.mux.Get(srv.config.Server.Endpoint, handlers.FeaturesHandler(srv.Client))

	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.PreHook(func() {
		printer.Log("received kill signal, gracefully stopping...")
	})

	printer.Log("pid: %d serving %s on %s", os.Getpid(), srv.config.Server.Endpoint, srv.config.Server.Host)
	err := graceful.ListenAndServe(srv.config.Server.Host, srv.mux)

	if err != nil {
		log.Println(err)
	}
}
