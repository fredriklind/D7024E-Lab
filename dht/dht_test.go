package dht

import (
	"fmt"

	log "github.com/cihub/seelog"
	//"sort"
	"testing"
)

// Initializer for the dht package
func init() {
	testConfig := `
		<seelog type="sync">
			<outputs>
				<file formatid="onlytime" path="logfile.log"/>
			</outputs>
			<formats>
				<format id="default" format="%Date %Time [%LEVEL] %Msg%n"/>
				<format id="onlytime" format="%Time [%LEVEL] %Msg%n"/>
			</formats>
		</seelog>
	`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}

// test cases can be run by calling e.g. go test -test.run TestRingSetup
// go run test will run all tests

func TestFinger160bits2(t *testing.T) {

	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	/*	node4 := makeDHTNode(nil, "localhost", "1114")
		node5 := makeDHTNode(nil, "localhost", "1115")
		node6 := makeDHTNode(nil, "localhost", "1116")
		node7 := makeDHTNode(nil, "localhost", "1117")
		node8 := makeDHTNode(nil, "localhost", "1118")
		node9 := makeDHTNode(nil, "localhost", "1119")*/

	for i := 0; i < 10; i++ {
		fmt.Println()
	}
	fmt.Println("------------------------------------------------------------------------------------------------")

	fmt.Println("1 joins <-----------------------")
	fmt.Printf("Node id = %s\n\n", node1.id)
	node1.join(nil)
	node1.printRing2()

	fmt.Println("2 joins <-----------------------")
	fmt.Printf("Node id = %s\n\n", node2.id)
	node2.join(node1)
	node1.printRing2()

	fmt.Println("3 joins <-----------------------")
	fmt.Printf("Node id = %s\n\n", node3.id)
	node3.join(node1)
	node1.printRing2()

	/*
		ids := []string{node1.id, node2.id, node3.id, node4.id, node5.id, node6.id, node7.id, node8.id, node9.id}
		sort.Strings(ids)

		fmt.Println("Id:s sorted in increasing order:\n")
		for i := 0; i < len(ids); i++ {
			fmt.Println(ids[i])
		}

		fmt.Println()

		s := generateNodeId()
		next := nextId(s)

		fmt.Printf("nextId(%s) = %s\n\n", s, next)
		fmt.Printf("lookup2(%s)=%s\n", s, node1.lookup2(s).id)

		for i:=0; i<5; i++ {
			fmt.Println()
		}
	*/

	fmt.Println(node1.id)
	fmt.Println(node2.id)
	fmt.Println(node3.id)

}

func TestFinger16bits(t *testing.T) {

	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id8 := "08"
	id9 := "09"
	idA := "0a"
	idB := "0b"
	idC := "0c"
	idD := "0d"
	idE := "0e"
	idF := "0f"

	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(&id1, "localhost", "1101")
	node2 := makeDHTNode(&id2, "localhost", "1102")
	node3 := makeDHTNode(&id3, "localhost", "1103")
	node4 := makeDHTNode(&id4, "localhost", "1104")
	node5 := makeDHTNode(&id5, "localhost", "1105")
	node6 := makeDHTNode(&id6, "localhost", "1106")
	node7 := makeDHTNode(&id7, "localhost", "1107")
	node8 := makeDHTNode(&id8, "localhost", "1108")
	node9 := makeDHTNode(&id9, "localhost", "1109")
	nodeA := makeDHTNode(&idA, "localhost", "1110")
	nodeB := makeDHTNode(&idB, "localhost", "1111")
	nodeC := makeDHTNode(&idC, "localhost", "1112")
	nodeD := makeDHTNode(&idD, "localhost", "1113")
	nodeE := makeDHTNode(&idE, "localhost", "1114")
	nodeF := makeDHTNode(&idF, "localhost", "1115")

	fmt.Println("------------------------------------------------------------------------------------------------")

	fmt.Println("3 joins <-----------------------")
	node3.join(nil)
	node3.printRing2()

	fmt.Println("4 joins <-----------------------")
	node4.join(node3)
	node3.printRing2()

	fmt.Println("1 joins <-----------------------")
	node1.join(node3)
	node3.printRing2()

	fmt.Println("a joins <-----------------------")
	nodeA.join(node1)
	node1.printRing2()

	fmt.Println("5 joins <-----------------------")
	//	fmt.Printf("nextId(%s)=%s\n", node1.id, nextId(node1.id))
	//	fmt.Printf("node1.lookup2(%s)=%s\n", node5.id, node1.lookup2(node5.id).id)
	node5.join(node1)
	node1.printRing2()
	/*
		fmt.Println(hexStringToByteArr("00ff"))
		bytearr := hexStringToByteArr("00ff")
		yHex := fmt.Sprintf("%x", bytearr)
		fmt.Println(yHex) */

	fmt.Println("6 joins <-----------------------")
	node6.join(node1)
	node1.printRing2()

	fmt.Println("7 joins <-----------------------")
	node7.join(node1)
	node1.printRing2()

	fmt.Println("2 joins <-----------------------")
	node2.join(node1)
	node1.printRing2()

	fmt.Println("8 joins <-----------------------")
	node8.join(node1)

	fmt.Println("9 joins <-----------------------")
	node9.join(node1)

	fmt.Println("6 joins <-----------------------")
	nodeE.join(node1)

	fmt.Println("B joins <-----------------------")
	nodeB.join(node1)

	fmt.Println("C joins <-----------------------")
	nodeC.join(node1)

	fmt.Println("D joins <-----------------------")
	nodeD.join(node1)

	fmt.Println("F joins <-----------------------")
	nodeF.join(node1)
	node1.printRing2()

	/*	fmt.Println("------------------------------------------------------------------------------------------------")

		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println("RING STRUCTURE")
		fmt.Println("------------------------------------------------------------------------------------------------")
		node1.printRing2()
		fmt.Println("------------------------------------------------------------------------------------------------")*/

	/*node1.printNodeWithFingers()
	node2.printNodeWithFingers()
	node3.printNodeWithFingers()
	node4.printNodeWithFingers()
	node5.printNodeWithFingers()
	node6.printNodeWithFingers()
	node7.printNodeWithFingers()
	node8.printNodeWithFingers()
	node9.printNodeWithFingers()
	node10.printNodeWithFingers()
	node11.printNodeWithFingers()
	node12.printNodeWithFingers()
	node13.printNodeWithFingers()
	node14.printNodeWithFingers()
	node15.printNodeWithFingers()
	node16.printNodeWithFingers()*/

	/*	node1.printNode2()
		node2.printNode2()
		node3.printNode2()
		node4.printNode2()
		node5.printNode2()
		node6.printNode2()
		node7.printNode2()
		node8.printNode2()
		node9.printNode2()
		node10.printNode2()
		node11.printNode2()
		node12.printNode2()
		node13.printNode2()
		node14.printNode2()
		node15.printNode2()
		node16.printNode2()*/

}

func TestDebug2(t *testing.T) {

	//	id0 := "00"
	id1 := "01"
	id2 := "02"
	//	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	//	id7 := "07"

	//	node0 := makeDHTNode(&id0, "localhost", "1111")
	node1 := makeDHTNode(&id1, "localhost", "1112")
	node2 := makeDHTNode(&id2, "localhost", "1113")
	//	node3 := makeDHTNode(&id3, "localhost", "1114")
	node4 := makeDHTNode(&id4, "localhost", "1115")
	node5 := makeDHTNode(&id5, "localhost", "1116")
	node6 := makeDHTNode(&id6, "localhost", "1117")
	//	node7 := makeDHTNode(&id7, "localhost", "1118")
	/*
		node0.printNode2()
		node1.printNode2()
		node2.printNode2()
		node3.printNode2()
		node4.printNode2()
		node5.printNode2()
		node6.printNode2()
		node7.printNode2()*/

	//	node0.join(nil)
	node1.join(nil)
	node5.join(node1)
	//	node3.join(node0)
	node4.join(node1)
	node2.join(node1)
	node6.join(node1)
	//	node7.join(node0)

	//	node0.printNodeWithFingers()
	node1.printNodeWithFingers()
	node2.printNodeWithFingers()
	//	node3.printNodeWithFingers()
	node4.printNodeWithFingers()
	node5.printNodeWithFingers()
	node6.printNodeWithFingers()
	//	node7.printNodeWithFingers()

	/*	fmt.Printf("PrevId to %s is %s\n", "00", prevId("00"))	// prevId och nextId testade!
		fmt.Printf("PrevId to %s is %s\n", "01", prevId("01"))
		fmt.Printf("PrevId to %s is %s\n", "02", prevId("02"))
		fmt.Printf("PrevId to %s is %s\n", "03", prevId("03"))
		fmt.Printf("PrevId to %s is %s\n", "04", prevId("04"))
		fmt.Printf("PrevId to %s is %s\n", "05", prevId("05"))
		fmt.Printf("PrevId to %s is %s\n", "06", prevId("06"))
		fmt.Printf("PrevId to %s is %s\n", "07", prevId("07"))*/

	node4.printRing2()
	fmt.Println("")
	/*
		idn := [8]string{"00","01","02","03","04","05","06","07"}
		nodes := [8]*DHTNode{node1,node2,node4,node5,node6,node1,node2,node4}*/
	/*
		for i:=0; i<8; i++ {
			for j:=0; j<8; j++ {
				if !((nodes[i].lookup2(idn[j])).id == idn[j]) {
				fmt.Println("ERROR")
				}
			}
		}
	*/ /*
		for i:=0; i<8; i++ {
			for j:=0; j<8; j++ {
				respNode := nodes[i].lookup2(idn[j])
				fmt.Printf("%s.lookup2(%s)->%s\n", nodes[i].id, idn[j], respNode.id)
			}
			fmt.Println("")
		}
	*/
	/*	fmt.Println([]byte(node0.id))
		fmt.Println([]byte(node1.id))
		fmt.Println([]byte(node2.id))
		fmt.Println([]byte(node3.id))
		fmt.Println([]byte(node4.id))
		fmt.Println([]byte(node5.id))
		fmt.Println([]byte(node6.id))
		fmt.Println([]byte(node7.id))

		fmt.Println(hexStringToByteArr(node0.id))
		fmt.Println(hexStringToByteArr(node1.id))
		fmt.Println(hexStringToByteArr(node2.id))
		fmt.Println(hexStringToByteArr(node3.id))
		fmt.Println(hexStringToByteArr(node4.id))
		fmt.Println(hexStringToByteArr(node5.id))
		fmt.Println(hexStringToByteArr(node6.id))
		fmt.Println(hexStringToByteArr(node7.id))

		byteArrTobigIntToString(hexStringToByteArr(node0.id))
		byteArrTobigIntToString(hexStringToByteArr(node1.id))
		byteArrTobigIntToString(hexStringToByteArr(node2.id))
		byteArrTobigIntToString(hexStringToByteArr(node3.id))
		byteArrTobigIntToString(hexStringToByteArr(node4.id))
		byteArrTobigIntToString(hexStringToByteArr(node5.id))
		byteArrTobigIntToString(hexStringToByteArr(node6.id))
		byteArrTobigIntToString(hexStringToByteArr(node7.id))*/
}

/*
 * Example of expected output of calling printRing().
 *
 * f38f3b2dcc69a2093f258e31902e40ad33148385 1390478919082870357587897783216576852537917080453
 * 10dc86630d9277a20e5f6176ff0786f66e781d97 96261723029167816257529941937491552490862681495
 * 35f2749bbe6fd0221a97ecf0df648bc8355c7a0e 307983449213776748112297858267528664243962149390
 * 3cb3aaec484f62c04dbab1512409b51887b28272 346546169131330955640073427806530491225644106354
 * 624778a652b23ebeb2ce133277ee8812fff87992 561074958520938864836545731942448707916353010066
 * a5a5dcfbd8c15e495242c4d7fe680fe986562ce2 945682350545587431465494866472073397640858316002
 * b94a0c51288cdaaa00cd5609faa2189f56251984 1057814620711304956240501530938795222302424635780
 * d8b6ac320d92fe71551bed2f702ba6ef2907283e 1237215742469423719453176640534983456657032816702
 * ee33f5aaf7cf6a7168a0f3a4449c19c9b4d1e399 1359898542148650805696846077009990511357036979097
 */
/*func TestRingSetup(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.join(nil)
	node2.join(node1)
	node3.join(node1)
	node4.join(node1)
	node5.join(node2)
	node6.join(node2)
	node7.join(node2)
	node8.join(node3)
	node9.join(node3)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
}
*/

/*
 * Example of expected output.
 *
 * str=hello students!
 * hashKey=cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 */

/*func TestLookup(t *testing.T) {
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.join(nil)
	node2.join(node1)
	node3.join(node1)
	node4.join(node2)
	node5.join(node2)
	node6.join(node3)
	node7.join(node3)
	node8.join(node4)
	node9.join(node4)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	str := "hello students!"
	hashKey := sha1hash(str)
	fmt.Println("str=" + str)
	fmt.Println("hashKey=" + hashKey)

	fmt.Println("node 1: " + node1.lookup(hashKey).id + " is respoinsible for " + hashKey)
	fmt.Println("node 5: " + node5.lookup(hashKey).id + " is respoinsible for " + hashKey)

	fmt.Println("------------------------------------------------------------------------------------------------")

}*/

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            1
 * m            3
 * 2^(k-1)      1
 * (n+2^(k-1))  1
 * 2^m          8
 * result       1
 * result (hex) 01
 * successor    01
 * distance     1
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            2
 * m            3
 * 2^(k-1)      2
 * (n+2^(k-1))  2
 * 2^m          8
 * result       2
 * result (hex) 02
 * successor    02
 * distance     2
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            3
 * m            3
 * 2^(k-1)      4
 * (n+2^(k-1))  4
 * 2^m          8
 * result       4
 * result (hex) 04
 * successor    04
 * distance     4
 */

func TestFinger3bits(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0 := makeDHTNode(&id0, "localhost", "1111")
	node1 := makeDHTNode(&id1, "localhost", "1112")
	node2 := makeDHTNode(&id2, "localhost", "1113")
	node3 := makeDHTNode(&id3, "localhost", "1114")
	node4 := makeDHTNode(&id4, "localhost", "1115")
	node5 := makeDHTNode(&id5, "localhost", "1116")
	node6 := makeDHTNode(&id6, "localhost", "1117")
	node7 := makeDHTNode(&id7, "localhost", "1118")

	node2.join(nil)
	node2.printNodeWithFingers()

	node0.join(node2)
	node0.printNodeWithFingers()

	fmt.Println("2 joined, 0 joined now the rest")

	node1.join(node2)
	node3.join(node1)
	node4.join(node1)
	node5.join(node1)
	node6.join(node1)
	node7.join(node2)

	node2.printNodeWithFingers()
	node0.printNodeWithFingers()
	node1.printNodeWithFingers()

	node3.printNodeWithFingers()
	node4.printNodeWithFingers()
	node5.printNodeWithFingers()
	node6.printNodeWithFingers()
	node7.printNodeWithFingers()

	//	node0.printRing2()
}

func TestDebug(t *testing.T) {

	//id0 := "00"
	/*	id1 := "01"
			id2 := "02"
			id3 := "03"
			id4 := "04"
			id5 := "05"
			id6 := "06"
			id7 := "07"

			//node0 := makeDHTNode(&id0, "localhost", "1111")
			node1 := makeDHTNode(&id1, "localhost", "1112")
			node2 := makeDHTNode(&id2, "localhost", "1113")
			node3 := makeDHTNode(&id3, "localhost", "1114")
			node4 := makeDHTNode(&id4, "localhost", "1115")
			node5 := makeDHTNode(&id5, "localhost", "1116")
			node6 := makeDHTNode(&id6, "localhost", "1117")
			node7 := makeDHTNode(&id7, "localhost", "1118")

			//node0.setSuccessor(node1)
			//node0.predecessor = node7

		/*	node1.setSuccessor(node2)
			node1.predecessor = node7

			node2.setSuccessor(node3)
			node2.predecessor = node1

			node3.setSuccessor(node4)
			node3.predecessor = node2

			node4.setSuccessor(node5)
			node4.predecessor = node3

			node5.setSuccessor(node6)
			node5.predecessor = node4

			node6.setSuccessor(node7)
			node6.predecessor = node5

			node7.setSuccessor(node1)
			node7.predecessor = node6

			testLookup(node6, "00") */
}

func testLookup(n *DHTNode, id string) {
	var result = n.lookup2(id)
	fmt.Printf("%s.lookup(%s) returns %s\n", n.id, id, result.id)
}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            0
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0

 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            1
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            80
 * m            160
 * 2^(k-1)      604462909807314587353088
 * (n+2^(k-1))  682874255151879437996523461382311327142222978674
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996523461382311327142222978674
 * finger (hex) 779d240121ed6d5e8bd14b6529b08e5c617b5e72
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            120
 * m            90
 * 2^(k-1)      664613997892457936451903530140172288
 * (n+2^(k-1))  682874255152544051994415314855853423357775797874
 * 2^m          1237940039285380274899124224
 * finger       1180872106465109536036052594
 * finger (hex) 03d0cb6529b08e5c617b5e72
 * successor    f880fb198b7059ae92a69968727d84da9c94dd15
 * distance     877444087302148207702277795
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            160
 * m            160
 * 2^(k-1)      730750818665451459101842416358141509827966271488
 * (n+2^(k-1))  1413625073817330897098365273277543029655601897074
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       1413625073817330897098365273277543029655601897074
 * finger (hex) f79d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * successor    d0a43af3a433353909e09739b964e64c107e5e92
 * distance     508258282811496687056817668076520806659544776736
 */

func TestFinger160bits(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.join(nil)
	node2.join(node1)
	node3.join(node1)
	node4.join(node2)
	node5.join(node2)
	node6.join(node3)
	node7.join(node3)
	node8.join(node4)
	node9.join(node4)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	/*
		node3.testCalcFingers(0, 160)
		fmt.Println("")
		node3.testCalcFingers(1, 160)
		fmt.Println("")
		node3.testCalcFingers(80, 160)
		fmt.Println("")
		node3.testCalcFingers(120, 90)
		fmt.Println("")
		node3.testCalcFingers(160, 160)
		fmt.Println("")
	*/
}

func TestHELLO(t *testing.T) {
	// Define a sequence of requests that are expected
	/*setupTest(t, []string{
		"Node 01 sent HELLO Request",
		"Node 02 got HELLO Request",
		"Node 02 sent ACK Response",
		"Node 01 got ACK Response",
	},
	)*/

	id1 := "01"
	id2 := "02"
	node1 := makeDHTNode(&id1, "127.0.0.1", "2000")
	node2 := makeDHTNode(&id2, "127.0.0.1", "3000")

	node1.sendRequest(Msg{
		Method: "HELLO",
		Dst:    node2.getAddress(),
	})
}

func Test3NodeForwarding(t *testing.T) {
	block := make(chan bool)

	id1 := "01"
	id2 := "02"
	id3 := "03"

	node1 := makeDHTNode(&id1, "127.0.0.1", "2000")
	node2 := makeDHTNode(&id2, "127.0.0.1", "3000")
	node3 := makeDHTNode(&id3, "127.0.0.1", "4000")

	node1.sendRequest(
		Msg{
			Method: "FORWARD",
			Values: map[string]string{
				"Method":             "HELLO",
				"FinalDestinationId": "03",
				"Sender":             node1.getAddress(),
			},
			Dst: node2.getAddress()},
	)

	// To prevent stupid warnings
	_ = node3
	<-block
}
