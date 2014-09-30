package dht

import (
	"fmt"
)

func (n *DHTNode) printNode2() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node        %s\n", n.id)
	if n.predecessor != nil {
		fmt.Printf("Predecessor  %s\n", n.predecessor.id)
	}
	//	fmt.Println("")
}

func (n *DHTNode) printNodeWithFingers() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node %s\n", n.id)
	if n.predecessor != nil {
		fmt.Printf("Predecessor %s\n", n.predecessor.id)
	}
	for i := 1; i <= m; i++ {
		fmt.Printf("Finger %s -> Node %s\n", n.fingerTable[i].startId, n.fingerTable[i].node.id)
	}
	fmt.Println("")
}

func (n *DHTNode) printRing2() {
	//fmt.Printf(".           %s\n", n.id)
	//	fmt.Printf("%s\n", n.id)
	n.printNodeWithFingers()
	newn := n.successor()
	for newn.id != n.id {
		//fmt.Printf(".           %s\n", newn.id)
		//		fmt.Printf("%s\n", newn.id)
		newn.printNodeWithFingers()
		newn = newn.successor()
	}
	//	fmt.Println()
}

// Returns the node who is responsible for the data corresponding to hashKey, traversing the ring linearly
func (n *DHTNode) linearLookup(hashKey string) *DHTNode {
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.predecessor.linearLookup(hashKey)
	}
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

// could be used in lookup
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
