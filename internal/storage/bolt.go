package storage

import "go.etcd.io/bbolt"

const (
	TASK_DB = "tasks"
)

var db *bbolt.DB

// OpenDatabase opens the BoltDB database at the given path
func OpenDatabase(path string) error {
	var err error
	db, err = bbolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}

	// Ensure the tasks bucket exists
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
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
