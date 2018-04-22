package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/msiadak/gophercises/task/util"
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

		dbPath, err := util.DefaultDBPath()
		if err != nil {
			log.Fatalf("Couldn't determine path to DB file\n%s\n", err)
		}

		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatalf("Couldn't open db: '%s'\n%s\n", dbPath, err)
		}
		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))
			if b == nil {
				return fmt.Errorf("Add a task using '%s %s' before trying to mark one done", rootCmd.Use, addCmd.Use)
			}

			n := b.Stats().KeyN
			if n < d {
				return fmt.Errorf("Couldn't delete task %d -- list only has %d tasks", d, n)
			}

			c := b.Cursor()

			i := 1
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if i == d {
					b.Delete(k)
					fmt.Printf("Marked '%s' done", v)
					return nil
				}
			}

			return fmt.Errorf("Unexpected error -- unable to mark task %d complete", d)
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}
