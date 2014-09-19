package dht

import (
	"fmt"
	"encoding/hex"
)

func (n *DHTNode) printNode2() {
	fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	if (n.predecessor != nil) {
		fmt.Printf("Predecessor %s\n", n.predecessor.id)
	}
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
	fmt.Println(n.id)
	newn := n.successor()
	for newn.id != n.id {
		fmt.Println(newn.id)
		newn = newn.successor()
	}
}

func hexStringToByteArr(hexId string) []byte {
	var hexbytes []byte
	hexbytes, _ = hex.DecodeString(hexId)
	return hexbytes
}

// Returns the node whose responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) linearLookup(hashKey string) *DHTNode {
//	fmt.Printf("Looking up %s\n", hashKey)
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.predecessor.lookup(hashKey)
	}
}

