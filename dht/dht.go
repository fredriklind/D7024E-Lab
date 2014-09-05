package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	successor        *DHTNode
	fingerTable      [][]string
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

func (n *DHTNode) findSuccessor(id string) *DHTNode {
	predecessor := n.findPredecessor(id)
	return predecessor.successor
}

func (n *DHTNode) findPredecessor(id string) *DHTNode {
	newn := n
	for !between([]byte(newn.id), []byte(newn.successor.id), []byte(id)) {
		newn = newn.closestPreceedingFinger(id)
	}
	return newn
}

func (n *DHTNode) closestPreceedingFinger(id string) *DHTNode {
	m := len([]byte(id))
	for m > 0 {
		if between([]byte(n.id), []byte(id), []byte(n.fingerTable[m][0].id)) {
			return n.fingerTable[m][1]
		}
		m -= 1
	}
	return n
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	if between([]byte(n.id), []byte(n.successor.id), []byte(hashKey)) {
		return n
	} else {
		return n.successor.lookup(hashKey)
	}
}

func (n *DHTNode) updateFingertable(k, m int) {
	calcFinger([]byte(n.id), k, m)
	//m := len(n.id)
	//for k:= 1; k < m; k++ {
		// or in place n.fingertable[i-1]
		//n.fingertable[k], _ = calcFinger([]byte(n.id), k, m)
	//}
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
