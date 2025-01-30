package addTask

import (
	"context"
	"encoding/json"
	"go-httpnet-todo-list/internal/database"
	"net/http"
)

type DB interface {
	AddTask(ctx context.Context, task database.Task) error
}

type AddTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func New(db DB, authUserIdKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(authUserIdKey).(int)

		var req AddTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
