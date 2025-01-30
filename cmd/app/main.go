package main

import (
	"go-httpnet-todo-list/internal/database/postgres"
	"go-httpnet-todo-list/internal/handlers/tasks/addTask"
	"go-httpnet-todo-list/internal/handlers/tasks/getTasks"
	"go-httpnet-todo-list/internal/handlers/tasks/markAsDeleted"
	"go-httpnet-todo-list/internal/handlers/tasks/markTask"
	"go-httpnet-todo-list/internal/httpserver"
	"go-httpnet-todo-list/internal/middlewares/logging"
	"go-httpnet-todo-list/internal/router"
	"log"
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
	db := postgres.New("")
	r.Get("/get-tasks", getTasks.New(db, "userId"))
	r.Put("/mark-task", markTask.New(db, "userId"))
	r.Post("/add-task", addTask.New(db, "userId"))
	r.Delete("/mark-as-deleted", markAsDeleted.New(db, "userId"))
}
