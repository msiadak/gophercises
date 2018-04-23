package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/msiadak/gophercises/task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.Join(args, " ")
		err := task.Add(name)
		if err != nil {
			log.Fatalf("Couldn't add task: %s\n", err)
		}

		fmt.Printf("Added task '%s'\n", name)
	},
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
