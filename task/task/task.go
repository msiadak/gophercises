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

// CloseDB closes the database.
func CloseDB() {
	db.Close()
}

// InitDB opens the database and creates the tasks bucket if it doesn't
// already exist.
func InitDB() error {
	err := openDB(defaultDBPath())
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

// Task represents something that needs doing and the time it was entered
// into the system and completed.
type Task struct {
	ID      uint64
	Name    string
	Created time.Time
	Done    time.Time
}

func (t *Task) idAsBytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, t.ID)
	return b
}

// Do marks a task as done in the database.
func (t *Task) Do() error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		t.Done = time.Now()

		rec, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(t.idAsBytes(), rec)
	})
}

// Rm removes a task from the database.
func (t *Task) Rm() error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.Delete(t.idAsBytes())
	})
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
		fmt.Println(err)
		return nil, err
	}
	return tasks, nil
}

// Add returns a new Task instance and saves it to the database.
func Add(name string) error {
	t := Task{Name: name, Created: time.Now()}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))

		id, _ := b.NextSequence()
		t.ID = id

		buf, err := json.Marshal(&t)
		if err != nil {
			return err
		}

		return b.Put(t.idAsBytes(), buf)
	})
	if err != nil {
		return err
	}

	return nil
}

// ListIncomplete returns a slice of the tasks that are incomplete (e.g. tasks
// that do not have a Done time)
func ListIncomplete() ([]Task, error) {
	tasks, err := list(func(t *Task) bool {
		return t.Done.IsZero()
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// ListCompletedToday returns a slice of the tasks that have been completed
// today.
func ListCompletedToday() ([]Task, error) {
	tasks, err := list(func(t *Task) bool {
		f := "2006 Jan _2"
		return t.Done.Format(f) == time.Now().Format(f)
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
