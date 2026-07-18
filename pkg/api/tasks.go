package api

import (
	"net/http"

	"go1fl-final-project/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается")
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeError(w, "Не удалось получить список задач")
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
