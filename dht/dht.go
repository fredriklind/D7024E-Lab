package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port *string
}

func (n DHTNode) printRing() {
	fmt.Println(n.id)
}


func makeDHTNode(id *string, adress *string, port *string) *DHTNode {

	if id == nil {
		id = generateNodeId()
	}

	return &DHTNode{id: id, adress: adress, port: port}
}
