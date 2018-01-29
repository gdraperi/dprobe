// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"io"
	"io/ioutil"
	"sync"
	"testing"
)

func muxPair() (*mux, *mux) ***REMOVED***
	a, b := memPipe()

	s := newMux(a)
	c := newMux(b)

	return s, c
***REMOVED***

// Returns both ends of a channel, and the mux for the the 2nd
// channel.
func channelPair(t *testing.T) (*channel, *channel, *mux) ***REMOVED***
	c, s := muxPair()

	res := make(chan *channel, 1)
	go func() ***REMOVED***
		newCh, ok := <-s.incomingChannels
		if !ok ***REMOVED***
			t.Fatalf("No incoming channel")
		***REMOVED***
		if newCh.ChannelType() != "chan" ***REMOVED***
			t.Fatalf("got type %q want chan", newCh.ChannelType())
		***REMOVED***
		ch, _, err := newCh.Accept()
		if err != nil ***REMOVED***
			t.Fatalf("Accept %v", err)
		***REMOVED***
		res <- ch.(*channel)
	***REMOVED***()

	ch, err := c.openChannel("chan", nil)
	if err != nil ***REMOVED***
		t.Fatalf("OpenChannel: %v", err)
	***REMOVED***

	return <-res, ch, c
***REMOVED***

// Test that stderr and stdout can be addressed from different
// goroutines. This is intended for use with the race detector.
func TestMuxChannelExtendedThreadSafety(t *testing.T) ***REMOVED***
	writer, reader, mux := channelPair(t)
	defer writer.Close()
	defer reader.Close()
	defer mux.Close()

	var wr, rd sync.WaitGroup
	magic := "hello world"

	wr.Add(2)
	go func() ***REMOVED***
		io.WriteString(writer, magic)
		wr.Done()
	***REMOVED***()
	go func() ***REMOVED***
		io.WriteString(writer.Stderr(), magic)
		wr.Done()
	***REMOVED***()

	rd.Add(2)
	go func() ***REMOVED***
		c, err := ioutil.ReadAll(reader)
		if string(c) != magic ***REMOVED***
			t.Fatalf("stdout read got %q, want %q (error %s)", c, magic, err)
		***REMOVED***
		rd.Done()
	***REMOVED***()
	go func() ***REMOVED***
		c, err := ioutil.ReadAll(reader.Stderr())
		if string(c) != magic ***REMOVED***
			t.Fatalf("stderr read got %q, want %q (error %s)", c, magic, err)
		***REMOVED***
		rd.Done()
	***REMOVED***()

	wr.Wait()
	writer.CloseWrite()
	rd.Wait()
***REMOVED***

func TestMuxReadWrite(t *testing.T) ***REMOVED***
	s, c, mux := channelPair(t)
	defer s.Close()
	defer c.Close()
	defer mux.Close()

	magic := "hello world"
	magicExt := "hello stderr"
	go func() ***REMOVED***
		_, err := s.Write([]byte(magic))
		if err != nil ***REMOVED***
			t.Fatalf("Write: %v", err)
		***REMOVED***
		_, err = s.Extended(1).Write([]byte(magicExt))
		if err != nil ***REMOVED***
			t.Fatalf("Write: %v", err)
		***REMOVED***
		err = s.Close()
		if err != nil ***REMOVED***
			t.Fatalf("Close: %v", err)
		***REMOVED***
	***REMOVED***()

	var buf [1024]byte
	n, err := c.Read(buf[:])
	if err != nil ***REMOVED***
		t.Fatalf("server Read: %v", err)
	***REMOVED***
	got := string(buf[:n])
	if got != magic ***REMOVED***
		t.Fatalf("server: got %q want %q", got, magic)
	***REMOVED***

	n, err = c.Extended(1).Read(buf[:])
	if err != nil ***REMOVED***
		t.Fatalf("server Read: %v", err)
	***REMOVED***

	got = string(buf[:n])
	if got != magicExt ***REMOVED***
		t.Fatalf("server: got %q want %q", got, magic)
	***REMOVED***
***REMOVED***

func TestMuxChannelOverflow(t *testing.T) ***REMOVED***
	reader, writer, mux := channelPair(t)
	defer reader.Close()
	defer writer.Close()
	defer mux.Close()

	wDone := make(chan int, 1)
	go func() ***REMOVED***
		if _, err := writer.Write(make([]byte, channelWindowSize)); err != nil ***REMOVED***
			t.Errorf("could not fill window: %v", err)
		***REMOVED***
		writer.Write(make([]byte, 1))
		wDone <- 1
	***REMOVED***()
	writer.remoteWin.waitWriterBlocked()

	// Send 1 byte.
	packet := make([]byte, 1+4+4+1)
	packet[0] = msgChannelData
	marshalUint32(packet[1:], writer.remoteId)
	marshalUint32(packet[5:], uint32(1))
	packet[9] = 42

	if err := writer.mux.conn.writePacket(packet); err != nil ***REMOVED***
		t.Errorf("could not send packet")
	***REMOVED***
	if _, err := reader.SendRequest("hello", true, nil); err == nil ***REMOVED***
		t.Errorf("SendRequest succeeded.")
	***REMOVED***
	<-wDone
***REMOVED***

func TestMuxChannelCloseWriteUnblock(t *testing.T) ***REMOVED***
	reader, writer, mux := channelPair(t)
	defer reader.Close()
	defer writer.Close()
	defer mux.Close()

	wDone := make(chan int, 1)
	go func() ***REMOVED***
		if _, err := writer.Write(make([]byte, channelWindowSize)); err != nil ***REMOVED***
			t.Errorf("could not fill window: %v", err)
		***REMOVED***
		if _, err := writer.Write(make([]byte, 1)); err != io.EOF ***REMOVED***
			t.Errorf("got %v, want EOF for unblock write", err)
		***REMOVED***
		wDone <- 1
	***REMOVED***()

	writer.remoteWin.waitWriterBlocked()
	reader.Close()
	<-wDone
***REMOVED***

func TestMuxConnectionCloseWriteUnblock(t *testing.T) ***REMOVED***
	reader, writer, mux := channelPair(t)
	defer reader.Close()
	defer writer.Close()
	defer mux.Close()

	wDone := make(chan int, 1)
	go func() ***REMOVED***
		if _, err := writer.Write(make([]byte, channelWindowSize)); err != nil ***REMOVED***
			t.Errorf("could not fill window: %v", err)
		***REMOVED***
		if _, err := writer.Write(make([]byte, 1)); err != io.EOF ***REMOVED***
			t.Errorf("got %v, want EOF for unblock write", err)
		***REMOVED***
		wDone <- 1
	***REMOVED***()

	writer.remoteWin.waitWriterBlocked()
	mux.Close()
	<-wDone
***REMOVED***

func TestMuxReject(t *testing.T) ***REMOVED***
	client, server := muxPair()
	defer server.Close()
	defer client.Close()

	go func() ***REMOVED***
		ch, ok := <-server.incomingChannels
		if !ok ***REMOVED***
			t.Fatalf("Accept")
		***REMOVED***
		if ch.ChannelType() != "ch" || string(ch.ExtraData()) != "extra" ***REMOVED***
			t.Fatalf("unexpected channel: %q, %q", ch.ChannelType(), ch.ExtraData())
		***REMOVED***
		ch.Reject(RejectionReason(42), "message")
	***REMOVED***()

	ch, err := client.openChannel("ch", []byte("extra"))
	if ch != nil ***REMOVED***
		t.Fatal("openChannel not rejected")
	***REMOVED***

	ocf, ok := err.(*OpenChannelError)
	if !ok ***REMOVED***
		t.Errorf("got %#v want *OpenChannelError", err)
	***REMOVED*** else if ocf.Reason != 42 || ocf.Message != "message" ***REMOVED***
		t.Errorf("got %#v, want ***REMOVED***Reason: 42, Message: %q***REMOVED***", ocf, "message")
	***REMOVED***

	want := "ssh: rejected: unknown reason 42 (message)"
	if err.Error() != want ***REMOVED***
		t.Errorf("got %q, want %q", err.Error(), want)
	***REMOVED***
***REMOVED***

func TestMuxChannelRequest(t *testing.T) ***REMOVED***
	client, server, mux := channelPair(t)
	defer server.Close()
	defer client.Close()
	defer mux.Close()

	var received int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() ***REMOVED***
		for r := range server.incomingRequests ***REMOVED***
			received++
			r.Reply(r.Type == "yes", nil)
		***REMOVED***
		wg.Done()
	***REMOVED***()
	_, err := client.SendRequest("yes", false, nil)
	if err != nil ***REMOVED***
		t.Fatalf("SendRequest: %v", err)
	***REMOVED***
	ok, err := client.SendRequest("yes", true, nil)
	if err != nil ***REMOVED***
		t.Fatalf("SendRequest: %v", err)
	***REMOVED***

	if !ok ***REMOVED***
		t.Errorf("SendRequest(yes): %v", ok)

	***REMOVED***

	ok, err = client.SendRequest("no", true, nil)
	if err != nil ***REMOVED***
		t.Fatalf("SendRequest: %v", err)
	***REMOVED***
	if ok ***REMOVED***
		t.Errorf("SendRequest(no): %v", ok)

	***REMOVED***

	client.Close()
	wg.Wait()

	if received != 3 ***REMOVED***
		t.Errorf("got %d requests, want %d", received, 3)
	***REMOVED***
***REMOVED***

func TestMuxGlobalRequest(t *testing.T) ***REMOVED***
	clientMux, serverMux := muxPair()
	defer serverMux.Close()
	defer clientMux.Close()

	var seen bool
	go func() ***REMOVED***
		for r := range serverMux.incomingRequests ***REMOVED***
			seen = seen || r.Type == "peek"
			if r.WantReply ***REMOVED***
				err := r.Reply(r.Type == "yes",
					append([]byte(r.Type), r.Payload...))
				if err != nil ***REMOVED***
					t.Errorf("AckRequest: %v", err)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	_, _, err := clientMux.SendRequest("peek", false, nil)
	if err != nil ***REMOVED***
		t.Errorf("SendRequest: %v", err)
	***REMOVED***

	ok, data, err := clientMux.SendRequest("yes", true, []byte("a"))
	if !ok || string(data) != "yesa" || err != nil ***REMOVED***
		t.Errorf("SendRequest(\"yes\", true, \"a\"): %v %v %v",
			ok, data, err)
	***REMOVED***
	if ok, data, err := clientMux.SendRequest("yes", true, []byte("a")); !ok || string(data) != "yesa" || err != nil ***REMOVED***
		t.Errorf("SendRequest(\"yes\", true, \"a\"): %v %v %v",
			ok, data, err)
	***REMOVED***

	if ok, data, err := clientMux.SendRequest("no", true, []byte("a")); ok || string(data) != "noa" || err != nil ***REMOVED***
		t.Errorf("SendRequest(\"no\", true, \"a\"): %v %v %v",
			ok, data, err)
	***REMOVED***

	if !seen ***REMOVED***
		t.Errorf("never saw 'peek' request")
	***REMOVED***
***REMOVED***

func TestMuxGlobalRequestUnblock(t *testing.T) ***REMOVED***
	clientMux, serverMux := muxPair()
	defer serverMux.Close()
	defer clientMux.Close()

	result := make(chan error, 1)
	go func() ***REMOVED***
		_, _, err := clientMux.SendRequest("hello", true, nil)
		result <- err
	***REMOVED***()

	<-serverMux.incomingRequests
	serverMux.conn.Close()
	err := <-result

	if err != io.EOF ***REMOVED***
		t.Errorf("want EOF, got %v", io.EOF)
	***REMOVED***
***REMOVED***

func TestMuxChannelRequestUnblock(t *testing.T) ***REMOVED***
	a, b, connB := channelPair(t)
	defer a.Close()
	defer b.Close()
	defer connB.Close()

	result := make(chan error, 1)
	go func() ***REMOVED***
		_, err := a.SendRequest("hello", true, nil)
		result <- err
	***REMOVED***()

	<-b.incomingRequests
	connB.conn.Close()
	err := <-result

	if err != io.EOF ***REMOVED***
		t.Errorf("want EOF, got %v", err)
	***REMOVED***
***REMOVED***

func TestMuxCloseChannel(t *testing.T) ***REMOVED***
	r, w, mux := channelPair(t)
	defer mux.Close()
	defer r.Close()
	defer w.Close()

	result := make(chan error, 1)
	go func() ***REMOVED***
		var b [1024]byte
		_, err := r.Read(b[:])
		result <- err
	***REMOVED***()
	if err := w.Close(); err != nil ***REMOVED***
		t.Errorf("w.Close: %v", err)
	***REMOVED***

	if _, err := w.Write([]byte("hello")); err != io.EOF ***REMOVED***
		t.Errorf("got err %v, want io.EOF after Close", err)
	***REMOVED***

	if err := <-result; err != io.EOF ***REMOVED***
		t.Errorf("got %v (%T), want io.EOF", err, err)
	***REMOVED***
***REMOVED***

func TestMuxCloseWriteChannel(t *testing.T) ***REMOVED***
	r, w, mux := channelPair(t)
	defer mux.Close()

	result := make(chan error, 1)
	go func() ***REMOVED***
		var b [1024]byte
		_, err := r.Read(b[:])
		result <- err
	***REMOVED***()
	if err := w.CloseWrite(); err != nil ***REMOVED***
		t.Errorf("w.CloseWrite: %v", err)
	***REMOVED***

	if _, err := w.Write([]byte("hello")); err != io.EOF ***REMOVED***
		t.Errorf("got err %v, want io.EOF after CloseWrite", err)
	***REMOVED***

	if err := <-result; err != io.EOF ***REMOVED***
		t.Errorf("got %v (%T), want io.EOF", err, err)
	***REMOVED***
***REMOVED***

func TestMuxInvalidRecord(t *testing.T) ***REMOVED***
	a, b := muxPair()
	defer a.Close()
	defer b.Close()

	packet := make([]byte, 1+4+4+1)
	packet[0] = msgChannelData
	marshalUint32(packet[1:], 29348723 /* invalid channel id */)
	marshalUint32(packet[5:], 1)
	packet[9] = 42

	a.conn.writePacket(packet)
	go a.SendRequest("hello", false, nil)
	// 'a' wrote an invalid packet, so 'b' has exited.
	req, ok := <-b.incomingRequests
	if ok ***REMOVED***
		t.Errorf("got request %#v after receiving invalid packet", req)
	***REMOVED***
***REMOVED***

func TestZeroWindowAdjust(t *testing.T) ***REMOVED***
	a, b, mux := channelPair(t)
	defer a.Close()
	defer b.Close()
	defer mux.Close()

	go func() ***REMOVED***
		io.WriteString(a, "hello")
		// bogus adjust.
		a.sendMessage(windowAdjustMsg***REMOVED******REMOVED***)
		io.WriteString(a, "world")
		a.Close()
	***REMOVED***()

	want := "helloworld"
	c, _ := ioutil.ReadAll(b)
	if string(c) != want ***REMOVED***
		t.Errorf("got %q want %q", c, want)
	***REMOVED***
***REMOVED***

func TestMuxMaxPacketSize(t *testing.T) ***REMOVED***
	a, b, mux := channelPair(t)
	defer a.Close()
	defer b.Close()
	defer mux.Close()

	large := make([]byte, a.maxRemotePayload+1)
	packet := make([]byte, 1+4+4+1+len(large))
	packet[0] = msgChannelData
	marshalUint32(packet[1:], a.remoteId)
	marshalUint32(packet[5:], uint32(len(large)))
	packet[9] = 42

	if err := a.mux.conn.writePacket(packet); err != nil ***REMOVED***
		t.Errorf("could not send packet")
	***REMOVED***

	go a.SendRequest("hello", false, nil)

	_, ok := <-b.incomingRequests
	if ok ***REMOVED***
		t.Errorf("connection still alive after receiving large packet.")
	***REMOVED***
***REMOVED***

// Don't ship code with debug=true.
func TestDebug(t *testing.T) ***REMOVED***
	if debugMux ***REMOVED***
		t.Error("mux debug switched on")
	***REMOVED***
	if debugHandshake ***REMOVED***
		t.Error("handshake debug switched on")
	***REMOVED***
	if debugTransport ***REMOVED***
		t.Error("transport debug switched on")
	***REMOVED***
***REMOVED***
