package main

import (
	"log"

	"go1fl-final-project/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
