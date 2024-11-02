package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"github.com/webbben/task/internal/util"
)

// AddTask creates a new task and stores it in the database
func AddTask(title, description, category string, dueDate time.Time) (types.Task, error) {
	task := types.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Category:    category,
		DueDate:     dueDate,
		Status:      "pending",
	}
	return task, storage.AddTaskToDB(task)
}

// GetTask retrieves a task by ID
func GetTask(id string) (*types.Task, error) {
	return storage.GetTaskFromDB(id)
}

func GetAllTasks() ([]types.Task, error) {
	return storage.GetAllTasksFromDB()
}

// DisplayTasks prints a list of tasks in a formatted table
func PrintListOfTasks(tasks []types.Task) {
	// Define headers and their width
	headers := []string{"ID", "Title", "Category", "Due Date", "Status", "Pr."}
	colWidths := []int{4, 15, 8, 12, 10, 3} // Width for each column

	// Create table borders and separators
	topBorder := "┌" + strings.Repeat("─", colWidths[0]+2) + "┬" +
		strings.Repeat("─", colWidths[1]+2) + "┬" +
		strings.Repeat("─", colWidths[2]+2) + "┬" +
		strings.Repeat("─", colWidths[3]+2) + "┬" +
		strings.Repeat("─", colWidths[4]+2) + "┬" +
		strings.Repeat("─", colWidths[5]+2) + "┐\n"

	headerSeparator := "├" + strings.Repeat("─", colWidths[0]+2) + "┼" +
		strings.Repeat("─", colWidths[1]+2) + "┼" +
		strings.Repeat("─", colWidths[2]+2) + "┼" +
		strings.Repeat("─", colWidths[3]+2) + "┼" +
		strings.Repeat("─", colWidths[4]+2) + "┼" +
		strings.Repeat("─", colWidths[5]+2) + "┤\n"

	bottomBorder := "└" + strings.Repeat("─", colWidths[0]+2) + "┴" +
		strings.Repeat("─", colWidths[1]+2) + "┴" +
		strings.Repeat("─", colWidths[2]+2) + "┴" +
		strings.Repeat("─", colWidths[3]+2) + "┴" +
		strings.Repeat("─", colWidths[4]+2) + "┴" +
		strings.Repeat("─", colWidths[5]+2) + "┘\n"

	// Print the top border
	fmt.Print(topBorder)

	// Print the headers
	fmt.Print("│")
	for i, header := range headers {
		fmt.Printf(" %-*s │", colWidths[i], header)
	}
	fmt.Print("\n" + headerSeparator)

	// Print each task row
	for _, task := range tasks {
		fmt.Print("│")
		values := []string{
			task.ID[:4],
			util.Truncate(task.Title, colWidths[1]),
			util.Truncate(task.Category, colWidths[2]),
			task.DueDate.Format("2006-01-02"),
			task.Status,
			fmt.Sprintf("%d", task.Priority),
		}
		for i, value := range values {
			fmt.Printf(" %-*s │", colWidths[i], value)
		}
		fmt.Println()
	}

	// Print the bottom border
	fmt.Print(bottomBorder)
}
