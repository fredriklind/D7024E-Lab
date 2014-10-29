package node

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"code.google.com/p/gorest"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
)

const (
	webServerRootPath = "./webserver"
	webServerPort     = "8080"
)

// Serves static files, such as index.html
func startWebServer() {
	// Write port number of API to a file
	port := theLocalNode.apiPort()
	err := ioutil.WriteFile(webServerRootPath+"/port.txt", []byte(port), 0777)

	if err != nil || port == "" {
		log.Errorf("Could not start web server: %s", err)
		return
	}

	// Start servig files
	fs := http.FileServer(http.Dir(webServerRootPath))
	http.ListenAndServe(":"+webServerPort, fs)
}

// Registers and starts the file api that serves all reqests to ip:apiPort/api
func startAPI() {
	serv := new(fileAPI)
	gorest.RegisterService(serv)
	serv.RestService.ResponseBuilder()
	http.Handle("/", gorest.Handle())
	if theLocalNode.apiPort() != "" {
		http.ListenAndServe(":"+theLocalNode.apiPort(), nil)
	} else {
		log.Error("API port not set, cannot start API")
	}
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

type Pairer interface {
	key() string
	value() string
}

// Needed to allow jQuery to do PUT/DELETE. (jQuery first sends OPTION)
func (serv fileAPI) OptionsRoute(_ string) {
	serv.setPerms()
}

// GET /storage (NOT IMPLEMENTED)
// Just there to give error in case the key searched for was empty
func (serv fileAPI) GetAll() string {
	// 400 Bad request
	serv.ResponseBuilder().SetResponseCode(400).Overide(true)
	return ""
}

// GET /storage/{key}
func (serv fileAPI) GetPair(key string) string {
	serv.setPerms()
	log.Tracef("Handling GET requset for key %s", key)

	// Validate request
	if key == "" {
		// 400 Bad request
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return ""
	}

	// Remvove things like %20 and stuff
	ogKey := key
	key, _ = url.QueryUnescape(key)

	/*responsibleNode, err := theLocalNode.lookup(key)
	if err != nil {
		log.Errorf("Lookup of responsible node failed: %s", err)
		return ""
	}*/
	responsibleNode := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	// If I'm responsible
	if responsibleNode.id() == theLocalNode.id() {
		log.Tracef("I am responsible, searching for %s", key)
		// Get the the value from primary db
		var value []byte
		primaryDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(mainBucket)
			value = b.Get([]byte(key))
			return nil
		})

		// The primaryDB did not have the value
		if value == nil {
			log.Trace("Didn't find value in primary db")
			// See if the replica has the value
			replicaDB.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(mainBucket)
				value = b.Get([]byte(key))
				return nil
			})

			if value == nil {
				// Replica didnt have the value either, 404 Not found
				log.Errorf("Didn't find value in replica db")
				serv.ResponseBuilder().SetResponseCode(404).Overide(true)
				return ""

			} else {
				// Replica did have the value
				log.Tracef("Found value \"%s\" in replica db", string(value))
				return string(value)
			}
		} else {
			log.Tracef("Found value \"%s\" in primary db", string(value))
			return string(value)
		}
		//If someone else is responsible, sent request to that guy
	} else {
		log.Tracef("%s is responsible for key %s, forwarding request", responsibleNode.address(), key)
		response, err := http.Get("http://" + responsibleNode.apiAddress() + "/api/storage/" + ogKey)
		defer response.Body.Close()

		if err != nil {
			log.Errorf("Forwarded request returned error: %s", err)
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
			return ""
		} else {
			// Parse bytes to string
			responseBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Errorf("Error parsing response: %s", err)
				serv.ResponseBuilder().SetResponseCode(500).Overide(true)
				return ""
			}
			if response.StatusCode == 200 {
				responseString := string(responseBytes)
				log.Tracef("Got value \"%s\" from responsible node", responseString)
				return responseString
			} else {
				serv.ResponseBuilder().SetResponseCode(response.StatusCode).Overide(true)
				return ""
			}
		}
	}

	// Everything fell through, didn't get any value
	serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	log.Tracef("Nothing found in GET, exiting")
	return ""
}

// POST /storage
func (serv fileAPI) SetPair(PostData KeyValuePair) {
	serv.setPerms()
	log.Tracef("Handling POST requset for key:%s value:%s ", PostData.Key, PostData.Value)

	if PostData.Key == "" || PostData.Value == "" {
		// 400 Bad request
		log.Tracef("Got bad POST request", PostData.Key)
		serv.ResponseBuilder().SetResponseCode(400).Overide(true)
		return
	}

	/*responsibleNode, err := theLocalNode.lookup(PostData.Key)
	if err != nil {
		log.Errorf("Lookup of responsible node failed: %s", err)
		return
	}*/
	responsibleNode := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	// If i'm responsible
	if responsibleNode.id() == theLocalNode.id() {
		log.Tracef("I am responsible, inserting into primary db...")

		didWrite := false
		var err error
		// Set the value, if the key does not already exist
		primaryDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(mainBucket)
			existingValue := b.Get([]byte(PostData.Key))
			if existingValue == nil {
				err := b.Put([]byte(PostData.Key), []byte(PostData.Value))
				if err == nil {
					didWrite = true
				}
				return err
			}
			return nil
		})

		if !didWrite {
			// 409 Conflict
			log.Tracef("Write conflict on key %s", PostData.Key)
			serv.ResponseBuilder().SetResponseCode(409).Overide(true)
			return
		}

		if err != nil {
			// 500 internal server error
			log.Tracef("Error saving value to primart db: %s", err)
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
			return
		}
		log.Tracef("Value was inserted")

		// Someone else is responsible, forward to that guy
	} else {
		log.Tracef("%s is responsible for key %s, forwarding request", responsibleNode.address(), PostData.Key)
		err, statusCode := sendPostRequest("http://"+responsibleNode.apiAddress()+"/api/storage", PostData)

		if err != nil {
			log.Errorf("Forward request error: %s", err)
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
			return
		}

		serv.ResponseBuilder().SetResponseCode(statusCode).Overide(true)
		return
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

	/*responsibleNode, err := theLocalNode.lookup(PostData.Key)
	if err != nil {
		log.Errorf("Lookup of responsible node failed: %s", err)
		return
	}*/
	responsibleNode := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	// If i'm responsible
	if responsibleNode.id() == theLocalNode.id() {

		didWrite := false
		var err error
		// Set the value, only if the key already exist
		primaryDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(mainBucket)
			existingValue := b.Get([]byte(key))
			if existingValue != nil {
				err := b.Put([]byte(key), []byte(PostData))
				if err == nil {
					didWrite = true
				}
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
	} else {
		log.Tracef("%s is responsible for key %s, forwarding request", responsibleNode.address(), key)
		err, statusCode := sendPutRequest("http://"+responsibleNode.apiAddress()+"/api/storage/"+url.QueryEscape(key), PostData)

		if err != nil {
			log.Errorf("Forward request error: %s", err)
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
			return
		}

		serv.ResponseBuilder().SetResponseCode(statusCode).Overide(true)
		return
	}
}

// DELETE /storage/{key}
func (serv fileAPI) DeletePair(key string) {
	serv.setPerms()

	/*responsibleNode, err := theLocalNode.lookup(PostData.Key)
	if err != nil {
		log.Errorf("Lookup of responsible node failed: %s", err)
		return
	}*/
	responsibleNode := newRemoteNode("03", "localhost", "5000", "5100", "5200")

	// If i'm responsible
	if responsibleNode.id() == theLocalNode.id() {

		var value []byte
		var err error

		// Check if value exists
		primaryDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(mainBucket)
			value = b.Get([]byte(key))
			return nil
		})

		if value == nil {
			// 404 Not found
			serv.ResponseBuilder().SetResponseCode(404).Overide(true)
		} else {
			// Delete the value
			primaryDB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket(mainBucket)
				err = b.Delete([]byte(key))
				return err
			})
		}

		if err != nil {
			// 500 Internal server error
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
		}
	} else {
		log.Tracef("%s is responsible for key %s, forwarding request", responsibleNode.address(), key)
		err, statusCode := sendDeleteRequest("http://" + responsibleNode.apiAddress() + "/api/storage/" + url.QueryEscape(key))

		if err != nil {
			log.Errorf("Forward request error: %s", err)
			serv.ResponseBuilder().SetResponseCode(500).Overide(true)
			return
		}

		serv.ResponseBuilder().SetResponseCode(statusCode).Overide(true)
		return
	}
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
}
