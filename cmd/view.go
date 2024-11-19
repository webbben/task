package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/completions"
	taskui "github.com/webbben/task/internal/ui/task-ui"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view the details of a single task",
	Long: `Launch a TUI application to view the details of a single task, such as description, notes, etc.
	
Example:

task view 9bf4`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("task ID required")
		}
		taskID := args[0]
		err := taskui.RunUI(taskID)
		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	viewCmd.ValidArgsFunction = completions.TaskIDCompletionFn(true)
	rootCmd.AddCommand(viewCmd)
}
