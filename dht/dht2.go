package dht

import (
	"fmt"
	"encoding/hex"
	"math/big"
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
	hexbytes, _ = hex.DecodeString(hexId) //([]byte, error)
	return hexbytes
}

func byteArrTobigIntToString(nId []byte) {

	nBigInt := big.Int{}
	nBigInt.SetBytes(nId)
	fmt.Println(nBigInt.String())

	resultBytes := nBigInt.Bytes()
	fmt.Printf("%q\n", resultBytes)

	resultHex := fmt.Sprintf("%x", resultBytes)
	fmt.Printf("%s\n", resultHex)
	fmt.Println("")
}


	/*var hex string
	var bytes []byte
	
	fmt.Println(hexbytes)

	fmt.Println([]byte(node7.id))
	hex, bytes = calcFinger([]byte(node7.id), 1, 3)

	fmt.Println(hex)
	fmt.Println(bytes)*/













