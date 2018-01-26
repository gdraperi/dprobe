package memberlist

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
)

type nodeStateType int

const (
	stateAlive nodeStateType = iota
	stateSuspect
	stateDead
)

// Node represents a node in the cluster.
type Node struct ***REMOVED***
	Name string
	Addr net.IP
	Port uint16
	Meta []byte // Metadata from the delegate for this node.
	PMin uint8  // Minimum protocol version this understands
	PMax uint8  // Maximum protocol version this understands
	PCur uint8  // Current version node is speaking
	DMin uint8  // Min protocol version for the delegate to understand
	DMax uint8  // Max protocol version for the delegate to understand
	DCur uint8  // Current version delegate is speaking
***REMOVED***

// Address returns the host:port form of a node's address, suitable for use
// with a transport.
func (n *Node) Address() string ***REMOVED***
	return joinHostPort(n.Addr.String(), n.Port)
***REMOVED***

// NodeState is used to manage our state view of another node
type nodeState struct ***REMOVED***
	Node
	Incarnation uint32        // Last known incarnation number
	State       nodeStateType // Current state
	StateChange time.Time     // Time last state change happened
***REMOVED***

// Address returns the host:port form of a node's address, suitable for use
// with a transport.
func (n *nodeState) Address() string ***REMOVED***
	return n.Node.Address()
***REMOVED***

// ackHandler is used to register handlers for incoming acks and nacks.
type ackHandler struct ***REMOVED***
	ackFn  func([]byte, time.Time)
	nackFn func()
	timer  *time.Timer
***REMOVED***

// NoPingResponseError is used to indicate a 'ping' packet was
// successfully issued but no response was received
type NoPingResponseError struct ***REMOVED***
	node string
***REMOVED***

func (f NoPingResponseError) Error() string ***REMOVED***
	return fmt.Sprintf("No response from node %s", f.node)
***REMOVED***

// Schedule is used to ensure the Tick is performed periodically. This
// function is safe to call multiple times. If the memberlist is already
// scheduled, then it won't do anything.
func (m *Memberlist) schedule() ***REMOVED***
	m.tickerLock.Lock()
	defer m.tickerLock.Unlock()

	// If we already have tickers, then don't do anything, since we're
	// scheduled
	if len(m.tickers) > 0 ***REMOVED***
		return
	***REMOVED***

	// Create the stop tick channel, a blocking channel. We close this
	// when we should stop the tickers.
	stopCh := make(chan struct***REMOVED******REMOVED***)

	// Create a new probeTicker
	if m.config.ProbeInterval > 0 ***REMOVED***
		t := time.NewTicker(m.config.ProbeInterval)
		go m.triggerFunc(m.config.ProbeInterval, t.C, stopCh, m.probe)
		m.tickers = append(m.tickers, t)
	***REMOVED***

	// Create a push pull ticker if needed
	if m.config.PushPullInterval > 0 ***REMOVED***
		go m.pushPullTrigger(stopCh)
	***REMOVED***

	// Create a gossip ticker if needed
	if m.config.GossipInterval > 0 && m.config.GossipNodes > 0 ***REMOVED***
		t := time.NewTicker(m.config.GossipInterval)
		go m.triggerFunc(m.config.GossipInterval, t.C, stopCh, m.gossip)
		m.tickers = append(m.tickers, t)
	***REMOVED***

	// If we made any tickers, then record the stopTick channel for
	// later.
	if len(m.tickers) > 0 ***REMOVED***
		m.stopTick = stopCh
	***REMOVED***
***REMOVED***

// triggerFunc is used to trigger a function call each time a
// message is received until a stop tick arrives.
func (m *Memberlist) triggerFunc(stagger time.Duration, C <-chan time.Time, stop <-chan struct***REMOVED******REMOVED***, f func()) ***REMOVED***
	// Use a random stagger to avoid syncronizing
	randStagger := time.Duration(uint64(rand.Int63()) % uint64(stagger))
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

// pushPullTrigger is used to periodically trigger a push/pull until
// a stop tick arrives. We don't use triggerFunc since the push/pull
// timer is dynamically scaled based on cluster size to avoid network
// saturation
func (m *Memberlist) pushPullTrigger(stop <-chan struct***REMOVED******REMOVED***) ***REMOVED***
	interval := m.config.PushPullInterval

	// Use a random stagger to avoid syncronizing
	randStagger := time.Duration(uint64(rand.Int63()) % uint64(interval))
	select ***REMOVED***
	case <-time.After(randStagger):
	case <-stop:
		return
	***REMOVED***

	// Tick using a dynamic timer
	for ***REMOVED***
		tickTime := pushPullScale(interval, m.estNumNodes())
		select ***REMOVED***
		case <-time.After(tickTime):
			m.pushPull()
		case <-stop:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Deschedule is used to stop the background maintenance. This is safe
// to call multiple times.
func (m *Memberlist) deschedule() ***REMOVED***
	m.tickerLock.Lock()
	defer m.tickerLock.Unlock()

	// If we have no tickers, then we aren't scheduled.
	if len(m.tickers) == 0 ***REMOVED***
		return
	***REMOVED***

	// Close the stop channel so all the ticker listeners stop.
	close(m.stopTick)

	// Explicitly stop all the tickers themselves so they don't take
	// up any more resources, and get rid of the list.
	for _, t := range m.tickers ***REMOVED***
		t.Stop()
	***REMOVED***
	m.tickers = nil
***REMOVED***

// Tick is used to perform a single round of failure detection and gossip
func (m *Memberlist) probe() ***REMOVED***
	// Track the number of indexes we've considered probing
	numCheck := 0
START:
	m.nodeLock.RLock()

	// Make sure we don't wrap around infinitely
	if numCheck >= len(m.nodes) ***REMOVED***
		m.nodeLock.RUnlock()
		return
	***REMOVED***

	// Handle the wrap around case
	if m.probeIndex >= len(m.nodes) ***REMOVED***
		m.nodeLock.RUnlock()
		m.resetNodes()
		m.probeIndex = 0
		numCheck++
		goto START
	***REMOVED***

	// Determine if we should probe this node
	skip := false
	var node nodeState

	node = *m.nodes[m.probeIndex]
	if node.Name == m.config.Name ***REMOVED***
		skip = true
	***REMOVED*** else if node.State == stateDead ***REMOVED***
		skip = true
	***REMOVED***

	// Potentially skip
	m.nodeLock.RUnlock()
	m.probeIndex++
	if skip ***REMOVED***
		numCheck++
		goto START
	***REMOVED***

	// Probe the specific node
	m.probeNode(&node)
***REMOVED***

// probeNode handles a single round of failure checking on a node.
func (m *Memberlist) probeNode(node *nodeState) ***REMOVED***
	defer metrics.MeasureSince([]string***REMOVED***"memberlist", "probeNode"***REMOVED***, time.Now())

	// We use our health awareness to scale the overall probe interval, so we
	// slow down if we detect problems. The ticker that calls us can handle
	// us running over the base interval, and will skip missed ticks.
	probeInterval := m.awareness.ScaleTimeout(m.config.ProbeInterval)
	if probeInterval > m.config.ProbeInterval ***REMOVED***
		metrics.IncrCounter([]string***REMOVED***"memberlist", "degraded", "probe"***REMOVED***, 1)
	***REMOVED***

	// Prepare a ping message and setup an ack handler.
	ping := ping***REMOVED***SeqNo: m.nextSeqNo(), Node: node.Name***REMOVED***
	ackCh := make(chan ackMessage, m.config.IndirectChecks+1)
	nackCh := make(chan struct***REMOVED******REMOVED***, m.config.IndirectChecks+1)
	m.setProbeChannels(ping.SeqNo, ackCh, nackCh, probeInterval)

	// Send a ping to the node. If this node looks like it's suspect or dead,
	// also tack on a suspect message so that it has a chance to refute as
	// soon as possible.
	deadline := time.Now().Add(probeInterval)
	addr := node.Address()
	if node.State == stateAlive ***REMOVED***
		if err := m.encodeAndSendMsg(addr, pingMsg, &ping); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to send ping: %s", err)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var msgs [][]byte
		if buf, err := encode(pingMsg, &ping); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to encode ping message: %s", err)
			return
		***REMOVED*** else ***REMOVED***
			msgs = append(msgs, buf.Bytes())
		***REMOVED***
		s := suspect***REMOVED***Incarnation: node.Incarnation, Node: node.Name, From: m.config.Name***REMOVED***
		if buf, err := encode(suspectMsg, &s); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to encode suspect message: %s", err)
			return
		***REMOVED*** else ***REMOVED***
			msgs = append(msgs, buf.Bytes())
		***REMOVED***

		compound := makeCompoundMessage(msgs)
		if err := m.rawSendMsgPacket(addr, &node.Node, compound.Bytes()); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to send compound ping and suspect message to %s: %s", addr, err)
			return
		***REMOVED***
	***REMOVED***

	// Mark the sent time here, which should be after any pre-processing and
	// system calls to do the actual send. This probably under-reports a bit,
	// but it's the best we can do.
	sent := time.Now()

	// Arrange for our self-awareness to get updated. At this point we've
	// sent the ping, so any return statement means the probe succeeded
	// which will improve our health until we get to the failure scenarios
	// at the end of this function, which will alter this delta variable
	// accordingly.
	awarenessDelta := -1
	defer func() ***REMOVED***
		m.awareness.ApplyDelta(awarenessDelta)
	***REMOVED***()

	// Wait for response or round-trip-time.
	select ***REMOVED***
	case v := <-ackCh:
		if v.Complete == true ***REMOVED***
			if m.config.Ping != nil ***REMOVED***
				rtt := v.Timestamp.Sub(sent)
				m.config.Ping.NotifyPingComplete(&node.Node, rtt, v.Payload)
			***REMOVED***
			return
		***REMOVED***

		// As an edge case, if we get a timeout, we need to re-enqueue it
		// here to break out of the select below.
		if v.Complete == false ***REMOVED***
			ackCh <- v
		***REMOVED***
	case <-time.After(m.config.ProbeTimeout):
		// Note that we don't scale this timeout based on awareness and
		// the health score. That's because we don't really expect waiting
		// longer to help get UDP through. Since health does extend the
		// probe interval it will give the TCP fallback more time, which
		// is more active in dealing with lost packets, and it gives more
		// time to wait for indirect acks/nacks.
		m.logger.Printf("[DEBUG] memberlist: Failed ping: %v (timeout reached)", node.Name)
	***REMOVED***

	// Get some random live nodes.
	m.nodeLock.RLock()
	kNodes := kRandomNodes(m.config.IndirectChecks, m.nodes, func(n *nodeState) bool ***REMOVED***
		return n.Name == m.config.Name ||
			n.Name == node.Name ||
			n.State != stateAlive
	***REMOVED***)
	m.nodeLock.RUnlock()

	// Attempt an indirect ping.
	expectedNacks := 0
	ind := indirectPingReq***REMOVED***SeqNo: ping.SeqNo, Target: node.Addr, Port: node.Port, Node: node.Name***REMOVED***
	for _, peer := range kNodes ***REMOVED***
		// We only expect nack to be sent from peers who understand
		// version 4 of the protocol.
		if ind.Nack = peer.PMax >= 4; ind.Nack ***REMOVED***
			expectedNacks++
		***REMOVED***

		if err := m.encodeAndSendMsg(peer.Address(), indirectPingMsg, &ind); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to send indirect ping: %s", err)
		***REMOVED***
	***REMOVED***

	// Also make an attempt to contact the node directly over TCP. This
	// helps prevent confused clients who get isolated from UDP traffic
	// but can still speak TCP (which also means they can possibly report
	// misinformation to other nodes via anti-entropy), avoiding flapping in
	// the cluster.
	//
	// This is a little unusual because we will attempt a TCP ping to any
	// member who understands version 3 of the protocol, regardless of
	// which protocol version we are speaking. That's why we've included a
	// config option to turn this off if desired.
	fallbackCh := make(chan bool, 1)
	if (!m.config.DisableTcpPings) && (node.PMax >= 3) ***REMOVED***
		go func() ***REMOVED***
			defer close(fallbackCh)
			didContact, err := m.sendPingAndWaitForAck(node.Address(), ping, deadline)
			if err != nil ***REMOVED***
				m.logger.Printf("[ERR] memberlist: Failed fallback ping: %s", err)
			***REMOVED*** else ***REMOVED***
				fallbackCh <- didContact
			***REMOVED***
		***REMOVED***()
	***REMOVED*** else ***REMOVED***
		close(fallbackCh)
	***REMOVED***

	// Wait for the acks or timeout. Note that we don't check the fallback
	// channel here because we want to issue a warning below if that's the
	// *only* way we hear back from the peer, so we have to let this time
	// out first to allow the normal UDP-based acks to come in.
	select ***REMOVED***
	case v := <-ackCh:
		if v.Complete == true ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	// Finally, poll the fallback channel. The timeouts are set such that
	// the channel will have something or be closed without having to wait
	// any additional time here.
	for didContact := range fallbackCh ***REMOVED***
		if didContact ***REMOVED***
			m.logger.Printf("[WARN] memberlist: Was able to connect to %s but other probes failed, network may be misconfigured", node.Name)
			return
		***REMOVED***
	***REMOVED***

	// Update our self-awareness based on the results of this failed probe.
	// If we don't have peers who will send nacks then we penalize for any
	// failed probe as a simple health metric. If we do have peers to nack
	// verify, then we can use that as a more sophisticated measure of self-
	// health because we assume them to be working, and they can help us
	// decide if the probed node was really dead or if it was something wrong
	// with ourselves.
	awarenessDelta = 0
	if expectedNacks > 0 ***REMOVED***
		if nackCount := len(nackCh); nackCount < expectedNacks ***REMOVED***
			awarenessDelta += (expectedNacks - nackCount)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		awarenessDelta += 1
	***REMOVED***

	// No acks received from target, suspect it as failed.
	m.logger.Printf("[INFO] memberlist: Suspect %s has failed, no acks received", node.Name)
	s := suspect***REMOVED***Incarnation: node.Incarnation, Node: node.Name, From: m.config.Name***REMOVED***
	m.suspectNode(&s)
***REMOVED***

// Ping initiates a ping to the node with the specified name.
func (m *Memberlist) Ping(node string, addr net.Addr) (time.Duration, error) ***REMOVED***
	// Prepare a ping message and setup an ack handler.
	ping := ping***REMOVED***SeqNo: m.nextSeqNo(), Node: node***REMOVED***
	ackCh := make(chan ackMessage, m.config.IndirectChecks+1)
	m.setProbeChannels(ping.SeqNo, ackCh, nil, m.config.ProbeInterval)

	// Send a ping to the node.
	if err := m.encodeAndSendMsg(addr.String(), pingMsg, &ping); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	// Mark the sent time here, which should be after any pre-processing and
	// system calls to do the actual send. This probably under-reports a bit,
	// but it's the best we can do.
	sent := time.Now()

	// Wait for response or timeout.
	select ***REMOVED***
	case v := <-ackCh:
		if v.Complete == true ***REMOVED***
			return v.Timestamp.Sub(sent), nil
		***REMOVED***
	case <-time.After(m.config.ProbeTimeout):
		// Timeout, return an error below.
	***REMOVED***

	m.logger.Printf("[DEBUG] memberlist: Failed UDP ping: %v (timeout reached)", node)
	return 0, NoPingResponseError***REMOVED***ping.Node***REMOVED***
***REMOVED***

// resetNodes is used when the tick wraps around. It will reap the
// dead nodes and shuffle the node list.
func (m *Memberlist) resetNodes() ***REMOVED***
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()

	// Move dead nodes, but respect gossip to the dead interval
	deadIdx := moveDeadNodes(m.nodes, m.config.GossipToTheDeadTime)

	// Deregister the dead nodes
	for i := deadIdx; i < len(m.nodes); i++ ***REMOVED***
		delete(m.nodeMap, m.nodes[i].Name)
		m.nodes[i] = nil
	***REMOVED***

	// Trim the nodes to exclude the dead nodes
	m.nodes = m.nodes[0:deadIdx]

	// Update numNodes after we've trimmed the dead nodes
	atomic.StoreUint32(&m.numNodes, uint32(deadIdx))

	// Shuffle live nodes
	shuffleNodes(m.nodes)
***REMOVED***

// gossip is invoked every GossipInterval period to broadcast our gossip
// messages to a few random nodes.
func (m *Memberlist) gossip() ***REMOVED***
	defer metrics.MeasureSince([]string***REMOVED***"memberlist", "gossip"***REMOVED***, time.Now())

	// Get some random live, suspect, or recently dead nodes
	m.nodeLock.RLock()
	kNodes := kRandomNodes(m.config.GossipNodes, m.nodes, func(n *nodeState) bool ***REMOVED***
		if n.Name == m.config.Name ***REMOVED***
			return true
		***REMOVED***

		switch n.State ***REMOVED***
		case stateAlive, stateSuspect:
			return false

		case stateDead:
			return time.Since(n.StateChange) > m.config.GossipToTheDeadTime

		default:
			return true
		***REMOVED***
	***REMOVED***)
	m.nodeLock.RUnlock()

	// Compute the bytes available
	bytesAvail := m.config.UDPBufferSize - compoundHeaderOverhead
	if m.config.EncryptionEnabled() ***REMOVED***
		bytesAvail -= encryptOverhead(m.encryptionVersion())
	***REMOVED***

	for _, node := range kNodes ***REMOVED***
		// Get any pending broadcasts
		msgs := m.getBroadcasts(compoundOverhead, bytesAvail)
		if len(msgs) == 0 ***REMOVED***
			return
		***REMOVED***

		addr := node.Address()
		if len(msgs) == 1 ***REMOVED***
			// Send single message as is
			if err := m.rawSendMsgPacket(addr, &node.Node, msgs[0]); err != nil ***REMOVED***
				m.logger.Printf("[ERR] memberlist: Failed to send gossip to %s: %s", addr, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Otherwise create and send a compound message
			compound := makeCompoundMessage(msgs)
			if err := m.rawSendMsgPacket(addr, &node.Node, compound.Bytes()); err != nil ***REMOVED***
				m.logger.Printf("[ERR] memberlist: Failed to send gossip to %s: %s", addr, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// pushPull is invoked periodically to randomly perform a complete state
// exchange. Used to ensure a high level of convergence, but is also
// reasonably expensive as the entire state of this node is exchanged
// with the other node.
func (m *Memberlist) pushPull() ***REMOVED***
	// Get a random live node
	m.nodeLock.RLock()
	nodes := kRandomNodes(1, m.nodes, func(n *nodeState) bool ***REMOVED***
		return n.Name == m.config.Name ||
			n.State != stateAlive
	***REMOVED***)
	m.nodeLock.RUnlock()

	// If no nodes, bail
	if len(nodes) == 0 ***REMOVED***
		return
	***REMOVED***
	node := nodes[0]

	// Attempt a push pull
	if err := m.pushPullNode(node.Address(), false); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Push/Pull with %s failed: %s", node.Name, err)
	***REMOVED***
***REMOVED***

// pushPullNode does a complete state exchange with a specific node.
func (m *Memberlist) pushPullNode(addr string, join bool) error ***REMOVED***
	defer metrics.MeasureSince([]string***REMOVED***"memberlist", "pushPullNode"***REMOVED***, time.Now())

	// Attempt to send and receive with the node
	remote, userState, err := m.sendAndReceiveState(addr, join)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := m.mergeRemoteState(join, remote, userState); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// verifyProtocol verifies that all the remote nodes can speak with our
// nodes and vice versa on both the core protocol as well as the
// delegate protocol level.
//
// The verification works by finding the maximum minimum and
// minimum maximum understood protocol and delegate versions. In other words,
// it finds the common denominator of protocol and delegate version ranges
// for the entire cluster.
//
// After this, it goes through the entire cluster (local and remote) and
// verifies that everyone's speaking protocol versions satisfy this range.
// If this passes, it means that every node can understand each other.
func (m *Memberlist) verifyProtocol(remote []pushNodeState) error ***REMOVED***
	m.nodeLock.RLock()
	defer m.nodeLock.RUnlock()

	// Maximum minimum understood and minimum maximum understood for both
	// the protocol and delegate versions. We use this to verify everyone
	// can be understood.
	var maxpmin, minpmax uint8
	var maxdmin, mindmax uint8
	minpmax = math.MaxUint8
	mindmax = math.MaxUint8

	for _, rn := range remote ***REMOVED***
		// If the node isn't alive, then skip it
		if rn.State != stateAlive ***REMOVED***
			continue
		***REMOVED***

		// Skip nodes that don't have versions set, it just means
		// their version is zero.
		if len(rn.Vsn) == 0 ***REMOVED***
			continue
		***REMOVED***

		if rn.Vsn[0] > maxpmin ***REMOVED***
			maxpmin = rn.Vsn[0]
		***REMOVED***

		if rn.Vsn[1] < minpmax ***REMOVED***
			minpmax = rn.Vsn[1]
		***REMOVED***

		if rn.Vsn[3] > maxdmin ***REMOVED***
			maxdmin = rn.Vsn[3]
		***REMOVED***

		if rn.Vsn[4] < mindmax ***REMOVED***
			mindmax = rn.Vsn[4]
		***REMOVED***
	***REMOVED***

	for _, n := range m.nodes ***REMOVED***
		// Ignore non-alive nodes
		if n.State != stateAlive ***REMOVED***
			continue
		***REMOVED***

		if n.PMin > maxpmin ***REMOVED***
			maxpmin = n.PMin
		***REMOVED***

		if n.PMax < minpmax ***REMOVED***
			minpmax = n.PMax
		***REMOVED***

		if n.DMin > maxdmin ***REMOVED***
			maxdmin = n.DMin
		***REMOVED***

		if n.DMax < mindmax ***REMOVED***
			mindmax = n.DMax
		***REMOVED***
	***REMOVED***

	// Now that we definitively know the minimum and maximum understood
	// version that satisfies the whole cluster, we verify that every
	// node in the cluster satisifies this.
	for _, n := range remote ***REMOVED***
		var nPCur, nDCur uint8
		if len(n.Vsn) > 0 ***REMOVED***
			nPCur = n.Vsn[2]
			nDCur = n.Vsn[5]
		***REMOVED***

		if nPCur < maxpmin || nPCur > minpmax ***REMOVED***
			return fmt.Errorf(
				"Node '%s' protocol version (%d) is incompatible: [%d, %d]",
				n.Name, nPCur, maxpmin, minpmax)
		***REMOVED***

		if nDCur < maxdmin || nDCur > mindmax ***REMOVED***
			return fmt.Errorf(
				"Node '%s' delegate protocol version (%d) is incompatible: [%d, %d]",
				n.Name, nDCur, maxdmin, mindmax)
		***REMOVED***
	***REMOVED***

	for _, n := range m.nodes ***REMOVED***
		nPCur := n.PCur
		nDCur := n.DCur

		if nPCur < maxpmin || nPCur > minpmax ***REMOVED***
			return fmt.Errorf(
				"Node '%s' protocol version (%d) is incompatible: [%d, %d]",
				n.Name, nPCur, maxpmin, minpmax)
		***REMOVED***

		if nDCur < maxdmin || nDCur > mindmax ***REMOVED***
			return fmt.Errorf(
				"Node '%s' delegate protocol version (%d) is incompatible: [%d, %d]",
				n.Name, nDCur, maxdmin, mindmax)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// nextSeqNo returns a usable sequence number in a thread safe way
func (m *Memberlist) nextSeqNo() uint32 ***REMOVED***
	return atomic.AddUint32(&m.sequenceNum, 1)
***REMOVED***

// nextIncarnation returns the next incarnation number in a thread safe way
func (m *Memberlist) nextIncarnation() uint32 ***REMOVED***
	return atomic.AddUint32(&m.incarnation, 1)
***REMOVED***

// skipIncarnation adds the positive offset to the incarnation number.
func (m *Memberlist) skipIncarnation(offset uint32) uint32 ***REMOVED***
	return atomic.AddUint32(&m.incarnation, offset)
***REMOVED***

// estNumNodes is used to get the current estimate of the number of nodes
func (m *Memberlist) estNumNodes() int ***REMOVED***
	return int(atomic.LoadUint32(&m.numNodes))
***REMOVED***

type ackMessage struct ***REMOVED***
	Complete  bool
	Payload   []byte
	Timestamp time.Time
***REMOVED***

// setProbeChannels is used to attach the ackCh to receive a message when an ack
// with a given sequence number is received. The `complete` field of the message
// will be false on timeout. Any nack messages will cause an empty struct to be
// passed to the nackCh, which can be nil if not needed.
func (m *Memberlist) setProbeChannels(seqNo uint32, ackCh chan ackMessage, nackCh chan struct***REMOVED******REMOVED***, timeout time.Duration) ***REMOVED***
	// Create handler functions for acks and nacks
	ackFn := func(payload []byte, timestamp time.Time) ***REMOVED***
		select ***REMOVED***
		case ackCh <- ackMessage***REMOVED***true, payload, timestamp***REMOVED***:
		default:
		***REMOVED***
	***REMOVED***
	nackFn := func() ***REMOVED***
		select ***REMOVED***
		case nackCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		default:
		***REMOVED***
	***REMOVED***

	// Add the handlers
	ah := &ackHandler***REMOVED***ackFn, nackFn, nil***REMOVED***
	m.ackLock.Lock()
	m.ackHandlers[seqNo] = ah
	m.ackLock.Unlock()

	// Setup a reaping routing
	ah.timer = time.AfterFunc(timeout, func() ***REMOVED***
		m.ackLock.Lock()
		delete(m.ackHandlers, seqNo)
		m.ackLock.Unlock()
		select ***REMOVED***
		case ackCh <- ackMessage***REMOVED***false, nil, time.Now()***REMOVED***:
		default:
		***REMOVED***
	***REMOVED***)
***REMOVED***

// setAckHandler is used to attach a handler to be invoked when an ack with a
// given sequence number is received. If a timeout is reached, the handler is
// deleted. This is used for indirect pings so does not configure a function
// for nacks.
func (m *Memberlist) setAckHandler(seqNo uint32, ackFn func([]byte, time.Time), timeout time.Duration) ***REMOVED***
	// Add the handler
	ah := &ackHandler***REMOVED***ackFn, nil, nil***REMOVED***
	m.ackLock.Lock()
	m.ackHandlers[seqNo] = ah
	m.ackLock.Unlock()

	// Setup a reaping routing
	ah.timer = time.AfterFunc(timeout, func() ***REMOVED***
		m.ackLock.Lock()
		delete(m.ackHandlers, seqNo)
		m.ackLock.Unlock()
	***REMOVED***)
***REMOVED***

// Invokes an ack handler if any is associated, and reaps the handler immediately
func (m *Memberlist) invokeAckHandler(ack ackResp, timestamp time.Time) ***REMOVED***
	m.ackLock.Lock()
	ah, ok := m.ackHandlers[ack.SeqNo]
	delete(m.ackHandlers, ack.SeqNo)
	m.ackLock.Unlock()
	if !ok ***REMOVED***
		return
	***REMOVED***
	ah.timer.Stop()
	ah.ackFn(ack.Payload, timestamp)
***REMOVED***

// Invokes nack handler if any is associated.
func (m *Memberlist) invokeNackHandler(nack nackResp) ***REMOVED***
	m.ackLock.Lock()
	ah, ok := m.ackHandlers[nack.SeqNo]
	m.ackLock.Unlock()
	if !ok || ah.nackFn == nil ***REMOVED***
		return
	***REMOVED***
	ah.nackFn()
***REMOVED***

// refute gossips an alive message in response to incoming information that we
// are suspect or dead. It will make sure the incarnation number beats the given
// accusedInc value, or you can supply 0 to just get the next incarnation number.
// This alters the node state that's passed in so this MUST be called while the
// nodeLock is held.
func (m *Memberlist) refute(me *nodeState, accusedInc uint32) ***REMOVED***
	// Make sure the incarnation number beats the accusation.
	inc := m.nextIncarnation()
	if accusedInc >= inc ***REMOVED***
		inc = m.skipIncarnation(accusedInc - inc + 1)
	***REMOVED***
	me.Incarnation = inc

	// Decrease our health because we are being asked to refute a problem.
	m.awareness.ApplyDelta(1)

	// Format and broadcast an alive message.
	a := alive***REMOVED***
		Incarnation: inc,
		Node:        me.Name,
		Addr:        me.Addr,
		Port:        me.Port,
		Meta:        me.Meta,
		Vsn: []uint8***REMOVED***
			me.PMin, me.PMax, me.PCur,
			me.DMin, me.DMax, me.DCur,
		***REMOVED***,
	***REMOVED***
	m.encodeAndBroadcast(me.Addr.String(), aliveMsg, a)
***REMOVED***

// aliveNode is invoked by the network layer when we get a message about a
// live node.
func (m *Memberlist) aliveNode(a *alive, notify chan struct***REMOVED******REMOVED***, bootstrap bool) ***REMOVED***
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	state, ok := m.nodeMap[a.Node]

	// It is possible that during a Leave(), there is already an aliveMsg
	// in-queue to be processed but blocked by the locks above. If we let
	// that aliveMsg process, it'll cause us to re-join the cluster. This
	// ensures that we don't.
	if m.leave && a.Node == m.config.Name ***REMOVED***
		return
	***REMOVED***

	// Invoke the Alive delegate if any. This can be used to filter out
	// alive messages based on custom logic. For example, using a cluster name.
	// Using a merge delegate is not enough, as it is possible for passive
	// cluster merging to still occur.
	if m.config.Alive != nil ***REMOVED***
		node := &Node***REMOVED***
			Name: a.Node,
			Addr: a.Addr,
			Port: a.Port,
			Meta: a.Meta,
			PMin: a.Vsn[0],
			PMax: a.Vsn[1],
			PCur: a.Vsn[2],
			DMin: a.Vsn[3],
			DMax: a.Vsn[4],
			DCur: a.Vsn[5],
		***REMOVED***
		if err := m.config.Alive.NotifyAlive(node); err != nil ***REMOVED***
			m.logger.Printf("[WARN] memberlist: ignoring alive message for '%s': %s",
				a.Node, err)
			return
		***REMOVED***
	***REMOVED***

	// Check if we've never seen this node before, and if not, then
	// store this node in our node map.
	if !ok ***REMOVED***
		state = &nodeState***REMOVED***
			Node: Node***REMOVED***
				Name: a.Node,
				Addr: a.Addr,
				Port: a.Port,
				Meta: a.Meta,
			***REMOVED***,
			State: stateDead,
		***REMOVED***

		// Add to map
		m.nodeMap[a.Node] = state

		// Get a random offset. This is important to ensure
		// the failure detection bound is low on average. If all
		// nodes did an append, failure detection bound would be
		// very high.
		n := len(m.nodes)
		offset := randomOffset(n)

		// Add at the end and swap with the node at the offset
		m.nodes = append(m.nodes, state)
		m.nodes[offset], m.nodes[n] = m.nodes[n], m.nodes[offset]

		// Update numNodes after we've added a new node
		atomic.AddUint32(&m.numNodes, 1)
	***REMOVED***

	// Check if this address is different than the existing node
	if !bytes.Equal([]byte(state.Addr), a.Addr) || state.Port != a.Port ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Conflicting address for %s. Mine: %v:%d Theirs: %v:%d",
			state.Name, state.Addr, state.Port, net.IP(a.Addr), a.Port)

		// Inform the conflict delegate if provided
		if m.config.Conflict != nil ***REMOVED***
			other := Node***REMOVED***
				Name: a.Node,
				Addr: a.Addr,
				Port: a.Port,
				Meta: a.Meta,
			***REMOVED***
			m.config.Conflict.NotifyConflict(&state.Node, &other)
		***REMOVED***
		return
	***REMOVED***

	// Bail if the incarnation number is older, and this is not about us
	isLocalNode := state.Name == m.config.Name
	if a.Incarnation <= state.Incarnation && !isLocalNode ***REMOVED***
		return
	***REMOVED***

	// Bail if strictly less and this is about us
	if a.Incarnation < state.Incarnation && isLocalNode ***REMOVED***
		return
	***REMOVED***

	// Clear out any suspicion timer that may be in effect.
	delete(m.nodeTimers, a.Node)

	// Store the old state and meta data
	oldState := state.State
	oldMeta := state.Meta

	// If this is us we need to refute, otherwise re-broadcast
	if !bootstrap && isLocalNode ***REMOVED***
		// Compute the version vector
		versions := []uint8***REMOVED***
			state.PMin, state.PMax, state.PCur,
			state.DMin, state.DMax, state.DCur,
		***REMOVED***

		// If the Incarnation is the same, we need special handling, since it
		// possible for the following situation to happen:
		// 1) Start with configuration C, join cluster
		// 2) Hard fail / Kill / Shutdown
		// 3) Restart with configuration C', join cluster
		//
		// In this case, other nodes and the local node see the same incarnation,
		// but the values may not be the same. For this reason, we always
		// need to do an equality check for this Incarnation. In most cases,
		// we just ignore, but we may need to refute.
		//
		if a.Incarnation == state.Incarnation &&
			bytes.Equal(a.Meta, state.Meta) &&
			bytes.Equal(a.Vsn, versions) ***REMOVED***
			return
		***REMOVED***

		m.refute(state, a.Incarnation)
		m.logger.Printf("[WARN] memberlist: Refuting an alive message")
	***REMOVED*** else ***REMOVED***
		m.encodeBroadcastNotify(a.Node, aliveMsg, a, notify)

		// Update protocol versions if it arrived
		if len(a.Vsn) > 0 ***REMOVED***
			state.PMin = a.Vsn[0]
			state.PMax = a.Vsn[1]
			state.PCur = a.Vsn[2]
			state.DMin = a.Vsn[3]
			state.DMax = a.Vsn[4]
			state.DCur = a.Vsn[5]
		***REMOVED***

		// Update the state and incarnation number
		state.Incarnation = a.Incarnation
		state.Meta = a.Meta
		if state.State != stateAlive ***REMOVED***
			state.State = stateAlive
			state.StateChange = time.Now()
		***REMOVED***
	***REMOVED***

	// Update metrics
	metrics.IncrCounter([]string***REMOVED***"memberlist", "msg", "alive"***REMOVED***, 1)

	// Notify the delegate of any relevant updates
	if m.config.Events != nil ***REMOVED***
		if oldState == stateDead ***REMOVED***
			// if Dead -> Alive, notify of join
			m.config.Events.NotifyJoin(&state.Node)

		***REMOVED*** else if !bytes.Equal(oldMeta, state.Meta) ***REMOVED***
			// if Meta changed, trigger an update notification
			m.config.Events.NotifyUpdate(&state.Node)
		***REMOVED***
	***REMOVED***
***REMOVED***

// suspectNode is invoked by the network layer when we get a message
// about a suspect node
func (m *Memberlist) suspectNode(s *suspect) ***REMOVED***
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	state, ok := m.nodeMap[s.Node]

	// If we've never heard about this node before, ignore it
	if !ok ***REMOVED***
		return
	***REMOVED***

	// Ignore old incarnation numbers
	if s.Incarnation < state.Incarnation ***REMOVED***
		return
	***REMOVED***

	// See if there's a suspicion timer we can confirm. If the info is new
	// to us we will go ahead and re-gossip it. This allows for multiple
	// independent confirmations to flow even when a node probes a node
	// that's already suspect.
	if timer, ok := m.nodeTimers[s.Node]; ok ***REMOVED***
		if timer.Confirm(s.From) ***REMOVED***
			m.encodeAndBroadcast(s.Node, suspectMsg, s)
		***REMOVED***
		return
	***REMOVED***

	// Ignore non-alive nodes
	if state.State != stateAlive ***REMOVED***
		return
	***REMOVED***

	// If this is us we need to refute, otherwise re-broadcast
	if state.Name == m.config.Name ***REMOVED***
		m.refute(state, s.Incarnation)
		m.logger.Printf("[WARN] memberlist: Refuting a suspect message (from: %s)", s.From)
		return // Do not mark ourself suspect
	***REMOVED*** else ***REMOVED***
		m.encodeAndBroadcast(s.Node, suspectMsg, s)
	***REMOVED***

	// Update metrics
	metrics.IncrCounter([]string***REMOVED***"memberlist", "msg", "suspect"***REMOVED***, 1)

	// Update the state
	state.Incarnation = s.Incarnation
	state.State = stateSuspect
	changeTime := time.Now()
	state.StateChange = changeTime

	// Setup a suspicion timer. Given that we don't have any known phase
	// relationship with our peers, we set up k such that we hit the nominal
	// timeout two probe intervals short of what we expect given the suspicion
	// multiplier.
	k := m.config.SuspicionMult - 2

	// If there aren't enough nodes to give the expected confirmations, just
	// set k to 0 to say that we don't expect any. Note we subtract 2 from n
	// here to take out ourselves and the node being probed.
	n := m.estNumNodes()
	if n-2 < k ***REMOVED***
		k = 0
	***REMOVED***

	// Compute the timeouts based on the size of the cluster.
	min := suspicionTimeout(m.config.SuspicionMult, n, m.config.ProbeInterval)
	max := time.Duration(m.config.SuspicionMaxTimeoutMult) * min
	fn := func(numConfirmations int) ***REMOVED***
		m.nodeLock.Lock()
		state, ok := m.nodeMap[s.Node]
		timeout := ok && state.State == stateSuspect && state.StateChange == changeTime
		m.nodeLock.Unlock()

		if timeout ***REMOVED***
			if k > 0 && numConfirmations < k ***REMOVED***
				metrics.IncrCounter([]string***REMOVED***"memberlist", "degraded", "timeout"***REMOVED***, 1)
			***REMOVED***

			m.logger.Printf("[INFO] memberlist: Marking %s as failed, suspect timeout reached (%d peer confirmations)",
				state.Name, numConfirmations)
			d := dead***REMOVED***Incarnation: state.Incarnation, Node: state.Name, From: m.config.Name***REMOVED***
			m.deadNode(&d)
		***REMOVED***
	***REMOVED***
	m.nodeTimers[s.Node] = newSuspicion(s.From, k, min, max, fn)
***REMOVED***

// deadNode is invoked by the network layer when we get a message
// about a dead node
func (m *Memberlist) deadNode(d *dead) ***REMOVED***
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	state, ok := m.nodeMap[d.Node]

	// If we've never heard about this node before, ignore it
	if !ok ***REMOVED***
		return
	***REMOVED***

	// Ignore old incarnation numbers
	if d.Incarnation < state.Incarnation ***REMOVED***
		return
	***REMOVED***

	// Clear out any suspicion timer that may be in effect.
	delete(m.nodeTimers, d.Node)

	// Ignore if node is already dead
	if state.State == stateDead ***REMOVED***
		return
	***REMOVED***

	// Check if this is us
	if state.Name == m.config.Name ***REMOVED***
		// If we are not leaving we need to refute
		if !m.leave ***REMOVED***
			m.refute(state, d.Incarnation)
			m.logger.Printf("[WARN] memberlist: Refuting a dead message (from: %s)", d.From)
			return // Do not mark ourself dead
		***REMOVED***

		// If we are leaving, we broadcast and wait
		m.encodeBroadcastNotify(d.Node, deadMsg, d, m.leaveBroadcast)
	***REMOVED*** else ***REMOVED***
		m.encodeAndBroadcast(d.Node, deadMsg, d)
	***REMOVED***

	// Update metrics
	metrics.IncrCounter([]string***REMOVED***"memberlist", "msg", "dead"***REMOVED***, 1)

	// Update the state
	state.Incarnation = d.Incarnation
	state.State = stateDead
	state.StateChange = time.Now()

	// Notify of death
	if m.config.Events != nil ***REMOVED***
		m.config.Events.NotifyLeave(&state.Node)
	***REMOVED***
***REMOVED***

// mergeState is invoked by the network layer when we get a Push/Pull
// state transfer
func (m *Memberlist) mergeState(remote []pushNodeState) ***REMOVED***
	for _, r := range remote ***REMOVED***
		switch r.State ***REMOVED***
		case stateAlive:
			a := alive***REMOVED***
				Incarnation: r.Incarnation,
				Node:        r.Name,
				Addr:        r.Addr,
				Port:        r.Port,
				Meta:        r.Meta,
				Vsn:         r.Vsn,
			***REMOVED***
			m.aliveNode(&a, nil, false)

		case stateDead:
			// If the remote node believes a node is dead, we prefer to
			// suspect that node instead of declaring it dead instantly
			fallthrough
		case stateSuspect:
			s := suspect***REMOVED***Incarnation: r.Incarnation, Node: r.Name, From: m.config.Name***REMOVED***
			m.suspectNode(&s)
		***REMOVED***
	***REMOVED***
***REMOVED***
