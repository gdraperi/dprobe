// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"errors"
	"io"
	"net"
	"testing"
)

type server struct ***REMOVED***
	*ServerConn
	chans <-chan NewChannel
***REMOVED***

func newServer(c net.Conn, conf *ServerConfig) (*server, error) ***REMOVED***
	sconn, chans, reqs, err := NewServerConn(c, conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go DiscardRequests(reqs)
	return &server***REMOVED***sconn, chans***REMOVED***, nil
***REMOVED***

func (s *server) Accept() (NewChannel, error) ***REMOVED***
	n, ok := <-s.chans
	if !ok ***REMOVED***
		return nil, io.EOF
	***REMOVED***
	return n, nil
***REMOVED***

func sshPipe() (Conn, *server, error) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	clientConf := ClientConfig***REMOVED***
		User:            "user",
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	serverConf := ServerConfig***REMOVED***
		NoClientAuth: true,
	***REMOVED***
	serverConf.AddHostKey(testSigners["ecdsa"])
	done := make(chan *server, 1)
	go func() ***REMOVED***
		server, err := newServer(c2, &serverConf)
		if err != nil ***REMOVED***
			done <- nil
		***REMOVED***
		done <- server
	***REMOVED***()

	client, _, reqs, err := NewClientConn(c1, "", &clientConf)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	server := <-done
	if server == nil ***REMOVED***
		return nil, nil, errors.New("server handshake failed.")
	***REMOVED***
	go DiscardRequests(reqs)

	return client, server, nil
***REMOVED***

func BenchmarkEndToEnd(b *testing.B) ***REMOVED***
	b.StopTimer()

	client, server, err := sshPipe()
	if err != nil ***REMOVED***
		b.Fatalf("sshPipe: %v", err)
	***REMOVED***

	defer client.Close()
	defer server.Close()

	size := (1 << 20)
	input := make([]byte, size)
	output := make([]byte, size)
	b.SetBytes(int64(size))
	done := make(chan int, 1)

	go func() ***REMOVED***
		newCh, err := server.Accept()
		if err != nil ***REMOVED***
			b.Fatalf("Client: %v", err)
		***REMOVED***
		ch, incoming, err := newCh.Accept()
		go DiscardRequests(incoming)
		for i := 0; i < b.N; i++ ***REMOVED***
			if _, err := io.ReadFull(ch, output); err != nil ***REMOVED***
				b.Fatalf("ReadFull: %v", err)
			***REMOVED***
		***REMOVED***
		ch.Close()
		done <- 1
	***REMOVED***()

	ch, in, err := client.OpenChannel("speed", nil)
	if err != nil ***REMOVED***
		b.Fatalf("OpenChannel: %v", err)
	***REMOVED***
	go DiscardRequests(in)

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := ch.Write(input); err != nil ***REMOVED***
			b.Fatalf("WriteFull: %v", err)
		***REMOVED***
	***REMOVED***
	ch.Close()
	b.StopTimer()

	<-done
***REMOVED***
