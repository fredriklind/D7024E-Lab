package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	successor        *DHTNode
	fingertable      []string
}

func (n *DHTNode) printRing() {
	fmt.Println(n.id)
	newn := n.successor
	for newn.id != n.id {
		fmt.Println(newn.id)
		newn = newn.successor
	}
}

func (n *DHTNode) addToRing(addedNode *DHTNode) {
	n.successor = addedNode
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	if between([]byte(n.id), []byte(n.successor.id), []byte(hashKey)) {
		return n
	} else {
		return n.successor.lookup(hashKey)
	}
}

func (n *DHTNode) updateFingertable() {
	// (n + 2^(k-1)) mod (2^m)
	// calcFinger(n []byte, k int, m int) (string, []byte) {
	m := len(n.id)
	for i:= 1; i < m; i++ {
		// or in place n.fingertable[i-1]
		n.fingertable[i], _ = calcFinger([]byte(n.id), i, m)
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
