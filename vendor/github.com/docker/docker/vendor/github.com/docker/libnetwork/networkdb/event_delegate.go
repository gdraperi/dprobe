package networkdb

import (
	"encoding/json"
	"net"

	"github.com/hashicorp/memberlist"
	"github.com/sirupsen/logrus"
)

type eventDelegate struct ***REMOVED***
	nDB *NetworkDB
***REMOVED***

func (e *eventDelegate) broadcastNodeEvent(addr net.IP, op opType) ***REMOVED***
	value, err := json.Marshal(&NodeAddr***REMOVED***addr***REMOVED***)
	if err == nil ***REMOVED***
		e.nDB.broadcaster.Write(makeEvent(op, NodeTable, "", "", value))
	***REMOVED*** else ***REMOVED***
		logrus.Errorf("Error marshalling node broadcast event %s", addr.String())
	***REMOVED***
***REMOVED***

func (e *eventDelegate) NotifyJoin(mn *memberlist.Node) ***REMOVED***
	logrus.Infof("Node %s/%s, joined gossip cluster", mn.Name, mn.Addr)
	e.broadcastNodeEvent(mn.Addr, opCreate)
	e.nDB.Lock()
	defer e.nDB.Unlock()
	// In case the node is rejoining after a failure or leave,
	// wait until an explicit join message arrives before adding
	// it to the nodes just to make sure this is not a stale
	// join. If you don't know about this node add it immediately.
	_, fOk := e.nDB.failedNodes[mn.Name]
	_, lOk := e.nDB.leftNodes[mn.Name]
	if fOk || lOk ***REMOVED***
		return
	***REMOVED***

	// Every node has a unique ID
	// Check on the base of the IP address if the new node that joined is actually a new incarnation of a previous
	// failed or shutdown one
	e.nDB.purgeReincarnation(mn)

	e.nDB.nodes[mn.Name] = &node***REMOVED***Node: *mn***REMOVED***
	logrus.Infof("Node %s/%s, added to nodes list", mn.Name, mn.Addr)
***REMOVED***

func (e *eventDelegate) NotifyLeave(mn *memberlist.Node) ***REMOVED***
	logrus.Infof("Node %s/%s, left gossip cluster", mn.Name, mn.Addr)
	e.broadcastNodeEvent(mn.Addr, opDelete)

	e.nDB.Lock()
	defer e.nDB.Unlock()

	n, currState, _ := e.nDB.findNode(mn.Name)
	if n == nil ***REMOVED***
		logrus.Errorf("Node %s/%s not found in the node lists", mn.Name, mn.Addr)
		return
	***REMOVED***
	// if the node was active means that did not send the leave cluster message, so it's probable that
	// failed. Else would be already in the left list so nothing else has to be done
	if currState == nodeActiveState ***REMOVED***
		moved, err := e.nDB.changeNodeState(mn.Name, nodeFailedState)
		if err != nil ***REMOVED***
			logrus.WithError(err).Errorf("impossible condition, node %s/%s not present in the list", mn.Name, mn.Addr)
			return
		***REMOVED***
		if moved ***REMOVED***
			logrus.Infof("Node %s/%s, added to failed nodes list", mn.Name, mn.Addr)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *eventDelegate) NotifyUpdate(n *memberlist.Node) ***REMOVED***
***REMOVED***
