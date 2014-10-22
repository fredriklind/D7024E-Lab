package dht

import (
	"net/http"
	"time"

	"code.google.com/p/gorest"
	"github.com/boltdb/bolt"
)

// For serving static files
func startWebServer() {
	fs := http.FileServer(http.Dir("."))
	http.ListenAndServe(":8080", fs)
}

func startAPI() {
	serv := new(fileAPI)
	gorest.RegisterService(serv)
	serv.RestService.ResponseBuilder()
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":8787", nil)
}

func (serv fileAPI) setPerms() {
	serv.RB().AddHeader("Access-Control-Allow-Origin", "*")
	serv.RB().AddHeader("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	// TODO Remove this
	time.Sleep(time.Millisecond * 500)
}

//Service Definition
type fileAPI struct {
	gorest.RestService `root:"/api/" consumes:"application/json" produces:"application/json"`
	getAll             gorest.EndPoint `method:"GET" path:"/storage/" output:"string"`
	getPair            gorest.EndPoint `method:"GET" path:"/storage/{key:string}" output:"string"`
	setPair            gorest.EndPoint `method:"POST" path:"/storage" postdata:"KeyValuePair"`
	updatePair         gorest.EndPoint `method:"PUT" path:"/storage/{key:string}" postdata:"string"`
	deletePair         gorest.EndPoint `method:"DELETE" path:"/storage/{key:string}"`
	optionsRoute       gorest.EndPoint `method:"OPTIONS" path:"/storage/{key:string}"`
}

type KeyValuePair struct {
	Key, Value string
}

// Needed to allow jQuery to do PUT/DELETE. (jQuery first sends OPTION)
func (serv fileAPI) OptionsRoute(_ string) {
	serv.setPerms()
}

// GET /storage (NOT IMPLEMENTED)
func (serv fileAPI) GetAll() string {
	// 400 Bad request
	serv.ResponseBuilder().SetResponseCode(400).Overide(true)
	return ""
}

// GET /storage/{key}
func (serv fileAPI) GetPair(key string) string {
	serv.setPerms()

	responsibleNode, _ := theLocalNode.lookup(key)

	// If I'm responsible
	if responsibleNode.id() == theLocalNode.id() {
		if key == "" {
			// 400 Bad request
			serv.ResponseBuilder().SetResponseCode(400).Overide(true)
			return ""
		}

		var value []byte

		// Start view transaction, get value
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("main"))
			value = b.Get([]byte(key))
			return nil
		})

		if value != nil {
			return string(value)
		} else {
			// 404 Not found
			serv.ResponseBuilder().SetResponseCode(404).Overide(true)
			return ""
		}

		// If someone else is responsible, sent request to that guy
	} else {
		//restClient, _ := gorest.NewRequestBuilder("http://" + responsibleNode.address())
	}

	return ""
}

// POST /storage
func (serv fileAPI) SetPair(PostData KeyValuePair) {
	serv.setPerms()

	if PostData.Key == "" || PostData.Value == "" {
		// 400 Bad request
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}

	didWrite := false
	var err error
	// Set the value, if the key does not already exist
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		existingValue := b.Get([]byte(PostData.Key))
		if existingValue == nil {
			err := b.Put([]byte(PostData.Key), []byte(PostData.Value))
			return err
		}
		return nil
	})

	if !didWrite {
		// 403 Forbidden
		serv.ResponseBuilder().SetResponseCode(403).Overide(true)
	}

	if err != nil {
		// 500 internal server error
		serv.ResponseBuilder().SetResponseCode(500).Overide(true)
	}
}

// PUT /storage/{key}
func (serv fileAPI) UpdatePair(PostData string, key string) {
	serv.setPerms()

	if key == "" || PostData == "" {
		// 400 Bad request
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}

	didWrite := false
	var err error
	// Set the value, only if the key already exist
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		existingValue := b.Get([]byte(key))
		if existingValue != nil {
			err := b.Put([]byte(key), []byte(PostData))
			return err
		}
		return nil
	})

	if !didWrite {
		// 404 Not found
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	}

	if err != nil {
		// 500 internal server error
		serv.ResponseBuilder().SetResponseCode(500).Overide(true)
	}
}

// DELETE /storage/{key}
func (serv fileAPI) DeletePair(key string) {
	serv.setPerms()
	var value []byte
	var err error

	// Check if value exists
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		value = b.Get([]byte(key))
		return nil
	})

	if value == nil {
		// 404 Not found
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	} else {
		// Delete the value
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("main"))
			err = b.Delete([]byte(key))
			return err
		})
	}

	if err != nil {
		// 500 Internal server error
		serv.ResponseBuilder().SetResponseCode(500).Overide(true)
	}
}
