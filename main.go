package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

var dbBucket = []byte("last_usage_time")

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		//os.ErrExist = fmt.Errorf("Two arguments expected")
		logrus.Errorf("Two arguments expected")
		os.Exit(1)
	}

	db, err := bolt.Open("pr00b121-tdd.db", 0600, nil)
	if err != nil {
		os.ErrExist = err
		logrus.Errorf("Failure caused by: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	switch argsWithoutProg[0] {
	case "read":
		id := argsWithoutProg[1]
		var writeTime string
		err = db.View(func(tx *bolt.Tx) error {
			bk := tx.Bucket(dbBucket)

			bs := bk.Get([]byte(id))
			if bs == nil {
				return fmt.Errorf("not found id: %s", id)
			}

			writeTime = string(bs)
			return nil
		})
		if err != nil {
			os.ErrExist = err
			logrus.Errorf("Write action failure caused by: %v", err)
			os.Exit(1)
		}

		logrus.Infof("Read: id=%s, time=%s", id, writeTime)

	case "write":
		id := argsWithoutProg[1]
		writeTime := time.Now().Format(time.RFC3339)
		logrus.Infof("Write: id=%s, time=%s", id, writeTime)

		err = db.Update(func(tx *bolt.Tx) error {
			bk, err := tx.CreateBucketIfNotExists(dbBucket)
			if err != nil {
				return fmt.Errorf("failed to create bucket: %v", err)
			}

			if err := bk.Put([]byte(id), []byte(writeTime)); err != nil {
				return fmt.Errorf("failed to insert '%s': %v", writeTime, err)
			}
			return nil
		})
		if err != nil {
			os.ErrExist = err
			logrus.Errorf("Write action failure caused by: %v", err)
			os.Exit(1)
		}
	default:
		logrus.Errorf("Unknown action: %s", argsWithoutProg[0])
		os.Exit(1)
	}

	fmt.Println(argsWithoutProg)

}
