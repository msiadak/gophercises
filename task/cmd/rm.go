package cmd

import (
	"log"
	"strconv"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use: "rm",
	Run: func(cmd *cobra.Command, args []string) {
		taskNum, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Please specify a task number to remove.")
		}

		tasks, err := task.ListIncomplete()
		if err != nil {
			log.Fatalf("Couldn't retrieve list of incomplete tasks: %s\n", err)
		}

		if taskNum > len(tasks) {
			log.Fatalln("Please specify a valid task number.")
		}

		for i, t := range tasks {
			if i+1 == taskNum {
				err := t.Rm()
				if err != nil {
					log.Fatalf("Couldn't remove task %d: %s\n", taskNum, err)
				}
				break
			}
		}
	},
}
