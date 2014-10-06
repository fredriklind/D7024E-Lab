package dht

const m = 160

// This must be used when the type of node is not known.
// For example when a method can return either a local or remote node
type node interface {
	// Getters
	id() string
	predecessor() node
	// Methods
	lookup(id string) node
	updateSuccessor(node)
	updatePredecessor(node)
}

type localNode struct {
	_id         string
	pred        node
	fingerTable [m + 1]finger
	isListening chan bool
}

type remoteNode struct {
	_id, address, port string
}

type finger struct {
	startId string
	node    node
}
