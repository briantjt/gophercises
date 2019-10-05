package db

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var Db *bolt.DB
var taskBucket = []byte("tasks")

type Task struct {
	Key   int
	Value string
}

func Init(dbDirectory string) {
	var err error
	Db, err = bolt.Open(filepath.Join(dbDirectory, "my.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return fmt.Errorf("Create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func CreateTask(task string) (int, error) {
	var id int
	err := Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}
	return id, nil
}

func ListTasks() ([]Task, error) {
	var tasks []Task
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func FindTaskID(index int) (int, error) {
	var key []byte
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		key, _ = c.First()
		for i := 0; i < index; i++ {
			key, _ = c.Next()
		}
		if key == nil {
			return errors.New("Could not find task with given id")
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return btoi(key), nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
