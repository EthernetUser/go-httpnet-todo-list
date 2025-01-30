package getTasks

import (
	"encoding/json"
	"go-httpnet-todo-list/internal/database"
	"net/http"
)

type DB interface {
	GetTasks(userId int) []database.Task
}

type GetTacksHandler struct {
	db DB
	authUserIdKey string
}

func (h *GetTacksHandler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(h.authUserIdKey).(int)

		tasks := h.db.GetTasks(userId)

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(tasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)
	}
}
