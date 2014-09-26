package dht

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

const timeoutSeconds = time.Second * 8

type Msg struct {
	Id, Type, Method, Src, Dst string
	Timestamp                  int64
	Values                     map[string]string
	DontWaitForResponse        bool
}

func (m *Msg) isRequest() bool  { return m.Type == "Request" }
func (m *Msg) isResponse() bool { return m.Type == "Response" }

func (n *DHTNode) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", n.getAddress())
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	defer n.listen()
	dec := json.NewDecoder(conn)

	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		log.Tracef("Node %s received msgId=%s %s %s from %s %+v", n.id, msg.Id, msg.Method, msg.Type, msg.Src, msg)
		time.Sleep(time.Second * 2)
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
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	t := time.Now()
	msg.Timestamp = t.Unix()

	if msg.Id == "" {
		msg.Id = generateNodeId()[0:4]
	}
	log.Tracef("Node %s sending msgId=%s %s %s to %s Timestamp: %d", n.id, msg.Id, msg.Method, msg.Type, msg.Dst, msg.Timestamp)
	time.Sleep(time.Second * 2)
	jsonMsg, err := json.Marshal(msg)
	_, err = conn.Write([]byte(jsonMsg))

	if msg.isRequest() && msg.DontWaitForResponse == false {
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
	msg := Msg{
		Type:   "Request",
		Method: request.Method,
		Src:    n.getAddress(),
		Dst:    request.Dst,
	}
	n.send(msg)
}

// When you get a request you need to handle (you = n)
func (n *DHTNode) handleRequest(request Msg) {

	// Construct the response to be sent
	response := Msg{Type: "Response", Src: n.getAddress()}

	switch request.Method {
	case "HELLO":
		response.Id = request.Id
		response.Method = "ACK"
		response.Dst = request.Src
		n.send(response)

	// Forwards a request to nextNode setting the method and Src depending
	// on the Values["Method"] and Values["Sender"]
	case "FORWARD":
		// If n is the final destination, answer the original sender
		if n.id == request.Values["FinalDestinationId"] {
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
				Id:                  request.Id,
				Method:              "FORWARD",
				Dst:                 nextNodeAddress,
				DontWaitForResponse: true,
			}
			n.sendRequest(forwardRequest)
		}

	default:
		log.Error("No request method specified!")
	}
}

// When you get something back from a request you sent (you = n)
func (n *DHTNode) handleResponse(response Msg) {
	switch response.Method {
	case "ACK":
		// Do nothing
	default:
		log.Debugf("No request method specified!")
	}
}
