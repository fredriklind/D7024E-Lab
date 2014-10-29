package node

import (
	"fmt"
)

func (n *localNode) printNode2() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node        %s\n", n.id)
	if n.predecessor() != nil {
		fmt.Printf("Predecessor  %s\n", n.predecessor().id)
	}
	//	fmt.Println("")
}

func (n *remoteNode) printNode2() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node        %s\n", n.id)
	if n.predecessor() != nil {
		fmt.Printf("Predecessor  %s\n", n.predecessor().id())
	}
	//	fmt.Println("")
}

func (n *localNode) printNodeWithFingers() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node %s\n", n.id())
	if n.predecessor() != nil {
		fmt.Printf("Predecessor %s\n", n.predecessor().id())
	}
	for i := 1; i <= m; i++ {
		fmt.Printf("Finger %s -> Node %s\n", n.fingerTable[i].startId, n.fingerTable[i].node.id())
	}
	fmt.Println("")
}

// Returns the node who is responsible for the data corresponding to hashKey, traversing the ring linearly
/*func (n *localNode) linearLookup(hashKey string) *node {
	if between(hexStringToByteArr(nextId(n.pred.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.pred.linearLookup(hashKey)
	}
}*/

func (n *localNode) printNode() {
	fmt.Println("------------------------")
	fmt.Printf("Node:        %s\n", n.id)
	fmt.Printf("Predecessor: %s\n", n.pred.id)
	n.printFingers()
	//fmt.Println("------------------------\n")
}

func (n *localNode) printFingers() {
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

/*
// could be used in lookup
func (n *node) findSuccessor(id string) *node {
	predecessor := n.findPredecessor(id)
	return predecessor.successor()
}

func (n *node) findPredecessor(id string) *node {
	n2 := n
	for !between(hexStringToByteArr(nextId(n2.id)), hexStringToByteArr(nextId(n2.successor().id)), hexStringToByteArr(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *node) closestPrecedingFinger(id string) *node {
	for i := m; i > 0; i-- {
		if between(hexStringToByteArr(nextId(n.id)), hexStringToByteArr(id), hexStringToByteArr(n.fingerTable[i].node.id)) {
			//			fmt.Printf(" %s\n", n.fingerTable[i].node.id)
			return n.fingerTable[i].node
		}
	}
	//	fmt.Printf(" %s\n", n.id)
	return n
}

// Turn the node into a JSON string containing id and address
func (n *node) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id      string `json:"id"`
		Address string `json:"address"`
	}{
		Address: n.address + ":" + n.port,
		Id:      n.id,
	})
}
*/
