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
	var initialLength = len(id)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Add(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)

	var returnString = y.String()

	for len(returnString) != initialLength {
		returnString = "0" + returnString
	}

	return returnString
	//	return id
}

func prevId(id string) string {
	nId := big.Int{}
	nId.SetString(id, 10)
	var initialLength = len(id)

	y := big.Int{}
	two := big.NewInt(2)
	one := big.NewInt(1)
	mbig := big.NewInt(m)

	y.Sub(&nId, one)
	// 2^m
	two.Exp(two, mbig, nil)
	y.Mod(&y, two)

	var returnString = y.String()

	for len(returnString) != initialLength {
		returnString = "0" + returnString
	}

	return returnString
	//	return id
}

func initTwoNodeRing(node1, node2 *DHTNode) {
	node1.predecessor = node2
	node2.predecessor = node1
	for i := 1; i <= m; i++ {
		node1.fingerTable[i].startId, _ = calcFinger([]byte(node1.id), i, m)

		node2.fingerTable[i].startId, _ = calcFinger([]byte(node2.id), i, m)
		node2.fingerTable[i].node = node1
	}
	node1.fingerTable[1].node = node2
	node1.fingerTable[2].node = node1
	node1.fingerTable[3].node = node1
}

func (n *DHTNode) printRing() {
	n.printNode()
	var visited []string
	visited = append(visited, n.id)
	newn := n.successor()

	for !stringInSlice(newn.id, visited) {
		newn.printNode()
		visited = append(visited, newn.id)
		newn = newn.successor()
	}
}

func (n *DHTNode) printNode() {
	fmt.Println("------------------------")
	fmt.Printf("Node:        %s\n", n.id)
	fmt.Printf("Predecessor: %s\n", n.predecessor.id)
	n.printFingers()
	//fmt.Println("------------------------\n")
}

func (n *DHTNode) printFingers() {
	fmt.Println("| startId  |   node.id |")
	for i := 1; i <= m; i++ {
		fmt.Printf("| %s       |        %s |\n", n.fingerTable[i].startId, n.fingerTable[i].node.id)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (nodeToAdd *DHTNode) join(n *DHTNode) {

	//	fmt.Printf("Adding node %s\n", nodeToAdd.id)
	// If nodeToAdd is the only node in the network
	if n == nil {
//		fmt.Printf("\nNode %s joins an empty ring\n", nodeToAdd.id)
		nodeToAdd.predecessor = nodeToAdd
		for i := 1; i <= m; i++ {
			nodeToAdd.fingerTable[i].startId, _ = calcFinger(hexStringToByteArr(nodeToAdd.id), i, m)
			nodeToAdd.fingerTable[i].node = nodeToAdd
		}
	} else {
		//fmt.Printf("\nNode %s joins, using node %s\n", nodeToAdd.id, n.id)
		nodeToAdd.initFingerTable(n)
		fmt.Printf("Node %s joined and initiated its finger now time for updating others\n", nodeToAdd.id)
		nodeToAdd.updateOthers()
		fmt.Println("")
	}
//	fmt.Printf("Ring structure after join, starting at %s: \n", nodeToAdd.id)
//	nodeToAdd.printRing()
//	fmt.Println("--- End ring\n")
}

// should be used in lookup and addToRing to find the right node / place in the ring
func (n *DHTNode) findSuccessor(id string) *DHTNode {
	return n.lookup(id)
	//predecessor := n.findPredecessor(id)
	//return predecessor.successor()
}

func (n *DHTNode) findPredecessor(id string) *DHTNode {
	n2 := n
	for !between([]byte(nextId(n2.id)), []byte(nextId(n2.successor().id)), []byte(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *DHTNode) closestPrecedingFinger(id string) *DHTNode {
	fmt.Printf("Finding closest preceding finger to to %s using %s", id, n.id)
	for i := m; i > 0; i-- {
		if between([]byte(nextId(n.id)), []byte(id), []byte(n.fingerTable[i].node.id)) {
			return n.fingerTable[i].node
		}
	}
	return n
}

func (nodeToUpdateTableOn *DHTNode) initFingerTable(n *DHTNode) {

	// Find successor node of nodeToUpdateTableOn using startId

	// Calculating first finger
	nodeToUpdateTableOn.fingerTable[1].startId, _ = calcFinger(hexStringToByteArr(nodeToUpdateTableOn.id), 1, m)
	// Successor to first finger

	nodeToUpdateTableOn.fingerTable[1].node = n.lookup(nodeToUpdateTableOn.fingerTable[1].startId)

	//	fmt.Printf("Due to initFingerTable: Finger 1 for node %s with startId = %s is set to %s\n", nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[1].startId, nodeToUpdateTableOn.fingerTable[1].node.id)

	//	fmt.Printf("Node %s first finger is %s \n", nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[1].node.id)

	// Set nodeToUpdateTableOns predecessor to the the node it's being inserted after
	nodeToUpdateTableOn.predecessor = nodeToUpdateTableOn.successor().predecessor

	// Set successor for node that´s before new node to new node
	//nodeToUpdateTableOn.predecessor.setSuccessor(nodeToUpdateTableOn)

	// Update the predecessor of the node that nodeToUpdateTableOn is inserted before
	nodeToUpdateTableOn.successor().predecessor = nodeToUpdateTableOn

	for i := 1; i <= (m - 1); i++ {
		// Calculating finger
		nodeToUpdateTableOn.fingerTable[i+1].startId, _ = calcFinger(hexStringToByteArr(nodeToUpdateTableOn.id), i+1, m)
		if between(
			hexStringToByteArr(nodeToUpdateTableOn.id),
			hexStringToByteArr(nodeToUpdateTableOn.fingerTable[i].node.id),
			hexStringToByteArr(nodeToUpdateTableOn.fingerTable[i+1].startId),
		) { // this happens when finger[k].interval does not contain any node
			// meaning [finger[k].startId, finger[k+1].startIf) does not contain any node! then finger[k+1].node = finger[k].node
			nodeToUpdateTableOn.fingerTable[i+1].node = nodeToUpdateTableOn.fingerTable[i].node

			//			fmt.Printf("First case: Due to initFingerTable: Finger %d for node %s with startId = %s is set to %s\n", i+1, nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[i+1].startId, nodeToUpdateTableOn.fingerTable[i+1].node.id)
		} else {
			nodeToUpdateTableOn.fingerTable[i+1].node = n.lookup(nodeToUpdateTableOn.fingerTable[i+1].startId)
			//			nodeToUpdateTableOn.fingerTable[i+1].node = nodeToUpdateTableOn.lookup(nodeToUpdateTableOn.fingerTable[i+1].startId).successor()
			//			fmt.Printf("Second case: Due to initFingerTable: Finger %d for node %s with startId = %s is set to %s\n", i+1, nodeToUpdateTableOn.id, nodeToUpdateTableOn.fingerTable[i+1].startId, nodeToUpdateTableOn.fingerTable[i+1].node.id)
		}
	}
}

// traverse the ring counter-clockwise to update all nodes whose finger table entries should refer to n
func (n *DHTNode) updateOthers() {
	for i := 1; i <= m; i++ {

		// find last node p whose i:th finger might be n
		nId := big.Int{}
		nId.SetString(n.id, 10)
		//fmt.Printf("nId som sträng %q\n", nId.String())
		var initialLength = len(n.id)

		y := big.Int{}
		two := big.NewInt(2)
		mbig := big.NewInt(m)

		y.Exp(two, big.NewInt(int64(i-1)), nil)
		y.Sub(&nId, &y)
		two.Exp(two, mbig, nil)
		y.Mod(&y, two)
		// y = nId - 2^(i-1)

		var returnString = y.String()

		for len(returnString) != initialLength {
			returnString = "0" + returnString
		}

		//fmt.Printf("Id to be or preceed %s\n", returnString)

		//fmt.Printf("in updateOthers: nId=%s, i=%d y=%s\n", nId.String(), i, y.String())
		//p := n.findPredecessor(y.String())
		p := n.lookup(returnString)
		if (p.id != returnString) {
			p = p.predecessor
		}

		//		fmt.Printf("p = %s\n", p.id)
		
		//fmt.Printf("Node %s inserted, update fingers of node %s since it is or preceeds %s\n", n.id, p.id, returnString)

		if (p.id != n.id) {
			p.updateFingerTable(n, i)
		}
	}
}

// if s should be the i:th finger of n -> update n's finger table entry i with n
func (n *DHTNode) updateFingerTable(s *DHTNode, i int) {

	//	if (s.id != n.id) {
	if between(
		hexStringToByteArr(n.id),
		hexStringToByteArr(n.fingerTable[i].node.id),
		hexStringToByteArr(s.id),
	) {
		fmt.Printf("Node %s between node %s and but not including %s\n", s.id, n.id, n.fingerTable[i].node.id)
		n.fingerTable[i].node = s
		fmt.Printf("Updated %s:s finger %d to %s\n", n.id, i, s.id)

		// get first node preceeding n
		p := n.predecessor
		if (p.id != s.id) {
			p.updateFingerTable(s, i)
		}
	}
	//	}
}

// returns a pointer to the node which is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
	//fmt.Printf("Looking up %s\n", hashKey)
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.predecessor.lookup(hashKey)
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
