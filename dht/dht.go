package dht

import (
	"fmt"
)

type DHTNode struct {
	id, adress, port string
	predecessor      *DHTNode
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
		nodeToAdd.predecessor = n
		nodeToAdd.successor = n
		n.predecessor = nodeToAdd
		n.successor = nodeToAdd
		m := len([]byte(n.id))
		fmt.Println(m)
		for k := 1; k < m; k++ {
			n.fingerTable[k].startId = n.id
			n.fingerTable[k].node = n
		}
	} else {
		node2 := n.findSuccessor(nodeToAdd.id)
		node1 := node2.predecessor

		nodeToAdd.predecessor = node1
		nodeToAdd.successor = node2

		node1.successor = nodeToAdd
		node2.predecessor = nodeToAdd
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

func (n *DHTNode) initFingerTable(nodeToUpdateTableOn *DHTNode) {

	// Find successor node of nodeToUpdateTableOn using startId,
	// ---- IS startID SET? ----
	nodeToUpdateTableOn.fingerTable[1].node = n.findSuccessor(nodeToUpdateTableOn.fingerTable[1].startId)

	// Set nodeToUpdateTableOn as the new predecessor the node it's being inserted after
	nodeToUpdateTableOn.predecessor = nodeToUpdateTableOn.fingerTable[1].node.predecessor

	// Update the old predecessor of the node that nodeToUpdateTableOn is inserted before
	nodeToUpdateTableOn.fingerTable[1].node.predecessor = nodeToUpdateTableOn

	m := len([]byte(n.id))

	for i := 1; i <= (m - 1); i++ {
		if between(
			[]byte(nodeToUpdateTableOn.id),
			[]byte(nodeToUpdateTableOn.fingerTable[i].node.id),
			[]byte(nodeToUpdateTableOn.fingerTable[i+1].startId),
		) {
		} else {
			nodeToUpdateTableOn.fingerTable[i+1].node = n.findSuccessor(nodeToUpdateTableOn.fingerTable[i+1].startId)
		}
	}
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

func (n *DHTNode) updateOthers() {

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
