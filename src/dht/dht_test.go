package dht

import "testing"

func TestReceive(t *testing.T) {
	node := makeLocalNode(nil, "localhost", "2000")
	_ = node
	block := make(chan bool)
	<-block
}

func TestHELLO(t *testing.T) {
	node1 := makeLocalNode(nil, "localhost", "3000")
	node2 := &remoteNode{address: "localhost", port: "2000"}

	node1.ping(node2)
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
