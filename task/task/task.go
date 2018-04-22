package task

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/boltdb/bolt"

	"github.com/mitchellh/go-homedir"
)

const taskBucket = "tasks"

var db *bolt.DB

// defaultDBPath returns the default boltdb file path ($HOME/.tasks.db)
func defaultDBPath() string {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Couldn't determine home dir: %s", err)
	}

	return path.Join(h, ".tasks.db")
}

func openDB(dbPath string) error {
	var err error

	o := bolt.Options{Timeout: 3 * time.Second}

	db, err = bolt.Open(dbPath, 0600, &o)
	if err != nil {
		return fmt.Errorf("Couldn't open db: %s", err)
	}

	return nil
}

func CloseDB() {
	db.Close()
}

// InitDB opens the database and creates the tasks bucket if it doesn't
// already exist.
func InitDB(dbPath string) error {
	err := openDB(dbPath)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(taskBucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Couldn't create tasks bucket: %s", err)
	}

	return nil
}

type jsonTime time.Time

func (jt jsonTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(jt).Format(time.RFC3339))
	return []byte(stamp), nil
}

func (jt *jsonTime) UnmarshalJSON(v []byte) error {
	t, err := time.Parse(time.RFC3339, string(v))
	if err != nil {
		return err
	}

	*jt = jsonTime(t)
	return nil
}

// Task represents something that needs doing and the time it was entered
// into the system and completed.
type Task struct {
	ID      int
	Name    string
	Created jsonTime
	Done    jsonTime
}

type taskPredicate func(t *Task) bool

// list returns a slice of tasks that satisfy the given predicate function
func list(fn taskPredicate) ([]Task, error) {
	tasks := make([]Task, 0)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := new(Task)
			err := json.Unmarshal(v, t)
			if err != nil {
				return err
			}

			if fn(t) {
				tasks = append(tasks, *t)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func Create(name string) *Task {
	t := Task{Name: name, Created: jsonTime(time.Now())}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))

		id, _ := b.NextSequence()

		buf, err := json.Marshal(&t)
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), buf)
	})

	return t
}

// ListIncomplete returns a slice of the tasks that are incomplete (e.g. tasks
// that do not have a Done time)
func ListIncomplete() ([]Task, error) {
	tasks, err := list(func(t *Task) bool {
		return time.Time(t.Done).IsZero()
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// ListCompleted returns a slice of the tasks that have been completed (e.g.
// tasks that have a Done time)
func ListCompleted() ([]Task, error) {
	tasks, err := list(func(t *Task) bool {
		return !time.Time(t.Done).IsZero()
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func MarkDone(taskNum uint) error {

}
