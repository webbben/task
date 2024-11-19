package completions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"github.com/webbben/task/internal/util"
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

// MatchFromListCompletionFn is a completion function for a given list of possible options.
//
// if the given string is a substring (or equal) to any of the options, those options will be returned.
func MatchFromListCompletionFn(s string, options []string, cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	matches := make([]string, 0)

	for _, op := range options {
		if len(s) > len(op) {
			continue
		}
		if s == op[:len(s)] {
			matches = append(matches, op)
		}
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}

// TaskIDCompletionFn provides a completion function for completing a task ID
//
// firstArgOnly - if true, only do completion if its the first arg after the command
func TaskIDCompletionFn(firstArgOnly bool) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// only do completion for the task ID arg, which is the first one
		if firstArgOnly && len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// TODO filter out already entered task IDs if args > 0

		taskPreviews, err := CompleteTaskID(toComplete)
		if err != nil {
			cmd.PrintErrln("Error finding task IDs:", err)
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var matches []string
		for _, preview := range taskPreviews {
			comp := fmt.Sprintf("%s\t(%s)", preview.ID, util.Truncate(preview.Title, 20))
			matches = append(matches, comp)
		}

		return matches, cobra.ShellCompDirectiveNoFileComp
	}
}
