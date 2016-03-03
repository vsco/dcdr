package server

import (
	"encoding/gob"
	"log"

	"github.com/pressly/cji"
	"github.com/vsco/dcdr/server/handlers"
	"github.com/vsco/dcdr/watcher"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func init() {
	gob.Register(&watcher.FeatureMap{})
}

type Server struct {
	mux        *web.Mux
	cfg        *watcher.Config
	middleware []interface{}
}

func New(cfg *watcher.Config, mux *web.Mux) (s *Server) {
	s = &Server{
		mux: mux,
		cfg: cfg,
	}

	return
}

func NewWithConfig(cfg *watcher.Config) (s *Server) {
	s = &Server{
		mux: goji.DefaultMux,
		cfg: cfg,
	}

	return
}

func NewDefault() (s *Server) {
	s = New(watcher.DefaultConfig(), goji.DefaultMux)

	return s
}

func (s *Server) Serve() {
	goji.Serve()
}

func (s *Server) Use(mdl ...web.MiddlewareType) *Server {
	for _, m := range mdl {
		s.middleware = append(s.middleware, m)
	}

	return s
}

func (s *Server) Mux() *web.Mux {
	return s.mux
}

func (s *Server) Init() *Server {
	ldr, err := watcher.NewWatcher(s.cfg)

	if err != nil {
		log.Fatal(err)
	}

	go ldr.WatchConfig()

	hdl := handlers.NewFeatureHandler(ldr, s.cfg).Serve
	s.mux.Get(s.cfg.FeatureEndpoint, cji.Use(s.middleware...).On(hdl))

	return s
}
