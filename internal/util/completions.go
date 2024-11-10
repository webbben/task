package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"go.etcd.io/bbolt"
)

type TaskPreview struct {
	ID    string
	Title string
}

// CompleteTaskID finds all tasks that have an ID starting with the given string.
func CompleteTaskID(s string) ([]TaskPreview, error) {
	var taskPreviews []TaskPreview

	if len(s) > 8 {
		return taskPreviews, errors.New("given ID prefix is too long")
	}

	db := storage.DB()
	if db == nil {
		return taskPreviews, errors.New("failed to get tasks db")
	}

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			id := string(k)

			if id[:len(s)] == s {
				var data types.Task
				err := json.Unmarshal(v, &data)
				if err != nil {
					return errors.New("failed to unmarshal task data")
				}
				taskPreviews = append(taskPreviews, TaskPreview{
					ID:    id,
					Title: data.Title,
				})
			}
			return nil
		})
	})

	return taskPreviews, err
}

// TaskIDCompletionFn is a completion function for task IDs. Used for completions with cobra commands.
func TaskIDCompletionFn(s string, cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	taskPreviews, err := CompleteTaskID(s)
	if err != nil {
		cmd.PrintErrln("Error finding task IDs:", err)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var matches []string
	for _, preview := range taskPreviews {
		comp := fmt.Sprintf("%s\t(%s)", preview.ID, Truncate(preview.Title, 10))
		matches = append(matches, comp)
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}
