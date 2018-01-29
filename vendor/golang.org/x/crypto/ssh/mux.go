// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
)

// debugMux, if set, causes messages in the connection protocol to be
// logged.
const debugMux = false

// chanList is a thread safe channel list.
type chanList struct ***REMOVED***
	// protects concurrent access to chans
	sync.Mutex

	// chans are indexed by the local id of the channel, which the
	// other side should send in the PeersId field.
	chans []*channel

	// This is a debugging aid: it offsets all IDs by this
	// amount. This helps distinguish otherwise identical
	// server/client muxes
	offset uint32
***REMOVED***

// Assigns a channel ID to the given channel.
func (c *chanList) add(ch *channel) uint32 ***REMOVED***
	c.Lock()
	defer c.Unlock()
	for i := range c.chans ***REMOVED***
		if c.chans[i] == nil ***REMOVED***
			c.chans[i] = ch
			return uint32(i) + c.offset
		***REMOVED***
	***REMOVED***
	c.chans = append(c.chans, ch)
	return uint32(len(c.chans)-1) + c.offset
***REMOVED***

// getChan returns the channel for the given ID.
func (c *chanList) getChan(id uint32) *channel ***REMOVED***
	id -= c.offset

	c.Lock()
	defer c.Unlock()
	if id < uint32(len(c.chans)) ***REMOVED***
		return c.chans[id]
	***REMOVED***
	return nil
***REMOVED***

func (c *chanList) remove(id uint32) ***REMOVED***
	id -= c.offset
	c.Lock()
	if id < uint32(len(c.chans)) ***REMOVED***
		c.chans[id] = nil
	***REMOVED***
	c.Unlock()
***REMOVED***

// dropAll forgets all channels it knows, returning them in a slice.
func (c *chanList) dropAll() []*channel ***REMOVED***
	c.Lock()
	defer c.Unlock()
	var r []*channel

	for _, ch := range c.chans ***REMOVED***
		if ch == nil ***REMOVED***
			continue
		***REMOVED***
		r = append(r, ch)
	***REMOVED***
	c.chans = nil
	return r
***REMOVED***

// mux represents the state for the SSH connection protocol, which
// multiplexes many channels onto a single packet transport.
type mux struct ***REMOVED***
	conn     packetConn
	chanList chanList

	incomingChannels chan NewChannel

	globalSentMu     sync.Mutex
	globalResponses  chan interface***REMOVED******REMOVED***
	incomingRequests chan *Request

	errCond *sync.Cond
	err     error
***REMOVED***

// When debugging, each new chanList instantiation has a different
// offset.
var globalOff uint32

func (m *mux) Wait() error ***REMOVED***
	m.errCond.L.Lock()
	defer m.errCond.L.Unlock()
	for m.err == nil ***REMOVED***
		m.errCond.Wait()
	***REMOVED***
	return m.err
***REMOVED***

// newMux returns a mux that runs over the given connection.
func newMux(p packetConn) *mux ***REMOVED***
	m := &mux***REMOVED***
		conn:             p,
		incomingChannels: make(chan NewChannel, chanSize),
		globalResponses:  make(chan interface***REMOVED******REMOVED***, 1),
		incomingRequests: make(chan *Request, chanSize),
		errCond:          newCond(),
	***REMOVED***
	if debugMux ***REMOVED***
		m.chanList.offset = atomic.AddUint32(&globalOff, 1)
	***REMOVED***

	go m.loop()
	return m
***REMOVED***

func (m *mux) sendMessage(msg interface***REMOVED******REMOVED***) error ***REMOVED***
	p := Marshal(msg)
	if debugMux ***REMOVED***
		log.Printf("send global(%d): %#v", m.chanList.offset, msg)
	***REMOVED***
	return m.conn.writePacket(p)
***REMOVED***

func (m *mux) SendRequest(name string, wantReply bool, payload []byte) (bool, []byte, error) ***REMOVED***
	if wantReply ***REMOVED***
		m.globalSentMu.Lock()
		defer m.globalSentMu.Unlock()
	***REMOVED***

	if err := m.sendMessage(globalRequestMsg***REMOVED***
		Type:      name,
		WantReply: wantReply,
		Data:      payload,
	***REMOVED***); err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	if !wantReply ***REMOVED***
		return false, nil, nil
	***REMOVED***

	msg, ok := <-m.globalResponses
	if !ok ***REMOVED***
		return false, nil, io.EOF
	***REMOVED***
	switch msg := msg.(type) ***REMOVED***
	case *globalRequestFailureMsg:
		return false, msg.Data, nil
	case *globalRequestSuccessMsg:
		return true, msg.Data, nil
	default:
		return false, nil, fmt.Errorf("ssh: unexpected response to request: %#v", msg)
	***REMOVED***
***REMOVED***

// ackRequest must be called after processing a global request that
// has WantReply set.
func (m *mux) ackRequest(ok bool, data []byte) error ***REMOVED***
	if ok ***REMOVED***
		return m.sendMessage(globalRequestSuccessMsg***REMOVED***Data: data***REMOVED***)
	***REMOVED***
	return m.sendMessage(globalRequestFailureMsg***REMOVED***Data: data***REMOVED***)
***REMOVED***

func (m *mux) Close() error ***REMOVED***
	return m.conn.Close()
***REMOVED***

// loop runs the connection machine. It will process packets until an
// error is encountered. To synchronize on loop exit, use mux.Wait.
func (m *mux) loop() ***REMOVED***
	var err error
	for err == nil ***REMOVED***
		err = m.onePacket()
	***REMOVED***

	for _, ch := range m.chanList.dropAll() ***REMOVED***
		ch.close()
	***REMOVED***

	close(m.incomingChannels)
	close(m.incomingRequests)
	close(m.globalResponses)

	m.conn.Close()

	m.errCond.L.Lock()
	m.err = err
	m.errCond.Broadcast()
	m.errCond.L.Unlock()

	if debugMux ***REMOVED***
		log.Println("loop exit", err)
	***REMOVED***
***REMOVED***

// onePacket reads and processes one packet.
func (m *mux) onePacket() error ***REMOVED***
	packet, err := m.conn.readPacket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if debugMux ***REMOVED***
		if packet[0] == msgChannelData || packet[0] == msgChannelExtendedData ***REMOVED***
			log.Printf("decoding(%d): data packet - %d bytes", m.chanList.offset, len(packet))
		***REMOVED*** else ***REMOVED***
			p, _ := decode(packet)
			log.Printf("decoding(%d): %d %#v - %d bytes", m.chanList.offset, packet[0], p, len(packet))
		***REMOVED***
	***REMOVED***

	switch packet[0] ***REMOVED***
	case msgChannelOpen:
		return m.handleChannelOpen(packet)
	case msgGlobalRequest, msgRequestSuccess, msgRequestFailure:
		return m.handleGlobalPacket(packet)
	***REMOVED***

	// assume a channel packet.
	if len(packet) < 5 ***REMOVED***
		return parseError(packet[0])
	***REMOVED***
	id := binary.BigEndian.Uint32(packet[1:])
	ch := m.chanList.getChan(id)
	if ch == nil ***REMOVED***
		return fmt.Errorf("ssh: invalid channel %d", id)
	***REMOVED***

	return ch.handlePacket(packet)
***REMOVED***

func (m *mux) handleGlobalPacket(packet []byte) error ***REMOVED***
	msg, err := decode(packet)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch msg := msg.(type) ***REMOVED***
	case *globalRequestMsg:
		m.incomingRequests <- &Request***REMOVED***
			Type:      msg.Type,
			WantReply: msg.WantReply,
			Payload:   msg.Data,
			mux:       m,
		***REMOVED***
	case *globalRequestSuccessMsg, *globalRequestFailureMsg:
		m.globalResponses <- msg
	default:
		panic(fmt.Sprintf("not a global message %#v", msg))
	***REMOVED***

	return nil
***REMOVED***

// handleChannelOpen schedules a channel to be Accept()ed.
func (m *mux) handleChannelOpen(packet []byte) error ***REMOVED***
	var msg channelOpenMsg
	if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
		return err
	***REMOVED***

	if msg.MaxPacketSize < minPacketLength || msg.MaxPacketSize > 1<<31 ***REMOVED***
		failMsg := channelOpenFailureMsg***REMOVED***
			PeersID:  msg.PeersID,
			Reason:   ConnectionFailed,
			Message:  "invalid request",
			Language: "en_US.UTF-8",
		***REMOVED***
		return m.sendMessage(failMsg)
	***REMOVED***

	c := m.newChannel(msg.ChanType, channelInbound, msg.TypeSpecificData)
	c.remoteId = msg.PeersID
	c.maxRemotePayload = msg.MaxPacketSize
	c.remoteWin.add(msg.PeersWindow)
	m.incomingChannels <- c
	return nil
***REMOVED***

func (m *mux) OpenChannel(chanType string, extra []byte) (Channel, <-chan *Request, error) ***REMOVED***
	ch, err := m.openChannel(chanType, extra)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return ch, ch.incomingRequests, nil
***REMOVED***

func (m *mux) openChannel(chanType string, extra []byte) (*channel, error) ***REMOVED***
	ch := m.newChannel(chanType, channelOutbound, extra)

	ch.maxIncomingPayload = channelMaxPacket

	open := channelOpenMsg***REMOVED***
		ChanType:         chanType,
		PeersWindow:      ch.myWindow,
		MaxPacketSize:    ch.maxIncomingPayload,
		TypeSpecificData: extra,
		PeersID:          ch.localId,
	***REMOVED***
	if err := m.sendMessage(open); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch msg := (<-ch.msg).(type) ***REMOVED***
	case *channelOpenConfirmMsg:
		return ch, nil
	case *channelOpenFailureMsg:
		return nil, &OpenChannelError***REMOVED***msg.Reason, msg.Message***REMOVED***
	default:
		return nil, fmt.Errorf("ssh: unexpected packet in response to channel open: %T", msg)
	***REMOVED***
***REMOVED***
