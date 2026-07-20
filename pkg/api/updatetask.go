package api

import (
	"encoding/json"
	"net/http"

	"go1fl-final-project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "Ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
