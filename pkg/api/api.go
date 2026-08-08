package api

import (
	"encoding/json"
	"net/http"
	"os"
)

var todoPassword string

func Init() {
	todoPassword = os.Getenv("TODO_PASSWORD")
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)

	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
