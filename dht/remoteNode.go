package dht

func (n *remoteNode) predecessor() *node {
	return n.owner.transport.predecessorRequest(n)
}

func (n *remoteNode) lookup() *node {
}

func (n *remoteNode) updateSuccessor(*node) {

}

func (n *remoteNode) updatePredecessor(*node) {

}
