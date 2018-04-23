package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List the tasks completed today",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := task.ListCompletedToday()
		if err != nil {
			log.Fatalf("Couldn't list completed tasks: %s\n", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks have been completed yet today.")
			return
		}

		fmt.Printf("Tasks completed today (%s):\n", time.Now().Format("Mon Jan 2 2006"))
		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i+1, t.Name)
		}
	},
}
