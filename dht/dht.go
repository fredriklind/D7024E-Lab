package dht

import (
	"transport"
)

const m = 160

type node interface {
	lookup(id string) *node
	predecessor() *node
	// Tell *node to check if node is its successor/predecessor
	updateSuccessor(*node)
	updatePredecessor(*node)
}

type localNode struct {
	id          string
	predecessor *localNode
	fingerTable [m + 1]Finger
	isListening chan bool
}

type remoteNode struct {
	id, address, port string
	owner             *node
}

type finger struct {
	startId string
	node    *node
}

func makelocalNode(idPointer *string, address string, port string) *localNode {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}
	node := localNode{id: id}
	transport.NewTransporter(address, port)
	return &node
}

func (n *localNode) successor() *localNode {
	return n.fingerTable[1].node
}

func (n *localNode) setSuccessor(successor *localNode) {
	n.fingerTable[1].node = successor
}
