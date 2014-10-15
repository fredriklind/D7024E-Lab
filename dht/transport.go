package dht

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"code.google.com/p/go-uuid/uuid"
	log "github.com/cihub/seelog"
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

// add function for periodically clean request-array, use timestamps in the requests. <--------

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
	}, nil
}

func (t *transporter) sendPredecessorResponse(request msg) {
	m := msg{
		Id:     request.Id,
		Type:   "Response",
		Method: "PREDECESSOR",
		Dst:    request.Src,
		Values: dictionary{},
	}

	n := theLocalNode.predecessor()
	m.Values["id"] = n.id() // "6899"

	switch n.(type) {
	case *remoteNode:
		m.Values["address"] = n.address()
	default:
		m.Values["address"] = t.Address
	}
	t.send(m)
}

func (t *transporter) sendUpdatePredecessorCall(destAddr, candidateId, candidateAddr string) {
	m := msg{
		Type:   "Request",
		Method: "UPDATE_PREDECESSOR",
		Dst:    destAddr,
		Values: dictionary{},
	}
	m.Values["id"] = candidateId
	m.Values["address"] = candidateAddr
	t.send(m)
}

func (_ *transporter) handleUpdatePredecessorCall(call msg) {
	if call.Values["id"] == theLocalNode.id() {
		theLocalNode.updatePredecessor(theLocalNode)
	} else {
		n := &remoteNode{_id: call.Values["id"], _address: call.Values["address"]}
		theLocalNode.updatePredecessor(n)
	}
}

func (t *transporter) sendUpdateSuccessorCall(destAddr, candidateId, candidateAddr string) {
	m := msg{
		Type:   "Request",
		Method: "UPDATE_SUCCESSOR",
		Dst:    destAddr,
		Values: dictionary{},
	}
	m.Values["id"] = candidateId
	m.Values["address"] = candidateAddr
	t.send(m)
}

func (_ *transporter) handleUpdateSuccessorCall(call msg) {
	if call.Values["id"] == theLocalNode.id() {
		theLocalNode.updateSuccessor(theLocalNode)
	} else {
		n := &remoteNode{_id: call.Values["id"], _address: call.Values["address"]}
		theLocalNode.updateSuccessor(n)
	}
}

func (t *transporter) sendLookupRequest(destAddr, key string) (dictionary, error) {

	m := msg{
		Type:   "Request",
		Method: "LOOKUP",
		Dst:    destAddr,
		Values: dictionary{},
		Sync:   true,
	}

	m.Values["key"] = key
	m.Id = uuid.New()[0:4]
	m.Values["original_msgid"] = m.Id
	m.Values["original_src"] = t.Address

	response, err := t.send(m)

	if err != nil {
		return dictionary{}, err
	}

	return dictionary{
		"id":      response.Values["id"],
		"address": response.Values["address"],
	}, nil
}

func (t *transporter) handleLookupRequest(request msg) {
	mg := msg{
		Id:     request.Id,
		Type:   "Response",
		Method: "LOOKUP_ACK",
		Dst:    request.Src,
		Values: dictionary{},
	}
	mg.Values["original_src"] = request.Values["original_src"]
	t.send(mg)

	if between(
		hexStringToByteArr(nextId(theLocalNode.predecessor().id())),
		hexStringToByteArr(nextId(theLocalNode.id())),
		hexStringToByteArr(request.Values["key"]),
	) {
		mg := msg{
			Id:     request.Values["original_msgid"],
			Type:   "Response",
			Method: "LOOKUP",
			Dst:    request.Values["original_src"],
			Values: dictionary{},
		}
		mg.Values["id"] = theLocalNode.id()
		mg.Values["address"] = t.Address
		t.send(mg)
	} else {

		for k := m; k > 0; {
			nextNode, i := theLocalNode.forwardingLookup(request.Values["key"], k)
			request.Dst = nextNode.address()
			request.Id = ""

			_, err := t.send(request)

			if err == nil {
				break
			} else {
				k = i - 1
			}
		}
	}
}

func (t *transporter) sendHelloRequest(destAddr string) {
	m := msg{
		Type:   "Request",
		Method: "HELLO",
		Dst:    destAddr,
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

func (t *transporter) ping(destAddr string) bool {
	m := msg{
		Type:   "Request",
		Method: "PING",
		Dst:    destAddr,
		Sync:   true,
	}
	_, err := t.send(m)
	if err == nil {
		return true
	} else {
		return false
	}
}

// ----------------------------------------------------------------------------------------
//										middle layer
// ----------------------------------------------------------------------------------------

// Parse requests and pass them on to their handler
func (t *transporter) handleRequest(request msg) {
	switch request.Method {
	case "HELLO":
		t.sendAck(request)

	case "PREDECESSOR":
		t.sendPredecessorResponse(request)

	case "UPDATE_SUCCESSOR":
		t.handleUpdateSuccessorCall(request)

	case "UPDATE_PREDECESSOR":
		t.handleUpdatePredecessorCall(request)

	case "LOOKUP":
		t.handleLookupRequest(request)
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

func (t *transporter) waitForLookupResponse(ackMsg msg) (msg, error) {

	// Got LOOKUP_ACK, wait for

	if ackMsg.Values["original_src"] == t.Address {
		select {
		// wait for real response
		case responseMsg := <-t.requests[ackMsg.Id]:
			return responseMsg, nil

		case <-time.After(time.Second * 10):
			return msg{}, errors.New("Timeout")

		}
	} else {
		// OK, got ACK, not original sender - don´t wait for any other response
		return msg{}, nil
	}
}

func (t *transporter) waitForResponse(msgId string, waitSeconds time.Duration) (msg, error) {

	// Save the channel so that the receive() method can un-block
	// this method when it receives a response with a matching id
	t.requests[msgId] = make(chan msg)

	// Wait for the msg-specific channel to get data, or time out
	select {

	case responseMsg := <-t.requests[msgId]:

		if responseMsg.Method == "LOOKUP_ACK" {
			return t.waitForLookupResponse(responseMsg)
		} else {
			return responseMsg, nil
		}

	case <-time.After(time.Second * waitSeconds):
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
	if m.Method == "LOOKUP" {
		if m.Type == "Request" {
			log.Tracef("%s: Sent %s %s: %s to %s key: %s origMsgId:%s origSrc:%s", t.Address, m.Method, m.Type, m.Id, m.Dst, m.Values["key"], m.Values["original_msgid"], m.Values["original_src"])
		} else if m.Type == "Response" {
			log.Tracef("%s: Sent %s %s: %s to %s respNodeId:%s", t.Address, m.Method, m.Type, m.Id, m.Dst, m.Values["id"])
		}
	} else if m.Method == "LOOKUP_ACK" {
		log.Tracef("%s: Sent %s %s: %s", t.Address, m.Method, m.Type, m.Id)
	} else {
		fmt.Printf("%s: Sent %s %s: %+v\n", t.Address, m.Method, m.Type, m)
		log.Tracef("%s: Sent %s %s: %s to %s", t.Address, m.Method, m.Type, m.Id, m.Dst)
	}

	time.Sleep(time.Second * 1)

	// Serialize and send the message (also wait to simulate network delay)
	jsonmsg, err := json.Marshal(m)
	_, err = conn.Write([]byte(jsonmsg))

	if err != nil {
		log.Errorf("%s: error in send: %s", t.Address, err.Error())
	}

	// Blocks until something is received on the channel that is associated with m.Id
	if m.isRequest() && m.Sync {
		return t.waitForResponse(m.Id, 5)
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
	m := msg{}

	err = dec.Decode(&m)
	fmt.Printf("%s: Got %s %s: %s\n", t.Address, m.Method, m.Type, m.Id)
	log.Tracef("%s: Got %s %s: %s", t.Address, m.Method, m.Type, m.Id)
	conn.Close()
	go t.receive()

	switch m.Type {
	case "Response":
		// Send response to the waiting request sender, or time out if no
		// one is waiting.
		select {
		case t.requests[m.Id] <- m:
		case <-time.After(time.Millisecond * 300):
		}

	case "Request":
		// Handle request
		t.handleRequest(m)
	default:
		log.Error("Message type not specified!")
	}
}
