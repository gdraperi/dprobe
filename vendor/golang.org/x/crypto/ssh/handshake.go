// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// debugHandshake, if set, prints messages sent and received.  Key
// exchange messages are printed as if DH were used, so the debug
// messages are wrong when using ECDH.
const debugHandshake = false

// chanSize sets the amount of buffering SSH connections. This is
// primarily for testing: setting chanSize=0 uncovers deadlocks more
// quickly.
const chanSize = 16

// keyingTransport is a packet based transport that supports key
// changes. It need not be thread-safe. It should pass through
// msgNewKeys in both directions.
type keyingTransport interface ***REMOVED***
	packetConn

	// prepareKeyChange sets up a key change. The key change for a
	// direction will be effected if a msgNewKeys message is sent
	// or received.
	prepareKeyChange(*algorithms, *kexResult) error
***REMOVED***

// handshakeTransport implements rekeying on top of a keyingTransport
// and offers a thread-safe writePacket() interface.
type handshakeTransport struct ***REMOVED***
	conn   keyingTransport
	config *Config

	serverVersion []byte
	clientVersion []byte

	// hostKeys is non-empty if we are the server. In that case,
	// it contains all host keys that can be used to sign the
	// connection.
	hostKeys []Signer

	// hostKeyAlgorithms is non-empty if we are the client. In that case,
	// we accept these key types from the server as host key.
	hostKeyAlgorithms []string

	// On read error, incoming is closed, and readError is set.
	incoming  chan []byte
	readError error

	mu             sync.Mutex
	writeError     error
	sentInitPacket []byte
	sentInitMsg    *kexInitMsg
	pendingPackets [][]byte // Used when a key exchange is in progress.

	// If the read loop wants to schedule a kex, it pings this
	// channel, and the write loop will send out a kex
	// message.
	requestKex chan struct***REMOVED******REMOVED***

	// If the other side requests or confirms a kex, its kexInit
	// packet is sent here for the write loop to find it.
	startKex chan *pendingKex

	// data for host key checking
	hostKeyCallback HostKeyCallback
	dialAddress     string
	remoteAddr      net.Addr

	// bannerCallback is non-empty if we are the client and it has been set in
	// ClientConfig. In that case it is called during the user authentication
	// dance to handle a custom server's message.
	bannerCallback BannerCallback

	// Algorithms agreed in the last key exchange.
	algorithms *algorithms

	readPacketsLeft uint32
	readBytesLeft   int64

	writePacketsLeft uint32
	writeBytesLeft   int64

	// The session ID or nil if first kex did not complete yet.
	sessionID []byte
***REMOVED***

type pendingKex struct ***REMOVED***
	otherInit []byte
	done      chan error
***REMOVED***

func newHandshakeTransport(conn keyingTransport, config *Config, clientVersion, serverVersion []byte) *handshakeTransport ***REMOVED***
	t := &handshakeTransport***REMOVED***
		conn:          conn,
		serverVersion: serverVersion,
		clientVersion: clientVersion,
		incoming:      make(chan []byte, chanSize),
		requestKex:    make(chan struct***REMOVED******REMOVED***, 1),
		startKex:      make(chan *pendingKex, 1),

		config: config,
	***REMOVED***
	t.resetReadThresholds()
	t.resetWriteThresholds()

	// We always start with a mandatory key exchange.
	t.requestKex <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
	return t
***REMOVED***

func newClientTransport(conn keyingTransport, clientVersion, serverVersion []byte, config *ClientConfig, dialAddr string, addr net.Addr) *handshakeTransport ***REMOVED***
	t := newHandshakeTransport(conn, &config.Config, clientVersion, serverVersion)
	t.dialAddress = dialAddr
	t.remoteAddr = addr
	t.hostKeyCallback = config.HostKeyCallback
	t.bannerCallback = config.BannerCallback
	if config.HostKeyAlgorithms != nil ***REMOVED***
		t.hostKeyAlgorithms = config.HostKeyAlgorithms
	***REMOVED*** else ***REMOVED***
		t.hostKeyAlgorithms = supportedHostKeyAlgos
	***REMOVED***
	go t.readLoop()
	go t.kexLoop()
	return t
***REMOVED***

func newServerTransport(conn keyingTransport, clientVersion, serverVersion []byte, config *ServerConfig) *handshakeTransport ***REMOVED***
	t := newHandshakeTransport(conn, &config.Config, clientVersion, serverVersion)
	t.hostKeys = config.hostKeys
	go t.readLoop()
	go t.kexLoop()
	return t
***REMOVED***

func (t *handshakeTransport) getSessionID() []byte ***REMOVED***
	return t.sessionID
***REMOVED***

// waitSession waits for the session to be established. This should be
// the first thing to call after instantiating handshakeTransport.
func (t *handshakeTransport) waitSession() error ***REMOVED***
	p, err := t.readPacket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if p[0] != msgNewKeys ***REMOVED***
		return fmt.Errorf("ssh: first packet should be msgNewKeys")
	***REMOVED***

	return nil
***REMOVED***

func (t *handshakeTransport) id() string ***REMOVED***
	if len(t.hostKeys) > 0 ***REMOVED***
		return "server"
	***REMOVED***
	return "client"
***REMOVED***

func (t *handshakeTransport) printPacket(p []byte, write bool) ***REMOVED***
	action := "got"
	if write ***REMOVED***
		action = "sent"
	***REMOVED***

	if p[0] == msgChannelData || p[0] == msgChannelExtendedData ***REMOVED***
		log.Printf("%s %s data (packet %d bytes)", t.id(), action, len(p))
	***REMOVED*** else ***REMOVED***
		msg, err := decode(p)
		log.Printf("%s %s %T %v (%v)", t.id(), action, msg, msg, err)
	***REMOVED***
***REMOVED***

func (t *handshakeTransport) readPacket() ([]byte, error) ***REMOVED***
	p, ok := <-t.incoming
	if !ok ***REMOVED***
		return nil, t.readError
	***REMOVED***
	return p, nil
***REMOVED***

func (t *handshakeTransport) readLoop() ***REMOVED***
	first := true
	for ***REMOVED***
		p, err := t.readOnePacket(first)
		first = false
		if err != nil ***REMOVED***
			t.readError = err
			close(t.incoming)
			break
		***REMOVED***
		if p[0] == msgIgnore || p[0] == msgDebug ***REMOVED***
			continue
		***REMOVED***
		t.incoming <- p
	***REMOVED***

	// Stop writers too.
	t.recordWriteError(t.readError)

	// Unblock the writer should it wait for this.
	close(t.startKex)

	// Don't close t.requestKex; it's also written to from writePacket.
***REMOVED***

func (t *handshakeTransport) pushPacket(p []byte) error ***REMOVED***
	if debugHandshake ***REMOVED***
		t.printPacket(p, true)
	***REMOVED***
	return t.conn.writePacket(p)
***REMOVED***

func (t *handshakeTransport) getWriteError() error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.writeError
***REMOVED***

func (t *handshakeTransport) recordWriteError(err error) ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.writeError == nil && err != nil ***REMOVED***
		t.writeError = err
	***REMOVED***
***REMOVED***

func (t *handshakeTransport) requestKeyExchange() ***REMOVED***
	select ***REMOVED***
	case t.requestKex <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
		// something already requested a kex, so do nothing.
	***REMOVED***
***REMOVED***

func (t *handshakeTransport) resetWriteThresholds() ***REMOVED***
	t.writePacketsLeft = packetRekeyThreshold
	if t.config.RekeyThreshold > 0 ***REMOVED***
		t.writeBytesLeft = int64(t.config.RekeyThreshold)
	***REMOVED*** else if t.algorithms != nil ***REMOVED***
		t.writeBytesLeft = t.algorithms.w.rekeyBytes()
	***REMOVED*** else ***REMOVED***
		t.writeBytesLeft = 1 << 30
	***REMOVED***
***REMOVED***

func (t *handshakeTransport) kexLoop() ***REMOVED***

write:
	for t.getWriteError() == nil ***REMOVED***
		var request *pendingKex
		var sent bool

		for request == nil || !sent ***REMOVED***
			var ok bool
			select ***REMOVED***
			case request, ok = <-t.startKex:
				if !ok ***REMOVED***
					break write
				***REMOVED***
			case <-t.requestKex:
				break
			***REMOVED***

			if !sent ***REMOVED***
				if err := t.sendKexInit(); err != nil ***REMOVED***
					t.recordWriteError(err)
					break
				***REMOVED***
				sent = true
			***REMOVED***
		***REMOVED***

		if err := t.getWriteError(); err != nil ***REMOVED***
			if request != nil ***REMOVED***
				request.done <- err
			***REMOVED***
			break
		***REMOVED***

		// We're not servicing t.requestKex, but that is OK:
		// we never block on sending to t.requestKex.

		// We're not servicing t.startKex, but the remote end
		// has just sent us a kexInitMsg, so it can't send
		// another key change request, until we close the done
		// channel on the pendingKex request.

		err := t.enterKeyExchange(request.otherInit)

		t.mu.Lock()
		t.writeError = err
		t.sentInitPacket = nil
		t.sentInitMsg = nil

		t.resetWriteThresholds()

		// we have completed the key exchange. Since the
		// reader is still blocked, it is safe to clear out
		// the requestKex channel. This avoids the situation
		// where: 1) we consumed our own request for the
		// initial kex, and 2) the kex from the remote side
		// caused another send on the requestKex channel,
	clear:
		for ***REMOVED***
			select ***REMOVED***
			case <-t.requestKex:
				//
			default:
				break clear
			***REMOVED***
		***REMOVED***

		request.done <- t.writeError

		// kex finished. Push packets that we received while
		// the kex was in progress. Don't look at t.startKex
		// and don't increment writtenSinceKex: if we trigger
		// another kex while we are still busy with the last
		// one, things will become very confusing.
		for _, p := range t.pendingPackets ***REMOVED***
			t.writeError = t.pushPacket(p)
			if t.writeError != nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		t.pendingPackets = t.pendingPackets[:0]
		t.mu.Unlock()
	***REMOVED***

	// drain startKex channel. We don't service t.requestKex
	// because nobody does blocking sends there.
	go func() ***REMOVED***
		for init := range t.startKex ***REMOVED***
			init.done <- t.writeError
		***REMOVED***
	***REMOVED***()

	// Unblock reader.
	t.conn.Close()
***REMOVED***

// The protocol uses uint32 for packet counters, so we can't let them
// reach 1<<32.  We will actually read and write more packets than
// this, though: the other side may send more packets, and after we
// hit this limit on writing we will send a few more packets for the
// key exchange itself.
const packetRekeyThreshold = (1 << 31)

func (t *handshakeTransport) resetReadThresholds() ***REMOVED***
	t.readPacketsLeft = packetRekeyThreshold
	if t.config.RekeyThreshold > 0 ***REMOVED***
		t.readBytesLeft = int64(t.config.RekeyThreshold)
	***REMOVED*** else if t.algorithms != nil ***REMOVED***
		t.readBytesLeft = t.algorithms.r.rekeyBytes()
	***REMOVED*** else ***REMOVED***
		t.readBytesLeft = 1 << 30
	***REMOVED***
***REMOVED***

func (t *handshakeTransport) readOnePacket(first bool) ([]byte, error) ***REMOVED***
	p, err := t.conn.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if t.readPacketsLeft > 0 ***REMOVED***
		t.readPacketsLeft--
	***REMOVED*** else ***REMOVED***
		t.requestKeyExchange()
	***REMOVED***

	if t.readBytesLeft > 0 ***REMOVED***
		t.readBytesLeft -= int64(len(p))
	***REMOVED*** else ***REMOVED***
		t.requestKeyExchange()
	***REMOVED***

	if debugHandshake ***REMOVED***
		t.printPacket(p, false)
	***REMOVED***

	if first && p[0] != msgKexInit ***REMOVED***
		return nil, fmt.Errorf("ssh: first packet should be msgKexInit")
	***REMOVED***

	if p[0] != msgKexInit ***REMOVED***
		return p, nil
	***REMOVED***

	firstKex := t.sessionID == nil

	kex := pendingKex***REMOVED***
		done:      make(chan error, 1),
		otherInit: p,
	***REMOVED***
	t.startKex <- &kex
	err = <-kex.done

	if debugHandshake ***REMOVED***
		log.Printf("%s exited key exchange (first %v), err %v", t.id(), firstKex, err)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	t.resetReadThresholds()

	// By default, a key exchange is hidden from higher layers by
	// translating it into msgIgnore.
	successPacket := []byte***REMOVED***msgIgnore***REMOVED***
	if firstKex ***REMOVED***
		// sendKexInit() for the first kex waits for
		// msgNewKeys so the authentication process is
		// guaranteed to happen over an encrypted transport.
		successPacket = []byte***REMOVED***msgNewKeys***REMOVED***
	***REMOVED***

	return successPacket, nil
***REMOVED***

// sendKexInit sends a key change message.
func (t *handshakeTransport) sendKexInit() error ***REMOVED***
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.sentInitMsg != nil ***REMOVED***
		// kexInits may be sent either in response to the other side,
		// or because our side wants to initiate a key change, so we
		// may have already sent a kexInit. In that case, don't send a
		// second kexInit.
		return nil
	***REMOVED***

	msg := &kexInitMsg***REMOVED***
		KexAlgos:                t.config.KeyExchanges,
		CiphersClientServer:     t.config.Ciphers,
		CiphersServerClient:     t.config.Ciphers,
		MACsClientServer:        t.config.MACs,
		MACsServerClient:        t.config.MACs,
		CompressionClientServer: supportedCompressions,
		CompressionServerClient: supportedCompressions,
	***REMOVED***
	io.ReadFull(rand.Reader, msg.Cookie[:])

	if len(t.hostKeys) > 0 ***REMOVED***
		for _, k := range t.hostKeys ***REMOVED***
			msg.ServerHostKeyAlgos = append(
				msg.ServerHostKeyAlgos, k.PublicKey().Type())
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		msg.ServerHostKeyAlgos = t.hostKeyAlgorithms
	***REMOVED***
	packet := Marshal(msg)

	// writePacket destroys the contents, so save a copy.
	packetCopy := make([]byte, len(packet))
	copy(packetCopy, packet)

	if err := t.pushPacket(packetCopy); err != nil ***REMOVED***
		return err
	***REMOVED***

	t.sentInitMsg = msg
	t.sentInitPacket = packet

	return nil
***REMOVED***

func (t *handshakeTransport) writePacket(p []byte) error ***REMOVED***
	switch p[0] ***REMOVED***
	case msgKexInit:
		return errors.New("ssh: only handshakeTransport can send kexInit")
	case msgNewKeys:
		return errors.New("ssh: only handshakeTransport can send newKeys")
	***REMOVED***

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.writeError != nil ***REMOVED***
		return t.writeError
	***REMOVED***

	if t.sentInitMsg != nil ***REMOVED***
		// Copy the packet so the writer can reuse the buffer.
		cp := make([]byte, len(p))
		copy(cp, p)
		t.pendingPackets = append(t.pendingPackets, cp)
		return nil
	***REMOVED***

	if t.writeBytesLeft > 0 ***REMOVED***
		t.writeBytesLeft -= int64(len(p))
	***REMOVED*** else ***REMOVED***
		t.requestKeyExchange()
	***REMOVED***

	if t.writePacketsLeft > 0 ***REMOVED***
		t.writePacketsLeft--
	***REMOVED*** else ***REMOVED***
		t.requestKeyExchange()
	***REMOVED***

	if err := t.pushPacket(p); err != nil ***REMOVED***
		t.writeError = err
	***REMOVED***

	return nil
***REMOVED***

func (t *handshakeTransport) Close() error ***REMOVED***
	return t.conn.Close()
***REMOVED***

func (t *handshakeTransport) enterKeyExchange(otherInitPacket []byte) error ***REMOVED***
	if debugHandshake ***REMOVED***
		log.Printf("%s entered key exchange", t.id())
	***REMOVED***

	otherInit := &kexInitMsg***REMOVED******REMOVED***
	if err := Unmarshal(otherInitPacket, otherInit); err != nil ***REMOVED***
		return err
	***REMOVED***

	magics := handshakeMagics***REMOVED***
		clientVersion: t.clientVersion,
		serverVersion: t.serverVersion,
		clientKexInit: otherInitPacket,
		serverKexInit: t.sentInitPacket,
	***REMOVED***

	clientInit := otherInit
	serverInit := t.sentInitMsg
	if len(t.hostKeys) == 0 ***REMOVED***
		clientInit, serverInit = serverInit, clientInit

		magics.clientKexInit = t.sentInitPacket
		magics.serverKexInit = otherInitPacket
	***REMOVED***

	var err error
	t.algorithms, err = findAgreedAlgorithms(clientInit, serverInit)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// We don't send FirstKexFollows, but we handle receiving it.
	//
	// RFC 4253 section 7 defines the kex and the agreement method for
	// first_kex_packet_follows. It states that the guessed packet
	// should be ignored if the "kex algorithm and/or the host
	// key algorithm is guessed wrong (server and client have
	// different preferred algorithm), or if any of the other
	// algorithms cannot be agreed upon". The other algorithms have
	// already been checked above so the kex algorithm and host key
	// algorithm are checked here.
	if otherInit.FirstKexFollows && (clientInit.KexAlgos[0] != serverInit.KexAlgos[0] || clientInit.ServerHostKeyAlgos[0] != serverInit.ServerHostKeyAlgos[0]) ***REMOVED***
		// other side sent a kex message for the wrong algorithm,
		// which we have to ignore.
		if _, err := t.conn.readPacket(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	kex, ok := kexAlgoMap[t.algorithms.kex]
	if !ok ***REMOVED***
		return fmt.Errorf("ssh: unexpected key exchange algorithm %v", t.algorithms.kex)
	***REMOVED***

	var result *kexResult
	if len(t.hostKeys) > 0 ***REMOVED***
		result, err = t.server(kex, t.algorithms, &magics)
	***REMOVED*** else ***REMOVED***
		result, err = t.client(kex, t.algorithms, &magics)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if t.sessionID == nil ***REMOVED***
		t.sessionID = result.H
	***REMOVED***
	result.SessionID = t.sessionID

	if err := t.conn.prepareKeyChange(t.algorithms, result); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = t.conn.writePacket([]byte***REMOVED***msgNewKeys***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	if packet, err := t.conn.readPacket(); err != nil ***REMOVED***
		return err
	***REMOVED*** else if packet[0] != msgNewKeys ***REMOVED***
		return unexpectedMessageError(msgNewKeys, packet[0])
	***REMOVED***

	return nil
***REMOVED***

func (t *handshakeTransport) server(kex kexAlgorithm, algs *algorithms, magics *handshakeMagics) (*kexResult, error) ***REMOVED***
	var hostKey Signer
	for _, k := range t.hostKeys ***REMOVED***
		if algs.hostKey == k.PublicKey().Type() ***REMOVED***
			hostKey = k
		***REMOVED***
	***REMOVED***

	r, err := kex.Server(t.conn, t.config.Rand, magics, hostKey)
	return r, err
***REMOVED***

func (t *handshakeTransport) client(kex kexAlgorithm, algs *algorithms, magics *handshakeMagics) (*kexResult, error) ***REMOVED***
	result, err := kex.Client(t.conn, t.config.Rand, magics)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hostKey, err := ParsePublicKey(result.HostKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := verifyHostKeySignature(hostKey, result); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = t.hostKeyCallback(t.dialAddress, t.remoteAddr, hostKey)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return result, nil
***REMOVED***
