/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/util"
)

// compCmd represents the comp command
var compCmd = &cobra.Command{
	Use:   "comp",
	Short: "Marks the given task as completed",
	Long: `Marks the given task as completed. The task is archived and removed from the active task list.
	
Example usage:

# mark task with ID beginning with 9bc3 as completed
task comp 9bc3`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("task ID required")
			return
		}
		taskID := args[0]

		if err := tasks.CompleteTask(taskID); err != nil {
			cmd.PrintErrln(err)
			return
		}

		todaysCompTasks, err := tasks.GetCompletedTasks(util.RoundDateDown(time.Now()))
		if err != nil {
			cmd.PrintErrln("task completed, but failed to get list of completed tasks today: ", err)
			return
		}
		tasks.PrintListOfTasks(todaysCompTasks)
	},
}

func init() {
	compCmd.ValidArgsFunction = completionFn
	rootCmd.AddCommand(compCmd)
}

func completionFn(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return taskIDCompletion(cmd, args, toComplete)
}
