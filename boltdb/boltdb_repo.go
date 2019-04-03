package boltdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var dbBucket = []byte("last_usage_time")

type BoltDb func() (*bolt.DB, error)

func NewBoltDb() BoltDb {
	return func() (db *bolt.DB, e error) {
		return bolt.Open("pr00b121-tdd.db", 0400, &bolt.Options{ReadOnly: true})
	}
}

func (dbFn BoltDb) Write(email string, connectTime time.Time) (err error) {

	return

}

func (dbFn BoltDb) Read(email string) (writeTime string, err error) {
	var db *bolt.DB
	if db, err = dbFn(); err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(dbBucket)

		bs := bk.Get([]byte(email))
		if bs == nil {
			return fmt.Errorf("email '%s' not found", email)
		}

		tmp := string(bs)
		writeTime = tmp
		return nil
	})
	return
}
