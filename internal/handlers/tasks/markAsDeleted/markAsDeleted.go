package markAsDeleted

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/consts"
	"log/slog"
	"net/http"
)

type DB interface {
	MarkAsDeleted(ctx context.Context, taskId int, userId int) error
}

type MarkAsDeletedRequest struct {
	TaskId int `json:"taskId"`
}

func New(logger *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.markAsDeleted.New"
		log := *logger.With(
			"op", op,
			"request_id", r.Header.Get(consts.RequestIdHeader),
		)

		userId := r.Context().Value(consts.AuthUserIdKey).(int)
		var req MarkAsDeletedRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.MarkAsDeleted(r.Context(), req.TaskId, userId)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
