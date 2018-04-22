package cmd

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/msiadak/gophercises/task/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the list of unfinished tasks",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := util.DefaultDBPath()
		if err != nil {
			fmt.Printf("Couldn't determine path to DB file\n%s\n", err)
			os.Exit(1)
		}

		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			fmt.Printf("Couldn't open db: '%s'\n%s\n", dbPath, err)
			os.Exit(1)
		}
		defer db.Close()

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))
			if b == nil {
				return fmt.Errorf("No tasks exist yet, add one with '%s %s' first", rootCmd.Use, addCmd.Use)
			}

			c := b.Cursor()

			i := 1
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("%d. %s\n", i, v)
				i++
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
