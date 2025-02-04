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
	"log/slog"
	"os"
)

func main() {
	cfg := config.New()
	logger := InitLogger(cfg.Env)

	logger.Debug("init db")
	db, err := postgres.New(cfg.Postgres.ConnString)
	if err != nil {
		logger.Error("Failed to connect to db", "err", err.Error())
		os.Exit(1)
	}

	logger.Debug("init router")
	v1 := router.New()
	loadRoutes(v1, db)

	logger.Debug("init middlewares")
	middlewareWrapper := router.CreateMiddlewaresWrapper(
		middlewares.RequestId,
		middlewares.Logging(logger),
		middlewares.Auth,
	)

	logger.Debug("init server")
	srvConfig := httpserver.ServerConfig{
		Addr:         cfg.HttpServer.Addr,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
		Handler:      middlewareWrapper(v1.GetMux()),
	}
	srv := httpserver.NewHttpServer(srvConfig)

	logger.Info("Starting server", "address", cfg.HttpServer.Addr)
	if err := srv.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func loadRoutes(r router.Router, db postgres.Postgres) {
	r.Get("/get-tasks", getTasks.New(db, consts.AuthUserIdKey))
	r.Put("/mark-task", markTask.New(db, consts.AuthUserIdKey))
	r.Post("/add-task", addTask.New(db, consts.AuthUserIdKey))
	r.Delete("/mark-as-deleted", markAsDeleted.New(db, consts.AuthUserIdKey))
}

func InitLogger(env string) *slog.Logger {
	switch env {
	case "dev":
		return slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case "prod":
		return slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	case "test":
		return slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	default:
		return slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
}
