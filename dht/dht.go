package dht

import (
	"fmt"
	"strconv"
	"math/big"
	"math"
)

const m = 3

type DHTNode struct {
	id, adress, port string
	predecessor      *DHTNode
	fingerTable      [m+1]Finger
}

type Finger struct {
	startId string
	node    *DHTNode
}

func (n *DHTNode) successor() *DHTNode{
	return n.fingerTable[1].node
}

func (n *DHTNode) setSuccessor(successor *DHTNode) {
	n.fingerTable[1].node = successor
}

func (n *DHTNode) printRing() {
	fmt.Println(n.id)
	newn := n.successor()
	for newn.id != n.id {
		fmt.Println(newn.id)
		newn = newn.successor()
	}
}

func (nodeToAdd *DHTNode) join(n *DHTNode) {

	// If nodeToAdd is the only node in the network
	if n == nil {
		nodeToAdd.predecessor = nodeToAdd
		for i := 1; i <= m; i++ {
			nodeToAdd.fingerTable[i].node = nodeToAdd
		}
	} else {
		nodeToAdd.initFingerTable(n)
		nodeToAdd.updateOthers()
	}
}

// should be used in lookup and addToRing to find the right node / place in the ring
func (n *DHTNode) findSuccessor(id string) *DHTNode {
	predecessor := n.findPredecessor(id)
	return predecessor.successor()
}

func (n *DHTNode) findPredecessor(id string) *DHTNode {
	n2 := n
	for !between([]byte(n2.id), []byte(n2.successor().id), []byte(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *DHTNode) closestPrecedingFinger(id string) *DHTNode {
	for i := m; i > 0; i-- {
		if between([]byte(n.id), []byte(id), []byte(n.fingerTable[i].node.id)) {
			return n.fingerTable[i].node
		}
	}
	return n
}

func (nodeToUpdateTableOn *DHTNode) initFingerTable(n *DHTNode) {

	// Find successor node of nodeToUpdateTableOn using startId

	// Calculating first finger
	nodeToUpdateTableOn.fingerTable[1].startId, _ = calcFinger([]byte(nodeToUpdateTableOn.id), 1, m)
	// Successor to first finger
	nodeToUpdateTableOn.fingerTable[1].node = n.findSuccessor(nodeToUpdateTableOn.fingerTable[1].startId)

	// Set nodeToUpdateTableOns predecessor to the the node it's being inserted after
	nodeToUpdateTableOn.predecessor = nodeToUpdateTableOn.successor().predecessor

	// Update the predecessor of the node that nodeToUpdateTableOn is inserted before
	nodeToUpdateTableOn.successor().predecessor = nodeToUpdateTableOn

	for i := 1; i <= (m - 1) ; i++ {
		// Calculating finger
		nodeToUpdateTableOn.fingerTable[i+1].startId, _ = calcFinger([]byte(nodeToUpdateTableOn.id), i+1, m)
		if between(
			[]byte(nodeToUpdateTableOn.id),
			[]byte(nodeToUpdateTableOn.fingerTable[i].node.id),
			[]byte(nodeToUpdateTableOn.fingerTable[i+1].startId),
		) {
			nodeToUpdateTableOn.fingerTable[i+1].node = nodeToUpdateTableOn.fingerTable[i].node
		} else {
			nodeToUpdateTableOn.fingerTable[i+1].node = n.findSuccessor(nodeToUpdateTableOn.fingerTable[i+1].startId)
		}
	}
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	if between([]byte(n.id), []byte(n.successor().id), []byte(hashKey)) {
		return n
	} else {

		return n.successor().lookup(hashKey)
	}
}
/*
func (n *DHTNode) updateFingertable(k, m int) {
	calcFinger([]byte(n.id), k, m)
	//m := len(n.id)
	//for k:= 1; k < m; k++ {
	// or in place n.fingertable[i-1]
	//n.fingertable[k], _ = calcFinger([]byte(n.id), k, m)
	//}
}
*/
func (n *DHTNode) updateOthers() {
	for i := 1; i <= m; i++ {
		// find last node p whose i:th finger might point to n
		nId := big.Int{}
		nId.SetBytes([]byte(n.id))
		someId := nId - Exp2(i-1)
		p := findPredecessor(FormatInt(someId, 10))
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
