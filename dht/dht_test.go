package dht

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
)

func TestReceive(t *testing.T) {
	id := "5"
	newLocalNode(&id, "localhost", "2000")

	node2 := &remoteNode{_id: "4", _address: "localhost:6600"}
	theLocalNode.pred = node2

	block := make(chan bool)
	<-block
}

func TestHELLO(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	_ = node2
	//theLocalNode.ping(node2)
	block := make(chan bool)
	<-block
}

func TestPredecessorRequest(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	_ = node2.predecessor()
	block := make(chan bool)
	<-block
}

func TestUpdateSuccessorCall(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	candidate := &remoteNode{_id: "8888", _address: "localhost:8877"}
	node2.updateSuccessor(candidate)
	block := make(chan bool)
	<-block
}

func TestUpdatePredecessorCall(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	candidate := &remoteNode{_id: "3", _address: "localhost:8877"}
	node2.updatePredecessor(candidate)
	block := make(chan bool)
	<-block
}

func TestNode0(t *testing.T) {
	id := "00"
	newLocalNode(&id, "localhost", "9000")
	node4 := &remoteNode{_id: "04", _address: "localhost:4000"}
	node6 := &remoteNode{_id: "06", _address: "localhost:6000"}
	theLocalNode.pred = node6
	theLocalNode.fingerTable[1].node = node4
	theLocalNode.fingerTable[2].node = node4
	theLocalNode.fingerTable[3].node = node4

	key := "05"
	n := theLocalNode.lookup(key)
	log.Tracef("%s.lookup(%s) = %s", theLocalNode.id(), key, n.id())

	block := make(chan bool)
	<-block
}

func TestNode4(t *testing.T) {
	id := "04"
	newLocalNode(&id, "localhost", "4000")
	node0 := &remoteNode{_id: "00", _address: "localhost:9000"}
	node6 := &remoteNode{_id: "06", _address: "localhost:6000"}
	theLocalNode.pred = node0
	theLocalNode.fingerTable[1].node = node6
	theLocalNode.fingerTable[2].node = node6
	theLocalNode.fingerTable[3].node = node0

	block := make(chan bool)
	<-block
}

func TestNode6(t *testing.T) {
	id := "06"
	newLocalNode(&id, "localhost", "6000")
	node4 := &remoteNode{_id: "04", _address: "localhost:4000"}
	node0 := &remoteNode{_id: "00", _address: "localhost:9000"}
	theLocalNode.pred = node4
	theLocalNode.fingerTable[1].node = node0
	theLocalNode.fingerTable[2].node = node0
	theLocalNode.fingerTable[3].node = node4

	block := make(chan bool)
	<-block
}

// 1. Sets up one primary db
// 2. Saves a value to it
// 3. Backs up the db
// 4. Reads the saved value from the backup db
func TestDB(t *testing.T) {
	db, err := bolt.Open("db/primary.db", 0600, nil)
	if err != nil {
		t.Errorf("Error opening db", err)
	}

	// Start a read-write transaction
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket1"))
		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	var value string

	if err != nil {
		t.Errorf("Failed to set value in db: %s", err)
	}

	// Start a read transaction
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bucket1"))
		value = string(b.Get([]byte("answer")))
		return nil
	})

	if err != nil {
		t.Errorf("Could not read from db: ", err)
	}

	if value != "42" {
		t.Errorf("Read wrong value from db. Expected 42, got %s", value)
	}

	err = db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile("db/replicas/primary.db", 0600)
		return err
	})

	db2, err := bolt.Open("db/replicas/primary.db", 0600, nil)
	if err != nil {
		t.Errorf("Error opening db2", err)
	}

	// Start a read transaction
	err = db2.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bucket1"))
		value = string(b.Get([]byte("answer")))
		return nil
	})

	if err != nil {
		t.Errorf("Could not read from db2: ", err)
	}

	if value != "42" {
		t.Errorf("Read wrong value from db2. Expected 42, got %s", value)
	}

	fmt.Println(db.GoString())
	db.Close()
}

func TestLocalDbBackup(t *testing.T) {
	id := "01"
	newLocalNode(&id, "localhost", "3000")
	theLocalNode.backupLocalDB()
}

/*
func Test3NodeForwarding(t *testing.T) {
	block := make(chan bool)

	id1 := "01"
	id2 := "02"
	id3 := "03"

	node1 := makeLocalNode(&id1, "127.0.0.1", "2000")
	node2 := makeLocalNode(&id2, "127.0.0.1", "3000")
	node3 := makeLocalNode(&id3, "127.0.0.1", "4000")

	node1.sendRequest(
		msg{
			Method: "FORWARD",
			Values: map[string]string{
				"Method":             "HELLO",
				"FinalDestinationId": "03",
				"Sender":             node1.getAddress(),
			},
			Dst: node2.getAddress()},
	)

	// To prevent stupid warnings
	_ = node3
	<-block
}*/
