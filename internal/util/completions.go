package util

import (
	"encoding/json"
	"errors"

	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"go.etcd.io/bbolt"
)

type TaskPreview struct {
	ID    string
	Title string
}

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
		b := tx.Bucket([]byte(storage.TASK_DB))
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
