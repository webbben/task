package main

import (
	"log"

	"github.com/webbben/task/cmd"
	"github.com/webbben/task/internal/storage"
)

func main() {
	if err := storage.OpenDatabase("tasks.db"); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer storage.CloseDatabase()

	cmd.Execute()
}
