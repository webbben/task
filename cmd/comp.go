package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/completions"
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
		// all arguments will be task IDs
		for i, taskID := range args {
			if err := tasks.CompleteTask(taskID); err != nil {
				cmd.PrintErrln(err)
				if i == 0 {
					return // no tasks were completed, so quit without showing summary
				}
			}
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
	compCmd.ValidArgsFunction = completions.TaskIDCompletionFn(false)
	rootCmd.AddCommand(compCmd)
}
