package main

import (
	"go-httpnet-todo-list/internal/config"
	"go-httpnet-todo-list/internal/consts"
	"go-httpnet-todo-list/internal/database/postgres"
	"go-httpnet-todo-list/internal/handlers/tasks/addTask"
	"go-httpnet-todo-list/internal/handlers/tasks/getTasks"
	"go-httpnet-todo-list/internal/handlers/tasks/markAsDeleted"
	"go-httpnet-todo-list/internal/handlers/tasks/markTask"
	"go-httpnet-todo-list/internal/httpserver"
	"go-httpnet-todo-list/internal/middlewares"
	"go-httpnet-todo-list/internal/router"
	"log"
)

func main() {
	cfg := config.New()
	v1 := router.New()
	loadRoutes(v1, cfg)

	middlewareWrapper := router.CreateMiddlewaresWrapper(
		middlewares.Logging,
		middlewares.Auth,
	)

	srvConfig := httpserver.ServerConfig{
		Addr:         cfg.HttpServer.Addr,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
		Handler:      middlewareWrapper(v1.GetMux()),
	}
	srv := httpserver.NewHttpServer(srvConfig)

	log.Printf("Starting server on %s", cfg.HttpServer.Addr)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadRoutes(r router.Router, cfg *config.Config) {
	db := postgres.New(cfg.Postgres.ConnString)
	r.Get("/get-tasks", getTasks.New(db, consts.AuthUserIdKey))
	r.Put("/mark-task", markTask.New(db, consts.AuthUserIdKey))
	r.Post("/add-task", addTask.New(db, consts.AuthUserIdKey))
	r.Delete("/mark-as-deleted", markAsDeleted.New(db, consts.AuthUserIdKey))
}
