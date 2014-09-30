package dht

import (
	"encoding/json"
	"fmt"
	"math/big"
)

const m = 160
const base = 16

// Turn the node into a JSON string containing id and address
func (n *DHTNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id      string `json:"id"`
		Address string `json:"address"`
	}{
		Address: n.address + ":" + n.port,
		Id:      n.id,
	})
}

func (n *DHTNode) getAddress() string {
	return n.address + ":" + n.port
}

// Returns the node who is responsible for the data corresponding to id, traversing the ring using finger tables
func (n *DHTNode) lookup2(id string) *DHTNode {
	//	fmt.Printf("Performing lookup from node %s\n", n.id)
	// n responsible for id
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(id)) {
		//		fmt.Printf("%s E (%s, %s], eg. [%s, %s) \n", id, n.predecessor.id, n.id, nextId(n.predecessor.id), nextId(n.id))
		return n
		// otherwise use fingers of n, starting with the one that is furthest away, to find responsible node
	} else {
		for i := m; i >= 1; i-- {

			//			fmt.Printf("i=%d\n", i)
			// special case - when n´s finger points to itself
			if n.fingerTable[i].node.id == n.id {

				// what to do? go to next finger...
				//				fmt.Println("Finger points to node itself<-------------------")
				// id between finger and node - got to that finger
			} else if between(hexStringToByteArr(n.fingerTable[i].node.id), hexStringToByteArr(n.id), hexStringToByteArr(id)) {
				//				fmt.Printf("%s E [%s,%s)\n", id, n.fingerTable[i].node.id, n.id)
				//				fmt.Printf("Go to node %s and perform lookup on %s\n", n.fingerTable[i].node.id, id)
				return ((n.fingerTable[i].node).lookup2(id))
			}
		}
		// if id is not between any finger and n - then id must be between n and its successor
		return n.successor()
	}
}

// lookup of finger.node for the case when a second node is added to a ring with only one node
func (n *DHTNode) specLookup(newNode *DHTNode, startId string) *DHTNode {
	if between(hexStringToByteArr(nextId(newNode.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(startId)) {
		return n
	}
	return newNode
}

func (newNode *DHTNode) initFingerTable(n *DHTNode) {

	oneNodeRing := false

	// Calculating first finger
	newNode.fingerTable[1].startId, _ = calcFinger(hexStringToByteArr(newNode.id), 1, m)

	// Successor to newNode
	newNode.fingerTable[1].node = n.lookup2(newNode.fingerTable[1].startId)

	// Set newNodes predecessor to the the node it is being inserted after
	newNode.predecessor = newNode.successor().predecessor

	// Set successor of newNode´s predecessor to newNode
	newNode.predecessor.setSuccessor(newNode)

	//	fmt.Printf("INIT: %s successor set to %s\n", newNode.id, newNode.fingerTable[1].node.id)
	//	fmt.Printf("INIT: %s predecessor set to %s\n", newNode.id, newNode.successor().predecessor.id)

	if n.predecessor.id == n.id {
		oneNodeRing = true
	} else {
		oneNodeRing = false
	}

	// Update the predecessor of the node that newNode is inserted before
	newNode.successor().predecessor = newNode

	//	fmt.Printf("INIT: %s, successor of %s set to %s\n", newNode.id, newNode.predecessor.id, newNode.id)
	//	fmt.Printf("INIT: %s, predecessor of %s set to %s\n", newNode.id, newNode.successor().id, newNode.id)
	//	fmt.Println("Hejpa")
	for i := 1; i <= (m - 1); i++ {
		// Calculating finger
		newNode.fingerTable[i+1].startId, _ = calcFinger(hexStringToByteArr(newNode.id), i+1, m)
		if between(
			hexStringToByteArr(newNode.id),
			hexStringToByteArr(nextId(newNode.fingerTable[i].node.id)),
			hexStringToByteArr(newNode.fingerTable[i+1].startId),
		) {

			//			fmt.Printf("%s between %s and %s\n", newNode.fingerTable[i+1].startId, newNode.id, nextId(newNode.fingerTable[i].node.id))
			// startId between node and previous finger.node
			newNode.fingerTable[i+1].node = newNode.fingerTable[i].node
		} else {
			//			fmt.Printf("Lookup instead, finger %d lookup2(%s)\n", i+1, newNode.fingerTable[i+1].startId)
			if oneNodeRing {
				newNode.fingerTable[i+1].node = n.specLookup(newNode, newNode.fingerTable[i+1].startId)
			} else {
				newNode.fingerTable[i+1].node = n.lookup2(newNode.fingerTable[i+1].startId)
				//			fmt.Printf("%s.lookup2(%s) = %s \n", n.id, newNode.fingerTable[i+1].startId, (n.lookup2(newNode.fingerTable[i+1].startId)).id)
				//			fmt.Printf("FingerNode %s \n", newNode.fingerTable[i+1].node.id)
			}
		}
		//		fmt.Printf("Node %s finger nr %d startId %s Node %s\n", newNode.id, i+1, newNode.fingerTable[i+1].startId, newNode.fingerTable[i+1].node.id)

	}
}

// Traverse the ring counter-clockwise to update all nodes whose finger table entries could/should refer to n
func (n *DHTNode) updateOthers() {
	//	fmt.Printf("%s.updateOthers()\n", n.id)
	for i := 1; i <= m; i++ {
		//		fmt.Printf("Loop %d of updateOthers()\n", i)

		// Find last preceeding node p whose i:th finger might be n
		nId := big.Int{}
		nId.SetString(n.id, base)

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

		//		fmt.Printf("in updateOthers: nId=%s, i=%d y=%s\n", nId.String(), i, y.String())

		p := n.lookup2(yHex)

		if p.id != yHex {
			p = p.predecessor
		}
		//		fmt.Printf("p = %s\n", p.id)
		if p.id != n.id {
			//			fmt.Printf("Going in to updateFingerTable, p is: %s\n", p.id)
			p.updateFingerTable(n, i)
		}
		//		fmt.Printf("%s.uptadeFingertable(%s,%d)\n", p.id, n.id, i)
	}
}

// If s should be the i:th finger of n -> update n's finger table entry nr i with s
func (n *DHTNode) updateFingerTable(s *DHTNode, i int) {
	// n´s finger.node points to n itself
	if n.id == n.fingerTable[i].node.id {
		if between(
			hexStringToByteArr(n.fingerTable[i].startId),
			hexStringToByteArr(n.fingerTable[i].node.id),
			hexStringToByteArr(s.id),
		) {
			n.fingerTable[i].node = s
			//			fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id, i, n.id, n.id, i, s.id)
		}
	} else if between(
		hexStringToByteArr(n.fingerTable[i].startId),
		hexStringToByteArr(n.fingerTable[i].node.id),
		hexStringToByteArr(s.id),
	) {
		if !(n.fingerTable[i].startId == n.fingerTable[i].node.id) {
			n.fingerTable[i].node = s
			//				fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id, i, n.id, n.id, i, s.id)
		}
	}
	// Get last node preceeding n and mayby update its finger i as well, check that it hasn´t come round to the node just added (s)
	p := n.predecessor
	//		fmt.Printf("p is: %s\n", p.id)
	if p.id != s.id {
		p.updateFingerTable(s, i)
	} else {
		//			fmt.Println("Not going to that node")
	}
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
		fmt.Println("After initFingerTable:\n")
		//		nodeToAdd.printNode2()
		//		nodeToAdd.printRing2()
		nodeToAdd.printNodeWithFingers()
		//		fmt.Printf("Node %s joined and initiated its finger now time for updating others\n", nodeToAdd.id)
		nodeToAdd.updateOthers()
		fmt.Println("After updateOthers:")
	}
	//	fmt.Printf("Ring structure after join, starting at %s: \n", nodeToAdd.id)
	//	nodeToAdd.printRing()
	//	fmt.Println("--- End ring\n")
}

type DHTNode struct {
	id, address, port string
	predecessor       *DHTNode
	fingerTable       [m + 1]Finger
	Requests          map[string]chan Msg
	isListening       chan bool
}

type Finger struct {
	startId string
	node    *DHTNode
}

func makeDHTNode(idPointer *string, address string, port string) *DHTNode {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}
	node := DHTNode{id: id, address: address, port: port}
	node.Requests = make(map[string]chan Msg)
	go node.listen()
	return &node
}

func (n *DHTNode) successor() *DHTNode {
	return n.fingerTable[1].node
}

func (n *DHTNode) setSuccessor(successor *DHTNode) {
	n.fingerTable[1].node = successor
}
