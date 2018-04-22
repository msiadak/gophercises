package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/msiadak/gophercises/task/util"
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

		dbPath, err := util.DefaultDBPath()
		if err != nil {
			log.Fatalf("couldn't determine path to DB file\n%s\n", err)
		}

		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatalf("couldn't open db: '%s'\n%s\n", dbPath, err)
		}
		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			id, _ := b.NextSequence()

			return b.Put(itob(int(id)), task)
		})
		if err != nil {
			log.Fatalf("couldn't add task: %s\n%s\n", task, err)
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
