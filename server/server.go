package server

import (
	"log"
	"syscall"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

type server struct {
	Client     client.ClientIFace
	middleware []interface{}
	config     *config.Config
	mux        *web.Mux
}

func NewServer(cfg *config.Config, mux *web.Mux, cl client.ClientIFace) (s *server) {
	s = &server{
		Client: cl,
		mux:    mux,
		config: cfg,
	}

	return
}

func (s *server) Use(mdl ...web.MiddlewareType) *server {
	for _, m := range mdl {
		s.middleware = append(s.middleware, m)
	}

	return s
}

func (s *server) Serve() {
	for _, md := range s.middleware {
		s.mux.Use(md)
	}

	s.mux.Get(s.config.Server.Endpoint, handlers.FeaturesHandler(s.config, s.Client))

	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	graceful.PreHook(func() {
		log.Print("[dcdr] received kill signal, gracefully stopping...")
	})

	log.Printf("[dcdr] serving %s on %s", s.config.Server.Endpoint, s.config.Server.Host)
	err := graceful.ListenAndServe(s.config.Server.Host, s.mux)

	if err != nil {
		log.Println(err)
	}
}
