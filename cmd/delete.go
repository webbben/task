package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/util"
)

var (
	all bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes tasks",
	Long: `Deletes tasks from the ongoing tasks database.
	
Example usage:

# delete a specific task
task delete <task-id>

# delete all tasks
task delete -a`,
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			// confirm deletion first
			if util.Confirm("Delete all tasks?") {
				tasks.DeleteAllTasks()
			}
			return
		}
		if len(args) == 0 {
			cmd.PrintErrln("task ID or --all flag is required")
			return
		}
		if err := tasks.DeleteTask(args[0]); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&all, "all", "a", false, "delete all tasks")
}
