// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

const (
	minPacketLength = 9
	// channelMaxPacket contains the maximum number of bytes that will be
	// sent in a single packet. As per RFC 4253, section 6.1, 32k is also
	// the minimum.
	channelMaxPacket = 1 << 15
	// We follow OpenSSH here.
	channelWindowSize = 64 * channelMaxPacket
)

// NewChannel represents an incoming request to a channel. It must either be
// accepted for use by calling Accept, or rejected by calling Reject.
type NewChannel interface ***REMOVED***
	// Accept accepts the channel creation request. It returns the Channel
	// and a Go channel containing SSH requests. The Go channel must be
	// serviced otherwise the Channel will hang.
	Accept() (Channel, <-chan *Request, error)

	// Reject rejects the channel creation request. After calling
	// this, no other methods on the Channel may be called.
	Reject(reason RejectionReason, message string) error

	// ChannelType returns the type of the channel, as supplied by the
	// client.
	ChannelType() string

	// ExtraData returns the arbitrary payload for this channel, as supplied
	// by the client. This data is specific to the channel type.
	ExtraData() []byte
***REMOVED***

// A Channel is an ordered, reliable, flow-controlled, duplex stream
// that is multiplexed over an SSH connection.
type Channel interface ***REMOVED***
	// Read reads up to len(data) bytes from the channel.
	Read(data []byte) (int, error)

	// Write writes len(data) bytes to the channel.
	Write(data []byte) (int, error)

	// Close signals end of channel use. No data may be sent after this
	// call.
	Close() error

	// CloseWrite signals the end of sending in-band
	// data. Requests may still be sent, and the other side may
	// still send data
	CloseWrite() error

	// SendRequest sends a channel request.  If wantReply is true,
	// it will wait for a reply and return the result as a
	// boolean, otherwise the return value will be false. Channel
	// requests are out-of-band messages so they may be sent even
	// if the data stream is closed or blocked by flow control.
	// If the channel is closed before a reply is returned, io.EOF
	// is returned.
	SendRequest(name string, wantReply bool, payload []byte) (bool, error)

	// Stderr returns an io.ReadWriter that writes to this channel
	// with the extended data type set to stderr. Stderr may
	// safely be read and written from a different goroutine than
	// Read and Write respectively.
	Stderr() io.ReadWriter
***REMOVED***

// Request is a request sent outside of the normal stream of
// data. Requests can either be specific to an SSH channel, or they
// can be global.
type Request struct ***REMOVED***
	Type      string
	WantReply bool
	Payload   []byte

	ch  *channel
	mux *mux
***REMOVED***

// Reply sends a response to a request. It must be called for all requests
// where WantReply is true and is a no-op otherwise. The payload argument is
// ignored for replies to channel-specific requests.
func (r *Request) Reply(ok bool, payload []byte) error ***REMOVED***
	if !r.WantReply ***REMOVED***
		return nil
	***REMOVED***

	if r.ch == nil ***REMOVED***
		return r.mux.ackRequest(ok, payload)
	***REMOVED***

	return r.ch.ackRequest(ok)
***REMOVED***

// RejectionReason is an enumeration used when rejecting channel creation
// requests. See RFC 4254, section 5.1.
type RejectionReason uint32

const (
	Prohibited RejectionReason = iota + 1
	ConnectionFailed
	UnknownChannelType
	ResourceShortage
)

// String converts the rejection reason to human readable form.
func (r RejectionReason) String() string ***REMOVED***
	switch r ***REMOVED***
	case Prohibited:
		return "administratively prohibited"
	case ConnectionFailed:
		return "connect failed"
	case UnknownChannelType:
		return "unknown channel type"
	case ResourceShortage:
		return "resource shortage"
	***REMOVED***
	return fmt.Sprintf("unknown reason %d", int(r))
***REMOVED***

func min(a uint32, b int) uint32 ***REMOVED***
	if a < uint32(b) ***REMOVED***
		return a
	***REMOVED***
	return uint32(b)
***REMOVED***

type channelDirection uint8

const (
	channelInbound channelDirection = iota
	channelOutbound
)

// channel is an implementation of the Channel interface that works
// with the mux class.
type channel struct ***REMOVED***
	// R/O after creation
	chanType          string
	extraData         []byte
	localId, remoteId uint32

	// maxIncomingPayload and maxRemotePayload are the maximum
	// payload sizes of normal and extended data packets for
	// receiving and sending, respectively. The wire packet will
	// be 9 or 13 bytes larger (excluding encryption overhead).
	maxIncomingPayload uint32
	maxRemotePayload   uint32

	mux *mux

	// decided is set to true if an accept or reject message has been sent
	// (for outbound channels) or received (for inbound channels).
	decided bool

	// direction contains either channelOutbound, for channels created
	// locally, or channelInbound, for channels created by the peer.
	direction channelDirection

	// Pending internal channel messages.
	msg chan interface***REMOVED******REMOVED***

	// Since requests have no ID, there can be only one request
	// with WantReply=true outstanding.  This lock is held by a
	// goroutine that has such an outgoing request pending.
	sentRequestMu sync.Mutex

	incomingRequests chan *Request

	sentEOF bool

	// thread-safe data
	remoteWin  window
	pending    *buffer
	extPending *buffer

	// windowMu protects myWindow, the flow-control window.
	windowMu sync.Mutex
	myWindow uint32

	// writeMu serializes calls to mux.conn.writePacket() and
	// protects sentClose and packetPool. This mutex must be
	// different from windowMu, as writePacket can block if there
	// is a key exchange pending.
	writeMu   sync.Mutex
	sentClose bool

	// packetPool has a buffer for each extended channel ID to
	// save allocations during writes.
	packetPool map[uint32][]byte
***REMOVED***

// writePacket sends a packet. If the packet is a channel close, it updates
// sentClose. This method takes the lock c.writeMu.
func (ch *channel) writePacket(packet []byte) error ***REMOVED***
	ch.writeMu.Lock()
	if ch.sentClose ***REMOVED***
		ch.writeMu.Unlock()
		return io.EOF
	***REMOVED***
	ch.sentClose = (packet[0] == msgChannelClose)
	err := ch.mux.conn.writePacket(packet)
	ch.writeMu.Unlock()
	return err
***REMOVED***

func (ch *channel) sendMessage(msg interface***REMOVED******REMOVED***) error ***REMOVED***
	if debugMux ***REMOVED***
		log.Printf("send(%d): %#v", ch.mux.chanList.offset, msg)
	***REMOVED***

	p := Marshal(msg)
	binary.BigEndian.PutUint32(p[1:], ch.remoteId)
	return ch.writePacket(p)
***REMOVED***

// WriteExtended writes data to a specific extended stream. These streams are
// used, for example, for stderr.
func (ch *channel) WriteExtended(data []byte, extendedCode uint32) (n int, err error) ***REMOVED***
	if ch.sentEOF ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	// 1 byte message type, 4 bytes remoteId, 4 bytes data length
	opCode := byte(msgChannelData)
	headerLength := uint32(9)
	if extendedCode > 0 ***REMOVED***
		headerLength += 4
		opCode = msgChannelExtendedData
	***REMOVED***

	ch.writeMu.Lock()
	packet := ch.packetPool[extendedCode]
	// We don't remove the buffer from packetPool, so
	// WriteExtended calls from different goroutines will be
	// flagged as errors by the race detector.
	ch.writeMu.Unlock()

	for len(data) > 0 ***REMOVED***
		space := min(ch.maxRemotePayload, len(data))
		if space, err = ch.remoteWin.reserve(space); err != nil ***REMOVED***
			return n, err
		***REMOVED***
		if want := headerLength + space; uint32(cap(packet)) < want ***REMOVED***
			packet = make([]byte, want)
		***REMOVED*** else ***REMOVED***
			packet = packet[:want]
		***REMOVED***

		todo := data[:space]

		packet[0] = opCode
		binary.BigEndian.PutUint32(packet[1:], ch.remoteId)
		if extendedCode > 0 ***REMOVED***
			binary.BigEndian.PutUint32(packet[5:], uint32(extendedCode))
		***REMOVED***
		binary.BigEndian.PutUint32(packet[headerLength-4:], uint32(len(todo)))
		copy(packet[headerLength:], todo)
		if err = ch.writePacket(packet); err != nil ***REMOVED***
			return n, err
		***REMOVED***

		n += len(todo)
		data = data[len(todo):]
	***REMOVED***

	ch.writeMu.Lock()
	ch.packetPool[extendedCode] = packet
	ch.writeMu.Unlock()

	return n, err
***REMOVED***

func (ch *channel) handleData(packet []byte) error ***REMOVED***
	headerLen := 9
	isExtendedData := packet[0] == msgChannelExtendedData
	if isExtendedData ***REMOVED***
		headerLen = 13
	***REMOVED***
	if len(packet) < headerLen ***REMOVED***
		// malformed data packet
		return parseError(packet[0])
	***REMOVED***

	var extended uint32
	if isExtendedData ***REMOVED***
		extended = binary.BigEndian.Uint32(packet[5:])
	***REMOVED***

	length := binary.BigEndian.Uint32(packet[headerLen-4 : headerLen])
	if length == 0 ***REMOVED***
		return nil
	***REMOVED***
	if length > ch.maxIncomingPayload ***REMOVED***
		// TODO(hanwen): should send Disconnect?
		return errors.New("ssh: incoming packet exceeds maximum payload size")
	***REMOVED***

	data := packet[headerLen:]
	if length != uint32(len(data)) ***REMOVED***
		return errors.New("ssh: wrong packet length")
	***REMOVED***

	ch.windowMu.Lock()
	if ch.myWindow < length ***REMOVED***
		ch.windowMu.Unlock()
		// TODO(hanwen): should send Disconnect with reason?
		return errors.New("ssh: remote side wrote too much")
	***REMOVED***
	ch.myWindow -= length
	ch.windowMu.Unlock()

	if extended == 1 ***REMOVED***
		ch.extPending.write(data)
	***REMOVED*** else if extended > 0 ***REMOVED***
		// discard other extended data.
	***REMOVED*** else ***REMOVED***
		ch.pending.write(data)
	***REMOVED***
	return nil
***REMOVED***

func (c *channel) adjustWindow(n uint32) error ***REMOVED***
	c.windowMu.Lock()
	// Since myWindow is managed on our side, and can never exceed
	// the initial window setting, we don't worry about overflow.
	c.myWindow += uint32(n)
	c.windowMu.Unlock()
	return c.sendMessage(windowAdjustMsg***REMOVED***
		AdditionalBytes: uint32(n),
	***REMOVED***)
***REMOVED***

func (c *channel) ReadExtended(data []byte, extended uint32) (n int, err error) ***REMOVED***
	switch extended ***REMOVED***
	case 1:
		n, err = c.extPending.Read(data)
	case 0:
		n, err = c.pending.Read(data)
	default:
		return 0, fmt.Errorf("ssh: extended code %d unimplemented", extended)
	***REMOVED***

	if n > 0 ***REMOVED***
		err = c.adjustWindow(uint32(n))
		// sendWindowAdjust can return io.EOF if the remote
		// peer has closed the connection, however we want to
		// defer forwarding io.EOF to the caller of Read until
		// the buffer has been drained.
		if n > 0 && err == io.EOF ***REMOVED***
			err = nil
		***REMOVED***
	***REMOVED***

	return n, err
***REMOVED***

func (c *channel) close() ***REMOVED***
	c.pending.eof()
	c.extPending.eof()
	close(c.msg)
	close(c.incomingRequests)
	c.writeMu.Lock()
	// This is not necessary for a normal channel teardown, but if
	// there was another error, it is.
	c.sentClose = true
	c.writeMu.Unlock()
	// Unblock writers.
	c.remoteWin.close()
***REMOVED***

// responseMessageReceived is called when a success or failure message is
// received on a channel to check that such a message is reasonable for the
// given channel.
func (ch *channel) responseMessageReceived() error ***REMOVED***
	if ch.direction == channelInbound ***REMOVED***
		return errors.New("ssh: channel response message received on inbound channel")
	***REMOVED***
	if ch.decided ***REMOVED***
		return errors.New("ssh: duplicate response received for channel")
	***REMOVED***
	ch.decided = true
	return nil
***REMOVED***

func (ch *channel) handlePacket(packet []byte) error ***REMOVED***
	switch packet[0] ***REMOVED***
	case msgChannelData, msgChannelExtendedData:
		return ch.handleData(packet)
	case msgChannelClose:
		ch.sendMessage(channelCloseMsg***REMOVED***PeersID: ch.remoteId***REMOVED***)
		ch.mux.chanList.remove(ch.localId)
		ch.close()
		return nil
	case msgChannelEOF:
		// RFC 4254 is mute on how EOF affects dataExt messages but
		// it is logical to signal EOF at the same time.
		ch.extPending.eof()
		ch.pending.eof()
		return nil
	***REMOVED***

	decoded, err := decode(packet)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch msg := decoded.(type) ***REMOVED***
	case *channelOpenFailureMsg:
		if err := ch.responseMessageReceived(); err != nil ***REMOVED***
			return err
		***REMOVED***
		ch.mux.chanList.remove(msg.PeersID)
		ch.msg <- msg
	case *channelOpenConfirmMsg:
		if err := ch.responseMessageReceived(); err != nil ***REMOVED***
			return err
		***REMOVED***
		if msg.MaxPacketSize < minPacketLength || msg.MaxPacketSize > 1<<31 ***REMOVED***
			return fmt.Errorf("ssh: invalid MaxPacketSize %d from peer", msg.MaxPacketSize)
		***REMOVED***
		ch.remoteId = msg.MyID
		ch.maxRemotePayload = msg.MaxPacketSize
		ch.remoteWin.add(msg.MyWindow)
		ch.msg <- msg
	case *windowAdjustMsg:
		if !ch.remoteWin.add(msg.AdditionalBytes) ***REMOVED***
			return fmt.Errorf("ssh: invalid window update for %d bytes", msg.AdditionalBytes)
		***REMOVED***
	case *channelRequestMsg:
		req := Request***REMOVED***
			Type:      msg.Request,
			WantReply: msg.WantReply,
			Payload:   msg.RequestSpecificData,
			ch:        ch,
		***REMOVED***

		ch.incomingRequests <- &req
	default:
		ch.msg <- msg
	***REMOVED***
	return nil
***REMOVED***

func (m *mux) newChannel(chanType string, direction channelDirection, extraData []byte) *channel ***REMOVED***
	ch := &channel***REMOVED***
		remoteWin:        window***REMOVED***Cond: newCond()***REMOVED***,
		myWindow:         channelWindowSize,
		pending:          newBuffer(),
		extPending:       newBuffer(),
		direction:        direction,
		incomingRequests: make(chan *Request, chanSize),
		msg:              make(chan interface***REMOVED******REMOVED***, chanSize),
		chanType:         chanType,
		extraData:        extraData,
		mux:              m,
		packetPool:       make(map[uint32][]byte),
	***REMOVED***
	ch.localId = m.chanList.add(ch)
	return ch
***REMOVED***

var errUndecided = errors.New("ssh: must Accept or Reject channel")
var errDecidedAlready = errors.New("ssh: can call Accept or Reject only once")

type extChannel struct ***REMOVED***
	code uint32
	ch   *channel
***REMOVED***

func (e *extChannel) Write(data []byte) (n int, err error) ***REMOVED***
	return e.ch.WriteExtended(data, e.code)
***REMOVED***

func (e *extChannel) Read(data []byte) (n int, err error) ***REMOVED***
	return e.ch.ReadExtended(data, e.code)
***REMOVED***

func (ch *channel) Accept() (Channel, <-chan *Request, error) ***REMOVED***
	if ch.decided ***REMOVED***
		return nil, nil, errDecidedAlready
	***REMOVED***
	ch.maxIncomingPayload = channelMaxPacket
	confirm := channelOpenConfirmMsg***REMOVED***
		PeersID:       ch.remoteId,
		MyID:          ch.localId,
		MyWindow:      ch.myWindow,
		MaxPacketSize: ch.maxIncomingPayload,
	***REMOVED***
	ch.decided = true
	if err := ch.sendMessage(confirm); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return ch, ch.incomingRequests, nil
***REMOVED***

func (ch *channel) Reject(reason RejectionReason, message string) error ***REMOVED***
	if ch.decided ***REMOVED***
		return errDecidedAlready
	***REMOVED***
	reject := channelOpenFailureMsg***REMOVED***
		PeersID:  ch.remoteId,
		Reason:   reason,
		Message:  message,
		Language: "en",
	***REMOVED***
	ch.decided = true
	return ch.sendMessage(reject)
***REMOVED***

func (ch *channel) Read(data []byte) (int, error) ***REMOVED***
	if !ch.decided ***REMOVED***
		return 0, errUndecided
	***REMOVED***
	return ch.ReadExtended(data, 0)
***REMOVED***

func (ch *channel) Write(data []byte) (int, error) ***REMOVED***
	if !ch.decided ***REMOVED***
		return 0, errUndecided
	***REMOVED***
	return ch.WriteExtended(data, 0)
***REMOVED***

func (ch *channel) CloseWrite() error ***REMOVED***
	if !ch.decided ***REMOVED***
		return errUndecided
	***REMOVED***
	ch.sentEOF = true
	return ch.sendMessage(channelEOFMsg***REMOVED***
		PeersID: ch.remoteId***REMOVED***)
***REMOVED***

func (ch *channel) Close() error ***REMOVED***
	if !ch.decided ***REMOVED***
		return errUndecided
	***REMOVED***

	return ch.sendMessage(channelCloseMsg***REMOVED***
		PeersID: ch.remoteId***REMOVED***)
***REMOVED***

// Extended returns an io.ReadWriter that sends and receives data on the given,
// SSH extended stream. Such streams are used, for example, for stderr.
func (ch *channel) Extended(code uint32) io.ReadWriter ***REMOVED***
	if !ch.decided ***REMOVED***
		return nil
	***REMOVED***
	return &extChannel***REMOVED***code, ch***REMOVED***
***REMOVED***

func (ch *channel) Stderr() io.ReadWriter ***REMOVED***
	return ch.Extended(1)
***REMOVED***

func (ch *channel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) ***REMOVED***
	if !ch.decided ***REMOVED***
		return false, errUndecided
	***REMOVED***

	if wantReply ***REMOVED***
		ch.sentRequestMu.Lock()
		defer ch.sentRequestMu.Unlock()
	***REMOVED***

	msg := channelRequestMsg***REMOVED***
		PeersID:             ch.remoteId,
		Request:             name,
		WantReply:           wantReply,
		RequestSpecificData: payload,
	***REMOVED***

	if err := ch.sendMessage(msg); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if wantReply ***REMOVED***
		m, ok := (<-ch.msg)
		if !ok ***REMOVED***
			return false, io.EOF
		***REMOVED***
		switch m.(type) ***REMOVED***
		case *channelRequestFailureMsg:
			return false, nil
		case *channelRequestSuccessMsg:
			return true, nil
		default:
			return false, fmt.Errorf("ssh: unexpected response to channel request: %#v", m)
		***REMOVED***
	***REMOVED***

	return false, nil
***REMOVED***

// ackRequest either sends an ack or nack to the channel request.
func (ch *channel) ackRequest(ok bool) error ***REMOVED***
	if !ch.decided ***REMOVED***
		return errUndecided
	***REMOVED***

	var msg interface***REMOVED******REMOVED***
	if !ok ***REMOVED***
		msg = channelRequestFailureMsg***REMOVED***
			PeersID: ch.remoteId,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		msg = channelRequestSuccessMsg***REMOVED***
			PeersID: ch.remoteId,
		***REMOVED***
	***REMOVED***
	return ch.sendMessage(msg)
***REMOVED***

func (ch *channel) ChannelType() string ***REMOVED***
	return ch.chanType
***REMOVED***

func (ch *channel) ExtraData() []byte ***REMOVED***
	return ch.extraData
***REMOVED***
