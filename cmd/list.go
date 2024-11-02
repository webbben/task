/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/types"
)

var (
	sortBy   string
	filterBy string
	limit    int
	todo     bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active tasks",
	Long: `List all the tasks that are currently in progress. Completed or deleted tasks are not shown.

Example usage:

# list all tasks
task list

# sort by a property (title, category, due date, priority, and status are supported)
task list -s duedate

# filter by a property value (status, category, priority, and due date are supported)
task list -f status=paused

# limit the number of results shown
task list -l 5

# sort the list to show the most important tasks for today (cannot be used with sort or filter)
task list -t
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// load all tasks
		tasks, err := tasks.GetAllTasks()
		if err != nil {
			cmd.PrintErrln("Error loading tasks:", err)
			return
		}

		// check for filtering
		// todo flag (-t) has priority over filter flag (-f) and sort flag (-s)
		if todo {
			showTodoTasks(tasks)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&sortBy, "sort", "s", "", "Sort the list by a property")
	listCmd.Flags().StringVarP(&filterBy, "filter", "f", "", "Filter the list by a property value")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of results shown")
	listCmd.Flags().BoolVarP(&todo, "todo", "t", false, "Show the most important tasks for today")
}

func showTodoTasks(t []types.Task) {
	// sort by due date, but for tasks that are the same due date, sort by priority
	sort.Slice(t, func(i, j int) bool {
		if t[i].DueDate.Equal(t[j].DueDate) {
			return t[i].Priority > t[j].Priority
		}
		return t[i].DueDate.Before(t[j].DueDate)
	})

	tasks.PrintListOfTasks(t)
}
