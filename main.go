package main

import (
	"log"
	"os"

	"go1fl-final-project/pkg/db"
	"go1fl-final-project/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	return server.Run()
}
