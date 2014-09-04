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
