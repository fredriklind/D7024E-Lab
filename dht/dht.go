package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	successor        *DHTNode
	fingerTable      []Finger
}

type Finger struct {
	startId string
	node    *DHTNode
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
		m := len(n.id)
		for k:= 1; k < m; k++ {
			n.fingerTable[k].startId = n.id
			n.fingerTable[k].node = n
		}
	} else if between([]byte(n.id), []byte(n.successor.id), []byte(nodeToAdd.id)) {
		nodeToAdd.successor = n.successor
		n.successor = nodeToAdd
	} else {
		n.successor.addToRing(nodeToAdd)
	}
}

// should be used in lookup and addToRing to find the right node / place in the ring
func (n *DHTNode) findSuccessor(id string) *DHTNode {
	predecessor := n.findPredecessor(id)
	return predecessor.successor
}

func (n *DHTNode) findPredecessor(id string) *DHTNode {
	n2 := n
	for !between([]byte(n2.id), []byte(n2.successor.id), []byte(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *DHTNode) closestPrecedingFinger(id string) *DHTNode {
	m := len([]byte(id))
	for m > 0 {
		if between([]byte(n.id), []byte(id), []byte(n.fingerTable[m].node.id)) {
			return n.fingerTable[m].node
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
