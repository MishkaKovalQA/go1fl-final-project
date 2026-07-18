package server

import (
	"net/http"
	"os"

	"go1fl-final-project/pkg/api"
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

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe(":"+port, nil)
}
