package node

import (
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	//	"time"
	"fmt"
	//"html"
	//"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var primaryDB *bolt.DB
var replicaDB *bolt.DB
var mainBucket = []byte("main")

const primary = 1
const replica = 2
const drop = 3
const rearrange = 4

///////////////////////////////////////////////////////////
// Node replication protocol on join
// 1. Node joins and gets successor and predecessor set
// 2. Node requests primary db from predecessor
// 3. Node initiates split with successor
// 4. Node gets whole db from successor, splits primary from precdecessor to self
// 5. Successor splits primary, moving  the range from Node to itself to replicas
///////////////////////////////////////////////////////////

func (n *localNode) initPrimaryAndReplicaDB(id string) {
	var err error

	// Open primary db
	primaryDB, err = bolt.Open("db/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}
	// Create one main bucket
	primaryDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(mainBucket)
		if err != nil {
			return log.Errorf("%s: %s", n.id(), err)
		}
		return nil
	})

	// Open replica db
	replicaDB, err = bolt.Open("db/replicas/primary"+id+".db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}
	// Create one main bucket
	replicaDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(mainBucket)
		if err != nil {
			return err // log.Errorf("%s: %s", n.id(), err)
		}
		return nil
	})

	go n.serveDBs()
}

func (n *localNode) serveDBs() {

	http.HandleFunc("/db/primary", getPrimaryDbHandleFunc)
	http.HandleFunc("/db/replica", getReplicaDbHandleFunc)
	http.HandleFunc("/split-dbs", splitDBHandleFunc)

	// Set up http server to listen for requests
	err := http.ListenAndServe(":"+n.dbPort(), nil)
	if err != nil {
		log.Errorf("%s: ListenAndServe: %s", n.dbAddress(), err)
	}
}

// handler for HTTP GET request for the primary db
func getPrimaryDbHandleFunc(w http.ResponseWriter, req *http.Request) {

	err := primaryDB.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="primary.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		return tx.Copy(w)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handler for HTTP GET request for the replica db
func getReplicaDbHandleFunc(w http.ResponseWriter, req *http.Request) {

	err := replicaDB.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="primary.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		return tx.Copy(w)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func splitDBHandleFunc(w http.ResponseWriter, req *http.Request) {

	log.Tracef("%s: in splitDBHandleFunc", theLocalNode.id())
	// drop previous replica and copy primary to replica
	err := copyFileContents("db/primary"+theLocalNode.id()+".db", "db/replicas/primary"+theLocalNode.id()+".db")
	if err != nil {
		log.Trace(err)
	}

	// run splitprimary
	theLocalNode.splitPrimaryDB()

	// run splitreplica
	theLocalNode.splitReplicaDB()

	io.WriteString(w, "split-dbs OK")
}

func (n *localNode) startReplication() {

	// backup predecessors db
	n.getDB(n.predecessor(), primary, replica)
	log.Trace("In startReplication, after get primary db from pred save as replica")
	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	// get successors db, set it as n´s own primary db
	n.getDB(n.successor(), primary, primary)
	log.Trace("In startReplication, after get primary db from successor save as primary")
	n.printMainBucket(primaryDB)
	n.printMainBucket(replicaDB)

	// take over only a part of it, A, drop the rest
	n.splitPrimaryDB()

}

// get DB from remote node n2
func (n *localNode) getDB(n2 node, getDb, setDbAs int) error {
	var url string
	// Get db from remote node
	if getDb == primary {
		url = "http://" + n2.dbAddress() + "/db/primary"
	} else if getDb == replica {
		url = "http://" + n2.dbAddress() + "/db/replica"
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("%s: %s", n.id(), err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	var saveFileTo string
	if setDbAs == primary {

		// save file as own primary.db
		saveFileTo = "db/primary" + n.id() + ".db"

	} else if setDbAs == replica {

		// save file as replica primary.db
		saveFileTo = "db/replicas/primary" + n.id() + ".db"
	} else if setDbAs == 88 {
		saveFileTo = "db/temp.db"
	}
	err = ioutil.WriteFile(saveFileTo, body, 0600)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (n *localNode) requestSplit(n2 node) {
	// send HTTP request - to invoke splitPrimaryDB(rearrange) on n2

	/*	req, err := http.NewRequest("DO", "http://"+n2.dbAddress()+"/split-rearrange", nil)
		if err != nil {
			log.Error(err)
		}*/
	resp, err2 := http.Get("http://" + n2.dbAddress() + "/split-dbs")
	if err2 != nil {
		log.Errorf("Error requesting split: %s", err2)
	}
	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		log.Errorf("Error reading split response: %s", err3)
	}
	log.Tracef("%s: %s", n.id(), body)
}

func (n *localNode) splitPrimaryDB() {

	err := primaryDB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(mainBucket)
		c := b.Cursor()
		var fail error
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			// keep keys in interval (n.pred, n], delete keys outside that interval
			if between(
				[]byte(nextId(n.predecessor().id())),
				[]byte(nextId(n.id())),
				k,
			) {
				// keep (k,v), in other words - do nothing
			} else {
				c.Delete()
			}
		}
		errorChan := tx.Check()
		fail = <-errorChan
		if fail != nil {
			log.Trace(fail)
		}
		return fail
	})
	if err != nil {
		log.Error(err)
		// doing this due to an error in bolt - page is freed twice and this causes panic
		n.splitPrimaryDB()
	}
}

func (n *localNode) splitReplicaDB() {
	for i := 0; i < 10000; i++ {
		if i == 10000 {
			fmt.Println(n.id())
		}
	}

	err := replicaDB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(mainBucket)
		c := b.Cursor()
		var fail error
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			// delete keys in interval (n.pred, n] - they belong in primaryDB
			if between(
				[]byte(nextId(n.predecessor().id())),
				[]byte(nextId(n.id())),
				k,
			) {
				c.Delete()
			} else {
				// keep (k,v)
			}
		}
		//str, err2 := tx.Page(2)
		//log.Tracef("pageinfo: %+v", str)
		//if err2 != nil {
		//	log.Error(err2)
		//}
		errorChan := tx.Check()
		fail = <-errorChan
		if fail != nil {
			log.Trace(fail)
		}
		return fail
	})
	if err != nil {
		log.Error(err)
		// doing this due to an error in bolt - page is freed twice and this causes panic
		n.splitReplicaDB()
	}
}

func (n *localNode) periodicReplication() {
	for {
		time.Sleep(time.Second * 20)
		// get predecessors db, if new data they will be backed up - otherwise the one got will be the same as the one had
		err := n.getDB(n.predecessor(), primary, replica)
		if err == nil {
			//n.printMainBucket(primaryDB)
			//n.printMainBucket(replicaDB)
		}
	}
}

/*
func (n *localNode) copyPrimaryDBtoReplica() {

	// kopiera primary skriv över replica

	err := ioutil.WriteFile("db/replica/primary.db", primDbBytes, 0600)
	if err != nil {
		log.Error(err)
	}
	// gå igenom replica och droppa allt som ska finnas i primary.
}
*/
// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (n *localNode) storeValue(key, value []byte) error {
	err := primaryDB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(mainBucket)
		err := b.Put(key, value)
		return err
	})
	return err
}

func (n *localNode) printDBsPeriodic() {
	for {
		log.Tracef("%s: primaryDB:", n.id())
		n.printMainBucket(primaryDB)
		log.Tracef("%s: replicaDB:", n.id())
		n.printMainBucket(replicaDB)
		time.Sleep(time.Second * 10)
	}
}

func (n *localNode) printMainBucket(db *bolt.DB) error {

	if db == primaryDB {
		log.Tracef("PrimaryDB:")
	} else if db == replicaDB {
		log.Tracef("ReplicaDB:")
	}
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(mainBucket)

		// Iterate over items in sorted key order
		err := b.ForEach(func(k, v []byte) error {
			log.Tracef("%s: (%s, %s)", n.id(), k, v)
			return nil
		})
		return err
	})
	return err
}

func (n *localNode) printDBs(n2 *remoteNode) {
	// Get remote nodes primary DB and print the contents
	n.getDB(n2, primary, 88)
	tempDB, err := bolt.Open("db/temp.db", 0600, nil)
	if err != nil {
		log.Errorf("Could not open db: %s", err)
	}
	log.Tracef("%s: PrimaryDB:", n2.id())
	n.printMainBucket(tempDB)

	// Get remote nodes replica DB and print the contents
	n.getDB(n2, replica, 88)
	log.Tracef("%s: ReplicaDB:", n2.id())
	n.printMainBucket(tempDB)
}

// Close the db after a time, to prevent it getting stucked as locked?
//db, err = bolt.Open("db/primary.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
//defer db.Close()
/*
func (n *localNode) getValue(key []byte) ([]byte, error) {
	// get corresponding value to key, lookup in both in both primary db and replica db
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(main)
		v = b.Get(key)
		fmt.Printf("In getValue: key %s value %s\n", key, v)
		return nil
	})
	return v, err
}
*/
