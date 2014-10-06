package dht

import "testing"

func TestReceive2000(t *testing.T) {
	id := "01"
	node := makeLocalNode(&id, "localhost", "2000")
	_ = node
	block := make(chan bool)
	<-block
}

func TestReceive3000(t *testing.T) {
	id := "02"
	node := makeLocalNode(&id, "localhost", "3000")
	_ = node
	block := make(chan bool)
	<-block
}

func TestHELLO(t *testing.T) {
	// Define a sequence of requests that are expected
	/*setupTest(t, []string{
		"Node 01 sent HELLO Request",
		"Node 02 got HELLO Request",
		"Node 02 sent ACK Response",
		"Node 01 got ACK Response",
	},
	)*/

	id1 := "01"
	id2 := "02"
	node1 := makeLocalNode(&id1, "127.0.0.1", "2000")
	node2 := makeLocalNode(&id2, "127.0.0.1", "3000")

	_ = node1
	_ = node2

	/*	sendRequest(msg{
		Method: "HELLO",
		Dst:    node2.getAddress(),
	})*/

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
