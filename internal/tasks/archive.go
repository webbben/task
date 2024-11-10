package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/webbben/task/internal/constants"
	"github.com/webbben/task/internal/storage"
	"github.com/webbben/task/internal/types"
	"go.etcd.io/bbolt"
)

func CompleteTask(id string) error {
	db := storage.DB()
	if db == nil {
		return errors.New("failed to get task database")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		// first, find the task in the active bucket
		b := tx.Bucket([]byte(storage.ACTIVE_BUCKET))
		if b == nil {
			return errors.New("active bucket does not exist")
		}
		taskData := b.Get([]byte(id))
		if taskData == nil {
			return errors.New("failed to get task from active bucket")
		}
		// delete it from the active bucket
		err := b.Delete([]byte(id))
		if err != nil {
			return fmt.Errorf("failed to delete task from active bucket: %s", err.Error())
		}
		// set the status to complete
		taskData, err = setTaskDataComplete(taskData)
		if err != nil {
			return err
		}
		// set the taskData in the archived bucket for the current month and year sub-bucket
		bucketName := monthBucketName(time.Now())
		monthBucket, err := getArchiveBucket(bucketName, tx)
		if err != nil {
			return err
		}
		return monthBucket.Put([]byte(id), taskData)
	})
}

func getArchiveBucket(bucketName string, tx *bbolt.Tx) (*bbolt.Bucket, error) {
	archiveBucket := tx.Bucket([]byte(storage.ARCHIVE_BUCKET))
	if archiveBucket == nil {
		return nil, fmt.Errorf("no bucket found: %s", storage.ARCHIVE_BUCKET)
	}
	monthBucket, err := archiveBucket.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return nil, fmt.Errorf("failed to get or create archive bucket: %s", bucketName)
	}
	return monthBucket, nil
}

func setTaskDataComplete(taskData []byte) ([]byte, error) {
	var task types.Task
	err := json.Unmarshal(taskData, &task)
	if err != nil {
		return []byte{}, err
	}
	task.Status = constants.TaskStatus.Complete
	task.LastUpdate = time.Now()
	return json.Marshal(task)
}

func GetCompletedTasks(lookbackDate time.Time) ([]types.Task, error) {
	db := storage.DB()
	if db == nil {
		return []types.Task{}, errors.New("failed to get database")
	}

	// get all the month buckets going back to the lookbackDate
	// ex: if today is Nov 10th and lookback date is Sept 5th, the month buckets
	// will be: Nov, Oct, Sept
	buckets := make([]string, 0)
	curMonth := time.Now()
	for curMonth.Year() >= lookbackDate.Year() && curMonth.Month() >= lookbackDate.Month() {
		buckets = append(buckets, monthBucketName(curMonth))
		curMonth = curMonth.AddDate(0, -1, 0)
	}

	var tasks []types.Task
	err := db.View(func(tx *bbolt.Tx) error {
		archiveBucket := tx.Bucket([]byte(storage.ARCHIVE_BUCKET))
		if archiveBucket == nil {
			return errors.New("failed to get archive bucket")
		}
		for _, month := range buckets {
			monthBucket := archiveBucket.Bucket([]byte(month))
			// some month buckets may not exist if no tasks were completed in that month
			if monthBucket == nil {
				continue
			}

			err := monthBucket.ForEach(func(k, v []byte) error {
				var t types.Task
				err := json.Unmarshal(v, &t)
				if err != nil {
					return err
				}
				if t.LastUpdate.After(lookbackDate) {
					tasks = append(tasks, t)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return tasks, err
}

func monthBucketName(date time.Time) string {
	return date.Format("2006-01")
}
