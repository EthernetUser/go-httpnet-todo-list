package addTask

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/consts"
	"go-httpnet-todo-list/internal/database"
	"log/slog"
	"net/http"
)

type DB interface {
	AddTask(ctx context.Context, task database.Task) error
}

type AddTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func New(logger *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.addTask.New"
		log := *logger.With(
			"op", op,
			"request_id", r.Header.Get(consts.RequestIdHeader),
		)

		userId := r.Context().Value(consts.AuthUserIdKey).(int)
		var req AddTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.AddTask(r.Context(), database.Task{
			Title:       req.Title,
			Description: req.Description,
			Done:        false,
			UserId:      userId,
		})
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
