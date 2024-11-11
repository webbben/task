package cmd

import (
	"github.com/spf13/cobra"
	taskui "github.com/webbben/task/internal/ui/task-ui"
	"github.com/webbben/task/internal/util"
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
	viewCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// only do completion for the task ID arg, which is the first one
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return util.TaskIDCompletionFn(toComplete, cmd)
	}
	rootCmd.AddCommand(viewCmd)
}
