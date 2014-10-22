package dht

import (
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
)

// RESKJFLEJFLS KEHFLS EKUFHL
// replication
// getPredReplica
// getPredPredReplica

// getOwnDB from successor
//notifySuccessor - of Takeover(of data) and about Drop previous predpredReplica
//notifySuccessorSuccessor - new division of data and about Drop previous predpredReplica

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
		log.Errorf("Could not open db: %s")
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("main"))
		return err
	})
}

func (n *localNode) serveDBs() {
	// Set up http server for both primary.db and predecessor.db
}

func (n *localNode) backupLocalDB() error {
	err := db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile("db/replicas/primary.db", 0600)
		return err
	})
	return err
}

func (n *localNode) backupPredecessorDB() error {
	// Get db from predecessor
	// Save it to db/replicas
	return nil
}
