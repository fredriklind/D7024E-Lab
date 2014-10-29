package node

import (
//log "github.com/cihub/seelog"
)

func newRemoteNode(id, ip, port, apiPort, dbPort string) *remoteNode {
	return &remoteNode{_id: id, _ip: ip, _port: port, _apiPort: apiPort, _dbPort: dbPort}
}

// ----------------------------------------------------------------------------------------
//										Getters + setters
// ----------------------------------------------------------------------------------------
func (n *remoteNode) id() string {
	return n._id
}

func (n *remoteNode) ip() string {
	return n._ip
}

func (n *remoteNode) port() string {
	return n._port
}

func (n *remoteNode) apiPort() string {
	return n._apiPort
}

func (n *remoteNode) dbPort() string {
	return n._dbPort
}

func (n *remoteNode) address() string {
	return n.ip() + ":" + n.port()
}

func (n *remoteNode) apiAddress() string {
	return n.ip() + ":" + n.apiPort()
}

func (n *remoteNode) dbAddress() string {
	return n.ip() + ":" + n.dbPort()
}

// TODO maybe return (node, error) here, to be able to handle errors better.
func (n *remoteNode) predecessor() node {
	pred, err := transport.sendPredecessorRequest(n.address())
	if err != nil {
		panic(err)
	}
	if pred["id"] == theLocalNode.id() {
		return theLocalNode
	} else {
		return newRemoteNode(pred["id"], pred["ip"], pred["port"], pred["apiPort"], pred["dbPort"])
	}
}

func (n *remoteNode) updatePredecessor(candidate node) {
	transport.sendUpdatePredecessorCall(n.address(), candidate)
}

func (n *remoteNode) updateSuccessor(candidate node) {
	transport.sendUpdateSuccessorCall(n.address(), candidate)
}

// ----------------------------------------------------------------------------------------
//										remoteNode Methods
// ----------------------------------------------------------------------------------------
func (n *remoteNode) lookup(key string) (node, error) {
	dict, err := transport.sendLookupRequest(n.address(), key)

	if err != nil {
		return nil, err
	}
	if dict["id"] == theLocalNode.id() {
		return theLocalNode, nil
	} else {
		return newRemoteNode(dict["id"], dict["ip"], dict["port"], dict["apiPort"], dict["dbPort"]), nil
	}
}

func (n *remoteNode) isAlive() bool {
	return transport.ping(n.address())
}
