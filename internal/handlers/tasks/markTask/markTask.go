package markTask

import (
	"context"
	"encoding/json"
	"net/http"
)

type DB interface {
	MarkTask(ctx context.Context, taskId int, userId int, done bool) error
}

type MarkTaskRequest struct {
	TaskId int  `json:"taskId"`
	Done   bool `json:"done"`
}

func New(db DB, authUserIdKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(authUserIdKey).(int)

		var req MarkTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.MarkTask(r.Context(), req.TaskId, userId, req.Done)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
