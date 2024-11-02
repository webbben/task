package types

import (
	"fmt"
	"time"

	"github.com/webbben/task/internal/util"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"` // "pending" or "done"
	Priority    int       `json:"priority"`
}

func (t Task) String() string {
	var result string

	headers := []string{"ID", "Title", "Desc", "Category", "Due Date", "Status"}
	values := []string{
		t.ID[:8],
		util.Truncate(t.Title, 10),
		util.Truncate(t.Description, 10),
		util.Truncate(t.Category, 10),
		t.DueDate.Format("1-2-2006"),
		t.Status,
	}

	// Create a border using Unicode box-drawing characters
	topBorder := "┌────────────┬────────────┬────────────┬────────────┬────────────┬────────────┐\n"
	headerSeparator := "├────────────┼────────────┼────────────┼────────────┼────────────┼────────────┤\n"
	bottomBorder := "└────────────┴────────────┴────────────┴────────────┴────────────┴────────────┘\n"

	// Start with the top border
	result += topBorder

	// Create header row with borders
	result += "│"
	for _, header := range headers {
		result += fmt.Sprintf(" %-10s │", header)
	}
	result += "\n" + headerSeparator

	// Create value row with borders
	result += "│"
	for _, value := range values {
		result += fmt.Sprintf(" %-10s │", value)
	}
	result += "\n" + bottomBorder

	return result
}
