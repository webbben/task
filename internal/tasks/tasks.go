package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/webbben/task/internal/constants"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"github.com/webbben/task/internal/util"
	"go.etcd.io/bbolt"
)

const (
	colID       = "ID"
	colTitle    = "Title"
	colCategory = "Cat."
	colDueDate  = "Due Date"
	colStatus   = "Status"
	colPriority = "Pr."
)

var (
	// colors for due dates
	veryLate = color.New(color.BgRed)
	late     = color.New(color.FgRed)
	today    = color.New(color.FgHiYellow)
	tomorrow = color.New(color.FgCyan)

	// columns that will be displayed in the table
	headers = []string{colID, colTitle, colDueDate, colStatus, colPriority}

	// widths for each column type
	colWidths = map[string]int{
		colID:       8,
		colTitle:    18,
		colCategory: 6,
		colDueDate:  10,
		colStatus:   8,
		colPriority: 3,
	}
)

// generates an ID that starts with the first 4 characters of the task title
func GenerateTaskID(title string) string {
	formatted := strings.ToLower(strings.ReplaceAll(title, " ", ""))
	genID := uuid.New().String()
	if len(formatted) > 4 {
		formatted = formatted[:4]
	}
	out := formatted + genID
	return out[:8]
}

// AddTask creates a new task and stores it in the database
func AddTask(title, description, category string, dueDate time.Time) (types.Task, error) {
	task := types.Task{
		ID:          GenerateTaskID(title),
		Title:       title,
		Description: description,
		Category:    category,
		DueDate:     dueDate,
		Status:      constants.TaskStatus.Pending,
		LastUpdate:  time.Now(),
	}

	db := storage.DB()
	if db == nil {
		return task, errors.New("failed to get task database")
	}

	return task, db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
		t.LastUpdate = time.Now()
		// advance to in progress if its still pending
		if t.Status == constants.TaskStatus.Pending {
			t.Status = constants.TaskStatus.InProgress
		}

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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
		if b == nil {
			// bucket doesn't exist yet
			return nil
		}
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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
		return b.Delete([]byte(id))
	})
}

func DeleteAllTasks() error {
	db := storage.DB()
	if db == nil {
		return errors.New("failed to get task database")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
	// Set up gray color for borders
	borderColor := color.New(color.FgHiBlack)
	totalWidth := totalTableWidth() + len(headers)*3 - 1

	// Create top border, header separator, and bottom border with lighter color
	topBorder := borderColor.Sprintf("┌%s┐\n", strings.Repeat("─", totalWidth))
	headerSeparator := borderColor.Sprintf("├%s┤\n", strings.Repeat("─", totalWidth))
	bottomBorder := borderColor.Sprintf("└%s┘\n", strings.Repeat("─", totalWidth))

	// Print the top border
	fmt.Print(topBorder)

	// Print the headers without vertical separators
	fmt.Print(borderColor.Sprintf("│"))
	for i, header := range headers {
		fmt.Printf(" %-*s", colWidths[header], header)
		if i < len(headers)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Print(borderColor.Sprintf(" │\n") + headerSeparator)

	// Print each task row
	for _, task := range tasks {
		fmt.Print(borderColor.Sprintf("│"))
		for i, header := range headers {
			var value string
			switch header {
			case colID:
				value = task.ID
			case colTitle:
				value = task.Title
			case colCategory:
				value = task.Category
			case colDueDate:
				value = formatDate(task.DueDate)
			case colStatus:
				value = task.Status
			case colPriority:
				value = fmt.Sprintf("%d", task.Priority)
			default:
				value = "?"
			}
			value = util.Truncate(value, colWidths[header])

			// if colors were used, we may need to add extra padding due to invisible ansi stuff
			value = addPadding(value, colWidths[header])
			fmt.Printf(" %-*s", colWidths[header], value)
			if i < len(headers)-1 {
				fmt.Print("  ")
			}
		}
		fmt.Print(borderColor.Sprintf(" │\n"))
	}

	// Print the bottom border
	fmt.Print(bottomBorder)
}

// sum calculates the total width of the columns and adds padding for borders
func totalTableWidth() int {
	total := 0
	for _, header := range headers {
		total += colWidths[header]
	}
	return total
}

// formatDate formats the given date to M-D. If year is not current, also shows year at the end in parentheses.
func formatDate(date time.Time) string {
	// Create a new time for the end of today
	now := time.Now()
	t := time.Date(
		now.Year(), now.Month(), now.Day(),
		23, 59, 59, 999999999, now.Location())

	dateYesterday := t.AddDate(0, 0, -1)
	dateTwoDaysAgo := t.AddDate(0, 0, -2)
	dateTomorrow := t.AddDate(0, 0, 1)

	out := date.Format("1-2")
	if date.Year() != t.Year() {
		out += "-" + date.Format("2006")
	}

	switch {
	case date.Before(dateTwoDaysAgo):
		out = veryLate.Sprint(out)
	case date.Before(dateYesterday):
		out = late.Sprint(out)
	case date.Before(t):
		out = today.Sprint(out)
	case date.Before(dateTomorrow):
		out = tomorrow.Sprint(out)
	}

	return out
}

func addPadding(s string, colWidth int) string {
	visibleLength := len(util.StripAnsi(s))
	padding := colWidth - visibleLength
	if padding > 0 {
		s += strings.Repeat(" ", padding)
	}
	return s
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
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
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
