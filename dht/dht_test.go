package dht

import (
	"github.com/boltdb/bolt"
	"testing"
)

func TestReceive(t *testing.T) {
	id := "5"
	newLocalNode(&id, "localhost", "2000", "", "")

	node2 := newRemoteNode("01", "localhost", "9000", "", "")
	theLocalNode.pred = node2

	block := make(chan bool)
	<-block
}

// 1. Sets up one primary db
// 2. Saves a value to it
// 3. Backs up the db to /replicas/primary.db
// 4. Reads the saved value from the backup db
func TestDB(t *testing.T) {
	id := "01"
	newLocalNode(&id, "localhost", "6000", "", "")

	// Start a read-write transaction
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("main"))
		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	var value string

	if err != nil {
		t.Errorf("Failed to set value in db: %s", err)
	}

	// Start a read transaction
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		value = string(b.Get([]byte("answer")))
		return nil
	})

	if err != nil {
		t.Errorf("Could not read from db: ", err)
	}

	if value != "42" {
		t.Errorf("Read wrong value from db. Expected 42, got %s", value)
	}

	// Backup the local DB
	theLocalNode.backupLocalDB()

	db2, err := bolt.Open("db/replicas/primary.db", 0600, nil)
	if err != nil {
		t.Errorf("Error opening db2", err)
	}

	// Start a read transaction
	err = db2.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		if b != nil {
			value = string(b.Get([]byte("answer")))
		}
		return nil
	})

	if err != nil {
		t.Errorf("Could not read from db2: ", err)
	}

	if value != "42" {
		t.Errorf("Read wrong value from db2. Expected 42, got %s", value)
	}
	db2.Close()
}

/*
// Run TestJoin3, TestJoin0 and TestJoin2 in that order from three separate tabs in terminal. (To test obj2).
func TestJoin3(t *testing.T) {
	id := "03"
	newLocalNode(&id, "localhost", "3000", "", "")

	theLocalNode.join(nil)

	block := make(chan bool)
	<-block
}

func TestJoin0(t *testing.T) {
	id := "00"
	newLocalNode(&id, "localhost", "9000", "", "")

	node3 := &remoteNode{_id: "03", _address: "localhost:3000"}

	theLocalNode.join(node3)
	block := make(chan bool)
	<-block
}

func TestJoin2(t *testing.T) {
	id := "02"
	newLocalNode(&id, "localhost", "2000", "", "")

	node3 := &remoteNode{_id: "03", _address: "localhost:3000"}

	theLocalNode.join(node3)

	block := make(chan bool)
	<-block
}*/
