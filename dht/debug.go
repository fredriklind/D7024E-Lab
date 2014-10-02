package dht

import (
	"fmt"
	log "github.com/cihub/seelog"
	"testing"
)

func (n *localNode) printNode2() {
	//fmt.Printf("Node %s, address %s, port %s\n", n.id, n.adress, n.port)
	fmt.Printf("Node        %s\n", n.id)
	if n.predecessor != nil {
		fmt.Printf("Predecessor  %s\n", n.predecessor.id)
	}
	//	fmt.Println("")
}

func (n *localNode) printNodeWithFingers() {
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

func (n *localNode) printRing2() {
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
func (n *localNode) linearLookup(hashKey string) *localNode {
	if between(hexStringToByteArr(nextId(n.predecessor.id)), hexStringToByteArr(nextId(n.id)), hexStringToByteArr(hashKey)) {
		return n
	} else {
		return n.predecessor.linearLookup(hashKey)
	}
}

func (n *localNode) printRing() {
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

func (n *localNode) printNode() {
	fmt.Println("------------------------")
	fmt.Printf("Node:        %s\n", n.id)
	fmt.Printf("Predecessor: %s\n", n.predecessor.id)
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

// Sets up a custom log receiver that compares log messages
// to the supplied valid logs and fails the test if it finds a
// mismatch.

type TestContext struct {
	validLogs []string
	test      *testing.T
}

type CustomReceiver struct {
	context TestContext
}

func (cr *CustomReceiver) ReceiveMessage(message string, level log.LogLevel, context log.LogContextInterface) error {
	message = message[:len(message)-1]
	if message == cr.context.validLogs[0] {
		cr.context.validLogs = append(cr.context.validLogs[:0], cr.context.validLogs[0+1:]...)
	} else {
		cr.context.test.Fatalf("Expecting '%s', got '%s'\n", cr.context.validLogs[0], message)
	}
	return nil
}
func (cr *CustomReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}
func (cr *CustomReceiver) Flush() {

}
func (cr *CustomReceiver) Close() error {
	return nil
}

func setupTest(t *testing.T, valids []string) {
	c := TestContext{test: t, validLogs: valids}
	testConfig := `
<seelog>
    <outputs>
        <custom name="myreceiver" formatid="test"/>
    </outputs>
    <formats>
        <format id="test" format="%Msg%n"/>
    </formats>
</seelog>
`
	parserParams := &log.CfgParseParams{
		CustomReceiverProducers: map[string]log.CustomReceiverProducer{
			"myreceiver": func(log.CustomReceiverInitArgs) (log.CustomReceiver, error) {
				return &CustomReceiver{c}, nil
			},
		},
	}
	logger, err := log.LoggerFromParamConfigAsString(testConfig, parserParams)
	if err != nil {
		panic(err)
	}
	defer logger.Flush()
	err = log.ReplaceLogger(logger)
	if err != nil {
		panic(err)
	}
}

// could be used in lookup
func (n *localNode) findSuccessor(id string) *localNode {
	predecessor := n.findPredecessor(id)
	return predecessor.successor()
}

func (n *localNode) findPredecessor(id string) *localNode {
	n2 := n
	for !between(hexStringToByteArr(nextId(n2.id)), hexStringToByteArr(nextId(n2.successor().id)), hexStringToByteArr(id)) {
		n2 = n2.closestPrecedingFinger(id)
	}
	return n2
}

func (n *localNode) closestPrecedingFinger(id string) *localNode {
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
func (n *localNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id      string `json:"id"`
		Address string `json:"address"`
	}{
		Address: n.address + ":" + n.port,
		Id:      n.id,
	})
}
