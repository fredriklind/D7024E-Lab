package node

import (
	//"fmt"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	"testing"
)

// run TestJoinSplitRequest first then this test to see if all key-value-pairs are assembled again in primary
func TestRecover(t *testing.T) {
	id := "04"
	var err error
	primaryDB, err = bolt.Open("db/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}
	replicaDB, err = bolt.Open("db/replicas/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}

	//recoverData()								// <---------- uncomment this when running test! and change the func stub in storage.go

	primaryDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mainBucket)

		// Iterate over items in sorted key order
		b.ForEach(func(k, v []byte) error {
			log.Tracef("(%s, %s)", k, v)
			return nil
		})
		return nil
	})
	log.Trace("replica")
	replicaDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mainBucket)

		// Iterate over items in sorted key order
		b.ForEach(func(k, v []byte) error {
			log.Tracef("(%s, %s)", k, v)
			return nil
		})
		return nil
	})
}

// this test fails every other time, in splitReplica: -> panic: assertion failed: page 2 already freed [recovered]
// something wrong in bolt, found some one that had related problem but that was said to be fixed?
func TestJoinSplitRequest(t *testing.T) {

	id := "04"
	NewLocalNode(&id, "localhost", "", "", "4444")
	n := theLocalNode
	node2 := newRemoteNode("02", "localhost", "", "", "")

	theLocalNode.storeValue([]byte("01"), []byte("02"))
	theLocalNode.storeValue([]byte("02"), []byte("02"))
	theLocalNode.storeValue([]byte("03"), []byte("04"))
	theLocalNode.storeValue([]byte("04"), []byte("04"))
	theLocalNode.storeValue([]byte("05"), []byte("02"))
	n.pred = node2

	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	err := copyFileContents("db/primary"+n.id()+".db", "db/replicas/primary"+n.id()+".db")
	if err != nil {
		log.Trace(err)
	}

	log.Trace("Copied primary to replica")

	// run splitprimary
	theLocalNode.splitPrimaryDB()

	log.Trace("After splitPrimary:")
	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	// run splitreplica
	theLocalNode.splitReplicaDB()

	log.Trace("After splitReplica:")
	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)
}

// Run TestJoin2_all, TestJoin4_all and TestJoin7_all in that order from three separate tabs in terminal. (Obj 3)
func TestJoin2_all(t *testing.T) {
	id := "02"
	NewLocalNode(&id, "localhost", "2001", "2002", "2003")
	n := theLocalNode

	n.join(nil)
	log.Tracef("NODE %s HAS JOINED THE RING!", n.id())

	n.storeValue([]byte("02"), []byte("02"))

	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	block := make(chan bool)
	<-block
}

func TestJoin4_all(t *testing.T) {
	id := "04"
	NewLocalNode(&id, "localhost", "4001", "4002", "4003")
	n := theLocalNode

	node2 := newRemoteNode("02", "localhost", "2001", "2002", "2003")

	n.join(node2)
	log.Tracef("NODE %s HAS JOINED THE RING!", n.id())

	n.storeValue([]byte("04"), []byte("04"))

	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	n.printDBs(node2)

	block := make(chan bool)
	<-block
}

func TestJoin7_all(t *testing.T) {
	id := "07"
	NewLocalNode(&id, "localhost", "7001", "7002", "7003")
	n := theLocalNode

	node2 := newRemoteNode("02", "localhost", "2001", "2002", "2003")
	node4 := newRemoteNode("04", "localhost", "4001", "4002", "4003")

	n.join(node2)
	log.Tracef("NODE %s HAS JOINED THE RING!", n.id())

	n.storeValue([]byte("07"), []byte("07"))

	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	n.printDBs(node4)
	n.printDBs(node2)

	block := make(chan bool)
	<-block
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func TestSplit(t *testing.T) {
	id := "04"
	NewLocalNode(&id, "localhost", "", "", "")
	n := theLocalNode
	node2 := newRemoteNode("02", "localhost", "", "", "")
	n.pred = node2

	n.storeValue([]byte("01"), []byte("OUTSIDE"))
	n.storeValue([]byte("02"), []byte("OUTSIDE"))
	n.storeValue([]byte("03"), []byte("IN"))
	n.storeValue([]byte("04"), []byte("IN"))

	log.Tracef("Before split")
	n.printMainBucket(primaryDB)

	n.splitPrimaryDB()
	log.Tracef("After split")
	n.printMainBucket(primaryDB)
}

func TestCopyFile(t *testing.T) {

	copyFileContents("db/primary04.db", "db/replicas/primary04.db")
}

func TestReplica(t *testing.T) {
	id := "04"
	NewLocalNode(&id, "localhost", "", "", "")
	n := theLocalNode

	var err error
	replicaDB, err = bolt.Open("db/replicas/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}

	n.printMainBucket(replicaDB)
}

func TestSplitReplica(t *testing.T) {
	id := "04"
	NewLocalNode(&id, "localhost", "", "", "")
	n := theLocalNode
	node2 := newRemoteNode("02", "localhost", "", "", "")
	n.pred = node2

	var err error
	replicaDB, err = bolt.Open("db/replicas/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}

	n.storeValue([]byte("01"), []byte("02"))
	n.storeValue([]byte("02"), []byte("02"))
	n.storeValue([]byte("03"), []byte("04"))
	n.storeValue([]byte("04"), []byte("04"))
	n.storeValue([]byte("05"), []byte("02"))

	n.splitReplicaDB()
	n.printMainBucket(replicaDB)
}

func TestReceive(t *testing.T) {
	id := "05"
	NewLocalNode(&id, "localhost", "5000", "", "")

	node4 := newRemoteNode("04", "localhost", "4000", "", "")
	theLocalNode.pred = node4

	block := make(chan bool)
	<-block
}

func TestPredecessorRequest(t *testing.T) {
	NewLocalNode(nil, "localhost", "2000", "", "")
	node5 := newRemoteNode("05", "localhost", "5000", "", "")

	_ = node5.predecessor()
	block := make(chan bool)
	<-block
}

// Run TestJoin1, TestJoin4 and TestJoin7 in that order from three separate tabs in terminal. (Obj 3)
func TestJoin2_db(t *testing.T) {
	id := "02"
	NewLocalNode(&id, "localhost", "2000", "", "2001")

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
	NewLocalNode(&id, "localhost", "4000", "", "4001")

	theLocalNode.storeValue([]byte("4444"), []byte("16"))

	node2 := newRemoteNode("02", "localhost", "2000", "", "2001")

	// print own db
	log.Tracef("%s: primaryDB:", theLocalNode.id())
	theLocalNode.printMainBucket(primaryDB)

	theLocalNode.getDB(node2, primary, primary)

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
	NewLocalNode(&id, "localhost", "7000", "", "7001")

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

// 1. Sets up one primary db
// 2. Saves a value to it
// 3. Backs up the db to /replicas/primary.db
// 4. Reads the saved value from the backup db
func TestDB(t *testing.T) {
	id := "01"
	NewLocalNode(&id, "localhost", "6000", "", "")

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

// Run TestJoin3, TestJoin0 and TestJoin2 in that order from three separate tabs in terminal. (To test obj2).
func TestJoin3(t *testing.T) {
	id := "03"
	NewLocalNode(&id, "localhost", "5000", "5100", "5200")

	theLocalNode.join(nil)

	block := make(chan bool)
	<-block
}

func TestJoin0(t *testing.T) {
	id := "00"
	NewLocalNode(&id, "localhost", "9000", "9100", "9200")

	node3 := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	theLocalNode.join(node3)
	block := make(chan bool)
	<-block
}

func TestJoin2(t *testing.T) {
	id := "02"
	NewLocalNode(&id, "localhost", "2000", "2100", "2200")

	node3 := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	theLocalNode.join(node3)

	block := make(chan bool)
	<-block
}

func TestForwardingReceiver(t *testing.T) {
	id := "03"
	NewLocalNode(&id, "localhost", "5000", "5100", "5200")

	block := make(chan bool)
	<-block
}

func TestForwardingSender(t *testing.T) {
	id := "01"
	NewLocalNode(&id, "localhost", "4000", "4100", "4200")

	block := make(chan bool)
	<-block
}
