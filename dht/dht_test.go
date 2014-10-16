package dht

import (
	//"fmt"
	log "github.com/cihub/seelog"
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

	node2.isAlive()
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

func TestNode0(t *testing.T) {
	id := "00"
	newLocalNode(&id, "localhost", "9000")
	node4 := &remoteNode{_id: "04", _address: "localhost:4000"}
	node6 := &remoteNode{_id: "06", _address: "localhost:6000"}
	theLocalNode.pred = node6
	//theLocalNode.fingerTable[1].node = node4
	//theLocalNode.fingerTable[2].node = node4
	theLocalNode.fingerTable[3].node = node4

	node2 := &remoteNode{_id: "02", _address: "localhost:2000"}
	theLocalNode.fingerTable[1].node = node2
	theLocalNode.fingerTable[2].node = node2

	key := "03"
	n, _ := theLocalNode.lookup(key)
	log.Tracef("%s.lookup(%s) = %s", theLocalNode.id(), key, n.id())

	block := make(chan bool)
	<-block
}

func TestNode4(t *testing.T) {
	id := "04"
	newLocalNode(&id, "localhost", "4000")
	node0 := &remoteNode{_id: "00", _address: "localhost:9000"}
	node6 := &remoteNode{_id: "06", _address: "localhost:6000"}
	theLocalNode.pred = node0
	theLocalNode.fingerTable[1].node = node6
	theLocalNode.fingerTable[2].node = node6
	theLocalNode.fingerTable[3].node = node0

	block := make(chan bool)
	<-block
}

func TestNode6(t *testing.T) {
	id := "06"
	newLocalNode(&id, "localhost", "6000")
	node4 := &remoteNode{_id: "04", _address: "localhost:4000"}
	node0 := &remoteNode{_id: "00", _address: "localhost:9000"}
	theLocalNode.pred = node4
	theLocalNode.fingerTable[1].node = node0
	theLocalNode.fingerTable[2].node = node0
	theLocalNode.fingerTable[3].node = node4

	block := make(chan bool)
	<-block
}

// Run TestJoin3, TestJoin0 and TestJoin2 in that order from three separate tabs in terminal. (To test obj2).

func TestJoin3(t *testing.T) {
	id := "03"
	newLocalNode(&id, "localhost", "3000")

	theLocalNode.join(nil)

	block := make(chan bool)
	<-block
}

func TestJoin0(t *testing.T) {
	id := "00"
	newLocalNode(&id, "localhost", "9000")

	node3 := &remoteNode{_id: "03", _address: "localhost:3000"}

	theLocalNode.join(node3)
	block := make(chan bool)
	<-block
}

func TestJoin2(t *testing.T) {
	id := "02"
	newLocalNode(&id, "localhost", "2000")

	node3 := &remoteNode{_id: "03", _address: "localhost:3000"}

	theLocalNode.join(node3)

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
