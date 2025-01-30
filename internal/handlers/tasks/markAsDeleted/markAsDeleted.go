package markAsDeleted

import (
	"context"
	"encoding/json"
	"net/http"
)

type DB interface {
	MarkAsDeleted(ctx context.Context, taskId int, userId int) error
}

type MarkAsDeletedRequest struct {
	TaskId int `json:"taskId"`
}

func New(db DB, authUserIdKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(authUserIdKey).(int)

		var req MarkAsDeletedRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.MarkAsDeleted(r.Context(), req.TaskId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
