/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/types"
)

var (
	description string
	category    string
	dueDate     string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long: `Add a new task to the list of ongoing tasks.

Example usage:

# Add a simple quick task that only has a title
task add "learn how to use cobra"

# Add a task with a title, description and due date M/D(/YYYY)
task add "work on task app" -d "finish the usage examples" -D 12/5/2024

# Add a task that is due in 2 days (d=days, w=weeks, m=months, y=years)
task add "get this done next week" -D 2d

the "title" argument is required, but all other arguments are optional. If no due date is provided, it defaults to today.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		// Parse the due date
		due, err := parseDueDate(dueDate)
		if err != nil {
			fmt.Println("Error parsing due date:", err)
			return
		}

		t, err := tasks.AddTask(title, description, category, due)
		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}
		tasks.PrintListOfTasks([]types.Task{t})
	},
}

func init() {
	addCmd.Flags().StringVarP(&description, "description", "d", "", "a description of the task")
	addCmd.Flags().StringVarP(&category, "category", "c", "", "a category for the task")
	addCmd.Flags().StringVarP(&dueDate, "due-date", "D", "", "the due date for the task")

	rootCmd.AddCommand(addCmd)
}

// parseDueDate parses the due date string and returns a time.Time
func parseDueDate(dueDate string) (time.Time, error) {
	if dueDate == "" {
		return time.Now(), nil
	}
	// check if the due date is a precise date (i.e. uses a slash delimiter)
	if strings.Contains(dueDate, "/") {
		// if year isn't defined, default to the current year and add it on
		parts := strings.Split(dueDate, "/")
		if len(parts) < 2 || len(parts) > 3 {
			return time.Time{}, fmt.Errorf("invalid date format")
		}
		if len(parts) == 2 {
			dueDate += "/" + time.Now().Format("2006")
		}
		return time.Parse("1/2/2006", dueDate)
	}
	// if its not a full date, then it should be a relative date
	// check if the format is correct (number followed by a letter)
	if len(dueDate) < 2 {
		return time.Time{}, fmt.Errorf("invalid date format")
	}
	// get the number and the unit and calculate the due date
	number, err := strconv.Atoi(dueDate[:len(dueDate)-1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid format for relative due date")
	}
	unit := dueDate[len(dueDate)-1:]
	today := time.Now()
	switch unit {
	case "d":
		return today.AddDate(0, 0, number), nil
	case "w":
		return today.AddDate(0, 0, number*7), nil
	case "m":
		return today.AddDate(0, number, 0), nil
	case "y":
		return today.AddDate(number, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid unit for relative due date")
	}
}
