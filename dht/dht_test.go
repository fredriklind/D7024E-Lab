package dht

import (
	//"fmt"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	"testing"
)

func TestReceive(t *testing.T) {
	id := "05"
	newLocalNode(&id, "localhost", "5000", "", "")

	node4 := newRemoteNode("04", "localhost", "4000", "", "")
	theLocalNode.pred = node4

	block := make(chan bool)
	<-block
}

func TestPredecessorRequest(t *testing.T) {
	newLocalNode(nil, "localhost", "2000", "", "")
	node5 := newRemoteNode("05", "localhost", "5000", "", "")

	_ = node5.predecessor()
	block := make(chan bool)
	<-block
}

// Run TestJoin1, TestJoin4 and TestJoin7 in that order from three separate tabs in terminal. (Obj 3)
func TestJoin2_db(t *testing.T) {
	id := "02"
	newLocalNode(&id, "localhost", "2000", "", "2001")

	theLocalNode.storeValue([]byte("2222"), []byte("8"))

	//	theLocalNode.join(nil)

	// print own db
	log.Tracef("%s: primaryDB:", theLocalNode.id())
	theLocalNode.printMainBucket(primaryDB)
	// print replica db
	//	log.Tracef("%s: replicaDB:", theLocalNode.id())
	//	theLocalNode.printMainBucket(replicaDB)

	block := make(chan bool)
	<-block
}

func TestJoin4_db(t *testing.T) {
	id := "04"
	newLocalNode(&id, "localhost", "4000", "", "4001")

	theLocalNode.storeValue([]byte("4444"), []byte("16"))

	node2 := newRemoteNode("02", "localhost", "2000", "", "2001")

	// print own db
	log.Tracef("%s: primaryDB:", theLocalNode.id())
	theLocalNode.printMainBucket(primaryDB)

	theLocalNode.getDB(node2, primary)

	// check if node 2´s key,value pair is in replica, primary8.db....!
	/*	newrepl, err := bolt.Open("db/primary04.db", 0600, nil)
		if err != nil {
			log.Errorf("Could not open db: %s", err)
		}

		theLocalNode.printMainBucket(newrepl)*/

	theLocalNode.printMainBucket(primaryDB)

	//	theLocalNode.join(node2)

	// print replica db
	//	log.Tracef("%s: replicaDB:", theLocalNode.id())
	//	theLocalNode.printMainBucket(replicaDB)

	block := make(chan bool)
	<-block
}

func TestJoin7_db(t *testing.T) {
	id := "07"
	newLocalNode(&id, "localhost", "7000", "", "7001")

	theLocalNode.storeValue([]byte("7777"), []byte("28"))

	//	node2 := newRemoteNode("02", "localhost", "2000", "", "2001")

	//	theLocalNode.join(node2)

	// print own db
	log.Tracef("%s: primaryDB:", theLocalNode.id())
	theLocalNode.printMainBucket(primaryDB)
	// print replica db
	log.Tracef("%s: replicaDB:", theLocalNode.id())
	theLocalNode.printMainBucket(replicaDB)

	block := make(chan bool)
	<-block
}

func TestBuild(t *testing.T) {
	// just test if the program compiles
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// 1. Sets up one primary db
// 2. Saves a value to it
// 3. Backs up the db to /replicas/primary.db
// 4. Reads the saved value from the backup db
func TestDB(t *testing.T) {
	id := "01"
	newLocalNode(&id, "localhost", "6000", "", "")

	// Start a read-write transaction
	err := primaryDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("main"))
		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	var value string

	if err != nil {
		t.Errorf("Failed to set value in primaryDB: %s", err)
	}

	// Start a read transaction
	err = primaryDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		value = string(b.Get([]byte("answer")))
		return nil
	})

	if err != nil {
		t.Errorf("Could not read from primaryDB: ", err)
	}

	if value != "42" {
		t.Errorf("Read wrong value from primaryDB. Expected 42, got %s", value)
	}
}

func TestMain(t *testing.T) {
	id := "01"
	newLocalNode(&id, "localhost", "3000", "", "")
	main()
	block := make(chan bool)
	<-block
}

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

	node3 := newRemoteNode("03", "localhost", "3000", "", "")

	theLocalNode.join(node3)
	block := make(chan bool)
	<-block
}

func TestJoin2(t *testing.T) {
	id := "02"
	newLocalNode(&id, "localhost", "2000", "", "")

	node3 := newRemoteNode("03", "localhost", "3000", "", "")

	theLocalNode.join(node3)

	block := make(chan bool)
	<-block
}
