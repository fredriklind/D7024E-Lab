package transport

import (
	"encoding/json"
	"net"
	"time"

	log "github.com/cihub/seelog"
)

const timeoutSeconds = time.Second * 8

type Transporter

type Msg struct {
	Id, Type, Method, Src, Dst string
	Timestamp                  int64
	Values                     map[string]string
	Sync                       bool
}

func (m *Msg) isRequest() bool  { return m.Type == "Request" }
func (m *Msg) isResponse() bool { return m.Type == "Response" }

func (n *localNode) getAddress() string {
	return n.address + ":" + n.port
}

// ----------------------------------------------------------------------------------------
//										public layer
// ----------------------------------------------------------------------------------------

func (t *transport) sendLookupRequest(destNode *remoteNode, id string) {
	// check queue
	// if lookup in queue - forward request
	// else send new lookupRequest
	lookupRequest := Msg{
		Method: "LOOKUP",
		Dst:    destNode.getAddress(),
	}
	t.sendRequest(lookupRequest)
}

// ----------------------------------------------------------------------------------------
//										middle layer
// ----------------------------------------------------------------------------------------

func (t *transport) sendRequest(request Msg) {
	// Construct the message
	msg := Msg{
		Type:   "Request",
		Method: request.Method,
		Src:    t.address,
		Dst:    request.Dst,
	}

	// This is to prevent send from generating a new msg.Id
	if request.Id != "" {
		msg.Id = request.Id
	}
	n.send(msg)
}

func (n *localNode) sendResponse(response Msg) {
	//Construct the message
	msg := Msg{
		Type:   "Response",
		Method: response.Method,
		Src:    n.getAddress(),
		Dst:    response.Dst,
	}

	// This is to prevent send from generating a new msg.Id
	if response.Id != "" {
		msg.Id = response.Id
	}
	n.send(msg)
}

func (n *localNode) waitForResponse(request Msg) {

	// Save the channel so that the receive() method can un-block
	// this method when it receives a response with a matching id
	n.Requests[request.Id] = make(chan Msg)

	// Wait for the Msg-specific channel to get data, or time out
	select {
	case response := <-n.Requests[request.Id]:
		n.handleResponse(response)
	case <-time.After(timeoutSeconds):
		log.Errorf("Node %s %s request with id %s timed out", n.id, request.Method, request.Id)
	}
}

// When you get a request you need to handle (you = n)
func (n *localNode) handleRequest(request Msg) {
	//log.Tracef("Got request: %+v", request)
	switch request.Method {
	case "HELLO":
		n.sendResponse(Msg{
			Id:     request.Id,
			Method: "ACK",
			Dst:    request.Src,
		})

	// Forwards a request to nextNode setting the method and Src depending
	// on the Values["Method"] and Values["Sender"]
	case "FORWARD":
		// If n is the final destination, answer the original sender
		if n.id == request.Values["FinalDestinationId"] {
			log.Tracef("Node %s is FinalDestination", n.id)
			newRequest := Msg{
				Method: request.Values["Method"],
				Src:    request.Values["Sender"],
			}
			// Handle the request contained in the FORWARD request
			n.handleRequest(newRequest)
		} else {
			// If n is not the searched for node, forward the request to the next node
			nextNodeAddress := "127.0.0.1:4000"
			forwardRequest := Msg{
				Id:     request.Id,
				Method: "FORWARD",
				Dst:    nextNodeAddress,
				Sync:   false,
			}
			n.sendRequest(forwardRequest)
		}

	default:
		log.Error("No request method specified!")
	}
}

// When you get something back from a request you sent (you = n)
func (n *localNode) handleResponse(response Msg) {
	switch response.Method {
	case "ACK":
		//log.Tracef("Node %s request satisfied", n.id)
	default:
		log.Debugf("No request method specified!")
	}
}

// ----------------------------------------------------------------------------------------
//										lowest layer
// ----------------------------------------------------------------------------------------

func (n *localNode) send(msg Msg) {

	// Start up network stuff
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	// Apply unix timestamp to the Msg
	t := time.Now()
	msg.Timestamp = t.Unix()

	// Assign a new Id to the Msg if it is not set
	if msg.Id == "" {
		msg.Id = generateNodeId()[0:4]
	}

	log.Tracef("Node %s sent %s %s", n.id, msg.Method, msg.Type)
	//log.Tracef("SENDING!: %+v", msg)

	// Serialize and send the message (also wait to simulate network delay)
	time.Sleep(time.Millisecond * 400)
	jsonMsg, err := json.Marshal(msg)
	_, err = conn.Write([]byte(jsonMsg))

	// Blocks until something is received on the channel that is associated with msg.Id
	if msg.isRequest() && msg.Sync {
		n.waitForResponse(msg)
	}

	if err != nil {
		log.Errorf("Node %s Error in send: %s", n.id, err.Error())
	}
}

func (n *localNode) receive() {

	// Start receiveing
	udpAddr, err := net.ResolveUDPAddr("udp", n.getAddress())
	conn, err := net.receiveUDP("udp", udpAddr)

	if err != nil {
		log.Errorf("Error in send on node %s, %s", n.id, err.Error())
		return
	}

	defer conn.Close()

	// receive again after this method finishes. TODO might be
	// that there is a better way to do this
	defer n.receive()
	dec := json.NewDecoder(conn)

	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		//log.Tracef("Node %s got %s %s", n.id, msg.Method, msg.Type)
		time.Sleep(time.Millisecond * 400)
		switch msg.Type {
		case "Response":
			// Pass message to sender
			n.Requests[msg.Id] <- msg
		case "Request":
			// Handle request
			n.handleRequest(msg)
		default:
			log.Error("Message type not specified!")
		}
	}
}
