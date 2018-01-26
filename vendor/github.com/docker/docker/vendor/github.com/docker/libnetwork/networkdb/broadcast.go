package networkdb

import (
	"errors"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

const broadcastTimeout = 5 * time.Second

type networkEventMessage struct ***REMOVED***
	id   string
	node string
	msg  []byte
***REMOVED***

func (m *networkEventMessage) Invalidates(other memberlist.Broadcast) bool ***REMOVED***
	otherm := other.(*networkEventMessage)
	return m.id == otherm.id && m.node == otherm.node
***REMOVED***

func (m *networkEventMessage) Message() []byte ***REMOVED***
	return m.msg
***REMOVED***

func (m *networkEventMessage) Finished() ***REMOVED***
***REMOVED***

func (nDB *NetworkDB) sendNetworkEvent(nid string, event NetworkEvent_Type, ltime serf.LamportTime) error ***REMOVED***
	nEvent := NetworkEvent***REMOVED***
		Type:      event,
		LTime:     ltime,
		NodeName:  nDB.config.NodeID,
		NetworkID: nid,
	***REMOVED***

	raw, err := encodeMessage(MessageTypeNetworkEvent, &nEvent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	nDB.networkBroadcasts.QueueBroadcast(&networkEventMessage***REMOVED***
		msg:  raw,
		id:   nid,
		node: nDB.config.NodeID,
	***REMOVED***)
	return nil
***REMOVED***

type nodeEventMessage struct ***REMOVED***
	msg    []byte
	notify chan<- struct***REMOVED******REMOVED***
***REMOVED***

func (m *nodeEventMessage) Invalidates(other memberlist.Broadcast) bool ***REMOVED***
	return false
***REMOVED***

func (m *nodeEventMessage) Message() []byte ***REMOVED***
	return m.msg
***REMOVED***

func (m *nodeEventMessage) Finished() ***REMOVED***
	if m.notify != nil ***REMOVED***
		close(m.notify)
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) sendNodeEvent(event NodeEvent_Type) error ***REMOVED***
	nEvent := NodeEvent***REMOVED***
		Type:     event,
		LTime:    nDB.networkClock.Increment(),
		NodeName: nDB.config.NodeID,
	***REMOVED***

	raw, err := encodeMessage(MessageTypeNodeEvent, &nEvent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	notifyCh := make(chan struct***REMOVED******REMOVED***)
	nDB.nodeBroadcasts.QueueBroadcast(&nodeEventMessage***REMOVED***
		msg:    raw,
		notify: notifyCh,
	***REMOVED***)

	nDB.RLock()
	noPeers := len(nDB.nodes) <= 1
	nDB.RUnlock()

	// Message enqueued, do not wait for a send if no peer is present
	if noPeers ***REMOVED***
		return nil
	***REMOVED***

	// Wait for the broadcast
	select ***REMOVED***
	case <-notifyCh:
	case <-time.After(broadcastTimeout):
		return errors.New("timed out broadcasting node event")
	***REMOVED***

	return nil
***REMOVED***

type tableEventMessage struct ***REMOVED***
	id    string
	tname string
	key   string
	msg   []byte
	node  string
***REMOVED***

func (m *tableEventMessage) Invalidates(other memberlist.Broadcast) bool ***REMOVED***
	otherm := other.(*tableEventMessage)
	return m.tname == otherm.tname && m.id == otherm.id && m.key == otherm.key
***REMOVED***

func (m *tableEventMessage) Message() []byte ***REMOVED***
	return m.msg
***REMOVED***

func (m *tableEventMessage) Finished() ***REMOVED***
***REMOVED***

func (nDB *NetworkDB) sendTableEvent(event TableEvent_Type, nid string, tname string, key string, entry *entry) error ***REMOVED***
	tEvent := TableEvent***REMOVED***
		Type:      event,
		LTime:     entry.ltime,
		NodeName:  nDB.config.NodeID,
		NetworkID: nid,
		TableName: tname,
		Key:       key,
		Value:     entry.value,
		// The duration in second is a float that below would be truncated
		ResidualReapTime: int32(entry.reapTime.Seconds()),
	***REMOVED***

	raw, err := encodeMessage(MessageTypeTableEvent, &tEvent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var broadcastQ *memberlist.TransmitLimitedQueue
	nDB.RLock()
	thisNodeNetworks, ok := nDB.networks[nDB.config.NodeID]
	if ok ***REMOVED***
		// The network may have been removed
		network, networkOk := thisNodeNetworks[nid]
		if !networkOk ***REMOVED***
			nDB.RUnlock()
			return nil
		***REMOVED***

		broadcastQ = network.tableBroadcasts
	***REMOVED***
	nDB.RUnlock()

	// The network may have been removed
	if broadcastQ == nil ***REMOVED***
		return nil
	***REMOVED***

	broadcastQ.QueueBroadcast(&tableEventMessage***REMOVED***
		msg:   raw,
		id:    nid,
		tname: tname,
		key:   key,
		node:  nDB.config.NodeID,
	***REMOVED***)
	return nil
***REMOVED***
