package cmd

import (
	"fmt"
	"log"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the list of unfinished tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := task.ListIncomplete()
		if err != nil {
			log.Fatalf("Couldn't retrieve incomplete tasks: %s\n", err)
		}

		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i+1, t.Name)
		}
	},
}
