package main

import (
	"net/http"
	"time"

	"code.google.com/p/gorest"
	"github.com/cznic/kv"
)

var db *kv.DB

// For serving static files
func startWebServer() {
	fs := http.FileServer(http.Dir("."))
	http.ListenAndServe(":8080", fs)
}

func startAPI() {
	opt := &kv.Options{}
	var err error

	db, err = kv.Open("keyValueStore", opt)

	if err != nil {
		db, err = kv.Create("keyValueStore", opt)
		if err != nil {
			panic("Could not open or create DB")
		}
	}

	db.Set([]byte("abc"), []byte("Testing"))
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

	if key == "" {
		// 400 Bad request
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return ""
	}

	value, _ := db.Get(nil, []byte(key))
	if value != nil {
		return string(value)
	} else {
		// 404 Not found
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)
		return ""
	}
}

// POST /storage
func (serv fileAPI) SetPair(PostData KeyValuePair) {
	serv.setPerms()

	if PostData.Key == "" || PostData.Value == "" {
		// 400 Bad request
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}

	// Set the value, if the key does not already exist
	_, didWrite, err := db.Put(nil, []byte(PostData.Key),
		func(k, existingValue []byte) ([]byte, bool, error) {
			if existingValue == nil {
				return []byte(PostData.Value), true, nil
			} else {
				return nil, false, nil
			}
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

	// Set the value, only if the key already exist
	_, didWrite, err := db.Put(nil, []byte(key),
		func(k, existingValue []byte) ([]byte, bool, error) {
			if existingValue == nil {
				return nil, false, nil
			} else {
				return []byte(PostData), true, nil
			}
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

// GET /storage/{key}
func (serv fileAPI) DeletePair(key string) {
	serv.setPerms()
	value, err := db.Extract(nil, []byte(key))

	if value == nil {
		// 404 Not found
		serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	}

	if err != nil {
		// 500 Internal server error
		serv.ResponseBuilder().SetResponseCode(500).Overide(true)
	}
}
