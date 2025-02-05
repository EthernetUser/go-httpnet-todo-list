package markTask

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/consts"
	"log/slog"
	"net/http"
)

type DB interface {
	MarkTask(ctx context.Context, taskId int, userId int, done bool) error
}

type MarkTaskRequest struct {
	TaskId int  `json:"taskId"`
	Done   bool `json:"done"`
}

func New(
	logger *slog.Logger,
	db DB,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.markTask.New"
		log := *logger.With(
			"op", op,
			"request_id", r.Header.Get(consts.RequestIdHeader),
		)

		userId := r.Context().Value(consts.AuthUserIdKey).(int)
		var req MarkTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.MarkTask(r.Context(), req.TaskId, userId, req.Done)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
