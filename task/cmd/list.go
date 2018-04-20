package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the list of unfinished tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1. Implement the list command")
	},
}
