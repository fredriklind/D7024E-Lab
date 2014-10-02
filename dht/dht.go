package dht

const m = 160
const base = 16

type node interface {
	lookup(id string) *node
	predecessor() *node
	// Tell *node to check if node is its successor/predecessor
	updateSuccessor(*node)
	updatePredecessor(*node)
}

type localNode struct {
	id, address, port string
	predecessor       *localNode
	fingerTable       [m + 1]Finger
	Requests          map[string]chan Msg
	isListening       chan bool
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
