package getTasks

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/database"
	"net/http"
)

type DB interface {
	GetTasks(ctx context.Context, userId int) ([]database.Task, error)
}

func New(db DB, authUserIdKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(authUserIdKey).(int)

		tasks, err := db.GetTasks(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonResp, err := json.Marshal(tasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
