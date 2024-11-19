package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/webbben/task/internal/completions"
	"github.com/webbben/task/internal/tasks"
	"github.com/webbben/task/internal/util"
)

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "create a new note update for a task",
	Long: `Create a new note for an existing task. You can specify a note, or leave it blank to launch an editor.

# add short note
task note 9bc3 "will follow-up next Monday"

# add a note that is composed in a terminal editor
task note 3bp4`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrln("task ID required")
			return
		}
		taskID := args[0]
		if len(taskID) > 8 {
			cmd.PrintErrln("invalid task ID")
			return
		}

		// if incomplete task ID, find a match if possible
		if len(taskID) < 8 {
			// find the full taskID
			matches, err := tasks.FindTasksByIDPrefix(taskID)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
			if len(matches) > 1 {
				tasks.ListPotentialTaskMatches(matches)
				return
			}
			taskID = matches[0]
		}

		// get note to add to task
		today := time.Now().Format("1-2-2006 15:04")
		note := ""
		if len(args) >= 2 {
			note = args[1]
		} else {
			// launch terminal text editor to get note from user
			note = util.OpenEditor()
		}
		if err := tasks.AddNote(taskID, note, today); err != nil {
			cmd.PrintErrln("Error adding note:", err)
			return
		}
		fmt.Printf("Added note \"%s\" to task %s: \n\"%s\"", today, taskID, note)
	},
}

func init() {
	noteCmd.ValidArgsFunction = completions.TaskIDCompletionFn(true)
	rootCmd.AddCommand(noteCmd)
}
