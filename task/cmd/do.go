package cmd

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task complete",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			fmt.Printf("Couldn't open db: '%s'\n%s\n", "tasks.db", err)
			os.Exit(1)
		}

		db.Update(func(tx *bolt.Tx) error {
			return nil
		})
	},
}
