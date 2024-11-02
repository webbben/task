package storage

import (
	"encoding/json"
	"fmt"

	"github.com/webbben/task/internal/types"
	"go.etcd.io/bbolt"
)

// AddTaskToDB saves a new task in the BoltDB database
func AddTaskToDB(task types.Task) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		data, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(task.ID), data)
	})
}

// GetTaskFromDB retrieves a task from the BoltDB tasks bucket
func GetTaskFromDB(id string) (*types.Task, error) {
	var task types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("task not found")
		}
		return json.Unmarshal(data, &task)
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasksFromDB retrieves all tasks from the BoltDB tasks bucket
func GetAllTasksFromDB() ([]types.Task, error) {
	var tasks []types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		return b.ForEach(func(k, v []byte) error {
			var task types.Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	return tasks, err
}
