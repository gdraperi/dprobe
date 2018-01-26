package memberlist

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-msgpack/codec"
)

// This is the minimum and maximum protocol version that we can
// _understand_. We're allowed to speak at any version within this
// range. This range is inclusive.
const (
	ProtocolVersionMin uint8 = 1

	// Version 3 added support for TCP pings but we kept the default
	// protocol version at 2 to ease transition to this new feature.
	// A memberlist speaking version 2 of the protocol will attempt
	// to TCP ping another memberlist who understands version 3 or
	// greater.
	//
	// Version 4 added support for nacks as part of indirect probes.
	// A memberlist speaking version 2 of the protocol will expect
	// nacks from another memberlist who understands version 4 or
	// greater, and likewise nacks will be sent to memberlists who
	// understand version 4 or greater.
	ProtocolVersion2Compatible = 2

	ProtocolVersionMax = 5
)

// messageType is an integer ID of a type of message that can be received
// on network channels from other members.
type messageType uint8

// The list of available message types.
const (
	pingMsg messageType = iota
	indirectPingMsg
	ackRespMsg
	suspectMsg
	aliveMsg
	deadMsg
	pushPullMsg
	compoundMsg
	userMsg // User mesg, not handled by us
	compressMsg
	encryptMsg
	nackRespMsg
	hasCrcMsg
)

// compressionType is used to specify the compression algorithm
type compressionType uint8

const (
	lzwAlgo compressionType = iota
)

const (
	MetaMaxSize            = 512 // Maximum size for node meta data
	compoundHeaderOverhead = 2   // Assumed header overhead
	compoundOverhead       = 2   // Assumed overhead per entry in compoundHeader
	userMsgOverhead        = 1
	blockingWarning        = 10 * time.Millisecond // Warn if a UDP packet takes this long to process
	maxPushStateBytes      = 10 * 1024 * 1024
)

// ping request sent directly to node
type ping struct ***REMOVED***
	SeqNo uint32

	// Node is sent so the target can verify they are
	// the intended recipient. This is to protect again an agent
	// restart with a new name.
	Node string
***REMOVED***

// indirect ping sent to an indirect ndoe
type indirectPingReq struct ***REMOVED***
	SeqNo  uint32
	Target []byte
	Port   uint16
	Node   string
	Nack   bool // true if we'd like a nack back
***REMOVED***

// ack response is sent for a ping
type ackResp struct ***REMOVED***
	SeqNo   uint32
	Payload []byte
***REMOVED***

// nack response is sent for an indirect ping when the pinger doesn't hear from
// the ping-ee within the configured timeout. This lets the original node know
// that the indirect ping attempt happened but didn't succeed.
type nackResp struct ***REMOVED***
	SeqNo uint32
***REMOVED***

// suspect is broadcast when we suspect a node is dead
type suspect struct ***REMOVED***
	Incarnation uint32
	Node        string
	From        string // Include who is suspecting
***REMOVED***

// alive is broadcast when we know a node is alive.
// Overloaded for nodes joining
type alive struct ***REMOVED***
	Incarnation uint32
	Node        string
	Addr        []byte
	Port        uint16
	Meta        []byte

	// The versions of the protocol/delegate that are being spoken, order:
	// pmin, pmax, pcur, dmin, dmax, dcur
	Vsn []uint8
***REMOVED***

// dead is broadcast when we confirm a node is dead
// Overloaded for nodes leaving
type dead struct ***REMOVED***
	Incarnation uint32
	Node        string
	From        string // Include who is suspecting
***REMOVED***

// pushPullHeader is used to inform the
// otherside how many states we are transferring
type pushPullHeader struct ***REMOVED***
	Nodes        int
	UserStateLen int  // Encodes the byte lengh of user state
	Join         bool // Is this a join request or a anti-entropy run
***REMOVED***

// userMsgHeader is used to encapsulate a userMsg
type userMsgHeader struct ***REMOVED***
	UserMsgLen int // Encodes the byte lengh of user state
***REMOVED***

// pushNodeState is used for pushPullReq when we are
// transferring out node states
type pushNodeState struct ***REMOVED***
	Name        string
	Addr        []byte
	Port        uint16
	Meta        []byte
	Incarnation uint32
	State       nodeStateType
	Vsn         []uint8 // Protocol versions
***REMOVED***

// compress is used to wrap an underlying payload
// using a specified compression algorithm
type compress struct ***REMOVED***
	Algo compressionType
	Buf  []byte
***REMOVED***

// msgHandoff is used to transfer a message between goroutines
type msgHandoff struct ***REMOVED***
	msgType messageType
	buf     []byte
	from    net.Addr
***REMOVED***

// encryptionVersion returns the encryption version to use
func (m *Memberlist) encryptionVersion() encryptionVersion ***REMOVED***
	switch m.ProtocolVersion() ***REMOVED***
	case 1:
		return 0
	default:
		return 1
	***REMOVED***
***REMOVED***

// streamListen is a long running goroutine that pulls incoming streams from the
// transport and hands them off for processing.
func (m *Memberlist) streamListen() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case conn := <-m.transport.StreamCh():
			go m.handleConn(conn)

		case <-m.shutdownCh:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// handleConn handles a single incoming stream connection from the transport.
func (m *Memberlist) handleConn(conn net.Conn) ***REMOVED***
	m.logger.Printf("[DEBUG] memberlist: Stream connection %s", LogConn(conn))

	defer conn.Close()
	metrics.IncrCounter([]string***REMOVED***"memberlist", "tcp", "accept"***REMOVED***, 1)

	conn.SetDeadline(time.Now().Add(m.config.TCPTimeout))
	msgType, bufConn, dec, err := m.readStream(conn)
	if err != nil ***REMOVED***
		if err != io.EOF ***REMOVED***
			m.logger.Printf("[ERR] memberlist: failed to receive: %s %s", err, LogConn(conn))
		***REMOVED***
		return
	***REMOVED***

	switch msgType ***REMOVED***
	case userMsg:
		if err := m.readUserMsg(bufConn, dec); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to receive user message: %s %s", err, LogConn(conn))
		***REMOVED***
	case pushPullMsg:
		join, remoteNodes, userState, err := m.readRemoteState(bufConn, dec)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to read remote state: %s %s", err, LogConn(conn))
			return
		***REMOVED***

		if err := m.sendLocalState(conn, join); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to push local state: %s %s", err, LogConn(conn))
			return
		***REMOVED***

		if err := m.mergeRemoteState(join, remoteNodes, userState); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed push/pull merge: %s %s", err, LogConn(conn))
			return
		***REMOVED***
	case pingMsg:
		var p ping
		if err := dec.Decode(&p); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to decode ping: %s %s", err, LogConn(conn))
			return
		***REMOVED***

		if p.Node != "" && p.Node != m.config.Name ***REMOVED***
			m.logger.Printf("[WARN] memberlist: Got ping for unexpected node %s %s", p.Node, LogConn(conn))
			return
		***REMOVED***

		ack := ackResp***REMOVED***p.SeqNo, nil***REMOVED***
		out, err := encode(ackRespMsg, &ack)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to encode ack: %s", err)
			return
		***REMOVED***

		err = m.rawSendMsgStream(conn, out.Bytes())
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to send ack: %s %s", err, LogConn(conn))
			return
		***REMOVED***
	default:
		m.logger.Printf("[ERR] memberlist: Received invalid msgType (%d) %s", msgType, LogConn(conn))
	***REMOVED***
***REMOVED***

// packetListen is a long running goroutine that pulls packets out of the
// transport and hands them off for processing.
func (m *Memberlist) packetListen() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case packet := <-m.transport.PacketCh():
			m.ingestPacket(packet.Buf, packet.From, packet.Timestamp)

		case <-m.shutdownCh:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (m *Memberlist) ingestPacket(buf []byte, from net.Addr, timestamp time.Time) ***REMOVED***
	// Check if encryption is enabled
	if m.config.EncryptionEnabled() ***REMOVED***
		// Decrypt the payload
		plain, err := decryptPayload(m.config.Keyring.GetKeys(), buf, nil)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Decrypt packet failed: %v %s", err, LogAddress(from))
			return
		***REMOVED***

		// Continue processing the plaintext buffer
		buf = plain
	***REMOVED***

	// See if there's a checksum included to verify the contents of the message
	if len(buf) >= 5 && messageType(buf[0]) == hasCrcMsg ***REMOVED***
		crc := crc32.ChecksumIEEE(buf[5:])
		expected := binary.BigEndian.Uint32(buf[1:5])
		if crc != expected ***REMOVED***
			m.logger.Printf("[WARN] memberlist: Got invalid checksum for UDP packet: %x, %x", crc, expected)
			return
		***REMOVED***
		m.handleCommand(buf[5:], from, timestamp)
	***REMOVED*** else ***REMOVED***
		m.handleCommand(buf, from, timestamp)
	***REMOVED***
***REMOVED***

func (m *Memberlist) handleCommand(buf []byte, from net.Addr, timestamp time.Time) ***REMOVED***
	// Decode the message type
	msgType := messageType(buf[0])
	buf = buf[1:]

	// Switch on the msgType
	switch msgType ***REMOVED***
	case compoundMsg:
		m.handleCompound(buf, from, timestamp)
	case compressMsg:
		m.handleCompressed(buf, from, timestamp)

	case pingMsg:
		m.handlePing(buf, from)
	case indirectPingMsg:
		m.handleIndirectPing(buf, from)
	case ackRespMsg:
		m.handleAck(buf, from, timestamp)
	case nackRespMsg:
		m.handleNack(buf, from)

	case suspectMsg:
		fallthrough
	case aliveMsg:
		fallthrough
	case deadMsg:
		fallthrough
	case userMsg:
		select ***REMOVED***
		case m.handoff <- msgHandoff***REMOVED***msgType, buf, from***REMOVED***:
		default:
			m.logger.Printf("[WARN] memberlist: handler queue full, dropping message (%d) %s", msgType, LogAddress(from))
		***REMOVED***

	default:
		m.logger.Printf("[ERR] memberlist: msg type (%d) not supported %s", msgType, LogAddress(from))
	***REMOVED***
***REMOVED***

// packetHandler is a long running goroutine that processes messages received
// over the packet interface, but is decoupled from the listener to avoid
// blocking the listener which may cause ping/ack messages to be delayed.
func (m *Memberlist) packetHandler() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case msg := <-m.handoff:
			msgType := msg.msgType
			buf := msg.buf
			from := msg.from

			switch msgType ***REMOVED***
			case suspectMsg:
				m.handleSuspect(buf, from)
			case aliveMsg:
				m.handleAlive(buf, from)
			case deadMsg:
				m.handleDead(buf, from)
			case userMsg:
				m.handleUser(buf, from)
			default:
				m.logger.Printf("[ERR] memberlist: Message type (%d) not supported %s (packet handler)", msgType, LogAddress(from))
			***REMOVED***

		case <-m.shutdownCh:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (m *Memberlist) handleCompound(buf []byte, from net.Addr, timestamp time.Time) ***REMOVED***
	// Decode the parts
	trunc, parts, err := decodeCompoundMessage(buf)
	if err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode compound request: %s %s", err, LogAddress(from))
		return
	***REMOVED***

	// Log any truncation
	if trunc > 0 ***REMOVED***
		m.logger.Printf("[WARN] memberlist: Compound request had %d truncated messages %s", trunc, LogAddress(from))
	***REMOVED***

	// Handle each message
	for _, part := range parts ***REMOVED***
		m.handleCommand(part, from, timestamp)
	***REMOVED***
***REMOVED***

func (m *Memberlist) handlePing(buf []byte, from net.Addr) ***REMOVED***
	var p ping
	if err := decode(buf, &p); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode ping request: %s %s", err, LogAddress(from))
		return
	***REMOVED***
	// If node is provided, verify that it is for us
	if p.Node != "" && p.Node != m.config.Name ***REMOVED***
		m.logger.Printf("[WARN] memberlist: Got ping for unexpected node '%s' %s", p.Node, LogAddress(from))
		return
	***REMOVED***
	var ack ackResp
	ack.SeqNo = p.SeqNo
	if m.config.Ping != nil ***REMOVED***
		ack.Payload = m.config.Ping.AckPayload()
	***REMOVED***
	if err := m.encodeAndSendMsg(from.String(), ackRespMsg, &ack); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to send ack: %s %s", err, LogAddress(from))
	***REMOVED***
***REMOVED***

func (m *Memberlist) handleIndirectPing(buf []byte, from net.Addr) ***REMOVED***
	var ind indirectPingReq
	if err := decode(buf, &ind); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode indirect ping request: %s %s", err, LogAddress(from))
		return
	***REMOVED***

	// For proto versions < 2, there is no port provided. Mask old
	// behavior by using the configured port.
	if m.ProtocolVersion() < 2 || ind.Port == 0 ***REMOVED***
		ind.Port = uint16(m.config.BindPort)
	***REMOVED***

	// Send a ping to the correct host.
	localSeqNo := m.nextSeqNo()
	ping := ping***REMOVED***SeqNo: localSeqNo, Node: ind.Node***REMOVED***

	// Setup a response handler to relay the ack
	cancelCh := make(chan struct***REMOVED******REMOVED***)
	respHandler := func(payload []byte, timestamp time.Time) ***REMOVED***
		// Try to prevent the nack if we've caught it in time.
		close(cancelCh)

		// Forward the ack back to the requestor.
		ack := ackResp***REMOVED***ind.SeqNo, nil***REMOVED***
		if err := m.encodeAndSendMsg(from.String(), ackRespMsg, &ack); err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to forward ack: %s %s", err, LogAddress(from))
		***REMOVED***
	***REMOVED***
	m.setAckHandler(localSeqNo, respHandler, m.config.ProbeTimeout)

	// Send the ping.
	addr := joinHostPort(net.IP(ind.Target).String(), ind.Port)
	if err := m.encodeAndSendMsg(addr, pingMsg, &ping); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to send ping: %s %s", err, LogAddress(from))
	***REMOVED***

	// Setup a timer to fire off a nack if no ack is seen in time.
	if ind.Nack ***REMOVED***
		go func() ***REMOVED***
			select ***REMOVED***
			case <-cancelCh:
				return
			case <-time.After(m.config.ProbeTimeout):
				nack := nackResp***REMOVED***ind.SeqNo***REMOVED***
				if err := m.encodeAndSendMsg(from.String(), nackRespMsg, &nack); err != nil ***REMOVED***
					m.logger.Printf("[ERR] memberlist: Failed to send nack: %s %s", err, LogAddress(from))
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

func (m *Memberlist) handleAck(buf []byte, from net.Addr, timestamp time.Time) ***REMOVED***
	var ack ackResp
	if err := decode(buf, &ack); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode ack response: %s %s", err, LogAddress(from))
		return
	***REMOVED***
	m.invokeAckHandler(ack, timestamp)
***REMOVED***

func (m *Memberlist) handleNack(buf []byte, from net.Addr) ***REMOVED***
	var nack nackResp
	if err := decode(buf, &nack); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode nack response: %s %s", err, LogAddress(from))
		return
	***REMOVED***
	m.invokeNackHandler(nack)
***REMOVED***

func (m *Memberlist) handleSuspect(buf []byte, from net.Addr) ***REMOVED***
	var sus suspect
	if err := decode(buf, &sus); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode suspect message: %s %s", err, LogAddress(from))
		return
	***REMOVED***
	m.suspectNode(&sus)
***REMOVED***

func (m *Memberlist) handleAlive(buf []byte, from net.Addr) ***REMOVED***
	var live alive
	if err := decode(buf, &live); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode alive message: %s %s", err, LogAddress(from))
		return
	***REMOVED***

	// For proto versions < 2, there is no port provided. Mask old
	// behavior by using the configured port
	if m.ProtocolVersion() < 2 || live.Port == 0 ***REMOVED***
		live.Port = uint16(m.config.BindPort)
	***REMOVED***

	m.aliveNode(&live, nil, false)
***REMOVED***

func (m *Memberlist) handleDead(buf []byte, from net.Addr) ***REMOVED***
	var d dead
	if err := decode(buf, &d); err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decode dead message: %s %s", err, LogAddress(from))
		return
	***REMOVED***
	m.deadNode(&d)
***REMOVED***

// handleUser is used to notify channels of incoming user data
func (m *Memberlist) handleUser(buf []byte, from net.Addr) ***REMOVED***
	d := m.config.Delegate
	if d != nil ***REMOVED***
		d.NotifyMsg(buf)
	***REMOVED***
***REMOVED***

// handleCompressed is used to unpack a compressed message
func (m *Memberlist) handleCompressed(buf []byte, from net.Addr, timestamp time.Time) ***REMOVED***
	// Try to decode the payload
	payload, err := decompressPayload(buf)
	if err != nil ***REMOVED***
		m.logger.Printf("[ERR] memberlist: Failed to decompress payload: %v %s", err, LogAddress(from))
		return
	***REMOVED***

	// Recursively handle the payload
	m.handleCommand(payload, from, timestamp)
***REMOVED***

// encodeAndSendMsg is used to combine the encoding and sending steps
func (m *Memberlist) encodeAndSendMsg(addr string, msgType messageType, msg interface***REMOVED******REMOVED***) error ***REMOVED***
	out, err := encode(msgType, msg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := m.sendMsg(addr, out.Bytes()); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// sendMsg is used to send a message via packet to another host. It will
// opportunistically create a compoundMsg and piggy back other broadcasts.
func (m *Memberlist) sendMsg(addr string, msg []byte) error ***REMOVED***
	// Check if we can piggy back any messages
	bytesAvail := m.config.UDPBufferSize - len(msg) - compoundHeaderOverhead
	if m.config.EncryptionEnabled() ***REMOVED***
		bytesAvail -= encryptOverhead(m.encryptionVersion())
	***REMOVED***
	extra := m.getBroadcasts(compoundOverhead, bytesAvail)

	// Fast path if nothing to piggypack
	if len(extra) == 0 ***REMOVED***
		return m.rawSendMsgPacket(addr, nil, msg)
	***REMOVED***

	// Join all the messages
	msgs := make([][]byte, 0, 1+len(extra))
	msgs = append(msgs, msg)
	msgs = append(msgs, extra...)

	// Create a compound message
	compound := makeCompoundMessage(msgs)

	// Send the message
	return m.rawSendMsgPacket(addr, nil, compound.Bytes())
***REMOVED***

// rawSendMsgPacket is used to send message via packet to another host without
// modification, other than compression or encryption if enabled.
func (m *Memberlist) rawSendMsgPacket(addr string, node *Node, msg []byte) error ***REMOVED***
	// Check if we have compression enabled
	if m.config.EnableCompression ***REMOVED***
		buf, err := compressPayload(msg)
		if err != nil ***REMOVED***
			m.logger.Printf("[WARN] memberlist: Failed to compress payload: %v", err)
		***REMOVED*** else ***REMOVED***
			// Only use compression if it reduced the size
			if buf.Len() < len(msg) ***REMOVED***
				msg = buf.Bytes()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Try to look up the destination node
	if node == nil ***REMOVED***
		toAddr, _, err := net.SplitHostPort(addr)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Failed to parse address %q: %v", addr, err)
			return err
		***REMOVED***
		m.nodeLock.RLock()
		nodeState, ok := m.nodeMap[toAddr]
		m.nodeLock.RUnlock()
		if ok ***REMOVED***
			node = &nodeState.Node
		***REMOVED***
	***REMOVED***

	// Add a CRC to the end of the payload if the recipient understands
	// ProtocolVersion >= 5
	if node != nil && node.PMax >= 5 ***REMOVED***
		crc := crc32.ChecksumIEEE(msg)
		header := make([]byte, 5, 5+len(msg))
		header[0] = byte(hasCrcMsg)
		binary.BigEndian.PutUint32(header[1:], crc)
		msg = append(header, msg...)
	***REMOVED***

	// Check if we have encryption enabled
	if m.config.EncryptionEnabled() ***REMOVED***
		// Encrypt the payload
		var buf bytes.Buffer
		primaryKey := m.config.Keyring.GetPrimaryKey()
		err := encryptPayload(m.encryptionVersion(), primaryKey, msg, nil, &buf)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERR] memberlist: Encryption of message failed: %v", err)
			return err
		***REMOVED***
		msg = buf.Bytes()
	***REMOVED***

	metrics.IncrCounter([]string***REMOVED***"memberlist", "udp", "sent"***REMOVED***, float32(len(msg)))
	_, err := m.transport.WriteTo(msg, addr)
	return err
***REMOVED***

// rawSendMsgStream is used to stream a message to another host without
// modification, other than applying compression and encryption if enabled.
func (m *Memberlist) rawSendMsgStream(conn net.Conn, sendBuf []byte) error ***REMOVED***
	// Check if compresion is enabled
	if m.config.EnableCompression ***REMOVED***
		compBuf, err := compressPayload(sendBuf)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERROR] memberlist: Failed to compress payload: %v", err)
		***REMOVED*** else ***REMOVED***
			sendBuf = compBuf.Bytes()
		***REMOVED***
	***REMOVED***

	// Check if encryption is enabled
	if m.config.EncryptionEnabled() ***REMOVED***
		crypt, err := m.encryptLocalState(sendBuf)
		if err != nil ***REMOVED***
			m.logger.Printf("[ERROR] memberlist: Failed to encrypt local state: %v", err)
			return err
		***REMOVED***
		sendBuf = crypt
	***REMOVED***

	// Write out the entire send buffer
	metrics.IncrCounter([]string***REMOVED***"memberlist", "tcp", "sent"***REMOVED***, float32(len(sendBuf)))

	if n, err := conn.Write(sendBuf); err != nil ***REMOVED***
		return err
	***REMOVED*** else if n != len(sendBuf) ***REMOVED***
		return fmt.Errorf("only %d of %d bytes written", n, len(sendBuf))
	***REMOVED***

	return nil
***REMOVED***

// sendUserMsg is used to stream a user message to another host.
func (m *Memberlist) sendUserMsg(addr string, sendBuf []byte) error ***REMOVED***
	conn, err := m.transport.DialTimeout(addr, m.config.TCPTimeout)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer conn.Close()

	bufConn := bytes.NewBuffer(nil)
	if err := bufConn.WriteByte(byte(userMsg)); err != nil ***REMOVED***
		return err
	***REMOVED***

	header := userMsgHeader***REMOVED***UserMsgLen: len(sendBuf)***REMOVED***
	hd := codec.MsgpackHandle***REMOVED******REMOVED***
	enc := codec.NewEncoder(bufConn, &hd)
	if err := enc.Encode(&header); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := bufConn.Write(sendBuf); err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.rawSendMsgStream(conn, bufConn.Bytes())
***REMOVED***

// sendAndReceiveState is used to initiate a push/pull over a stream with a
// remote host.
func (m *Memberlist) sendAndReceiveState(addr string, join bool) ([]pushNodeState, []byte, error) ***REMOVED***
	// Attempt to connect
	conn, err := m.transport.DialTimeout(addr, m.config.TCPTimeout)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer conn.Close()
	m.logger.Printf("[DEBUG] memberlist: Initiating push/pull sync with: %s", conn.RemoteAddr())
	metrics.IncrCounter([]string***REMOVED***"memberlist", "tcp", "connect"***REMOVED***, 1)

	// Send our state
	if err := m.sendLocalState(conn, join); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	conn.SetDeadline(time.Now().Add(m.config.TCPTimeout))
	msgType, bufConn, dec, err := m.readStream(conn)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	// Quit if not push/pull
	if msgType != pushPullMsg ***REMOVED***
		err := fmt.Errorf("received invalid msgType (%d), expected pushPullMsg (%d) %s", msgType, pushPullMsg, LogConn(conn))
		return nil, nil, err
	***REMOVED***

	// Read remote state
	_, remoteNodes, userState, err := m.readRemoteState(bufConn, dec)
	return remoteNodes, userState, err
***REMOVED***

// sendLocalState is invoked to send our local state over a stream connection.
func (m *Memberlist) sendLocalState(conn net.Conn, join bool) error ***REMOVED***
	// Setup a deadline
	conn.SetDeadline(time.Now().Add(m.config.TCPTimeout))

	// Prepare the local node state
	m.nodeLock.RLock()
	localNodes := make([]pushNodeState, len(m.nodes))
	for idx, n := range m.nodes ***REMOVED***
		localNodes[idx].Name = n.Name
		localNodes[idx].Addr = n.Addr
		localNodes[idx].Port = n.Port
		localNodes[idx].Incarnation = n.Incarnation
		localNodes[idx].State = n.State
		localNodes[idx].Meta = n.Meta
		localNodes[idx].Vsn = []uint8***REMOVED***
			n.PMin, n.PMax, n.PCur,
			n.DMin, n.DMax, n.DCur,
		***REMOVED***
	***REMOVED***
	m.nodeLock.RUnlock()

	// Get the delegate state
	var userData []byte
	if m.config.Delegate != nil ***REMOVED***
		userData = m.config.Delegate.LocalState(join)
	***REMOVED***

	// Create a bytes buffer writer
	bufConn := bytes.NewBuffer(nil)

	// Send our node state
	header := pushPullHeader***REMOVED***Nodes: len(localNodes), UserStateLen: len(userData), Join: join***REMOVED***
	hd := codec.MsgpackHandle***REMOVED******REMOVED***
	enc := codec.NewEncoder(bufConn, &hd)

	// Begin state push
	if _, err := bufConn.Write([]byte***REMOVED***byte(pushPullMsg)***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := enc.Encode(&header); err != nil ***REMOVED***
		return err
	***REMOVED***
	for i := 0; i < header.Nodes; i++ ***REMOVED***
		if err := enc.Encode(&localNodes[i]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Write the user state as well
	if userData != nil ***REMOVED***
		if _, err := bufConn.Write(userData); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Get the send buffer
	return m.rawSendMsgStream(conn, bufConn.Bytes())
***REMOVED***

// encryptLocalState is used to help encrypt local state before sending
func (m *Memberlist) encryptLocalState(sendBuf []byte) ([]byte, error) ***REMOVED***
	var buf bytes.Buffer

	// Write the encryptMsg byte
	buf.WriteByte(byte(encryptMsg))

	// Write the size of the message
	sizeBuf := make([]byte, 4)
	encVsn := m.encryptionVersion()
	encLen := encryptedLength(encVsn, len(sendBuf))
	binary.BigEndian.PutUint32(sizeBuf, uint32(encLen))
	buf.Write(sizeBuf)

	// Write the encrypted cipher text to the buffer
	key := m.config.Keyring.GetPrimaryKey()
	err := encryptPayload(encVsn, key, sendBuf, buf.Bytes()[:5], &buf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***

// decryptRemoteState is used to help decrypt the remote state
func (m *Memberlist) decryptRemoteState(bufConn io.Reader) ([]byte, error) ***REMOVED***
	// Read in enough to determine message length
	cipherText := bytes.NewBuffer(nil)
	cipherText.WriteByte(byte(encryptMsg))
	_, err := io.CopyN(cipherText, bufConn, 4)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Ensure we aren't asked to download too much. This is to guard against
	// an attack vector where a huge amount of state is sent
	moreBytes := binary.BigEndian.Uint32(cipherText.Bytes()[1:5])
	if moreBytes > maxPushStateBytes ***REMOVED***
		return nil, fmt.Errorf("Remote node state is larger than limit (%d)", moreBytes)
	***REMOVED***

	// Read in the rest of the payload
	_, err = io.CopyN(cipherText, bufConn, int64(moreBytes))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Decrypt the cipherText
	dataBytes := cipherText.Bytes()[:5]
	cipherBytes := cipherText.Bytes()[5:]

	// Decrypt the payload
	keys := m.config.Keyring.GetKeys()
	return decryptPayload(keys, cipherBytes, dataBytes)
***REMOVED***

// readStream is used to read from a stream connection, decrypting and
// decompressing the stream if necessary.
func (m *Memberlist) readStream(conn net.Conn) (messageType, io.Reader, *codec.Decoder, error) ***REMOVED***
	// Created a buffered reader
	var bufConn io.Reader = bufio.NewReader(conn)

	// Read the message type
	buf := [1]byte***REMOVED***0***REMOVED***
	if _, err := bufConn.Read(buf[:]); err != nil ***REMOVED***
		return 0, nil, nil, err
	***REMOVED***
	msgType := messageType(buf[0])

	// Check if the message is encrypted
	if msgType == encryptMsg ***REMOVED***
		if !m.config.EncryptionEnabled() ***REMOVED***
			return 0, nil, nil,
				fmt.Errorf("Remote state is encrypted and encryption is not configured")
		***REMOVED***

		plain, err := m.decryptRemoteState(bufConn)
		if err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***

		// Reset message type and bufConn
		msgType = messageType(plain[0])
		bufConn = bytes.NewReader(plain[1:])
	***REMOVED*** else if m.config.EncryptionEnabled() ***REMOVED***
		return 0, nil, nil,
			fmt.Errorf("Encryption is configured but remote state is not encrypted")
	***REMOVED***

	// Get the msgPack decoders
	hd := codec.MsgpackHandle***REMOVED******REMOVED***
	dec := codec.NewDecoder(bufConn, &hd)

	// Check if we have a compressed message
	if msgType == compressMsg ***REMOVED***
		var c compress
		if err := dec.Decode(&c); err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***
		decomp, err := decompressBuffer(&c)
		if err != nil ***REMOVED***
			return 0, nil, nil, err
		***REMOVED***

		// Reset the message type
		msgType = messageType(decomp[0])

		// Create a new bufConn
		bufConn = bytes.NewReader(decomp[1:])

		// Create a new decoder
		dec = codec.NewDecoder(bufConn, &hd)
	***REMOVED***

	return msgType, bufConn, dec, nil
***REMOVED***

// readRemoteState is used to read the remote state from a connection
func (m *Memberlist) readRemoteState(bufConn io.Reader, dec *codec.Decoder) (bool, []pushNodeState, []byte, error) ***REMOVED***
	// Read the push/pull header
	var header pushPullHeader
	if err := dec.Decode(&header); err != nil ***REMOVED***
		return false, nil, nil, err
	***REMOVED***

	// Allocate space for the transfer
	remoteNodes := make([]pushNodeState, header.Nodes)

	// Try to decode all the states
	for i := 0; i < header.Nodes; i++ ***REMOVED***
		if err := dec.Decode(&remoteNodes[i]); err != nil ***REMOVED***
			return false, nil, nil, err
		***REMOVED***
	***REMOVED***

	// Read the remote user state into a buffer
	var userBuf []byte
	if header.UserStateLen > 0 ***REMOVED***
		userBuf = make([]byte, header.UserStateLen)
		bytes, err := io.ReadAtLeast(bufConn, userBuf, header.UserStateLen)
		if err == nil && bytes != header.UserStateLen ***REMOVED***
			err = fmt.Errorf(
				"Failed to read full user state (%d / %d)",
				bytes, header.UserStateLen)
		***REMOVED***
		if err != nil ***REMOVED***
			return false, nil, nil, err
		***REMOVED***
	***REMOVED***

	// For proto versions < 2, there is no port provided. Mask old
	// behavior by using the configured port
	for idx := range remoteNodes ***REMOVED***
		if m.ProtocolVersion() < 2 || remoteNodes[idx].Port == 0 ***REMOVED***
			remoteNodes[idx].Port = uint16(m.config.BindPort)
		***REMOVED***
	***REMOVED***

	return header.Join, remoteNodes, userBuf, nil
***REMOVED***

// mergeRemoteState is used to merge the remote state with our local state
func (m *Memberlist) mergeRemoteState(join bool, remoteNodes []pushNodeState, userBuf []byte) error ***REMOVED***
	if err := m.verifyProtocol(remoteNodes); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Invoke the merge delegate if any
	if join && m.config.Merge != nil ***REMOVED***
		nodes := make([]*Node, len(remoteNodes))
		for idx, n := range remoteNodes ***REMOVED***
			nodes[idx] = &Node***REMOVED***
				Name: n.Name,
				Addr: n.Addr,
				Port: n.Port,
				Meta: n.Meta,
				PMin: n.Vsn[0],
				PMax: n.Vsn[1],
				PCur: n.Vsn[2],
				DMin: n.Vsn[3],
				DMax: n.Vsn[4],
				DCur: n.Vsn[5],
			***REMOVED***
		***REMOVED***
		if err := m.config.Merge.NotifyMerge(nodes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Merge the membership state
	m.mergeState(remoteNodes)

	// Invoke the delegate for user state
	if userBuf != nil && m.config.Delegate != nil ***REMOVED***
		m.config.Delegate.MergeRemoteState(userBuf, join)
	***REMOVED***
	return nil
***REMOVED***

// readUserMsg is used to decode a userMsg from a stream.
func (m *Memberlist) readUserMsg(bufConn io.Reader, dec *codec.Decoder) error ***REMOVED***
	// Read the user message header
	var header userMsgHeader
	if err := dec.Decode(&header); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Read the user message into a buffer
	var userBuf []byte
	if header.UserMsgLen > 0 ***REMOVED***
		userBuf = make([]byte, header.UserMsgLen)
		bytes, err := io.ReadAtLeast(bufConn, userBuf, header.UserMsgLen)
		if err == nil && bytes != header.UserMsgLen ***REMOVED***
			err = fmt.Errorf(
				"Failed to read full user message (%d / %d)",
				bytes, header.UserMsgLen)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		d := m.config.Delegate
		if d != nil ***REMOVED***
			d.NotifyMsg(userBuf)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// sendPingAndWaitForAck makes a stream connection to the given address, sends
// a ping, and waits for an ack. All of this is done as a series of blocking
// operations, given the deadline. The bool return parameter is true if we
// we able to round trip a ping to the other node.
func (m *Memberlist) sendPingAndWaitForAck(addr string, ping ping, deadline time.Time) (bool, error) ***REMOVED***
	conn, err := m.transport.DialTimeout(addr, m.config.TCPTimeout)
	if err != nil ***REMOVED***
		// If the node is actually dead we expect this to fail, so we
		// shouldn't spam the logs with it. After this point, errors
		// with the connection are real, unexpected errors and should
		// get propagated up.
		return false, nil
	***REMOVED***
	defer conn.Close()
	conn.SetDeadline(deadline)

	out, err := encode(pingMsg, &ping)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if err = m.rawSendMsgStream(conn, out.Bytes()); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	msgType, _, dec, err := m.readStream(conn)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if msgType != ackRespMsg ***REMOVED***
		return false, fmt.Errorf("Unexpected msgType (%d) from ping %s", msgType, LogConn(conn))
	***REMOVED***

	var ack ackResp
	if err = dec.Decode(&ack); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if ack.SeqNo != ping.SeqNo ***REMOVED***
		return false, fmt.Errorf("Sequence number from ack (%d) doesn't match ping (%d)", ack.SeqNo, ping.SeqNo, LogConn(conn))
	***REMOVED***

	return true, nil
***REMOVED***
