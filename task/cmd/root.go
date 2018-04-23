package cmd

import (
	"fmt"
	"os"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI task manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		task.InitDB()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		task.CloseDB()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
