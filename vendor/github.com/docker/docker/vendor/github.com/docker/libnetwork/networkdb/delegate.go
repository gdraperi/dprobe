package networkdb

import (
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

type delegate struct ***REMOVED***
	nDB *NetworkDB
***REMOVED***

func (d *delegate) NodeMeta(limit int) []byte ***REMOVED***
	return []byte***REMOVED******REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleNodeEvent(nEvent *NodeEvent) bool ***REMOVED***
	// Update our local clock if the received messages has newer
	// time.
	nDB.networkClock.Witness(nEvent.LTime)

	nDB.RLock()
	defer nDB.RUnlock()

	// check if the node exists
	n, _, _ := nDB.findNode(nEvent.NodeName)
	if n == nil ***REMOVED***
		return false
	***REMOVED***

	// check if the event is fresh
	if n.ltime >= nEvent.LTime ***REMOVED***
		return false
	***REMOVED***

	// If we are here means that the event is fresher and the node is known. Update the laport time
	n.ltime = nEvent.LTime

	// If it is a node leave event for a manager and this is the only manager we
	// know of we want the reconnect logic to kick in. In a single manager
	// cluster manager's gossip can't be bootstrapped unless some other node
	// connects to it.
	if len(nDB.bootStrapIP) == 1 && nEvent.Type == NodeEventTypeLeave ***REMOVED***
		for _, ip := range nDB.bootStrapIP ***REMOVED***
			if ip.Equal(n.Addr) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch nEvent.Type ***REMOVED***
	case NodeEventTypeJoin:
		moved, err := nDB.changeNodeState(n.Name, nodeActiveState)
		if err != nil ***REMOVED***
			logrus.WithError(err).Error("unable to find the node to move")
			return false
		***REMOVED***
		if moved ***REMOVED***
			logrus.Infof("%v(%v): Node join event for %s/%s", nDB.config.Hostname, nDB.config.NodeID, n.Name, n.Addr)
		***REMOVED***
		return moved
	case NodeEventTypeLeave:
		moved, err := nDB.changeNodeState(n.Name, nodeLeftState)
		if err != nil ***REMOVED***
			logrus.WithError(err).Error("unable to find the node to move")
			return false
		***REMOVED***
		if moved ***REMOVED***
			logrus.Infof("%v(%v): Node leave event for %s/%s", nDB.config.Hostname, nDB.config.NodeID, n.Name, n.Addr)
		***REMOVED***
		return moved
	***REMOVED***

	return false
***REMOVED***

func (nDB *NetworkDB) handleNetworkEvent(nEvent *NetworkEvent) bool ***REMOVED***
	// Update our local clock if the received messages has newer
	// time.
	nDB.networkClock.Witness(nEvent.LTime)

	nDB.Lock()
	defer nDB.Unlock()

	if nEvent.NodeName == nDB.config.NodeID ***REMOVED***
		return false
	***REMOVED***

	nodeNetworks, ok := nDB.networks[nEvent.NodeName]
	if !ok ***REMOVED***
		// We haven't heard about this node at all.  Ignore the leave
		if nEvent.Type == NetworkEventTypeLeave ***REMOVED***
			return false
		***REMOVED***

		nodeNetworks = make(map[string]*network)
		nDB.networks[nEvent.NodeName] = nodeNetworks
	***REMOVED***

	if n, ok := nodeNetworks[nEvent.NetworkID]; ok ***REMOVED***
		// We have the latest state. Ignore the event
		// since it is stale.
		if n.ltime >= nEvent.LTime ***REMOVED***
			return false
		***REMOVED***

		n.ltime = nEvent.LTime
		n.leaving = nEvent.Type == NetworkEventTypeLeave
		if n.leaving ***REMOVED***
			n.reapTime = nDB.config.reapNetworkInterval

			// The remote node is leaving the network, but not the gossip cluster.
			// Mark all its entries in deleted state, this will guarantee that
			// if some node bulk sync with us, the deleted state of
			// these entries will be propagated.
			nDB.deleteNodeNetworkEntries(nEvent.NetworkID, nEvent.NodeName)
		***REMOVED***

		if nEvent.Type == NetworkEventTypeLeave ***REMOVED***
			nDB.deleteNetworkNode(nEvent.NetworkID, nEvent.NodeName)
		***REMOVED*** else ***REMOVED***
			nDB.addNetworkNode(nEvent.NetworkID, nEvent.NodeName)
		***REMOVED***

		return true
	***REMOVED***

	if nEvent.Type == NetworkEventTypeLeave ***REMOVED***
		return false
	***REMOVED***

	// If the node is not known from memberlist we cannot process save any state of it else if it actually
	// dies we won't receive any notification and we will remain stuck with it
	if _, ok := nDB.nodes[nEvent.NodeName]; !ok ***REMOVED***
		return false
	***REMOVED***

	// This remote network join is being seen the first time.
	nodeNetworks[nEvent.NetworkID] = &network***REMOVED***
		id:    nEvent.NetworkID,
		ltime: nEvent.LTime,
	***REMOVED***

	nDB.addNetworkNode(nEvent.NetworkID, nEvent.NodeName)
	return true
***REMOVED***

func (nDB *NetworkDB) handleTableEvent(tEvent *TableEvent) bool ***REMOVED***
	// Update our local clock if the received messages has newer time.
	nDB.tableClock.Witness(tEvent.LTime)

	// Ignore the table events for networks that are in the process of going away
	nDB.RLock()
	networks := nDB.networks[nDB.config.NodeID]
	network, ok := networks[tEvent.NetworkID]
	// Check if the owner of the event is still part of the network
	nodes := nDB.networkNodes[tEvent.NetworkID]
	var nodePresent bool
	for _, node := range nodes ***REMOVED***
		if node == tEvent.NodeName ***REMOVED***
			nodePresent = true
			break
		***REMOVED***
	***REMOVED***
	nDB.RUnlock()
	if !ok || network.leaving || !nodePresent ***REMOVED***
		// I'm out of the network OR the event owner is not anymore part of the network so do not propagate
		return false
	***REMOVED***

	e, err := nDB.getEntry(tEvent.TableName, tEvent.NetworkID, tEvent.Key)
	if err == nil ***REMOVED***
		// We have the latest state. Ignore the event
		// since it is stale.
		if e.ltime >= tEvent.LTime ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	e = &entry***REMOVED***
		ltime:    tEvent.LTime,
		node:     tEvent.NodeName,
		value:    tEvent.Value,
		deleting: tEvent.Type == TableEventTypeDelete,
		reapTime: time.Duration(tEvent.ResidualReapTime) * time.Second,
	***REMOVED***

	// All the entries marked for deletion should have a reapTime set greater than 0
	// This case can happen if the cluster is running different versions of the engine where the old version does not have the
	// field. If that is not the case, this can be a BUG
	if e.deleting && e.reapTime == 0 ***REMOVED***
		logrus.Warnf("%v(%v) handleTableEvent object %+v has a 0 reapTime, is the cluster running the same docker engine version?",
			nDB.config.Hostname, nDB.config.NodeID, tEvent)
		e.reapTime = nDB.config.reapEntryInterval
	***REMOVED***

	nDB.Lock()
	nDB.createOrUpdateEntry(tEvent.NetworkID, tEvent.TableName, tEvent.Key, e)
	nDB.Unlock()

	if err != nil && tEvent.Type == TableEventTypeDelete ***REMOVED***
		// If it is a delete event and we did not have a state for it, don't propagate to the application
		// If the residual reapTime is lower or equal to 1/6 of the total reapTime don't bother broadcasting it around
		// most likely the cluster is already aware of it, if not who will sync with this node will catch the state too.
		// This also avoids that deletion of entries close to their garbage collection ends up circuling around forever
		return e.reapTime > nDB.config.reapEntryInterval/6
	***REMOVED***

	var op opType
	switch tEvent.Type ***REMOVED***
	case TableEventTypeCreate:
		op = opCreate
	case TableEventTypeUpdate:
		op = opUpdate
	case TableEventTypeDelete:
		op = opDelete
	***REMOVED***

	nDB.broadcaster.Write(makeEvent(op, tEvent.TableName, tEvent.NetworkID, tEvent.Key, tEvent.Value))
	return true
***REMOVED***

func (nDB *NetworkDB) handleCompound(buf []byte, isBulkSync bool) ***REMOVED***
	// Decode the parts
	parts, err := decodeCompoundMessage(buf)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to decode compound request: %v", err)
		return
	***REMOVED***

	// Handle each message
	for _, part := range parts ***REMOVED***
		nDB.handleMessage(part, isBulkSync)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleTableMessage(buf []byte, isBulkSync bool) ***REMOVED***
	var tEvent TableEvent
	if err := proto.Unmarshal(buf, &tEvent); err != nil ***REMOVED***
		logrus.Errorf("Error decoding table event message: %v", err)
		return
	***REMOVED***

	// Ignore messages that this node generated.
	if tEvent.NodeName == nDB.config.NodeID ***REMOVED***
		return
	***REMOVED***

	if rebroadcast := nDB.handleTableEvent(&tEvent); rebroadcast ***REMOVED***
		var err error
		buf, err = encodeRawMessage(MessageTypeTableEvent, buf)
		if err != nil ***REMOVED***
			logrus.Errorf("Error marshalling gossip message for network event rebroadcast: %v", err)
			return
		***REMOVED***

		nDB.RLock()
		n, ok := nDB.networks[nDB.config.NodeID][tEvent.NetworkID]
		nDB.RUnlock()

		// if the network is not there anymore, OR we are leaving the network OR the broadcast queue is not present
		if !ok || n.leaving || n.tableBroadcasts == nil ***REMOVED***
			return
		***REMOVED***

		n.tableBroadcasts.QueueBroadcast(&tableEventMessage***REMOVED***
			msg:   buf,
			id:    tEvent.NetworkID,
			tname: tEvent.TableName,
			key:   tEvent.Key,
			node:  tEvent.NodeName,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleNodeMessage(buf []byte) ***REMOVED***
	var nEvent NodeEvent
	if err := proto.Unmarshal(buf, &nEvent); err != nil ***REMOVED***
		logrus.Errorf("Error decoding node event message: %v", err)
		return
	***REMOVED***

	if rebroadcast := nDB.handleNodeEvent(&nEvent); rebroadcast ***REMOVED***
		var err error
		buf, err = encodeRawMessage(MessageTypeNodeEvent, buf)
		if err != nil ***REMOVED***
			logrus.Errorf("Error marshalling gossip message for node event rebroadcast: %v", err)
			return
		***REMOVED***

		nDB.nodeBroadcasts.QueueBroadcast(&nodeEventMessage***REMOVED***
			msg: buf,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleNetworkMessage(buf []byte) ***REMOVED***
	var nEvent NetworkEvent
	if err := proto.Unmarshal(buf, &nEvent); err != nil ***REMOVED***
		logrus.Errorf("Error decoding network event message: %v", err)
		return
	***REMOVED***

	if rebroadcast := nDB.handleNetworkEvent(&nEvent); rebroadcast ***REMOVED***
		var err error
		buf, err = encodeRawMessage(MessageTypeNetworkEvent, buf)
		if err != nil ***REMOVED***
			logrus.Errorf("Error marshalling gossip message for network event rebroadcast: %v", err)
			return
		***REMOVED***

		nDB.networkBroadcasts.QueueBroadcast(&networkEventMessage***REMOVED***
			msg:  buf,
			id:   nEvent.NetworkID,
			node: nEvent.NodeName,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleBulkSync(buf []byte) ***REMOVED***
	var bsm BulkSyncMessage
	if err := proto.Unmarshal(buf, &bsm); err != nil ***REMOVED***
		logrus.Errorf("Error decoding bulk sync message: %v", err)
		return
	***REMOVED***

	if bsm.LTime > 0 ***REMOVED***
		nDB.tableClock.Witness(bsm.LTime)
	***REMOVED***

	nDB.handleMessage(bsm.Payload, true)

	// Don't respond to a bulk sync which was not unsolicited
	if !bsm.Unsolicited ***REMOVED***
		nDB.Lock()
		ch, ok := nDB.bulkSyncAckTbl[bsm.NodeName]
		if ok ***REMOVED***
			close(ch)
			delete(nDB.bulkSyncAckTbl, bsm.NodeName)
		***REMOVED***
		nDB.Unlock()

		return
	***REMOVED***

	var nodeAddr net.IP
	nDB.RLock()
	if node, ok := nDB.nodes[bsm.NodeName]; ok ***REMOVED***
		nodeAddr = node.Addr
	***REMOVED***
	nDB.RUnlock()

	if err := nDB.bulkSyncNode(bsm.Networks, bsm.NodeName, false); err != nil ***REMOVED***
		logrus.Errorf("Error in responding to bulk sync from node %s: %v", nodeAddr, err)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) handleMessage(buf []byte, isBulkSync bool) ***REMOVED***
	mType, data, err := decodeMessage(buf)
	if err != nil ***REMOVED***
		logrus.Errorf("Error decoding gossip message to get message type: %v", err)
		return
	***REMOVED***

	switch mType ***REMOVED***
	case MessageTypeNodeEvent:
		nDB.handleNodeMessage(data)
	case MessageTypeNetworkEvent:
		nDB.handleNetworkMessage(data)
	case MessageTypeTableEvent:
		nDB.handleTableMessage(data, isBulkSync)
	case MessageTypeBulkSync:
		nDB.handleBulkSync(data)
	case MessageTypeCompound:
		nDB.handleCompound(data, isBulkSync)
	default:
		logrus.Errorf("%v(%v): unknown message type %d", nDB.config.Hostname, nDB.config.NodeID, mType)
	***REMOVED***
***REMOVED***

func (d *delegate) NotifyMsg(buf []byte) ***REMOVED***
	if len(buf) == 0 ***REMOVED***
		return
	***REMOVED***

	d.nDB.handleMessage(buf, false)
***REMOVED***

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte ***REMOVED***
	msgs := d.nDB.networkBroadcasts.GetBroadcasts(overhead, limit)
	msgs = append(msgs, d.nDB.nodeBroadcasts.GetBroadcasts(overhead, limit)...)
	return msgs
***REMOVED***

func (d *delegate) LocalState(join bool) []byte ***REMOVED***
	if join ***REMOVED***
		// Update all the local node/network state to a new time to
		// force update on the node we are trying to rejoin, just in
		// case that node has these in leaving state still. This is
		// facilitate fast convergence after recovering from a gossip
		// failure.
		d.nDB.updateLocalNetworkTime()
	***REMOVED***

	d.nDB.RLock()
	defer d.nDB.RUnlock()

	pp := NetworkPushPull***REMOVED***
		LTime:    d.nDB.networkClock.Time(),
		NodeName: d.nDB.config.NodeID,
	***REMOVED***

	for name, nn := range d.nDB.networks ***REMOVED***
		for _, n := range nn ***REMOVED***
			pp.Networks = append(pp.Networks, &NetworkEntry***REMOVED***
				LTime:     n.ltime,
				NetworkID: n.id,
				NodeName:  name,
				Leaving:   n.leaving,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	buf, err := encodeMessage(MessageTypePushPull, &pp)
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to encode local network state: %v", err)
		return nil
	***REMOVED***

	return buf
***REMOVED***

func (d *delegate) MergeRemoteState(buf []byte, isJoin bool) ***REMOVED***
	if len(buf) == 0 ***REMOVED***
		logrus.Error("zero byte remote network state received")
		return
	***REMOVED***

	var gMsg GossipMessage
	err := proto.Unmarshal(buf, &gMsg)
	if err != nil ***REMOVED***
		logrus.Errorf("Error unmarshalling push pull message: %v", err)
		return
	***REMOVED***

	if gMsg.Type != MessageTypePushPull ***REMOVED***
		logrus.Errorf("Invalid message type %v received from remote", buf[0])
	***REMOVED***

	pp := NetworkPushPull***REMOVED******REMOVED***
	if err := proto.Unmarshal(gMsg.Data, &pp); err != nil ***REMOVED***
		logrus.Errorf("Failed to decode remote network state: %v", err)
		return
	***REMOVED***

	nodeEvent := &NodeEvent***REMOVED***
		LTime:    pp.LTime,
		NodeName: pp.NodeName,
		Type:     NodeEventTypeJoin,
	***REMOVED***
	d.nDB.handleNodeEvent(nodeEvent)

	for _, n := range pp.Networks ***REMOVED***
		nEvent := &NetworkEvent***REMOVED***
			LTime:     n.LTime,
			NodeName:  n.NodeName,
			NetworkID: n.NetworkID,
			Type:      NetworkEventTypeJoin,
		***REMOVED***

		if n.Leaving ***REMOVED***
			nEvent.Type = NetworkEventTypeLeave
		***REMOVED***

		d.nDB.handleNetworkEvent(nEvent)
	***REMOVED***

***REMOVED***
