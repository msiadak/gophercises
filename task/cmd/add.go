package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		task := []byte(strings.Join(args, " "))

		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			fmt.Printf("Couldn't open db: '%s'\n%s\n", "tasks.db", err)
			os.Exit(1)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
			if err != nil {
				return fmt.Errorf("create bucket: %s\n", err)
			}

			id, _ := b.NextSequence()

			return b.Put(itob(int(id)), task)
		})
		if err != nil {
			fmt.Printf("Couldn't add task: %s\n%s\n", task, err)
			os.Exit(1)
		}

		fmt.Printf("Added '%s' to tasks\n", task)
	},
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
