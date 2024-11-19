package storage

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

const (
	TASK_DB        = "tasks.db"
	ACTIVE_BUCKET  = "active"
	ARCHIVE_BUCKET = "archive"
)

func ConfigPathUnix() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("failed to get config path:", err)
	}
	path := filepath.Join(usr.HomeDir, ".config/task")
	if err := ensureDir(path); err != nil {
		log.Fatal(err)
	}
	return path
}

func AppDataPathUnix() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("failed to get config path:", err)
	}
	path := filepath.Join(usr.HomeDir, ".local/share/task")
	if err := ensureDir(path); err != nil {
		log.Fatal(err)
	}
	return path
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// if it doesn't exist, create it
			return os.MkdirAll(path, os.ModePerm)
		}
		return err
	}
	return nil
}

// OpenDatabase opens the BoltDB database of the given name
func OpenDatabase(name string) error {
	fullpath := filepath.Join(AppDataPathUnix(), name)
	var err error
	db, err = bbolt.Open(fullpath, 0600, nil)
	if err != nil {
		return err
	}

	// Ensure the tasks bucket exists
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ACTIVE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ARCHIVE_BUCKET))
		return err
	})
}

// CloseDatabase closes the BoltDB database
func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}

// DB returns the BoltDB database
func DB() *bbolt.DB {
	return db
}
