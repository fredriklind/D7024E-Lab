package dht

import (
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	//	"time"
	"fmt"
	"io/ioutil"
	"net/http"
)

var bucket1 = []byte("Bucket1")

const primary = 1
const replica = 2

// replication
// getPredReplica

// getOwnDB from successor
//notifySuccessor - of Takeover(of data) and about Drop previous predReplica

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.

///////////////////////////////////////////////////////////
// Node replication protocol on join
// 1. Node joins and gets successor and predecessor set
// 2. Node requests primary db from predecessor
// 3. Node initiates split with successor
// 4. Node gets whole db from successor, splits primary from precdecessor to self
// 5. Successor splits primary, moving  the range from Node to itself to replicas
///////////////////////////////////////////////////////////

func (n *localNode) initPrimaryDB() {
	var err error
	db, err = bolt.Open("db/primary.db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}
	// Create one main bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucket1)
		if err != nil {
			return log.Errorf("Error on create bucket: %s", err)
		}
		return nil
	})
}

func (n *localNode) startReplication() {

	// as long as not one node ring
	if n.id() != n.successor().id() {

		// get successors db
		n.getDB(n.successor(), primary)
		// take over only a part of it, A, drop the rest
		// ...
		// request successor to split its primary and replace its previos replica with part A
		// ...
	}
	// still as long as not one node ring
	if n.id() != n.predecessor().id() {

		// backup predecessors db
		n.getDB(n.predecessor(), replica)
	}
}

// get DB from remote node,
func (n *localNode) getDB(n2 node, useDbAs int) {

	// Get db from remote node
	resp, err := http.Get(n2.dbAddress())
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("%s", body)

	if useDbAs == primary {
		// save file as own primary.db
	} else if useDbAs == replica {
		// save file as replica primary.db
	}
}

func (n *localNode) serveDBs() {
	// Set up http server for both primary.db and predecessor.db

	// create responseWriter
	// call
	/*
	   	func BackupHandleFunc(w http.ResponseWriter, req *http.Request) {
	       err := db.View(func(tx bolt.Tx) error {
	           w.Header().Set("Content-Type", "application/octet-stream")
	           w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
	           w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
	           return tx.Copy(w)
	       })
	       if err != nil {
	           http.Error(w, err.Error(), http.StatusInternalServerError)
	       }
	   }*/
}

// Close the db after a time, to prevent it getting stucked as locked?
//db, err = bolt.Open("db/primary.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
//defer db.Close()
/*
func (n *localNode) getValue(key []byte) ([]byte, error) {
	// get corresponding value to key, lookup in both in both primary db and replica db
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket1)
		v = b.Get(key)
		fmt.Printf("In getValue: key %s value %s\n", key, v)
		return nil
	})
	return v, err
}

func (n *localNode) storeValue(key, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucket1)
		err := b.Put(key, value)
		return err
	})
	return err
}
*/

func (n *localNode) backupLocalDB() error {
	err := db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile("db/replicas/primary.db", 0600)
		return err
	})
	return err
}
