package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"code.google.com/p/go-uuid/uuid"
	log "github.com/cihub/seelog"
)

// ----------------------------------------------------------------------------------------
//										Init, types and variables
// ----------------------------------------------------------------------------------------

var t *transport

const timeoutSeconds = time.Second * 8

type transport struct {
	address  string
	requests map[string]chan msg
}

type msg struct {
	Id, Type, Method, Src, Dst string
	Timestamp                  int64
	Values                     map[string]string
	Sync                       bool
}

// Initializer for the package, sets up the logger
/*func init() {
	testConfig := `
		<seelog type="sync">
			<outputs>
				<file formatid="onlytime" path="logfile.log"/>
			</outputs>
			<formats>
				<format id="onlytime" format="%Time [%LEVEL] %Msg%n"/>
			</formats>
		</seelog>
	`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}
*/

// ----------------------------------------------------------------------------------------
//										public layer
// ----------------------------------------------------------------------------------------

// Instantiates a new transporter that listens on the provided ip and port.
// This step is required to be able to use the transport package.
func NewTransporter(ip, port string) {
	t = &transport{}
	t.address = ip + ":" + port
	t.requests = make(map[string]chan msg)
	go receive()
}

// TODO define a return type for all methods like this one, maybe map[string]string?
func SendLookupRequest(address, id string) {
	// check queue
	// if lookup in queue - forward request
	// else send new lookupRequest
	lookupRequest := msg{
		Method: "LOOKUP",
		Dst:    t.address,
	}
	sendRequest(lookupRequest)
}

func SendHelloRequest(ip string) {
	request := msg{}
	request.Dst = ip
	request.Method = "HELLO"
	sendRequest(request)
}

func PredecessorRequest(address string) {
	// TODO
}

// ----------------------------------------------------------------------------------------
//										middle layer
// ----------------------------------------------------------------------------------------

func sendRequest(request msg) {
	// Construct the message
	msg := msg{
		Type:   "Request",
		Method: request.Method,
		Src:    t.address,
		Dst:    request.Dst,
	}

	// This is to prevent send from generating a new msg.Id
	if request.Id != "" {
		msg.Id = request.Id
	}
	send(msg)
}

func sendResponse(response msg) {
	//Construct the message
	msg := msg{
		Type:   "Response",
		Method: response.Method,
		Src:    t.address,
		Dst:    response.Dst,
	}

	// This is to prevent send from generating a new msg.Id
	if response.Id != "" {
		msg.Id = response.Id
	}
	send(msg)
}

func waitForResponse(request msg) {

	// Save the channel so that the receive() method can un-block
	// this method when it receives a response with a matching id
	t.requests[request.Id] = make(chan msg)

	// Wait for the msg-specific channel to get data, or time out
	select {
	case response := <-t.requests[request.Id]:
		handleResponse(response)
	case <-time.After(timeoutSeconds):
		log.Errorf("%s: %s request with id %s timed out", t.address, request.Method, request.Id)
	}
}

// When you get a request you need to handle
func handleRequest(request msg) {
	switch request.Method {
	case "HELLO":
		sendResponse(msg{
			Id:     request.Id,
			Method: "ACK",
			Dst:    request.Src,
		})

		// Forwards a request to nextNode setting the method and Src depending
		// on the Values["Method"] and Values["Sender"]
		/*	case "FORWARD":
				// If n is the final destination, answer the original sender
				if n.id == request.Values["FinalDestinationId"] {
					//			log.Tracef("Node %s is FinalDestination", n.id)
					newRequest := msg{
						Method: request.Values["Method"],
						Src:    request.Values["Sender"],
					}
					// Handle the request contained in the FORWARD request
					handleRequest(newRequest)
				} else {
					// If n is not the searched for node, forward the request to the next node
					nextNodeAddress := "127.0.0.1:4000"
					forwardRequest := msg{
						Id:     request.Id,
						Method: "FORWARD",
						Dst:    nextNodeAddress,
						Sync:   false,
					}
					sendRequest(forwardRequest)
				}

			default:
				log.Error("No request method specified!")*/
	}
}

// When you get something back from a request you sent (you = n)
func handleResponse(response msg) {
	switch response.Method {
	case "ACK":
		log.Tracef("%s HELLO request satisfied", t.address)
	default:
		log.Errorf("No request method specified!")
	}
}

func (m *msg) isRequest() bool  { return m.Type == "Request" }
func (m *msg) isResponse() bool { return m.Type == "Response" }

// ----------------------------------------------------------------------------------------
//										lowest layer
// ----------------------------------------------------------------------------------------

func send(msg msg) {

	// Start up network stuff
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	// Apply unix timestamp to the msg
	ti := time.Now()
	msg.Timestamp = ti.Unix()

	// Assign a new Id to the msg if it is not set
	if msg.Id == "" {
		msg.Id = uuid.New()[0:4]
	}

	//	log.Tracef("%s: Sent %s %s", n.id, msg.Method, msg.Type)
	fmt.Printf("%s: Sent %s %s\n", t.address, msg.Method, msg.Type)
	log.Tracef("%s: Sent %s %s", t.address, msg.Method, msg.Type)

	// Serialize and send the message (also wait to simulate network delay)
	jsonmsg, err := json.Marshal(msg)
	_, err = conn.Write([]byte(jsonmsg))

	// Blocks until something is received on the channel that is associated with msg.Id
	if msg.isRequest() && msg.Sync {
		waitForResponse(msg)
	}

	if err != nil {
		log.Errorf("%s: Error in send: %s", t.address, err.Error())
	}
}

func receive() {
	// Start receiveing
	udpAddr, err := net.ResolveUDPAddr("udp", t.address)
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Errorf("%s: Error in receive, %s", t.address, err.Error())
		return
	}

	defer conn.Close()

	// receive again after this method finishes. TODO might be
	// that there is a better way to do this
	defer receive()
	dec := json.NewDecoder(conn)

	for {
		msg := msg{}
		err = dec.Decode(&msg)
		fmt.Printf("%s: Got %s %s\n", t.address, msg.Method, msg.Type)
		log.Tracef("%s: Got %s %s", t.address, msg.Method, msg.Type)
		switch msg.Type {
		case "Response":
			// Pass message to sender
			t.requests[msg.Id] <- msg
		case "Request":
			// Handle request
			handleRequest(msg)
		default:
			log.Error("Message type not specified!")
		}
	}
}
