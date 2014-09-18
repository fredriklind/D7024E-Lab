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
//		fmt.Printf("\nNode %s joins, using node %s\n", nodeToAdd.id, n.id)
		nodeToAdd.initFingerTable(n)
//		fmt.Printf("Node %s joined and initiated its finger now time for updating others\n", nodeToAdd.id)
		nodeToAdd.updateOthers()
//		fmt.Println("")
	}
//	fmt.Printf("Ring structure after join, starting at %s: \n", nodeToAdd.id)
//	nodeToAdd.printRing()
//	fmt.Println("--- End ring\n")
}

	// should be used in lookup
func (n *DHTNode) findSuccessor(id string) *DHTNode {
	predecessor := n.findPredecessor(id)
	return predecessor.successor()
}

func (n *DHTNode) findPredecessor(id string) *DHTNode {
	n2 := n
	for !between(hexStringToByteArr(nextId(n2.id)), hexStringToByteArr(nextId(n2.successor().id)), hexStringToByteArr(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *DHTNode) closestPrecedingFinger(id string) *DHTNode {
	for i := m; i > 0; i-- {
		if between(hexStringToByteArr(nextId(n.id)), hexStringToByteArr(id), hexStringToByteArr(n.fingerTable[i].node.id)) {
//			fmt.Printf(" %s\n", n.fingerTable[i].node.id)
			return n.fingerTable[i].node
		}
	}
//	fmt.Printf(" %s\n", n.id)
	return n
}

func (newNode *DHTNode) initFingerTable(n *DHTNode) {

	// Calculating first finger
	newNode.fingerTable[1].startId, _ = calcFinger(hexStringToByteArr(newNode.id), 1, m)
	
	// Successor to first finger
	newNode.fingerTable[1].node = n.lookup(newNode.fingerTable[1].startId)

	// Set newNodes predecessor to the the node it is being inserted after
	newNode.predecessor = newNode.successor().predecessor

//	fmt.Printf("INIT: %s successor set to %s\n", newNode.id, newNode.fingerTable[1].node.id)
//	fmt.Printf("INIT: %s predecessor set to %s\n", newNode.id, newNode.successor().predecessor.id)

	// Update the predecessor of the node that newNode is inserted before
	newNode.successor().predecessor = newNode

//	fmt.Printf("INIT: %s, successor of %s set to %s\n", newNode.id, newNode.predecessor.id, newNode.id)
//	fmt.Printf("INIT: %s, predecessor of %s set to %s\n", newNode.id, newNode.successor().id, newNode.id)

	for i:=1; i<=(m-1); i++ {
		// Calculating finger
		newNode.fingerTable[i+1].startId, _ = calcFinger(hexStringToByteArr(newNode.id), i+1, m)
		if between(
			hexStringToByteArr(newNode.id),
			hexStringToByteArr(newNode.fingerTable[i].node.id),
			hexStringToByteArr(newNode.fingerTable[i+1].startId),
		) { 
			newNode.fingerTable[i+1].node = newNode.fingerTable[i].node
		} else {
			newNode.fingerTable[i+1].node = n.lookup(newNode.fingerTable[i+1].startId)
		}
	}
}

// Traverse the ring counter-clockwise to update all nodes whose finger table entries could/should refer to n
func (n *DHTNode) updateOthers() {
//	fmt.Printf("%s.updateOthers()\n", n.id)
	for i := 1; i <= m; i++ {
//		fmt.Printf("Loop %d of updateOthers()\n", i)

		// Find last preceeding node p whose i:th finger might be n
		nId := big.Int{}
		nId.SetString(n.id, 10)

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

//		fmt.Printf("in updateOthers: nId=%s, i=%d y=%s\n", nId.String(), i, y.String())

		p := n.lookup(returnString)
		if (p.id != returnString) {
			p = p.predecessor
		}

//		fmt.Printf("p = %s\n", p.id)

		if (p.id != n.id) {
			p.updateFingerTable(n, i)
		}

//		fmt.Printf("%s.uptadeFingertable(%s,%d)\n", p.id, n.id, i)
	}
}

// If s should be the i:th finger of n -> update n's finger table entry i with n
func (n *DHTNode) updateFingerTable(s *DHTNode, i int) {

	if between(
		hexStringToByteArr(n.id),
		hexStringToByteArr(n.fingerTable[i].node.id),
		hexStringToByteArr(s.id),
	) {
		n.fingerTable[i].node = s

		// Get last node preceeding n, check that it hasn´t come round to the node just added (s)
		p := n.predecessor
		if (p.id != s.id) {
			p.updateFingerTable(s, i)
		}
	}
}

// Returns the node whose responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) lookup(hashKey string) *DHTNode {
//	fmt.Printf("Looking up %s\n", hashKey)
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
