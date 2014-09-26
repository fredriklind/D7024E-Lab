package dht

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Msg struct {
	key, src, dst string
}

type Transport struct {
	listenAddress string
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.listenAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)
		fmt.Printf("%+v\n", msg)
	}

	if err != nil {
		fmt.Printf("Error, %s", err.Error())
	}
}

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.dst)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	jsonMsg, err := json.Marshal(&msg)
	fmt.Printf("%+v\n", msg)
	os.Stdout.Write(jsonMsg)

	_, err = conn.Write([]byte(jsonMsg))

	if err != nil {
		fmt.Printf("Error, %s", err.Error())
	}
}
