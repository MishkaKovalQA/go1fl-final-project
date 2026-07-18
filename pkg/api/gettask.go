package api

import (
	"database/sql"
	"errors"
	"net/http"

	"go1fl-final-project/pkg/db"
)

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
