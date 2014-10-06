package dht

func (n *remoteNode) predecessor() *node {
	return transport.predecessorRequest(n)
}

func (n *remoteNode) lookup(id string) *node {
	return transport.SendLookupRequest(n.getAddress(), id)

}

func (n *remoteNode) updateSuccessor(*node) {

}

func (n *remoteNode) updatePredecessor(*node) {

}

func (n *remoteNode) getAddress() string {
	return n.address + ":" + n.port
}
