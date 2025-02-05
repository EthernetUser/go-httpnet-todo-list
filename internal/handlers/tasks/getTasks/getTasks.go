package getTasks

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/consts"
	"go-httpnet-todo-list/internal/database"
	"log/slog"
	"net/http"
)

type DB interface {
	GetTasks(ctx context.Context, userId int) ([]database.Task, error)
}

func New(logger *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getTasks.New"
		log := *logger.With(
			"op", op,
			"request_id", r.Header.Get(consts.RequestIdHeader),
		)

		userId := r.Context().Value(consts.AuthUserIdKey).(int)
		tasks, err := db.GetTasks(r.Context(), userId)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonResp, err := json.Marshal(tasks)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
