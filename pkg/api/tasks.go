package api

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	search := r.FormValue("search")

	tasks, err := db.Tasks(search, tasksLimit)
	if err != nil {
		writeError(w, "Не удалось получить список задач")
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, "Задача не найдена")
			return
		}

		writeError(w, "Не удалось получить задачу")
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
