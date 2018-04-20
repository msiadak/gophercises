package cmd

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the list of unfinished tasks",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			fmt.Printf("Couldn't open db: '%s'\n%s\n", "tasks.db", err)
			os.Exit(1)
		}

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))

			c := b.Cursor()

			i := 1
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("%d. %s\n", i, v)
				i++
			}

			return nil
		})
		if err != nil {
			fmt.Printf("Couldn't create db view\n%s\n", err)
			os.Exit(1)
		}
	},
}
