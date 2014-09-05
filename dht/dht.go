package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	successor        *DHTNode
}

func (n *DHTNode) printRing() {
	fmt.Println(n.id)
	newn := n.successor
	for newn.id != n.id {
		fmt.Println(newn.id)
		newn = newn.successor
	}
}

func (n *DHTNode) addToRing(nodeToAdd *DHTNode) {

	// If n is alone
	if n.successor == nil {
		nodeToAdd.successor = n
		n.successor = nodeToAdd
	} else if between([]byte(n.id), []byte(n.successor.id), []byte(nodeToAdd.id)) {
		nodeToAdd.successor = n.successor
		n.successor = nodeToAdd
	} else {
		n.successor.addToRing(nodeToAdd)
	}
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	if between([]byte(n.id), []byte(n.successor.id), []byte(hashKey)) {
		return n
	} else {
		return n.successor.lookup(hashKey)
	}
}

func makeDHTNode(idPointer *string, adress string, port string) *DHTNode {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}

	return &DHTNode{id: id, adress: adress, port: port}
}
