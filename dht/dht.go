package dht

import (
	"fmt"
	"math/big"
)

const m = 3

type DHTNode struct {
	id, adress, port string
	predecessor      *DHTNode
	fingerTable      [m + 1]Finger
}

type Finger struct {
	startId string
	node    *DHTNode
}

func (n *DHTNode) successor() *DHTNode {
	return n.fingerTable[1].node
}

func (n *DHTNode) setSuccessor(successor *DHTNode) {
	n.fingerTable[1].node = successor
}

func nextId(id string) string {
	nId := big.Int{}
	nId.SetString(id, 10)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Add(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)
	return y.String()
	//return id
}

func prevId(id string) string {
	nId := big.Int{}
	nId.SetString(id, 10)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Add(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)
	return y.String()
	//return id
}

func (n *DHTNode) printRing() {
	fmt.Println(n.id)
	newn := n.successor()
	for newn.id != n.id {
		fmt.Println(newn.id)
		newn = newn.successor()
	}
}

func (n *DHTNode) printFingers() {
	fmt.Printf("NodeId = %s\n", n.id)
	for i := 1; i <= m; i++ {
		fmt.Printf("FingerStartId = %s Node = %s\n", n.fingerTable[i].startId, n.fingerTable[i].node.id)
	}
}

func (nodeToAdd *DHTNode) join(n *DHTNode) {

	// If nodeToAdd is the only node in the network
	if n == nil {
		nodeToAdd.predecessor = nodeToAdd
		for i := 1; i <= m; i++ {
			nodeToAdd.fingerTable[i].startId, _ = calcFinger([]byte(nodeToAdd.id), i, m)
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
	for !between([]byte(nextId(n2.id)), []byte(n2.successor().id), []byte(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *DHTNode) closestPrecedingFinger(id string) *DHTNode {
	for i := m; i > 0; i-- {
		if between([]byte(nextId(n.id)), []byte(prevId(id)), []byte(n.fingerTable[i].node.id)) {
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
	fmt.Printf("Due to initFingerTable: Finger 1 for node %s is set to %s\n", nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[1].node.id)

	// Set nodeToUpdateTableOns predecessor to the the node it's being inserted after
	nodeToUpdateTableOn.predecessor = nodeToUpdateTableOn.successor().predecessor

	// Update the predecessor of the node that nodeToUpdateTableOn is inserted before
	nodeToUpdateTableOn.successor().predecessor = nodeToUpdateTableOn

	for i := 1; i <= (m - 1); i++ {
		// Calculating finger
		nodeToUpdateTableOn.fingerTable[i+1].startId, _ = calcFinger([]byte(nodeToUpdateTableOn.id), i+1, m)
		if between(
			[]byte(nodeToUpdateTableOn.id),
			[]byte(prevId(nodeToUpdateTableOn.fingerTable[i].node.id)),
			[]byte(nodeToUpdateTableOn.fingerTable[i+1].startId),
		) {
			nodeToUpdateTableOn.fingerTable[i+1].node = nodeToUpdateTableOn.fingerTable[i].node
			fmt.Printf("Due to initFingerTable: Finger %d for node %s is set to %s\n", i+1, nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[i+1].node.id)
		} else {
			nodeToUpdateTableOn.fingerTable[i+1].node = n.findSuccessor(nodeToUpdateTableOn.fingerTable[i+1].startId)
			fmt.Printf("Due to initFingerTable: Finger %d for node %s is set to %s\n", i+1, nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[i+1].node.id)
		}
	}
}

// traverse the ring counter-clockwise to update all nodes whose finger table entries should refer to n
func (n *DHTNode) updateOthers() {
	for i := 1; i <= m; i++ {

		// find last node p whose i:th finger might be n
		nId := big.Int{}
		nId.SetString(n.id, 10)

		y := big.Int{}
		two := big.NewInt(2)
		mbig := big.NewInt(m)

		y.Exp(two, big.NewInt(int64(i-1)), nil)
		y.Sub(&nId, &y)
		two.Exp(two, mbig, nil)
		y.Mod(&y, two)
		// y = nId - 2^(i-1)
		fmt.Printf("in updateOthers: nId=%s, i=%d y=%s\n", nId.String(), i, y.String())
		p := n.findPredecessor(nextId(y.String()))
		fmt.Printf("p = %s\n", p.id)

		p.updateFingerTable(n, i)
	}
}

// if s should be the i:th finger of n -> update n's finger table entry i with n
func (n *DHTNode) updateFingerTable(s *DHTNode, i int) {
	//	if (s.id != n.id) {
	if between(
		[]byte(n.id),
		[]byte(prevId(n.fingerTable[i].node.id)),
		[]byte(s.id),
	) {
		n.fingerTable[i].node = s
		fmt.Printf("Due to updateOthers: Node %s finger %d is set to %s\n", n.id, i, s.id)
		// get first node preceeding n
		p := n.predecessor
		p.updateFingerTable(s, i)
	}
	//	}
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	if between([]byte(n.id), []byte(n.successor().id), []byte(hashKey)) {
		return n
	} else {

		return n.successor().lookup(hashKey)
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
