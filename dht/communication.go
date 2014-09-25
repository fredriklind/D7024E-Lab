package dht

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const timeoutSeconds = time.Second * 4

type Msg struct {
	Id, Method, Src, Dst string
	Values               []string
	Channel              chan Msg
	WaitForResponse      bool
}

type Communication struct {
	listenAddress string
	channel       chan Msg
	stopListening chan bool
}

func (n *DHTNode) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", n.adress+":"+n.port)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		// Got message!
		n.Requests[msg.Id] <- msg
	}

	if err != nil {
		fmt.Printf("Error, %s", err.Error())
	}
}

func (n *DHTNode) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	msg.Id = generateNodeId()

	jsonMsg, err := json.Marshal(*msg)
	_, err = conn.Write([]byte(jsonMsg))

	if msg.WaitForResponse {
		// Blocks until something is received on the channel that is associated with msg.Id
		n.waitForResponse(msg)
	}

	if err != nil {
		fmt.Printf("Error, %s", err.Error())
	}
}

func (n *DHTNode) waitForResponse(msg *Msg) {
	n.Requests[msg.Id] = make(chan Msg)

	select {
	case response := <-n.Requests[msg.Id]:
		fmt.Printf("Received: %+v\n", response)
	case <-time.After(timeoutSeconds):
		fmt.Printf("Timed out message %s", msg.Id)
	}
}
