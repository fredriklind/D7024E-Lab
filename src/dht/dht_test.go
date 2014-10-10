package dht

import (
	//log "github.com/cihub/seelog"
	"testing"
)

func TestReceive(t *testing.T) {
	id := "5"
	newLocalNode(&id, "localhost", "2000")

	node2 := &remoteNode{_id: "4", _address: "localhost:6600"}
	theLocalNode.pred = node2

	block := make(chan bool)
	<-block
}

func TestHELLO(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	theLocalNode.ping(node2)
	block := make(chan bool)
	<-block
}

func TestPredecessorRequest(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	_ = node2.predecessor()
	block := make(chan bool)
	<-block
}

func TestUpdateSuccessorCall(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	candidate := &remoteNode{_id: "8888", _address: "localhost:8877"}
	node2.updateSuccessor(candidate)
	block := make(chan bool)
	<-block
}

func TestUpdatePredecessorCall(t *testing.T) {
	newLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{_address: "localhost:2000"}

	candidate := &remoteNode{_id: "3", _address: "localhost:8877"}
	node2.updatePredecessor(candidate)
	block := make(chan bool)
	<-block
}

/*
func Test3NodeForwarding(t *testing.T) {
	block := make(chan bool)

	id1 := "01"
	id2 := "02"
	id3 := "03"

	node1 := makeLocalNode(&id1, "127.0.0.1", "2000")
	node2 := makeLocalNode(&id2, "127.0.0.1", "3000")
	node3 := makeLocalNode(&id3, "127.0.0.1", "4000")

	node1.sendRequest(
		msg{
			Method: "FORWARD",
			Values: map[string]string{
				"Method":             "HELLO",
				"FinalDestinationId": "03",
				"Sender":             node1.getAddress(),
			},
			Dst: node2.getAddress()},
	)

	// To prevent stupid warnings
	_ = node3
	<-block
}*/
