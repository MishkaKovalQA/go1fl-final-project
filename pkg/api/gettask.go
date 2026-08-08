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
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "задача не найдена")
			return
		}

		writeError(w, http.StatusInternalServerError, "не удалось получить задачу")
		return
	}

	writeJSON(w, task)
}
