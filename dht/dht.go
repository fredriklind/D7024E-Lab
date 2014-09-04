package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	successor        *DHTNode
}

func (n DHTNode) printRing() {
	fmt.Println(n.id)
	id = n.successor.id
	for id != n.id {
		fmt.Println(id)
		id = n.successor.id
	}
}

func (n DHTNode) addToRing(addedNode *DHTNode) {
	n.successor = addedNode
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
