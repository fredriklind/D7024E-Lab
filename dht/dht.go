package dht

import (
	"encoding/json"
	"fmt"
	"math/big"
)

const m = 160
const base = 16

type node interface {
}

type localNode struct {
	id, address, port string
	predecessor       *localNode
	fingerTable       [m + 1]Finger
	Requests          map[string]chan Msg
	isListening       chan bool
}

type finger struct {
	startId string
	node    *localNode
}

func makelocalNode(idPointer *string, address string, port string) *localNode {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}
	node := localNode{id: id, address: address, port: port}
	node.Requests = make(map[string]chan Msg)
	go node.listen()
	return &node
}

func (n *localNode) successor() *localNode {
	return n.fingerTable[1].node
}

func (n *localNode) setSuccessor(successor *localNode) {
	n.fingerTable[1].node = successor
}
