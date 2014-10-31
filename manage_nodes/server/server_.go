package server

// this file should contain code for starting up a webserver,
// the REST-API with its routes to functions in manage_dht.go

import (
	"net/http"
	//"net/url"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"code.google.com/p/gorest"
	"github.com/boltdb/bolt"
)

var nodeDB *bolt.DB

const (
	webServerRootPath = "./server/webserver"
	webServerPort     = "8080"
	apiPort           = "8787"
	logPort           = "8888"
)

type Node struct {
	Id, Ip, Port string
}

// Serves static files, such as index.html
func StartWebServer() {
	// Write port number of API to a file
	err := ioutil.WriteFile(webServerRootPath+"/port.txt", []byte(apiPort), 0777)

	if err != nil {
		fmt.Printf("Could not start web server: %s", err)
		return
	}

	// Start servig files
	fs := http.FileServer(http.Dir(webServerRootPath))
	http.ListenAndServe(":"+webServerPort, fs)
}

// Registers and starts the file api that serves all reqests to ip:apiPort/api
func StartAPI() {
	initNodeDB()
	serv := new(fileAPI)
	gorest.RegisterService(serv)
	serv.RestService.ResponseBuilder()
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+apiPort, nil)
}

// Sets headers so that the API can be accessed from other domains, and so that
// jQuery works with PUT and DELETE
func (serv fileAPI) setPerms() {
	serv.RB().AddHeader("Access-Control-Allow-Origin", "*")
	serv.RB().AddHeader("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	serv.RB().ConnectionKeepAlive()
}

//Service Definition
type fileAPI struct {
	gorest.RestService `root:"/nodes/" consumes:"application/json" produces:"application/json"`
	getAllNodes        gorest.EndPoint `method:"GET" path:"/" output:"[]Node"`
	/*getPair            gorest.EndPoint `method:"GET" path:"/storage/{key:string}" output:"string"`
	setPair            gorest.EndPoint `method:"POST" path:"/storage" postdata:"KeyValuePair"`
	updatePair         gorest.EndPoint `method:"PUT" path:"/storage/{key:string}" postdata:"string"`
	deletePair         gorest.EndPoint `method:"DELETE" path:"/storage/{key:string}"`
	optionsRoute       gorest.EndPoint `method:"OPTIONS" path:"/storage/{key:string}"`*/
}

// Needed to allow jQuery to do PUT/DELETE. (jQuery first sends OPTION)
func (serv fileAPI) OptionsRoute(_ string) {
	serv.setPerms()
}

// GET /storage (NOT IMPLEMENTED)
// Just there to give error in case the key searched for was empty
func (serv fileAPI) GetAllNodes() []Node {
	// 400 Bad request
	//serv.ResponseBuilder().SetResponseCode(400).Overide(true)
	var nodes []Node

	nodeDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(v))
			var n Node
			err := json.Unmarshal(v, &n)
			if err != nil {
				fmt.Printf("Error unmarshaling: %s\n", err)
				serv.ResponseBuilder().SetResponseCode(500).Overide(true)
				return err
			}
			nodes = append(nodes, n)
			return nil
		})
		return nil
	})
	return nodes
}

/*
// GET /storage/{key}
func (serv fileAPI) GetPair(key string) string {
	serv.setPerms()
}

// POST /storage
func (serv fileAPI) SetPair(PostData KeyValuePair) {
	serv.setPerms()
}

// PUT /storage/{key}
func (serv fileAPI) UpdatePair(PostData string, key string) {
	serv.setPerms()

}

// DELETE /storage/{key}
func (serv fileAPI) DeletePair(key string) {
	serv.setPerms()

}

func sendPostRequest(url string, data KeyValuePair) (error, int) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Could not parse JSON: %s", err)
		return err, 0
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending POST request: %s", err)
		return err, 0
	}
	defer resp.Body.Close()

	return nil, resp.StatusCode
}

func sendPutRequest(url string, value string) (error, int) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(value)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending PUT request: %s", err)
		return err, 0
	}
	defer resp.Body.Close()

	return nil, resp.StatusCode
}

func sendDeleteRequest(url string) (error, int) {
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending DELETE request: %s", err)
		return err, 0
	}
	defer resp.Body.Close()

	return nil, resp.StatusCode
}*/

func initNodeDB() {
	var err error
	nodeDB, err = bolt.Open("nodes.db", 0600, nil)
	if err != nil {
		fmt.Printf("Could not open db: %s", err)
	}
	// Create one main bucket
	nodeDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("main"))
		if err != nil {
			fmt.Printf("Error creating bucket: %s", err)
			return err
		}
		return nil
	})

	nodeDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		node1 := Node{Id: "AaNNeDa8ehcB3iwc3boa"}
		node2 := Node{Id: "CJHakuwhc8k28hKHkfhk"}
		jsonBytes1, err := json.Marshal(node1)
		jsonBytes2, err := json.Marshal(node2)
		if err != nil {
			fmt.Printf("Error parsing json: %s", err)
		}
		b.Put([]byte(node1.Id), jsonBytes1)
		b.Put([]byte(node2.Id), jsonBytes2)
		return nil
	})
}
