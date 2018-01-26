package networkdb

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	rnd "math/rand"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/sirupsen/logrus"
)

const (
	reapPeriod       = 5 * time.Second
	retryInterval    = 1 * time.Second
	nodeReapInterval = 24 * time.Hour
	nodeReapPeriod   = 2 * time.Hour
)

type logWriter struct***REMOVED******REMOVED***

func (l *logWriter) Write(p []byte) (int, error) ***REMOVED***
	str := string(p)
	str = strings.TrimSuffix(str, "\n")

	switch ***REMOVED***
	case strings.HasPrefix(str, "[WARN] "):
		str = strings.TrimPrefix(str, "[WARN] ")
		logrus.Warn(str)
	case strings.HasPrefix(str, "[DEBUG] "):
		str = strings.TrimPrefix(str, "[DEBUG] ")
		logrus.Debug(str)
	case strings.HasPrefix(str, "[INFO] "):
		str = strings.TrimPrefix(str, "[INFO] ")
		logrus.Info(str)
	case strings.HasPrefix(str, "[ERR] "):
		str = strings.TrimPrefix(str, "[ERR] ")
		logrus.Warn(str)
	***REMOVED***

	return len(p), nil
***REMOVED***

// SetKey adds a new key to the key ring
func (nDB *NetworkDB) SetKey(key []byte) ***REMOVED***
	logrus.Debugf("Adding key %s", hex.EncodeToString(key)[0:5])
	nDB.Lock()
	defer nDB.Unlock()
	for _, dbKey := range nDB.config.Keys ***REMOVED***
		if bytes.Equal(key, dbKey) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	nDB.config.Keys = append(nDB.config.Keys, key)
	if nDB.keyring != nil ***REMOVED***
		nDB.keyring.AddKey(key)
	***REMOVED***
***REMOVED***

// SetPrimaryKey sets the given key as the primary key. This should have
// been added apriori through SetKey
func (nDB *NetworkDB) SetPrimaryKey(key []byte) ***REMOVED***
	logrus.Debugf("Primary Key %s", hex.EncodeToString(key)[0:5])
	nDB.RLock()
	defer nDB.RUnlock()
	for _, dbKey := range nDB.config.Keys ***REMOVED***
		if bytes.Equal(key, dbKey) ***REMOVED***
			if nDB.keyring != nil ***REMOVED***
				nDB.keyring.UseKey(dbKey)
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// RemoveKey removes a key from the key ring. The key being removed
// can't be the primary key
func (nDB *NetworkDB) RemoveKey(key []byte) ***REMOVED***
	logrus.Debugf("Remove Key %s", hex.EncodeToString(key)[0:5])
	nDB.Lock()
	defer nDB.Unlock()
	for i, dbKey := range nDB.config.Keys ***REMOVED***
		if bytes.Equal(key, dbKey) ***REMOVED***
			nDB.config.Keys = append(nDB.config.Keys[:i], nDB.config.Keys[i+1:]...)
			if nDB.keyring != nil ***REMOVED***
				nDB.keyring.RemoveKey(dbKey)
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) clusterInit() error ***REMOVED***
	nDB.lastStatsTimestamp = time.Now()
	nDB.lastHealthTimestamp = nDB.lastStatsTimestamp

	config := memberlist.DefaultLANConfig()
	config.Name = nDB.config.NodeID
	config.BindAddr = nDB.config.BindAddr
	config.AdvertiseAddr = nDB.config.AdvertiseAddr
	config.UDPBufferSize = nDB.config.PacketBufferSize

	if nDB.config.BindPort != 0 ***REMOVED***
		config.BindPort = nDB.config.BindPort
	***REMOVED***

	config.ProtocolVersion = memberlist.ProtocolVersion2Compatible
	config.Delegate = &delegate***REMOVED***nDB: nDB***REMOVED***
	config.Events = &eventDelegate***REMOVED***nDB: nDB***REMOVED***
	// custom logger that does not add time or date, so they are not
	// duplicated by logrus
	config.Logger = log.New(&logWriter***REMOVED******REMOVED***, "", 0)

	var err error
	if len(nDB.config.Keys) > 0 ***REMOVED***
		for i, key := range nDB.config.Keys ***REMOVED***
			logrus.Debugf("Encryption key %d: %s", i+1, hex.EncodeToString(key)[0:5])
		***REMOVED***
		nDB.keyring, err = memberlist.NewKeyring(nDB.config.Keys, nDB.config.Keys[0])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config.Keyring = nDB.keyring
	***REMOVED***

	nDB.networkBroadcasts = &memberlist.TransmitLimitedQueue***REMOVED***
		NumNodes: func() int ***REMOVED***
			nDB.RLock()
			num := len(nDB.nodes)
			nDB.RUnlock()
			return num
		***REMOVED***,
		RetransmitMult: config.RetransmitMult,
	***REMOVED***

	nDB.nodeBroadcasts = &memberlist.TransmitLimitedQueue***REMOVED***
		NumNodes: func() int ***REMOVED***
			nDB.RLock()
			num := len(nDB.nodes)
			nDB.RUnlock()
			return num
		***REMOVED***,
		RetransmitMult: config.RetransmitMult,
	***REMOVED***

	mlist, err := memberlist.Create(config)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to create memberlist: %v", err)
	***REMOVED***

	nDB.stopCh = make(chan struct***REMOVED******REMOVED***)
	nDB.memberlist = mlist

	for _, trigger := range []struct ***REMOVED***
		interval time.Duration
		fn       func()
	***REMOVED******REMOVED***
		***REMOVED***reapPeriod, nDB.reapState***REMOVED***,
		***REMOVED***config.GossipInterval, nDB.gossip***REMOVED***,
		***REMOVED***config.PushPullInterval, nDB.bulkSyncTables***REMOVED***,
		***REMOVED***retryInterval, nDB.reconnectNode***REMOVED***,
		***REMOVED***nodeReapPeriod, nDB.reapDeadNode***REMOVED***,
	***REMOVED*** ***REMOVED***
		t := time.NewTicker(trigger.interval)
		go nDB.triggerFunc(trigger.interval, t.C, nDB.stopCh, trigger.fn)
		nDB.tickers = append(nDB.tickers, t)
	***REMOVED***

	return nil
***REMOVED***

func (nDB *NetworkDB) retryJoin(members []string, stop <-chan struct***REMOVED******REMOVED***) ***REMOVED***
	t := time.NewTicker(retryInterval)
	defer t.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case <-t.C:
			if _, err := nDB.memberlist.Join(members); err != nil ***REMOVED***
				logrus.Errorf("Failed to join memberlist %s on retry: %v", members, err)
				continue
			***REMOVED***
			if err := nDB.sendNodeEvent(NodeEventTypeJoin); err != nil ***REMOVED***
				logrus.Errorf("failed to send node join on retry: %v", err)
				continue
			***REMOVED***
			return
		case <-stop:
			return
		***REMOVED***
	***REMOVED***

***REMOVED***

func (nDB *NetworkDB) clusterJoin(members []string) error ***REMOVED***
	mlist := nDB.memberlist

	if _, err := mlist.Join(members); err != nil ***REMOVED***
		// In case of failure, keep retrying join until it succeeds or the cluster is shutdown.
		go nDB.retryJoin(members, nDB.stopCh)
		return fmt.Errorf("could not join node to memberlist: %v", err)
	***REMOVED***

	if err := nDB.sendNodeEvent(NodeEventTypeJoin); err != nil ***REMOVED***
		return fmt.Errorf("failed to send node join: %v", err)
	***REMOVED***

	return nil
***REMOVED***

func (nDB *NetworkDB) clusterLeave() error ***REMOVED***
	mlist := nDB.memberlist

	if err := nDB.sendNodeEvent(NodeEventTypeLeave); err != nil ***REMOVED***
		logrus.Errorf("failed to send node leave: %v", err)
	***REMOVED***

	if err := mlist.Leave(time.Second); err != nil ***REMOVED***
		return err
	***REMOVED***

	close(nDB.stopCh)

	for _, t := range nDB.tickers ***REMOVED***
		t.Stop()
	***REMOVED***

	return mlist.Shutdown()
***REMOVED***

func (nDB *NetworkDB) triggerFunc(stagger time.Duration, C <-chan time.Time, stop <-chan struct***REMOVED******REMOVED***, f func()) ***REMOVED***
	// Use a random stagger to avoid syncronizing
	randStagger := time.Duration(uint64(rnd.Int63()) % uint64(stagger))
	select ***REMOVED***
	case <-time.After(randStagger):
	case <-stop:
		return
	***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-C:
			f()
		case <-stop:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) reapDeadNode() ***REMOVED***
	nDB.Lock()
	defer nDB.Unlock()
	for _, nodeMap := range []map[string]*node***REMOVED***
		nDB.failedNodes,
		nDB.leftNodes,
	***REMOVED*** ***REMOVED***
		for id, n := range nodeMap ***REMOVED***
			if n.reapTime > nodeReapPeriod ***REMOVED***
				n.reapTime -= nodeReapPeriod
				continue
			***REMOVED***
			logrus.Debugf("Garbage collect node %v", n.Name)
			delete(nodeMap, id)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) reconnectNode() ***REMOVED***
	nDB.RLock()
	if len(nDB.failedNodes) == 0 ***REMOVED***
		nDB.RUnlock()
		return
	***REMOVED***

	nodes := make([]*node, 0, len(nDB.failedNodes))
	for _, n := range nDB.failedNodes ***REMOVED***
		nodes = append(nodes, n)
	***REMOVED***
	nDB.RUnlock()

	node := nodes[randomOffset(len(nodes))]
	addr := net.UDPAddr***REMOVED***IP: node.Addr, Port: int(node.Port)***REMOVED***

	if _, err := nDB.memberlist.Join([]string***REMOVED***addr.String()***REMOVED***); err != nil ***REMOVED***
		return
	***REMOVED***

	if err := nDB.sendNodeEvent(NodeEventTypeJoin); err != nil ***REMOVED***
		return
	***REMOVED***

	logrus.Debugf("Initiating bulk sync with node %s after reconnect", node.Name)
	nDB.bulkSync([]string***REMOVED***node.Name***REMOVED***, true)
***REMOVED***

// For timing the entry deletion in the repaer APIs that doesn't use monotonic clock
// source (time.Now, Sub etc.) should be avoided. Hence we use reapTime in every
// entry which is set initially to reapInterval and decremented by reapPeriod every time
// the reaper runs. NOTE nDB.reapTableEntries updates the reapTime with a readlock. This
// is safe as long as no other concurrent path touches the reapTime field.
func (nDB *NetworkDB) reapState() ***REMOVED***
	// The reapTableEntries leverage the presence of the network so garbage collect entries first
	nDB.reapTableEntries()
	nDB.reapNetworks()
***REMOVED***

func (nDB *NetworkDB) reapNetworks() ***REMOVED***
	nDB.Lock()
	for _, nn := range nDB.networks ***REMOVED***
		for id, n := range nn ***REMOVED***
			if n.leaving ***REMOVED***
				if n.reapTime <= 0 ***REMOVED***
					delete(nn, id)
					continue
				***REMOVED***
				n.reapTime -= reapPeriod
			***REMOVED***
		***REMOVED***
	***REMOVED***
	nDB.Unlock()
***REMOVED***

func (nDB *NetworkDB) reapTableEntries() ***REMOVED***
	var nodeNetworks []string
	// This is best effort, if the list of network changes will be picked up in the next cycle
	nDB.RLock()
	for nid := range nDB.networks[nDB.config.NodeID] ***REMOVED***
		nodeNetworks = append(nodeNetworks, nid)
	***REMOVED***
	nDB.RUnlock()

	cycleStart := time.Now()
	// In order to avoid blocking the database for a long time, apply the garbage collection logic by network
	// The lock is taken at the beginning of the cycle and the deletion is inline
	for _, nid := range nodeNetworks ***REMOVED***
		nDB.Lock()
		nDB.indexes[byNetwork].WalkPrefix(fmt.Sprintf("/%s", nid), func(path string, v interface***REMOVED******REMOVED***) bool ***REMOVED***
			// timeCompensation compensate in case the lock took some time to be released
			timeCompensation := time.Since(cycleStart)
			entry, ok := v.(*entry)
			if !ok || !entry.deleting ***REMOVED***
				return false
			***REMOVED***

			// In this check we are adding an extra 1 second to guarantee that when the number is truncated to int32 to fit the packet
			// for the tableEvent the number is always strictly > 1 and never 0
			if entry.reapTime > reapPeriod+timeCompensation+time.Second ***REMOVED***
				entry.reapTime -= reapPeriod + timeCompensation
				return false
			***REMOVED***

			params := strings.Split(path[1:], "/")
			nid := params[0]
			tname := params[1]
			key := params[2]

			okTable, okNetwork := nDB.deleteEntry(nid, tname, key)
			if !okTable ***REMOVED***
				logrus.Errorf("Table tree delete failed, entry with key:%s does not exists in the table:%s network:%s", key, tname, nid)
			***REMOVED***
			if !okNetwork ***REMOVED***
				logrus.Errorf("Network tree delete failed, entry with key:%s does not exists in the network:%s table:%s", key, nid, tname)
			***REMOVED***

			return false
		***REMOVED***)
		nDB.Unlock()
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) gossip() ***REMOVED***
	networkNodes := make(map[string][]string)
	nDB.RLock()
	thisNodeNetworks := nDB.networks[nDB.config.NodeID]
	for nid := range thisNodeNetworks ***REMOVED***
		networkNodes[nid] = nDB.networkNodes[nid]
	***REMOVED***
	printStats := time.Since(nDB.lastStatsTimestamp) >= nDB.config.StatsPrintPeriod
	printHealth := time.Since(nDB.lastHealthTimestamp) >= nDB.config.HealthPrintPeriod
	nDB.RUnlock()

	if printHealth ***REMOVED***
		healthScore := nDB.memberlist.GetHealthScore()
		if healthScore != 0 ***REMOVED***
			logrus.Warnf("NetworkDB stats %v(%v) - healthscore:%d (connectivity issues)", nDB.config.Hostname, nDB.config.NodeID, healthScore)
		***REMOVED***
		nDB.lastHealthTimestamp = time.Now()
	***REMOVED***

	for nid, nodes := range networkNodes ***REMOVED***
		mNodes := nDB.mRandomNodes(3, nodes)
		bytesAvail := nDB.config.PacketBufferSize - compoundHeaderOverhead

		nDB.RLock()
		network, ok := thisNodeNetworks[nid]
		nDB.RUnlock()
		if !ok || network == nil ***REMOVED***
			// It is normal for the network to be removed
			// between the time we collect the network
			// attachments of this node and processing
			// them here.
			continue
		***REMOVED***

		broadcastQ := network.tableBroadcasts

		if broadcastQ == nil ***REMOVED***
			logrus.Errorf("Invalid broadcastQ encountered while gossiping for network %s", nid)
			continue
		***REMOVED***

		msgs := broadcastQ.GetBroadcasts(compoundOverhead, bytesAvail)
		// Collect stats and print the queue info, note this code is here also to have a view of the queues empty
		network.qMessagesSent += len(msgs)
		if printStats ***REMOVED***
			logrus.Infof("NetworkDB stats %v(%v) - netID:%s leaving:%t netPeers:%d entries:%d Queue qLen:%d netMsg/s:%d",
				nDB.config.Hostname, nDB.config.NodeID,
				nid, network.leaving, broadcastQ.NumNodes(), network.entriesNumber, broadcastQ.NumQueued(),
				network.qMessagesSent/int((nDB.config.StatsPrintPeriod/time.Second)))
			network.qMessagesSent = 0
		***REMOVED***

		if len(msgs) == 0 ***REMOVED***
			continue
		***REMOVED***

		// Create a compound message
		compound := makeCompoundMessage(msgs)

		for _, node := range mNodes ***REMOVED***
			nDB.RLock()
			mnode := nDB.nodes[node]
			nDB.RUnlock()

			if mnode == nil ***REMOVED***
				break
			***REMOVED***

			// Send the compound message
			if err := nDB.memberlist.SendBestEffort(&mnode.Node, compound); err != nil ***REMOVED***
				logrus.Errorf("Failed to send gossip to %s: %s", mnode.Addr, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Reset the stats
	if printStats ***REMOVED***
		nDB.lastStatsTimestamp = time.Now()
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) bulkSyncTables() ***REMOVED***
	var networks []string
	nDB.RLock()
	for nid, network := range nDB.networks[nDB.config.NodeID] ***REMOVED***
		if network.leaving ***REMOVED***
			continue
		***REMOVED***
		networks = append(networks, nid)
	***REMOVED***
	nDB.RUnlock()

	for ***REMOVED***
		if len(networks) == 0 ***REMOVED***
			break
		***REMOVED***

		nid := networks[0]
		networks = networks[1:]

		nDB.RLock()
		nodes := nDB.networkNodes[nid]
		nDB.RUnlock()

		// No peer nodes on this network. Move on.
		if len(nodes) == 0 ***REMOVED***
			continue
		***REMOVED***

		completed, err := nDB.bulkSync(nodes, false)
		if err != nil ***REMOVED***
			logrus.Errorf("periodic bulk sync failure for network %s: %v", nid, err)
			continue
		***REMOVED***

		// Remove all the networks for which we have
		// successfully completed bulk sync in this iteration.
		updatedNetworks := make([]string, 0, len(networks))
		for _, nid := range networks ***REMOVED***
			var found bool
			for _, completedNid := range completed ***REMOVED***
				if nid == completedNid ***REMOVED***
					found = true
					break
				***REMOVED***
			***REMOVED***

			if !found ***REMOVED***
				updatedNetworks = append(updatedNetworks, nid)
			***REMOVED***
		***REMOVED***

		networks = updatedNetworks
	***REMOVED***
***REMOVED***

func (nDB *NetworkDB) bulkSync(nodes []string, all bool) ([]string, error) ***REMOVED***
	if !all ***REMOVED***
		// Get 2 random nodes. 2nd node will be tried if the bulk sync to
		// 1st node fails.
		nodes = nDB.mRandomNodes(2, nodes)
	***REMOVED***

	if len(nodes) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	var err error
	var networks []string
	for _, node := range nodes ***REMOVED***
		if node == nDB.config.NodeID ***REMOVED***
			continue
		***REMOVED***
		logrus.Debugf("%v(%v): Initiating bulk sync with node %v", nDB.config.Hostname, nDB.config.NodeID, node)
		networks = nDB.findCommonNetworks(node)
		err = nDB.bulkSyncNode(networks, node, true)
		// if its periodic bulksync stop after the first successful sync
		if !all && err == nil ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			err = fmt.Errorf("bulk sync to node %s failed: %v", node, err)
			logrus.Warn(err.Error())
		***REMOVED***
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return networks, nil
***REMOVED***

// Bulk sync all the table entries belonging to a set of networks to a
// single peer node. It can be unsolicited or can be in response to an
// unsolicited bulk sync
func (nDB *NetworkDB) bulkSyncNode(networks []string, node string, unsolicited bool) error ***REMOVED***
	var msgs [][]byte

	var unsolMsg string
	if unsolicited ***REMOVED***
		unsolMsg = "unsolicited"
	***REMOVED***

	logrus.Debugf("%v(%v): Initiating %s bulk sync for networks %v with node %s",
		nDB.config.Hostname, nDB.config.NodeID, unsolMsg, networks, node)

	nDB.RLock()
	mnode := nDB.nodes[node]
	if mnode == nil ***REMOVED***
		nDB.RUnlock()
		return nil
	***REMOVED***

	for _, nid := range networks ***REMOVED***
		nDB.indexes[byNetwork].WalkPrefix(fmt.Sprintf("/%s", nid), func(path string, v interface***REMOVED******REMOVED***) bool ***REMOVED***
			entry, ok := v.(*entry)
			if !ok ***REMOVED***
				return false
			***REMOVED***

			eType := TableEventTypeCreate
			if entry.deleting ***REMOVED***
				eType = TableEventTypeDelete
			***REMOVED***

			params := strings.Split(path[1:], "/")
			tEvent := TableEvent***REMOVED***
				Type:      eType,
				LTime:     entry.ltime,
				NodeName:  entry.node,
				NetworkID: nid,
				TableName: params[1],
				Key:       params[2],
				Value:     entry.value,
				// The duration in second is a float that below would be truncated
				ResidualReapTime: int32(entry.reapTime.Seconds()),
			***REMOVED***

			msg, err := encodeMessage(MessageTypeTableEvent, &tEvent)
			if err != nil ***REMOVED***
				logrus.Errorf("Encode failure during bulk sync: %#v", tEvent)
				return false
			***REMOVED***

			msgs = append(msgs, msg)
			return false
		***REMOVED***)
	***REMOVED***
	nDB.RUnlock()

	// Create a compound message
	compound := makeCompoundMessage(msgs)

	bsm := BulkSyncMessage***REMOVED***
		LTime:       nDB.tableClock.Time(),
		Unsolicited: unsolicited,
		NodeName:    nDB.config.NodeID,
		Networks:    networks,
		Payload:     compound,
	***REMOVED***

	buf, err := encodeMessage(MessageTypeBulkSync, &bsm)
	if err != nil ***REMOVED***
		return fmt.Errorf("failed to encode bulk sync message: %v", err)
	***REMOVED***

	nDB.Lock()
	ch := make(chan struct***REMOVED******REMOVED***)
	nDB.bulkSyncAckTbl[node] = ch
	nDB.Unlock()

	err = nDB.memberlist.SendReliable(&mnode.Node, buf)
	if err != nil ***REMOVED***
		nDB.Lock()
		delete(nDB.bulkSyncAckTbl, node)
		nDB.Unlock()

		return fmt.Errorf("failed to send a TCP message during bulk sync: %v", err)
	***REMOVED***

	// Wait on a response only if it is unsolicited.
	if unsolicited ***REMOVED***
		startTime := time.Now()
		t := time.NewTimer(30 * time.Second)
		select ***REMOVED***
		case <-t.C:
			logrus.Errorf("Bulk sync to node %s timed out", node)
		case <-ch:
			logrus.Debugf("%v(%v): Bulk sync to node %s took %s", nDB.config.Hostname, nDB.config.NodeID, node, time.Since(startTime))
		***REMOVED***
		t.Stop()
	***REMOVED***

	return nil
***REMOVED***

// Returns a random offset between 0 and n
func randomOffset(n int) int ***REMOVED***
	if n == 0 ***REMOVED***
		return 0
	***REMOVED***

	val, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil ***REMOVED***
		logrus.Errorf("Failed to get a random offset: %v", err)
		return 0
	***REMOVED***

	return int(val.Int64())
***REMOVED***

// mRandomNodes is used to select up to m random nodes. It is possible
// that less than m nodes are returned.
func (nDB *NetworkDB) mRandomNodes(m int, nodes []string) []string ***REMOVED***
	n := len(nodes)
	mNodes := make([]string, 0, m)
OUTER:
	// Probe up to 3*n times, with large n this is not necessary
	// since k << n, but with small n we want search to be
	// exhaustive
	for i := 0; i < 3*n && len(mNodes) < m; i++ ***REMOVED***
		// Get random node
		idx := randomOffset(n)
		node := nodes[idx]

		if node == nDB.config.NodeID ***REMOVED***
			continue
		***REMOVED***

		// Check if we have this node already
		for j := 0; j < len(mNodes); j++ ***REMOVED***
			if node == mNodes[j] ***REMOVED***
				continue OUTER
			***REMOVED***
		***REMOVED***

		// Append the node
		mNodes = append(mNodes, node)
	***REMOVED***

	return mNodes
***REMOVED***
