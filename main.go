package main

import (
	"log"

	"github.com/webbben/task/cmd"
	"github.com/webbben/task/internal/storage"
)

func main() {
	if err := storage.OpenDatabase(storage.TASK_DB); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer storage.CloseDatabase()

	cmd.Execute()
}
