package cmd

import (
	"log"
	"strconv"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task complete",
	Run: func(cmd *cobra.Command, args []string) {
		d, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Couldn't parse arg '%s', please provide an integer\n", args[0])
		}

		tasks, err := task.ListIncomplete()
		if err != nil {
			log.Fatalf("Couldn't retrieve incomplete tasks: %s\n", args[0])
		}

		for i, t := range tasks {
			if i == d {
				t.Do()
			}
		}
	},
}
