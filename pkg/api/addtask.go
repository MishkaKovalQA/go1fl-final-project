package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go1fl-final-project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	default:
		writeError(w, "Метод не поддерживается")
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Ошибка десериализации JSON")
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Не удалось добавить задачу")
		return
	}

	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "" {
		task.Date = today
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("Неверный формат даты")
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, date) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}

func afterNow(now, date time.Time) bool {
	return now.Format(dateFormat) > date.Format(dateFormat)
}
