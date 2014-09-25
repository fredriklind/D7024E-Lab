package dht

import (
	"fmt"
	"encoding/hex"
)

func (n *DHTNode) printNode2() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node        %s\n", n.id)
	if (n.predecessor != nil) {
		fmt.Printf("Predecessor  %s\n", n.predecessor.id)
	}
	fmt.Println("")
}

func (n *DHTNode) printNodeWithFingers() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node %s\n", n.id)
	if (n.predecessor != nil) {
		fmt.Printf("Predecessor %s\n", n.predecessor.id)
	}
	for i:=1; i<=m; i++ {
		fmt.Printf("Finger %s -> Node %s\n", n.fingerTable[i].startId, n.fingerTable[i].node.id)
	}
	fmt.Println("")
}

func (n *DHTNode) printRing2() {
	fmt.Printf(".           %s\n", n.id)
	newn := n.successor()
	for newn.id != n.id {
		fmt.Printf(".           %s\n", newn.id)
		newn = newn.successor()
	}
	fmt.Println()
}

func hexStringToByteArr(hexId string) []byte {
	var hexbytes []byte
	hexbytes, _ = hex.DecodeString(hexId)
	return hexbytes
}

// Returns the node who is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) linearLookup(hashKey string) *DHTNode {
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.predecessor.lookup(hashKey)
	}
}

// Returns the node who is responsible for the data corresponding to id, traversing the ring using finger tables
func (n *DHTNode) lookup2(id string) *DHTNode {

	// n responsible for id
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(id)) {
		return n
	// otherwise use fingers of n, starting with the one that is furthest away, to find responsible node
	} else {
		for i:=m; i>=1; i-- {


			// special case - when n´s finger points to itself
			if (n.fingerTable[i].node.id == n.id) {

				// what to do?
				// go to next finger...
//				fmt.Println("Finger points to node itself<-------------------")


			} else if between(hexStringToByteArr(n.fingerTable[i].node.id), hexStringToByteArr(n.id), hexStringToByteArr(id)) {
				return n.fingerTable[i].node.lookup2(id)
			}
		}
		/*// if id is not between any finger and n - then id must be between n and its successor
		if !(n.fingerTable[1].node.id == n.id) {
			return n.fingerTable[1].node
		}*/
		return n.fingerTable[1].node
	}
}	

// lookup of finger.node for the case when a second node is added to a ring with only one node 
func (n *DHTNode) specLookup(newNode *DHTNode, startId string) * DHTNode {
	if between(hexStringToByteArr(nextId(newNode.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(startId)) {
		return n
	}
	return newNode
}









