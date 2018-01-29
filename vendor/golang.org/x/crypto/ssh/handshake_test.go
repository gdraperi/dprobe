// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type testChecker struct ***REMOVED***
	calls []string
***REMOVED***

func (t *testChecker) Check(dialAddr string, addr net.Addr, key PublicKey) error ***REMOVED***
	if dialAddr == "bad" ***REMOVED***
		return fmt.Errorf("dialAddr is bad")
	***REMOVED***

	if tcpAddr, ok := addr.(*net.TCPAddr); !ok || tcpAddr == nil ***REMOVED***
		return fmt.Errorf("testChecker: got %T want *net.TCPAddr", addr)
	***REMOVED***

	t.calls = append(t.calls, fmt.Sprintf("%s %v %s %x", dialAddr, addr, key.Type(), key.Marshal()))

	return nil
***REMOVED***

// netPipe is analogous to net.Pipe, but it uses a real net.Conn, and
// therefore is buffered (net.Pipe deadlocks if both sides start with
// a write.)
func netPipe() (net.Conn, net.Conn, error) ***REMOVED***
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		listener, err = net.Listen("tcp", "[::1]:0")
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***
	defer listener.Close()
	c1, err := net.Dial("tcp", listener.Addr().String())
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	c2, err := listener.Accept()
	if err != nil ***REMOVED***
		c1.Close()
		return nil, nil, err
	***REMOVED***

	return c1, c2, nil
***REMOVED***

// noiseTransport inserts ignore messages to check that the read loop
// and the key exchange filters out these messages.
type noiseTransport struct ***REMOVED***
	keyingTransport
***REMOVED***

func (t *noiseTransport) writePacket(p []byte) error ***REMOVED***
	ignore := []byte***REMOVED***msgIgnore***REMOVED***
	if err := t.keyingTransport.writePacket(ignore); err != nil ***REMOVED***
		return err
	***REMOVED***
	debug := []byte***REMOVED***msgDebug, 1, 2, 3***REMOVED***
	if err := t.keyingTransport.writePacket(debug); err != nil ***REMOVED***
		return err
	***REMOVED***

	return t.keyingTransport.writePacket(p)
***REMOVED***

func addNoiseTransport(t keyingTransport) keyingTransport ***REMOVED***
	return &noiseTransport***REMOVED***t***REMOVED***
***REMOVED***

// handshakePair creates two handshakeTransports connected with each
// other. If the noise argument is true, both transports will try to
// confuse the other side by sending ignore and debug messages.
func handshakePair(clientConf *ClientConfig, addr string, noise bool) (client *handshakeTransport, server *handshakeTransport, err error) ***REMOVED***
	a, b, err := netPipe()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var trC, trS keyingTransport

	trC = newTransport(a, rand.Reader, true)
	trS = newTransport(b, rand.Reader, false)
	if noise ***REMOVED***
		trC = addNoiseTransport(trC)
		trS = addNoiseTransport(trS)
	***REMOVED***
	clientConf.SetDefaults()

	v := []byte("version")
	client = newClientTransport(trC, v, v, clientConf, addr, a.RemoteAddr())

	serverConf := &ServerConfig***REMOVED******REMOVED***
	serverConf.AddHostKey(testSigners["ecdsa"])
	serverConf.AddHostKey(testSigners["rsa"])
	serverConf.SetDefaults()
	server = newServerTransport(trS, v, v, serverConf)

	if err := server.waitSession(); err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("server.waitSession: %v", err)
	***REMOVED***
	if err := client.waitSession(); err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("client.waitSession: %v", err)
	***REMOVED***

	return client, server, nil
***REMOVED***

func TestHandshakeBasic(t *testing.T) ***REMOVED***
	if runtime.GOOS == "plan9" ***REMOVED***
		t.Skip("see golang.org/issue/7237")
	***REMOVED***

	checker := &syncChecker***REMOVED***
		waitCall: make(chan int, 10),
		called:   make(chan int, 10),
	***REMOVED***

	checker.waitCall <- 1
	trC, trS, err := handshakePair(&ClientConfig***REMOVED***HostKeyCallback: checker.Check***REMOVED***, "addr", false)
	if err != nil ***REMOVED***
		t.Fatalf("handshakePair: %v", err)
	***REMOVED***

	defer trC.Close()
	defer trS.Close()

	// Let first kex complete normally.
	<-checker.called

	clientDone := make(chan int, 0)
	gotHalf := make(chan int, 0)
	const N = 20

	go func() ***REMOVED***
		defer close(clientDone)
		// Client writes a bunch of stuff, and does a key
		// change in the middle. This should not confuse the
		// handshake in progress. We do this twice, so we test
		// that the packet buffer is reset correctly.
		for i := 0; i < N; i++ ***REMOVED***
			p := []byte***REMOVED***msgRequestSuccess, byte(i)***REMOVED***
			if err := trC.writePacket(p); err != nil ***REMOVED***
				t.Fatalf("sendPacket: %v", err)
			***REMOVED***
			if (i % 10) == 5 ***REMOVED***
				<-gotHalf
				// halfway through, we request a key change.
				trC.requestKeyExchange()

				// Wait until we can be sure the key
				// change has really started before we
				// write more.
				<-checker.called
			***REMOVED***
			if (i % 10) == 7 ***REMOVED***
				// write some packets until the kex
				// completes, to test buffering of
				// packets.
				checker.waitCall <- 1
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Server checks that client messages come in cleanly
	i := 0
	err = nil
	for ; i < N; i++ ***REMOVED***
		var p []byte
		p, err = trS.readPacket()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if (i % 10) == 5 ***REMOVED***
			gotHalf <- 1
		***REMOVED***

		want := []byte***REMOVED***msgRequestSuccess, byte(i)***REMOVED***
		if bytes.Compare(p, want) != 0 ***REMOVED***
			t.Errorf("message %d: got %v, want %v", i, p, want)
		***REMOVED***
	***REMOVED***
	<-clientDone
	if err != nil && err != io.EOF ***REMOVED***
		t.Fatalf("server error: %v", err)
	***REMOVED***
	if i != N ***REMOVED***
		t.Errorf("received %d messages, want 10.", i)
	***REMOVED***

	close(checker.called)
	if _, ok := <-checker.called; ok ***REMOVED***
		// If all went well, we registered exactly 2 key changes: one
		// that establishes the session, and one that we requested
		// additionally.
		t.Fatalf("got another host key checks after 2 handshakes")
	***REMOVED***
***REMOVED***

func TestForceFirstKex(t *testing.T) ***REMOVED***
	// like handshakePair, but must access the keyingTransport.
	checker := &testChecker***REMOVED******REMOVED***
	clientConf := &ClientConfig***REMOVED***HostKeyCallback: checker.Check***REMOVED***
	a, b, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***

	var trC, trS keyingTransport

	trC = newTransport(a, rand.Reader, true)

	// This is the disallowed packet:
	trC.writePacket(Marshal(&serviceRequestMsg***REMOVED***serviceUserAuth***REMOVED***))

	// Rest of the setup.
	trS = newTransport(b, rand.Reader, false)
	clientConf.SetDefaults()

	v := []byte("version")
	client := newClientTransport(trC, v, v, clientConf, "addr", a.RemoteAddr())

	serverConf := &ServerConfig***REMOVED******REMOVED***
	serverConf.AddHostKey(testSigners["ecdsa"])
	serverConf.AddHostKey(testSigners["rsa"])
	serverConf.SetDefaults()
	server := newServerTransport(trS, v, v, serverConf)

	defer client.Close()
	defer server.Close()

	// We setup the initial key exchange, but the remote side
	// tries to send serviceRequestMsg in cleartext, which is
	// disallowed.

	if err := server.waitSession(); err == nil ***REMOVED***
		t.Errorf("server first kex init should reject unexpected packet")
	***REMOVED***
***REMOVED***

func TestHandshakeAutoRekeyWrite(t *testing.T) ***REMOVED***
	checker := &syncChecker***REMOVED***
		called:   make(chan int, 10),
		waitCall: nil,
	***REMOVED***
	clientConf := &ClientConfig***REMOVED***HostKeyCallback: checker.Check***REMOVED***
	clientConf.RekeyThreshold = 500
	trC, trS, err := handshakePair(clientConf, "addr", false)
	if err != nil ***REMOVED***
		t.Fatalf("handshakePair: %v", err)
	***REMOVED***
	defer trC.Close()
	defer trS.Close()

	input := make([]byte, 251)
	input[0] = msgRequestSuccess

	done := make(chan int, 1)
	const numPacket = 5
	go func() ***REMOVED***
		defer close(done)
		j := 0
		for ; j < numPacket; j++ ***REMOVED***
			if p, err := trS.readPacket(); err != nil ***REMOVED***
				break
			***REMOVED*** else if !bytes.Equal(input, p) ***REMOVED***
				t.Errorf("got packet type %d, want %d", p[0], input[0])
			***REMOVED***
		***REMOVED***

		if j != numPacket ***REMOVED***
			t.Errorf("got %d, want 5 messages", j)
		***REMOVED***
	***REMOVED***()

	<-checker.called

	for i := 0; i < numPacket; i++ ***REMOVED***
		p := make([]byte, len(input))
		copy(p, input)
		if err := trC.writePacket(p); err != nil ***REMOVED***
			t.Errorf("writePacket: %v", err)
		***REMOVED***
		if i == 2 ***REMOVED***
			// Make sure the kex is in progress.
			<-checker.called
		***REMOVED***

	***REMOVED***
	<-done
***REMOVED***

type syncChecker struct ***REMOVED***
	waitCall chan int
	called   chan int
***REMOVED***

func (c *syncChecker) Check(dialAddr string, addr net.Addr, key PublicKey) error ***REMOVED***
	c.called <- 1
	if c.waitCall != nil ***REMOVED***
		<-c.waitCall
	***REMOVED***
	return nil
***REMOVED***

func TestHandshakeAutoRekeyRead(t *testing.T) ***REMOVED***
	sync := &syncChecker***REMOVED***
		called:   make(chan int, 2),
		waitCall: nil,
	***REMOVED***
	clientConf := &ClientConfig***REMOVED***
		HostKeyCallback: sync.Check,
	***REMOVED***
	clientConf.RekeyThreshold = 500

	trC, trS, err := handshakePair(clientConf, "addr", false)
	if err != nil ***REMOVED***
		t.Fatalf("handshakePair: %v", err)
	***REMOVED***
	defer trC.Close()
	defer trS.Close()

	packet := make([]byte, 501)
	packet[0] = msgRequestSuccess
	if err := trS.writePacket(packet); err != nil ***REMOVED***
		t.Fatalf("writePacket: %v", err)
	***REMOVED***

	// While we read out the packet, a key change will be
	// initiated.
	done := make(chan int, 1)
	go func() ***REMOVED***
		defer close(done)
		if _, err := trC.readPacket(); err != nil ***REMOVED***
			t.Fatalf("readPacket(client): %v", err)
		***REMOVED***

	***REMOVED***()

	<-done
	<-sync.called
***REMOVED***

// errorKeyingTransport generates errors after a given number of
// read/write operations.
type errorKeyingTransport struct ***REMOVED***
	packetConn
	readLeft, writeLeft int
***REMOVED***

func (n *errorKeyingTransport) prepareKeyChange(*algorithms, *kexResult) error ***REMOVED***
	return nil
***REMOVED***

func (n *errorKeyingTransport) getSessionID() []byte ***REMOVED***
	return nil
***REMOVED***

func (n *errorKeyingTransport) writePacket(packet []byte) error ***REMOVED***
	if n.writeLeft == 0 ***REMOVED***
		n.Close()
		return errors.New("barf")
	***REMOVED***

	n.writeLeft--
	return n.packetConn.writePacket(packet)
***REMOVED***

func (n *errorKeyingTransport) readPacket() ([]byte, error) ***REMOVED***
	if n.readLeft == 0 ***REMOVED***
		n.Close()
		return nil, errors.New("barf")
	***REMOVED***

	n.readLeft--
	return n.packetConn.readPacket()
***REMOVED***

func TestHandshakeErrorHandlingRead(t *testing.T) ***REMOVED***
	for i := 0; i < 20; i++ ***REMOVED***
		testHandshakeErrorHandlingN(t, i, -1, false)
	***REMOVED***
***REMOVED***

func TestHandshakeErrorHandlingWrite(t *testing.T) ***REMOVED***
	for i := 0; i < 20; i++ ***REMOVED***
		testHandshakeErrorHandlingN(t, -1, i, false)
	***REMOVED***
***REMOVED***

func TestHandshakeErrorHandlingReadCoupled(t *testing.T) ***REMOVED***
	for i := 0; i < 20; i++ ***REMOVED***
		testHandshakeErrorHandlingN(t, i, -1, true)
	***REMOVED***
***REMOVED***

func TestHandshakeErrorHandlingWriteCoupled(t *testing.T) ***REMOVED***
	for i := 0; i < 20; i++ ***REMOVED***
		testHandshakeErrorHandlingN(t, -1, i, true)
	***REMOVED***
***REMOVED***

// testHandshakeErrorHandlingN runs handshakes, injecting errors. If
// handshakeTransport deadlocks, the go runtime will detect it and
// panic.
func testHandshakeErrorHandlingN(t *testing.T, readLimit, writeLimit int, coupled bool) ***REMOVED***
	msg := Marshal(&serviceRequestMsg***REMOVED***strings.Repeat("x", int(minRekeyThreshold)/4)***REMOVED***)

	a, b := memPipe()
	defer a.Close()
	defer b.Close()

	key := testSigners["ecdsa"]
	serverConf := Config***REMOVED***RekeyThreshold: minRekeyThreshold***REMOVED***
	serverConf.SetDefaults()
	serverConn := newHandshakeTransport(&errorKeyingTransport***REMOVED***a, readLimit, writeLimit***REMOVED***, &serverConf, []byte***REMOVED***'a'***REMOVED***, []byte***REMOVED***'b'***REMOVED***)
	serverConn.hostKeys = []Signer***REMOVED***key***REMOVED***
	go serverConn.readLoop()
	go serverConn.kexLoop()

	clientConf := Config***REMOVED***RekeyThreshold: 10 * minRekeyThreshold***REMOVED***
	clientConf.SetDefaults()
	clientConn := newHandshakeTransport(&errorKeyingTransport***REMOVED***b, -1, -1***REMOVED***, &clientConf, []byte***REMOVED***'a'***REMOVED***, []byte***REMOVED***'b'***REMOVED***)
	clientConn.hostKeyAlgorithms = []string***REMOVED***key.PublicKey().Type()***REMOVED***
	clientConn.hostKeyCallback = InsecureIgnoreHostKey()
	go clientConn.readLoop()
	go clientConn.kexLoop()

	var wg sync.WaitGroup

	for _, hs := range []packetConn***REMOVED***serverConn, clientConn***REMOVED*** ***REMOVED***
		if !coupled ***REMOVED***
			wg.Add(2)
			go func(c packetConn) ***REMOVED***
				for i := 0; ; i++ ***REMOVED***
					str := fmt.Sprintf("%08x", i) + strings.Repeat("x", int(minRekeyThreshold)/4-8)
					err := c.writePacket(Marshal(&serviceRequestMsg***REMOVED***str***REMOVED***))
					if err != nil ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				wg.Done()
				c.Close()
			***REMOVED***(hs)
			go func(c packetConn) ***REMOVED***
				for ***REMOVED***
					_, err := c.readPacket()
					if err != nil ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				wg.Done()
			***REMOVED***(hs)
		***REMOVED*** else ***REMOVED***
			wg.Add(1)
			go func(c packetConn) ***REMOVED***
				for ***REMOVED***
					_, err := c.readPacket()
					if err != nil ***REMOVED***
						break
					***REMOVED***
					if err := c.writePacket(msg); err != nil ***REMOVED***
						break
					***REMOVED***

				***REMOVED***
				wg.Done()
			***REMOVED***(hs)
		***REMOVED***
	***REMOVED***
	wg.Wait()
***REMOVED***

func TestDisconnect(t *testing.T) ***REMOVED***
	if runtime.GOOS == "plan9" ***REMOVED***
		t.Skip("see golang.org/issue/7237")
	***REMOVED***
	checker := &testChecker***REMOVED******REMOVED***
	trC, trS, err := handshakePair(&ClientConfig***REMOVED***HostKeyCallback: checker.Check***REMOVED***, "addr", false)
	if err != nil ***REMOVED***
		t.Fatalf("handshakePair: %v", err)
	***REMOVED***

	defer trC.Close()
	defer trS.Close()

	trC.writePacket([]byte***REMOVED***msgRequestSuccess, 0, 0***REMOVED***)
	errMsg := &disconnectMsg***REMOVED***
		Reason:  42,
		Message: "such is life",
	***REMOVED***
	trC.writePacket(Marshal(errMsg))
	trC.writePacket([]byte***REMOVED***msgRequestSuccess, 0, 0***REMOVED***)

	packet, err := trS.readPacket()
	if err != nil ***REMOVED***
		t.Fatalf("readPacket 1: %v", err)
	***REMOVED***
	if packet[0] != msgRequestSuccess ***REMOVED***
		t.Errorf("got packet %v, want packet type %d", packet, msgRequestSuccess)
	***REMOVED***

	_, err = trS.readPacket()
	if err == nil ***REMOVED***
		t.Errorf("readPacket 2 succeeded")
	***REMOVED*** else if !reflect.DeepEqual(err, errMsg) ***REMOVED***
		t.Errorf("got error %#v, want %#v", err, errMsg)
	***REMOVED***

	_, err = trS.readPacket()
	if err == nil ***REMOVED***
		t.Errorf("readPacket 3 succeeded")
	***REMOVED***
***REMOVED***

func TestHandshakeRekeyDefault(t *testing.T) ***REMOVED***
	clientConf := &ClientConfig***REMOVED***
		Config: Config***REMOVED***
			Ciphers: []string***REMOVED***"aes128-ctr"***REMOVED***,
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	trC, trS, err := handshakePair(clientConf, "addr", false)
	if err != nil ***REMOVED***
		t.Fatalf("handshakePair: %v", err)
	***REMOVED***
	defer trC.Close()
	defer trS.Close()

	trC.writePacket([]byte***REMOVED***msgRequestSuccess, 0, 0***REMOVED***)
	trC.Close()

	rgb := (1024 + trC.readBytesLeft) >> 30
	wgb := (1024 + trC.writeBytesLeft) >> 30

	if rgb != 64 ***REMOVED***
		t.Errorf("got rekey after %dG read, want 64G", rgb)
	***REMOVED***
	if wgb != 64 ***REMOVED***
		t.Errorf("got rekey after %dG write, want 64G", wgb)
	***REMOVED***
***REMOVED***
