package dht

const m = 3

var transport *transporter

// This must be used when the type of node is not known.
// For example when a method can return either a local or remote node.
type node interface {
	// Getters
	id() string
	predecessor() node
	// Address for the transport layer
	address() string
	// Address for the REST API
	apiAddress() string
	// Address for accessing DB
	dbAddress() string

	ip() string
	port() string
	apiPort() string
	dbPort() string

	// Methods
	lookup(id string) (node, error)
	updateSuccessor(node)
	updatePredecessor(node)
}

type localNode struct {
	_id            string
	pred           node
	fingerTable    [m + 1]finger
	fixFingersChan chan bool
}

type remoteNode struct {
	_id, _ip, _port, _apiPort, _dbPort string
}

type finger struct {
	startId string
	node    node
}

func main() {
	go startWebServer()
	go startAPI()

	/*for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		cmd, _ := reader.ReadString('\n')
		if cmd == "e\n" {
			db.Close()
			break
		}
	}
	*/
}
