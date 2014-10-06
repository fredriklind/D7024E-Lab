package dht

import (
	"fmt"
	"transport"
)

// ----------------------------------------------------------------------------------------
//										Initializer
// ----------------------------------------------------------------------------------------
func makeLocalNode(idPointer *string, address string, port string) *localNode {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}
	node := localNode{_id: id}
	transport.NewTransporter(address, port)
	return &node
}

// ----------------------------------------------------------------------------------------
//										Getters + setters
// ----------------------------------------------------------------------------------------
func (n *localNode) id() string {
	return n._id
}

func (n *localNode) predecessor() node {
	return n.pred
}

func (n *localNode) successor() node {
	return n.fingerTable[1].node
}

func (n *localNode) updatePredecessor(n2 node) {
	n.pred = n2
}

func (n *localNode) updateSuccessor(n2 node) {
	n.fingerTable[1].node = n2
}

// ----------------------------------------------------------------------------------------
//										localNode methods
// ----------------------------------------------------------------------------------------

// Returns the node who is responsible for the data corresponding to id, traversing the ring using finger tables
func (n *localNode) lookup(id string) node {
	// n responsible for id
	if between(
		hexStringToByteArr(nextId(n.predecessor().id())),
		hexStringToByteArr(nextId(n.id())),
		hexStringToByteArr(id),
	) {
		return n
		// otherwise use fingers of n, starting with the one that is furthest away, to find responsible node
	} else {
		for i := m; i >= 1; i-- {
			// special case - when n´s finger points to itself
			if n.fingerTable[i].node.id() == n.id() {

				// what to do? go to next finger...
				// id between finger and node - got to that finger
			} else if between(
				hexStringToByteArr(n.fingerTable[i].node.id()),
				hexStringToByteArr(n.id()),
				hexStringToByteArr(id),
			) {
				return ((n.fingerTable[i].node).lookup(id))
			}
		}
		// if id is not between any finger and n - then id must be between n and its successor
		return n.successor()
	}
}

// lookup of finger.node for the case when a second node is added to a ring with only one node
func (n *localNode) specLookup(newNode *localNode, startId string) *localNode {
	if between(
		hexStringToByteArr(nextId(newNode.id())),
		hexStringToByteArr(nextId(n.id())),
		hexStringToByteArr(startId),
	) {
		return n
	}
	return newNode
}

func (newNode *localNode) initFingerTable(n *localNode) {
	oneNodeRing := false

	// Calculating first finger
	newNode.fingerTable[1].startId, _ = calcFinger(hexStringToByteArr(newNode.id()), 1, m)

	// Successor to newNode
	newNode.fingerTable[1].node = n.lookup(newNode.fingerTable[1].startId)

	// Set newNodes predecessor to the the node it is being inserted after
	newNode.pred = newNode.successor().predecessor()

	// Set successor of newNode´s predecessor to newNode
	newNode.predecessor().updateSuccessor(newNode)

	if n.predecessor().id() == n.id() {
		oneNodeRing = true
	} else {
		oneNodeRing = false
	}

	// Update the predecessor of the node that newNode is inserted before
	newNode.successor().updatePredecessor(newNode)

	for i := 1; i <= (m - 1); i++ {
		// Calculating finger
		newNode.fingerTable[i+1].startId, _ = calcFinger(hexStringToByteArr(newNode.id()), i+1, m)
		if between(
			hexStringToByteArr(newNode.id()),
			hexStringToByteArr(nextId(newNode.fingerTable[i].node.id())),
			hexStringToByteArr(newNode.fingerTable[i+1].startId),
		) {
			// startId between node and previous finger.node
			newNode.fingerTable[i+1].node = newNode.fingerTable[i].node
		} else {
			if oneNodeRing {
				newNode.fingerTable[i+1].node = n.specLookup(newNode, newNode.fingerTable[i+1].startId)
			} else {
				newNode.fingerTable[i+1].node = n.lookup(newNode.fingerTable[i+1].startId)
			}
		}
	}
}

// Traverse the ring counter-clockwise to update all nodes whose finger table entries could/should refer to n
/*func (n *localNode) updateOthers() {
	for i := 1; i <= m; i++ {
		// Find last preceeding node p whose i:th finger might be n
		nId := big.Int{}
		nId.SetString(n.id(), base)

		y := big.Int{}
		two := big.NewInt(2)
		mbig := big.NewInt(m)

		y.Exp(two, big.NewInt(int64(i-1)), nil)
		y.Sub(&nId, &y)

		two.Exp(two, mbig, nil)
		y.Mod(&y, two)
		// y = nId - 2^(i-1)

		yBytes := y.Bytes()
		yHex := fmt.Sprintf("%x", yBytes)

		p := n.lookup(yHex)

		if p.id() != yHex {
			p = p.predecessor()
		}
		if p.id() != n.id() {
			p.updateFingerTable(n, i)
		}
	}
}

// If s should be the i:th finger of n -> update n's finger table entry nr i with s
func (n *localNode) updateFingerTable(s *localNode, i int) {
	// n´s finger.node points to n itself
	if n.id() == n.fingerTable[i].node.id() {
		if between(
			hexStringToByteArr(n.fingerTable[i].startId),
			hexStringToByteArr(n.fingerTable[i].node.id()),
			hexStringToByteArr(s.id()),
		) {
			n.fingerTable[i].node = s
			//			fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id(), i, n.id(), n.id(), i, s.id())
		}
	} else if between(
		hexStringToByteArr(n.fingerTable[i].startId),
		hexStringToByteArr(n.fingerTable[i].node.id()),
		hexStringToByteArr(s.id()),
	) {
		if !(n.fingerTable[i].startId == n.fingerTable[i].node.id()) {
			n.fingerTable[i].node = s
			//				fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id(), i, n.id(), n.id(), i, s.id())
		}
	}
	// Get last node preceeding n and mayby update its finger i as well, check that it hasn´t come round to the node just added (s)
	p := n.predecessor()
	//		fmt.Printf("p is: %s\n", p.id())
	if p.id() != s.id() {
		p.updateFingerTable(s, i)
	} else {
		//			fmt.Println("Not going to that node")
	}
}*/

func (nodeToAdd *localNode) join(n *localNode) {

	//	fmt.Printf("Adding node %s\n", nodeToAdd.id())

	// If nodeToAdd is the only node in the network
	if n == nil {
		//		fmt.Printf("\nNode %s joins an empty ring\n", nodeToAdd.id())
		nodeToAdd.updatePredecessor(nodeToAdd)
		for i := 1; i <= m; i++ {
			nodeToAdd.fingerTable[i].startId, _ = calcFinger(hexStringToByteArr(nodeToAdd.id()), i, m)
			nodeToAdd.fingerTable[i].node = nodeToAdd
		}
	} else {
		//		fmt.Printf("\nNode %s joins, using node %s\n", nodeToAdd.id(), n.id())
		nodeToAdd.initFingerTable(n)
		fmt.Println("After initFingerTable:\n")
		//		nodeToAdd.printNode2()
		//		nodeToAdd.printRing2()
		nodeToAdd.printNodeWithFingers()
		//		fmt.Printf("Node %s joined and initiated its finger now time for updating others\n", nodeToAdd.id())
		//nodeToAdd.updateOthers()
		fmt.Println("After updateOthers:")
	}
	//	fmt.Printf("Ring structure after join, starting at %s: \n", nodeToAdd.id())
	//	nodeToAdd.printRing()
	//	fmt.Println("--- End ring\n")
}

// Caleld periodically to update fingers
func (n *localNode) fixFingers() {

}
