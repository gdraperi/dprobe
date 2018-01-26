package networkdb

import (
	"fmt"

	"github.com/hashicorp/memberlist"
	"github.com/sirupsen/logrus"
)

type nodeState int

const (
	nodeNotFound    nodeState = -1
	nodeActiveState nodeState = 0
	nodeLeftState   nodeState = 1
	nodeFailedState nodeState = 2
)

var nodeStateName = map[nodeState]string***REMOVED***
	-1: "NodeNotFound",
	0:  "NodeActive",
	1:  "NodeLeft",
	2:  "NodeFailed",
***REMOVED***

// findNode search the node into the 3 node lists and returns the node pointer and the list
// where it got found
func (nDB *NetworkDB) findNode(nodeName string) (*node, nodeState, map[string]*node) ***REMOVED***
	for i, nodes := range []map[string]*node***REMOVED***
		nDB.nodes,
		nDB.leftNodes,
		nDB.failedNodes,
	***REMOVED*** ***REMOVED***
		if n, ok := nodes[nodeName]; ok ***REMOVED***
			return n, nodeState(i), nodes
		***REMOVED***
	***REMOVED***
	return nil, nodeNotFound, nil
***REMOVED***

// changeNodeState changes the state of the node specified, returns true if the node was moved,
// false if there was no need to change the node state. Error will be returned if the node does not
// exists
func (nDB *NetworkDB) changeNodeState(nodeName string, newState nodeState) (bool, error) ***REMOVED***
	n, currState, m := nDB.findNode(nodeName)
	if n == nil ***REMOVED***
		return false, fmt.Errorf("node %s not found", nodeName)
	***REMOVED***

	switch newState ***REMOVED***
	case nodeActiveState:
		if currState == nodeActiveState ***REMOVED***
			return false, nil
		***REMOVED***

		delete(m, nodeName)
		// reset the node reap time
		n.reapTime = 0
		nDB.nodes[nodeName] = n
	case nodeLeftState:
		if currState == nodeLeftState ***REMOVED***
			return false, nil
		***REMOVED***

		delete(m, nodeName)
		nDB.leftNodes[nodeName] = n
	case nodeFailedState:
		if currState == nodeFailedState ***REMOVED***
			return false, nil
		***REMOVED***

		delete(m, nodeName)
		nDB.failedNodes[nodeName] = n
	***REMOVED***

	logrus.Infof("Node %s change state %s --> %s", nodeName, nodeStateName[currState], nodeStateName[newState])

	if newState == nodeLeftState || newState == nodeFailedState ***REMOVED***
		// set the node reap time, if not already set
		// It is possible that a node passes from failed to left and the reaptime was already set so keep that value
		if n.reapTime == 0 ***REMOVED***
			n.reapTime = nodeReapInterval
		***REMOVED***
		// The node leave or fails, delete all the entries created by it.
		// If the node was temporary down, deleting the entries will guarantee that the CREATE events will be accepted
		// If the node instead left because was going down, then it makes sense to just delete all its state
		nDB.deleteNodeFromNetworks(n.Name)
		nDB.deleteNodeTableEntries(n.Name)
	***REMOVED***

	return true, nil
***REMOVED***

func (nDB *NetworkDB) purgeReincarnation(mn *memberlist.Node) bool ***REMOVED***
	for name, node := range nDB.nodes ***REMOVED***
		if node.Addr.Equal(mn.Addr) && node.Port == mn.Port && mn.Name != name ***REMOVED***
			logrus.Infof("Node %s/%s, is the new incarnation of the active node %s/%s", mn.Name, mn.Addr, name, node.Addr)
			nDB.changeNodeState(name, nodeLeftState)
			return true
		***REMOVED***
	***REMOVED***

	for name, node := range nDB.failedNodes ***REMOVED***
		if node.Addr.Equal(mn.Addr) && node.Port == mn.Port && mn.Name != name ***REMOVED***
			logrus.Infof("Node %s/%s, is the new incarnation of the failed node %s/%s", mn.Name, mn.Addr, name, node.Addr)
			nDB.changeNodeState(name, nodeLeftState)
			return true
		***REMOVED***
	***REMOVED***

	for name, node := range nDB.leftNodes ***REMOVED***
		if node.Addr.Equal(mn.Addr) && node.Port == mn.Port && mn.Name != name ***REMOVED***
			logrus.Infof("Node %s/%s, is the new incarnation of the shutdown node %s/%s", mn.Name, mn.Addr, name, node.Addr)
			nDB.changeNodeState(name, nodeLeftState)
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
