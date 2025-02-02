package httpserver

import (
	"net/http"
	"time"
)

type HttpServer interface {
	Run() error
}

type httpServer struct {
	srv *http.Server
}

type ServerConfig struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Handler      http.Handler
}

func NewHttpServer(config ServerConfig) HttpServer {
	return &httpServer{
		srv: &http.Server{
			Addr:         config.Addr,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout,
			IdleTimeout:  config.IdleTimeout,
			Handler:      config.Handler,
		},
	}
}

func (s *httpServer) Run() error {
	return s.srv.ListenAndServe()
}
