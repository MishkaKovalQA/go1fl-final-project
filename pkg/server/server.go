package server

import (
	"net/http"
	"os"
)

const (
	defaultPort = "7540"
	webDir      = "web"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe(":"+port, nil)
}
