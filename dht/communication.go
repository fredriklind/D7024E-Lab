package dht

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

const timeoutSeconds = time.Second * 4

type Msg struct {
	Id, Type, Method, Src, Dst string
	Values                     []string
	Channel                    chan Msg
	WaitForResponse            bool
}

func (m *Msg) isRequest() bool  { return m.Type == "Request" }
func (m *Msg) isResponse() bool { return m.Type == "Response" }

func (n *DHTNode) listen() {
	log.Debugf("Node %s listening on %s\n", n.id, n.getAddress())
	udpAddr, err := net.ResolveUDPAddr("udp", n.getAddress())
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		// Got a message!
		log.Debug("Got message!")
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

	if err != nil {
		log.Errorf("Error in send on node %s, %s", n.id, err.Error())
	}
}

func (n *DHTNode) send(msg Msg) {
	log.Debugf("Sending %s %s to %s", msg.Type, msg.Method, msg.Dst)
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	msg.Id = generateNodeId() //[0:4]

	jsonMsg, err := json.Marshal(msg)
	_, err = conn.Write([]byte(jsonMsg))

	if msg.isRequest() {
		// Blocks until something is received on the channel that is associated with msg.Id
		n.waitForResponse(msg)
	}

	if err != nil {
		log.Errorf("Error in send on node %s, %s", n.id, err.Error())
	}
}

func (n *DHTNode) waitForResponse(request Msg) {
	n.Requests[request.Id] = make(chan Msg)

	select {
	case response := <-n.Requests[request.Id]:
		n.handleResponse(response)
	case <-time.After(timeoutSeconds):
		log.Debugf("%s request with id %s timed out", request.Method, request.Id)
	}
}

func (n *DHTNode) sendRequest(request Msg) {
	msg := &Msg{
		Type:   "Request",
		Method: request.Method,
		Src:    n.getAddress(),
		Dst:    request.Dst,
	}
	n.send(*msg)
}

// When you get a request you need to handle (you = n)
func (n *DHTNode) handleRequest(request Msg) {

	// Construct the response to be sent
	response := Msg{Type: "Response", Src: n.getAddress()}

	switch request.Method {
	case "HELLO":
		response.Method = "ACK"
		response.Dst = request.Src
		n.send(response)

	default:
		log.Error("No request method specified!")
	}
}

// When you get something back from a request you sent (you = n)
func (n *DHTNode) handleResponse(response Msg) {
	switch response.Method {
	case "ACK":
		log.Debugf("Received ACK from %s", response.Src)
	default:
		log.Debugf("No request method specified!")
	}
}
