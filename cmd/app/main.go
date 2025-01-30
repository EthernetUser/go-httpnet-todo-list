package main

import (
	"go-httpnet-todo-list/internal/httpserver"
	"go-httpnet-todo-list/internal/middlewares/logging"
	"go-httpnet-todo-list/internal/router"
	"log"
	"net/http"
	"time"
)

func main() {
	v1 := router.New()
	loadRoutes(v1)

	middlewareWrapper := router.AddMiddlewares(logging.LoggingMiddleware)

	srvConfig := httpserver.ServerConfig{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      middlewareWrapper(v1.GetMux()),
	}
	srv := httpserver.NewHttpServer(srvConfig)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadRoutes(r router.Router) {
	r.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}
