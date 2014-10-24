package dht

import (
	//	"fmt"
	"time"

	log "github.com/cihub/seelog"
)

// ----------------------------------------------------------------------------------------
//										Initializer
// ----------------------------------------------------------------------------------------
func newLocalNode(idPointer *string, ip, port, apiPort, dbPort string) {
	var id string

	if idPointer == nil {
		id = generateNodeId()
	} else {
		id = *idPointer
	}
	theLocalNode = &localNode{_id: id}
	transport = newTransporter(ip, port, apiPort, dbPort)
	theLocalNode.initPrimaryAndReplicaDB(theLocalNode.id())
}

// ----------------------------------------------------------------------------------------
//										Getters + setters
// ----------------------------------------------------------------------------------------
func (n *localNode) id() string {
	return n._id
}

func (n *localNode) ip() string {
	return transport.Ip
}

func (n *localNode) port() string {
	return transport.Port
}

func (n *localNode) apiPort() string {
	return transport.ApiPort
}

func (n *localNode) dbPort() string {
	return transport.DbPort
}

func (n *localNode) predecessor() node {
	return n.pred
}

func (n *localNode) successor() node {
	return n.fingerTable[1].node
}

func (n *localNode) address() string {
	return transport.Ip + ":" + transport.Port
}

func (n *localNode) apiAddress() string {
	return transport.Ip + ":" + transport.ApiPort
}

func (n *localNode) dbAddress() string {
	return transport.Ip + ":" + transport.DbPort
}

func (n *localNode) updatePredecessor(candidate node) {
	if between(hexStringToByteArr(n.predecessor().id()), hexStringToByteArr(n.id()), hexStringToByteArr(candidate.id())) {
		n.pred = candidate
		log.Tracef("%s: Predecessor updated to: %s", theLocalNode.address(), candidate.id())
	} else {
		log.Tracef("%s: Predecessor NOT updated to: %s", theLocalNode.address(), candidate.id())
	}
	theLocalNode.fixFingersChan <- true
}

func (n *localNode) updateSuccessor(candidate node) {
	if between(hexStringToByteArr(n.id()), hexStringToByteArr(n.successor().id()), hexStringToByteArr(candidate.id())) {
		n.fingerTable[1].node = candidate
		log.Tracef("%s: Successor updated to: %s", theLocalNode.address(), candidate.id())
	} else {
		log.Tracef("%s: Successor NOT updated to: %s", theLocalNode.address(), candidate.id())
	}
	theLocalNode.fixFingersChan <- true
}

// ----------------------------------------------------------------------------------------
//										localNode methods
// ----------------------------------------------------------------------------------------

// Returns the node who is responsible for the data corresponding to id, traversing the ring using finger tables
func (n *localNode) lookup(id string) (node, error) {
	// n responsible for id
	if between(
		hexStringToByteArr(nextId(n.predecessor().id())),
		hexStringToByteArr(nextId(n.id())),
		hexStringToByteArr(id),
	) {
		return n, nil
		// otherwise use fingers of n, starting with the one that is furthest away, to find responsible node
	} else {
		for k := m; k > 0; {

			nextNode, i := theLocalNode.forwardingLookup(id, k)

			responsibleNode, err := nextNode.lookup(id)

			if err == nil {
				return responsibleNode, nil
			} else {
				k = i - 1

			}
		}
		// all fingers dead... what to do? fixfingers! and fix predAndsucc! don´t send ACK?
		return nil, nil
	}
}

func (n *localNode) forwardingLookup(id string, j int) (node, int) {
	for i := j; i > 0; i-- {
		// special case - when n´s finger points to itself
		if n.fingerTable[i].node.id() == n.id() {

			// what to do? go to next finger...
			// id between finger and node - got to that finger
		} else if between(
			hexStringToByteArr(n.fingerTable[i].node.id()),
			hexStringToByteArr(n.id()),
			hexStringToByteArr(id),
		) {
			return n.fingerTable[i].node, i
		}
	}
	// if id is not between any finger and n - then id must be between n and its successor
	return n.successor(), 0
}

// lookup of finger.node for the case when a second node is added to a ring with only one node
func (newNode *localNode) specLookup(n *remoteNode, startId string) node {
	if between(
		hexStringToByteArr(nextId(n.id())),
		hexStringToByteArr(nextId(newNode.id())),
		hexStringToByteArr(startId),
	) {
		// newNodes first finger/successor is newNode itself
		return newNode
	}
	// newNodes first finger/successor is n
	return n
}

// n needs to be an up and running remote node in the ring
func (newNode *localNode) join(n *remoteNode) {

	// If newNode is the only node in the network
	if n == nil {
		newNode.pred = newNode
		for i := 1; i <= m; i++ {
			newNode.fingerTable[i].startId, _ = calcFinger(hexStringToByteArr(newNode.id()), i, m)
			newNode.fingerTable[i].node = newNode
		}

	} else {
		newNode.initFingers(n)
	}
	go newNode.startFixFingers()
}

func (newNode *localNode) initFingers(n *remoteNode) {
	oneNodeRing := false

	// Calculating first finger
	newNode.fingerTable[1].startId, _ = calcFinger(hexStringToByteArr(newNode.id()), 1, m)

	// Successor to newNode
	newNode.fingerTable[1].node, _ = n.lookup(newNode.fingerTable[1].startId)
	log.Tracef("%s: Set successor to: %s", theLocalNode.address(), newNode.successor().id())

	// Predecessor to newNode
	newNode.pred = newNode.successor().predecessor()
	log.Tracef("%s: Set predecessor to: %s", theLocalNode.address(), newNode.predecessor().id())

	if newNode.successor().id() == newNode.predecessor().id() { // n.predecessor().id() == n.id() {
		oneNodeRing = true
	}

	// backup predecessors db and takeover correct part of successors db
	newNode.startReplication()

	// Set successor of newNode´s predecessor to newNode  						<----------- should be made sync!
	newNode.predecessor().updateSuccessor(newNode)

	// Update the predecessor of the node that newNode is inserted before  	  	<---------- should be made sync!
	newNode.successor().updatePredecessor(newNode)

	// request successor node to split its primary and replace its previous replica with part from its primary
	newNode.requestSplit(newNode.successor())

	for i := 1; i < m; i++ {

		// Calculating finger
		newNode.fingerTable[i+1].startId, _ = calcFinger(hexStringToByteArr(newNode.id()), i+1, m)

		if between(
			hexStringToByteArr(newNode.id()),
			hexStringToByteArr(nextId(newNode.fingerTable[i].node.id())),
			hexStringToByteArr(newNode.fingerTable[i+1].startId),
		) {
			// startId between newNode and previous finger.node
			newNode.fingerTable[i+1].node = newNode.fingerTable[i].node

		} else {

			if oneNodeRing {
				newNode.fingerTable[i+1].node = newNode.specLookup(n, newNode.fingerTable[i+1].startId)

			} else {
				newNode.fingerTable[i+1].node, _ = n.lookup(newNode.fingerTable[i+1].startId)
				log.Tracef("%s: In join, set finger %s to %s", theLocalNode.address(), newNode.fingerTable[i+1].startId, newNode.fingerTable[i+1].node.id())
			}
		}
	}
}

func (n *localNode) startFixFingers() {
	for {
		select {
		case <-time.After(time.Second * 10):
			n.fixFingers()
		case <-theLocalNode.fixFingersChan:
			n.fixFingers()
		}
	}
}

// an optimized fixFingers could be called from updateSuccessor and updatePredecessor.
// In those functions you get a candidate. Use the candidates id
// to determine for each finger if it should be updated to the candidate or not.
// fixFingers without any remote lookup requests.
// se old updateFingerTable.

// Called periodically to update fingers
func (n *localNode) fixFingers() {

	//log.Tracef("%s: Running fixFingers", theLocalNode.address())

	succ, _ := n.lookup(n.fingerTable[1].startId)
	if succ.id() != n.successor().id() {
		n.updateSuccessor(succ)
	}

	for i := 1; i < m; i++ {

		if between(
			hexStringToByteArr(n.id()),
			hexStringToByteArr(nextId(n.fingerTable[i].node.id())),
			hexStringToByteArr(n.fingerTable[i+1].startId),
		) {
			// startId between n and previous finger.node
			n.fingerTable[i+1].node = n.fingerTable[i].node
		} else {
			n.fingerTable[i+1].node, _ = n.lookup(n.fingerTable[i+1].startId)
			if theLocalNode.address() == "localhost:2000" {
				//				log.Tracef("%s: In fixFingers: Lookuped and updated finger %s to: %s", theLocalNode.address(), n.fingerTable[i+1].startId, n.fingerTable[i+1].node.id())
			}
		}
	}
}

/*
// If s should be the i:th finger of n -> update n's finger table entry nr i with s
func (n *localNode) updateFingerTable(s *localNode, i int) {
	// n´s finger.node points to n itself
	if n.id() == n.fingerTable[i].node.id() {
		if between(
			hexStringToByteArr(n.fingerTable[i].startId),
			hexStringToByteArr(n.fingerTable[i].node.id()),
			hexStringToByteArr(s.id()),
		) {
			n.fingerTable[i].node = s
			//			fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id(), i, n.id(), n.id(), i, s.id())
		}
	} else if between(
		hexStringToByteArr(n.fingerTable[i].startId),
		hexStringToByteArr(n.fingerTable[i].node.id()),
		hexStringToByteArr(s.id()),
	) {
		if !(n.fingerTable[i].startId == n.fingerTable[i].node.id()) {
			n.fingerTable[i].node = s
			//				fmt.Printf("%s should be the %d:th finger of %s -> update %s's finger table entry nr %d with %s\n", s.id(), i, n.id(), n.id(), i, s.id())
		}
	}
	// Get last node preceeding n and mayby update its finger i as well, check that it hasn´t come round to the node just added (s)
	p := n.predecessor()
	//		fmt.Printf("p is: %s\n", p.id())
	if p.id() != s.id() {
		p.updateFingerTable(s, i)
	} else {
		//			fmt.Println("Not going to that node")
	}
}*/
