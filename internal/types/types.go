package types

import (
	"time"
)

type Task struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	DueDate     time.Time         `json:"due_date"`
	Status      string            `json:"status"` // "pending" or "done"
	Priority    int               `json:"priority"`
	ChildTasks  []Task            `json:"child_tasks"`
	Notes       map[string]string `json:"notes"`
}
