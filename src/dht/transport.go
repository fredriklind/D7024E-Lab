package dht

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

// ----------------------------------------------------------------------------------------
//										Init, types and variables
// ----------------------------------------------------------------------------------------

const timeoutSeconds = time.Second * 8

var theLocalNode *localNode

type transporter struct {
	Address  string
	requests map[string]chan msg
}

type msg struct {
	Id, Type, Method, Src, Dst string
	Timestamp                  int64
	Values                     dictionary
	Sync                       bool
}

type dictionary map[string]string

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
func newTransporter(address string) *transporter {
	t := &transporter{}
	t.Address = address
	t.requests = make(map[string]chan msg)
	go t.receive()
	return t
}

func (t *transporter) sendPredecessorRequest(destAddr string) (dictionary, error) {
	m := msg{
		Type:   "Request",
		Method: "PREDECESSOR",
		Dst:    destAddr,
		Sync:   true,
	}
	response, err := t.send(m)

	if err != nil {
		return dictionary{}, err
	}

	return dictionary{
		"id":      response.Values["id"],
		"address": response.Values["address"],
		"port":    response.Values["port"],
	}, nil
}

func (t *transporter) sendPredecessorResponse(request msg) {
	m := msg{
		Id:     request.Id,
		Type:   "Response",
		Method: "PREDECESSOR",
		Dst:    request.Src,
	}
	//n := theLocalNode.predecessor()
	m.Values["id"] = "6899" //n.id()
	if false {              //n.id() == theLocalNode.id() {
		m.Values["address"] = t.Address
	} else {
		m.Values["address"] = "localhost:6000" //n.address()
	}
	log.Trace("sendPredecessorResponse")
	_, err := t.send(m)
	log.Errorf("Error in send: %s", err)
}

func (t *transporter) sendHelloRequest(address string) {
	m := msg{
		Type:   "Request",
		Method: "HELLO",
		Dst:    address,
	}
	t.send(m)
}

func (t *transporter) sendAck(request msg) {
	m := msg{
		Type:   "Response",
		Method: "ACK",
		Dst:    request.Src,
	}
	t.send(m)
}

// ----------------------------------------------------------------------------------------
//										middle layer
// ----------------------------------------------------------------------------------------

// When you get a request you need to handle
func (t *transporter) handleRequest(request msg) {
	switch request.Method {
	case "HELLO":
		t.sendAck(request)

	case "PREDECESSOR":
		log.Trace("Handlerequest")
		t.sendPredecessorResponse(request)
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
func (t *transporter) handleResponse(response msg) {
	switch response.Method {
	case "ACK":
		log.Tracef("%s HELLO request satisfied", t.Address)
	default:
		log.Errorf("No request method specified!")
	}
}

func (m *msg) isRequest() bool  { return m.Type == "Request" }
func (m *msg) isResponse() bool { return m.Type == "Response" }

// ----------------------------------------------------------------------------------------
//										lowest layer
// ----------------------------------------------------------------------------------------

func (t *transporter) waitForResponse(msgId string) (msg, error) {

	// Save the channel so that the receive() method can un-block
	// this method when it receives a response with a matching id
	t.requests[msgId] = make(chan msg)

	// Wait for the msg-specific channel to get data, or time out
	select {
	case responseMsg := <-t.requests[msgId]:
		return responseMsg, nil
	case <-time.After(timeoutSeconds):
		log.Errorf("%s: request with id %s timed out", t.Address, msgId)
		return msg{}, errors.New("Timeout")
	}
}

func (t *transporter) send(m msg) (msg, error) {
	m.Src = t.Address
	// Start up network stuff
	udpAddr, err := net.ResolveUDPAddr("udp", m.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	// Apply unix timestamp to the msg
	ti := time.Now()
	m.Timestamp = ti.Unix()

	// Assign a new Id to the msg if it is not set
	if m.Id == "" {
		m.Id = uuid.New()[0:4]
	}

	//	log.Tracef("%s: Sent %s %s", n.id, m.Method, m.Type)
	fmt.Printf("%s: Sent %s %s\n", t.Address, m.Method, m.Type)
	log.Tracef("%s: Sent %s %s", t.Address, m.Method, m.Type)

	// Serialize and send the message (also wait to simulate network delay)
	jsonmsg, err := json.Marshal(m)
	_, err = conn.Write([]byte(jsonmsg))

	if err != nil {
		log.Errorf("%s: Error in send: %s", t.Address, err.Error())
		return msg{}, err
	}

	// Blocks until something is received on the channel that is associated with m.Id
	if m.isRequest() && m.Sync {
		return t.waitForResponse(m.Id)
	} else {
		return msg{}, nil
	}
}

func (t *transporter) receive() {
	// Start receiveing
	udpAddr, err := net.ResolveUDPAddr("udp", t.Address)
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Errorf("%s: Error in receive, %s", t.Address, err.Error())
		return
	}

	dec := json.NewDecoder(conn)
	msg := msg{}

	err = dec.Decode(&msg)
	fmt.Printf("%s: Got %s %s\n", t.Address, msg.Method, msg.Type)
	log.Tracef("%s: Got %s %s", t.Address, msg.Method, msg.Type)
	conn.Close()
	go t.receive()

	switch msg.Type {
	case "Response":
		// Pass message to sender
		t.requests[msg.Id] <- msg
	case "Request":
		// Handle request
		log.Trace("In receive, a request")
		t.handleRequest(msg)
	default:
		log.Error("Message type not specified!")
	}
}
