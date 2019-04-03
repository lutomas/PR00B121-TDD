package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"os"
	"regexp"
	"time"
)

var dbBucket = []byte("last_usage_time")

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		logrus.Errorf("Two arguments expected")
		os.Exit(1)
	}

	db, err := bolt.Open("pr00b121-tdd.db", 0600, nil)
	if err != nil {
		logrus.Errorf("Failure caused by: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	action := argsWithoutProg[0]
	email := argsWithoutProg[1]

	matched, err := regexp.Match(`.*@.*\..*`, []byte(email))
	if matched == false {
		logrus.Errorf("Non email provided: '%s'", email)
		os.Exit(1)
	}

	switch action {
	case "read":
		var writeTime string
		err = db.View(func(tx *bolt.Tx) error {
			bk := tx.Bucket(dbBucket)

			bs := bk.Get([]byte(email))
			if bs == nil {
				return fmt.Errorf("email '%s' not found", email)
			}

			writeTime = string(bs)
			return nil
		})
		if err != nil {
			os.ErrExist = err
			logrus.Errorf("Read failed: %v", err)
			os.Exit(1)
		}

		logrus.Infof("Read: email=%s, time=%s", email, writeTime)

	case "write":
		writeTime := time.Now().Format(time.RFC3339)
		logrus.Infof("Write: email=%s, time=%s", email, writeTime)

		err = db.Update(func(tx *bolt.Tx) error {
			bk, err := tx.CreateBucketIfNotExists(dbBucket)
			if err != nil {
				return fmt.Errorf("failed to create bucket: %v", err)
			}

			if err := bk.Put([]byte(email), []byte(writeTime)); err != nil {
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
