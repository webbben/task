package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"github.com/webbben/task/internal/util"
	"go.etcd.io/bbolt"
)

const (
	TASK_DB = "tasks"
)

// AddTask creates a new task and stores it in the database
func AddTask(title, description, category string, dueDate time.Time) (types.Task, error) {
	task := types.Task{
		ID:          uuid.New().String()[:8],
		Title:       title,
		Description: description,
		Category:    category,
		DueDate:     dueDate,
		Status:      "pending",
	}

	db := storage.DB()
	if db == nil {
		return task, errors.New("failed to get task database")
	}

	return task, db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		data, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(task.ID), data)
	})
}

func AddNote(taskID, note, noteName string) error {
	db := storage.DB()
	if db == nil {
		return errors.New("failed to get task database")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		if b == nil {
			return errors.New("failed to get task database")
		}
		data := b.Get([]byte(taskID))
		var t types.Task
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		if t.Notes == nil {
			t.Notes = make(map[string]string)
		}
		t.Notes[noteName] = note

		// put back into json and put back into db
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(taskID), data)
	})
}

// GetTask retrieves a task by ID
func GetTask(id string) (*types.Task, error) {
	db := storage.DB()
	if db == nil {
		return nil, errors.New("failed to get task database")
	}
	var task types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
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

func GetTasks(ids []string) ([]types.Task, error) {
	db := storage.DB()
	if db == nil {
		return nil, errors.New("failed to get task database")
	}

	var tasks []types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		if b == nil {
			return errors.New("task bucket not found")
		}
		for _, id := range ids {
			data := b.Get([]byte(id))
			if data == nil {
				return errors.New("task not found")
			}
			var t types.Task
			if err := json.Unmarshal(data, &t); err != nil {
				return err
			}
			tasks = append(tasks, t)
		}
		return nil
	})

	return tasks, err
}

func GetAllTasks() ([]types.Task, error) {
	db := storage.DB()
	if db == nil {
		return []types.Task{}, errors.New("failed to get task database")
	}
	var tasks []types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
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

func DeleteTask(id string) error {
	db := storage.DB()
	if db == nil {
		return errors.New("failed to get task database")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		return b.Delete([]byte(id))
	})
}

func DeleteAllTasks() error {
	db := storage.DB()
	if db == nil {
		return errors.New("failed to get task database")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
	})
}

// DisplayTasks prints a list of tasks in a formatted table
func PrintListOfTasks(tasks []types.Task) {
	// Define headers and their width
	headers := []string{"ID", "Title", "Category", "Due Date", "Status", "Pr."}
	colWidths := []int{4, 15, 8, 12, 10, 3} // Width for each column

	// Set up gray color for borders
	borderColor := color.New(color.FgHiBlack)
	totalWidth := sum(colWidths) + len(colWidths)*3 - 1

	// Create top border, header separator, and bottom border with lighter color
	topBorder := borderColor.Sprintf("┌%s┐\n", strings.Repeat("─", totalWidth))
	headerSeparator := borderColor.Sprintf("├%s┤\n", strings.Repeat("─", totalWidth))
	bottomBorder := borderColor.Sprintf("└%s┘\n", strings.Repeat("─", totalWidth))

	// Print the top border
	fmt.Print(topBorder)

	// Print the headers without vertical separators
	fmt.Print(borderColor.Sprintf("│"))
	for i, header := range headers {
		fmt.Printf(" %-*s", colWidths[i], header)
		if i < len(headers)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Print(borderColor.Sprintf(" │\n") + headerSeparator)

	// Print each task row
	for _, task := range tasks {
		fmt.Print(borderColor.Sprintf("│"))
		values := []string{
			task.ID[:4],
			util.Truncate(task.Title, colWidths[1]),
			util.Truncate(task.Category, colWidths[2]),
			formatDate(task.DueDate),
			task.Status,
			fmt.Sprintf("%d", task.Priority),
		}
		for i, value := range values {
			fmt.Printf(" %-*s", colWidths[i], value)
			if i < len(values)-1 {
				fmt.Print("  ")
			}
		}
		fmt.Print(borderColor.Sprintf(" │\n"))
	}

	// Print the bottom border
	fmt.Print(bottomBorder)
}

// sum calculates the total width of the columns and adds padding for borders
func sum(colWidths []int) int {
	total := 0
	for _, width := range colWidths {
		total += width
	}
	return total
}

// formatDate formats the given date to M-D. If year is not current, also shows year at the end in parentheses.
func formatDate(date time.Time) string {
	out := date.Format("1-2")
	if date.Year() != time.Now().Year() {
		out += fmt.Sprintf(" (%s)", date.Format("2006"))
	}
	return out
}

// FindTasksByIDPrefix finds a list of potential ID matches for a given ID prefix string.
func FindTasksByIDPrefix(prefix string) ([]string, error) {
	var matchingIDs []string
	if len(prefix) > 8 {
		return matchingIDs, errors.New("given ID prefix is too long")
	}

	db := storage.DB()
	if db == nil {
		return []string{}, errors.New("failed to get tasks db")
	}

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TASK_DB))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			id := string(k)
			if len(prefix) == len(id) {
				if prefix == id {
					matchingIDs = append(matchingIDs, id)
				}
				return nil
			}

			if id[:len(prefix)] == prefix {
				matchingIDs = append(matchingIDs, id)
			}
			return nil
		})
	})
	return matchingIDs, err
}

func ListPotentialTaskMatches(taskIDs []string) {
	db := storage.DB()
	if db == nil {
		fmt.Println()
	}
	tasks, err := GetTasks(taskIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, t := range tasks {
		taskTitle := t.Title
		if len(taskTitle) > 15 {
			taskTitle = taskTitle[:12] + "..."
		}
		fmt.Printf("%s (%s)", t.ID, taskTitle)
	}
}
