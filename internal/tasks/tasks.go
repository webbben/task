package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/webbben/task/internal/constants"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"github.com/webbben/task/internal/util"
	"go.etcd.io/bbolt"
)

const (
	colID         = "ID"
	colTitle      = "Title"
	colCategory   = "Cat."
	colDueDate    = "Due Date"
	colStatus     = "Status"
	colPriority   = "Pr."
	colLastUpdate = "Upd."
)

var (
	// colors for due dates
	veryLate = color.New(color.BgRed)
	late     = color.New(color.FgHiRed)
	today    = color.New(color.FgHiYellow)
	tomorrow = color.New(color.FgCyan)

	// status colors
	comp = color.New(color.BgGreen, color.FgBlack)
	prog = color.New(color.FgCyan)

	// columns that will be displayed in the table
	headers = []string{colID, colTitle, colDueDate, colStatus, colPriority, colLastUpdate}

	// widths for each column type
	colWidths = map[string]int{
		colID:         8,
		colTitle:      18,
		colCategory:   6,
		colDueDate:    10,
		colStatus:     8,
		colPriority:   3,
		colLastUpdate: 4,
	}
)

// generates an ID that includes the task title as much as possible.
// doing this so it's easy to predict the task ID when typing them in the CLI, since it's a lot
// easier than memorizing completely randomized IDs.
//
// If a given title is already taken (e.g. you use the same or similar title for tasks often), then
// the task ID will progressively become more full of random digits, starting from the end.
//
// pass in an empty string to get a random ID that isn't based on any title text.
func GenerateTaskID(title string) string {
	maxIDLen := 12
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	formatted := strings.ToLower(re.ReplaceAllString(title, ""))
	genID := strings.ReplaceAll(uuid.New().String(), "-", "")
	out := (formatted + genID)[:maxIDLen]
	for {
		inUse, err := idAlreadyUsed(out)
		if err != nil {
			log.Println("error occurred during ID generation:", err)
			log.Println("proceeding with random ID")
			return genID[:maxIDLen]
		}
		if !inUse {
			return out
		}
		if len(formatted) == 0 {
			return genID[:maxIDLen]
		}
		// keep reducing the title and replacing it with random digits until an unused ID is found
		formatted = formatted[:len(formatted)-1]
		out = (formatted + genID)[:maxIDLen]
	}
}

func idAlreadyUsed(id string) (bool, error) {
	db := storage.DB()
	if db == nil {
		return true, errors.New("failed to get task database")
	}

	inUse := false
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
		if b == nil {
			return errors.New("failed to get active tasks bucket")
		}
		if b.Get([]byte(id)) != nil {
			inUse = true
		}
		return nil
	})
	return inUse, err
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
			t, err := unpackTaskJson(data)
			if err != nil {
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
			task, err := unpackTaskJson(v)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	return tasks, err
}

func unpackTaskJson(v []byte) (types.Task, error) {
	var task types.Task
	if err := json.Unmarshal(v, &task); err != nil {
		return task, err
	}
	// handle any data updates or corrections
	// e.g. a new property was added and doesn't exist on a task yet, so add a default
	if task.LastUpdate.IsZero() {
		task.LastUpdate = time.Now()
	}
	return task, nil
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
	totalWidth, _, err := term.GetSize(os.Stdin.Fd())
	if err != nil {
		log.Println("failed to get terminal size:", err)
		totalWidth = 80
	}
	totalWidth -= 2 // make space for the borders on the sides
	titleWidth := totalWidth - baseTableWidth()
	if titleWidth < colWidths[colTitle] {
		// terminal is too small for the current configuration; consider removing columns
		overflow := colWidths[colTitle] - titleWidth
		removeHeader(colPriority) // hide priority col
		if overflow >= 7 {
			// also hide last update col as a final measure
			removeHeader(colLastUpdate)
		}
		// recalculate title width after column removals
		titleWidth = totalWidth - baseTableWidth()
	}
	colWidths[colTitle] = titleWidth

	// Create top border, header separator, and bottom border with lighter color
	borderColor := color.New(color.FgHiBlack) // Set up gray color for borders
	topBorder := borderColor.Sprintf("┌%s┐\n", strings.Repeat("─", totalWidth))
	headerSeparator := borderColor.Sprintf("├%s┤\n", strings.Repeat("─", totalWidth))
	bottomBorder := borderColor.Sprintf("└%s┘\n", strings.Repeat("─", totalWidth))

	// Print the top border
	fmt.Print(topBorder)

	// Print the headers without vertical separators
	fmt.Print(borderColor.Sprintf("│"))
	for i, header := range headers {
		if header == "X" {
			// ignore removed headers
			continue
		}
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
				value = task.ID[:8]
			case colTitle:
				value = task.Title
			case colCategory:
				value = task.Category
			case colDueDate:
				value = formatDate(task.DueDate, task.Status == constants.TaskStatus.Complete)
			case colStatus:
				value = formatStatus(task.Status)
			case colPriority:
				value = fmt.Sprintf("%d", task.Priority)
			case colLastUpdate:
				value = timeSinceDateFormat(task.LastUpdate)
			case "X":
				continue // deleted header due to terminal being too small
			default:
				value = "?" // unknown header
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

// sum calculates the total width of the fixed-size columns (i.e. those besides the title)
// including padding
func baseTableWidth() int {
	total := 2 // left border of table
	for i, header := range headers {
		// exclude title from calc. also exclude any removed headers (if terminal is too small, some may be removed)
		if header == colTitle || header == "X" {
			continue
		}
		total += colWidths[header] + 1
		if i < len(headers)-1 {
			total += 2 // gap between each column
		}
	}
	return total + 2 // right border of table
}

// formatDate formats the given date to M-D. If year is not current, also shows year at the end in parentheses.
func formatDate(date time.Time, skipColor bool) string {
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

	if !skipColor {
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
	}

	return out
}

// timeSinceDateFormat returns the number of days, weeks, or months since the given date.
//
// the string is formatted as a number followed by a letter which represents the unit ("d", "w", or "m").
func timeSinceDateFormat(date time.Time) string {
	duration := time.Since(date)

	days := int(duration.Hours() / 24)
	if days < 14 {
		return fmt.Sprintf("%vd", days)
	}
	weeks := days / 7
	if weeks < 8 {
		return fmt.Sprintf("%vw", weeks)
	}
	months := days / 30
	return fmt.Sprintf("%vm", months)
}

func formatStatus(status int) string {
	out := constants.TaskStatusDisplay[status]
	if status == constants.TaskStatus.Complete {
		out = comp.Sprint(out)
	} else if status == constants.TaskStatus.InProgress {
		out = prog.Sprint(out)
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

func removeHeader(h string) {
	for i := 0; i < len(headers); i++ {
		if headers[i] == h {
			headers[i] = "X"
		}
	}
}
