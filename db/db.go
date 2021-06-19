package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	defaultBucket = []byte("default")
)

// Database ....
type Database struct {
	db *bolt.DB
}

func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDB, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	db = &Database{db: boltDB}
	if err := db.createDefaultBucket(); err != nil {
		db.db.Close()
		return nil, nil, fmt.Errorf("failed to create bucket: %v", err)
	}

	return db, boltDB.Close, nil
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

func (d *Database) SetKey(key string, value []byte) error {
	err := d.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
	return err
}

func (d *Database) GetKey(key string) (value []byte, err error) {
	err = d.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(defaultBucket)
		value = b.Get([]byte(key))
		if value == nil {
			return fmt.Errorf("key %s not found", key)
		}
		return nil
	})
	return value, err
}
