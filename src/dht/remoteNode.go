package dht

import (
//	"transport"
)

// ----------------------------------------------------------------------------------------
//										Getters + setters
// ----------------------------------------------------------------------------------------
func (n *remoteNode) id() string {
	return n._id
}

// TODO maybe return (node, error) here, to be able to handle errors better.
func (n *remoteNode) predecessor() node {
	// TODO add conversion from what transport returns and what
	// this method should return
	dict, err := transport.PredecessorRequest(n.getAddress())
	if err != nil {
		panic(err)
	}
	if dict["id"] == theLocalNode.id() {
		return theLocalNode
	} else {
		return &remoteNode{_id: dict["id"], address: dict["address"], port: dict["port"]}
	}
}

func (n *remoteNode) updateSuccessor(node) {

}

func (n *remoteNode) updatePredecessor(node) {

}

// ----------------------------------------------------------------------------------------
//										remoteNode Methods
// ----------------------------------------------------------------------------------------
func (n *remoteNode) lookup(id string) node {
	// TODO add conversion from what transport returns and what
	// this method should return
	//transport.SendLookupRequest(n.getAddress(), id)
	return n
}

func (n *remoteNode) getAddress() string {
	return n.address + ":" + n.port
}
