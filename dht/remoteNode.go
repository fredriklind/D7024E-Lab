package dht

// ----------------------------------------------------------------------------------------
//										Getters + setters
// ----------------------------------------------------------------------------------------
func (n *remoteNode) id() string {
	return n._id
}

func (n *remoteNode) address() string {
	return n._address
}

// TODO maybe return (node, error) here, to be able to handle errors better.
func (n *remoteNode) predecessor() node {
	// TODO add conversion from what transport returns and what
	// this method should return
	dict, err := transport.sendPredecessorRequest(n.address())
	if err != nil {
		panic(err)
	}
	if dict["id"] == theLocalNode.id() {
		return theLocalNode
	} else {
		return &remoteNode{_id: dict["id"], _address: dict["address"]}
	}
}

func (n *remoteNode) updatePredecessor(candidate node) {
	transport.sendUpdatePredecessorCall(n.address(), candidate.id(), candidate.address())
}

func (n *remoteNode) updateSuccessor(candidate node) {
	transport.sendUpdateSuccessorCall(n.address(), candidate.id(), candidate.address())
}

// ----------------------------------------------------------------------------------------
//										remoteNode Methods
// ----------------------------------------------------------------------------------------
func (n *remoteNode) lookup(key string) (node, error) {
	// TODO add conversion from what transport returns and what
	// this method should return
	dict, err := transport.sendLookupRequest(n.address(), key)

	if err != nil {
		return nil, err
	}
	if dict["id"] == theLocalNode.id() {
		return theLocalNode, nil
	} else {
		return &remoteNode{_id: dict["id"], _address: dict["address"]}, nil
	}
}

func (n *remoteNode) isAlive() bool {
	return transport.ping(n.address())
}
