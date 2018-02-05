// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package websocket

import (
	"io"
	"io/ioutil"
	"sync/atomic"
	"testing"
)

// broadcastBench allows to run broadcast benchmarks.
// In every broadcast benchmark we create many connections, then send the same
// message into every connection and wait for all writes complete. This emulates
// an application where many connections listen to the same data - i.e. PUB/SUB
// scenarios with many subscribers in one channel.
type broadcastBench struct ***REMOVED***
	w           io.Writer
	message     *broadcastMessage
	closeCh     chan struct***REMOVED******REMOVED***
	doneCh      chan struct***REMOVED******REMOVED***
	count       int32
	conns       []*broadcastConn
	compression bool
	usePrepared bool
***REMOVED***

type broadcastMessage struct ***REMOVED***
	payload  []byte
	prepared *PreparedMessage
***REMOVED***

type broadcastConn struct ***REMOVED***
	conn  *Conn
	msgCh chan *broadcastMessage
***REMOVED***

func newBroadcastConn(c *Conn) *broadcastConn ***REMOVED***
	return &broadcastConn***REMOVED***
		conn:  c,
		msgCh: make(chan *broadcastMessage, 1),
	***REMOVED***
***REMOVED***

func newBroadcastBench(usePrepared, compression bool) *broadcastBench ***REMOVED***
	bench := &broadcastBench***REMOVED***
		w:           ioutil.Discard,
		doneCh:      make(chan struct***REMOVED******REMOVED***),
		closeCh:     make(chan struct***REMOVED******REMOVED***),
		usePrepared: usePrepared,
		compression: compression,
	***REMOVED***
	msg := &broadcastMessage***REMOVED***
		payload: textMessages(1)[0],
	***REMOVED***
	if usePrepared ***REMOVED***
		pm, _ := NewPreparedMessage(TextMessage, msg.payload)
		msg.prepared = pm
	***REMOVED***
	bench.message = msg
	bench.makeConns(10000)
	return bench
***REMOVED***

func (b *broadcastBench) makeConns(numConns int) ***REMOVED***
	conns := make([]*broadcastConn, numConns)

	for i := 0; i < numConns; i++ ***REMOVED***
		c := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: b.w***REMOVED***, true, 1024, 1024)
		if b.compression ***REMOVED***
			c.enableWriteCompression = true
			c.newCompressionWriter = compressNoContextTakeover
		***REMOVED***
		conns[i] = newBroadcastConn(c)
		go func(c *broadcastConn) ***REMOVED***
			for ***REMOVED***
				select ***REMOVED***
				case msg := <-c.msgCh:
					if b.usePrepared ***REMOVED***
						c.conn.WritePreparedMessage(msg.prepared)
					***REMOVED*** else ***REMOVED***
						c.conn.WriteMessage(TextMessage, msg.payload)
					***REMOVED***
					val := atomic.AddInt32(&b.count, 1)
					if val%int32(numConns) == 0 ***REMOVED***
						b.doneCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
				case <-b.closeCh:
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***(conns[i])
	***REMOVED***
	b.conns = conns
***REMOVED***

func (b *broadcastBench) close() ***REMOVED***
	close(b.closeCh)
***REMOVED***

func (b *broadcastBench) runOnce() ***REMOVED***
	for _, c := range b.conns ***REMOVED***
		c.msgCh <- b.message
	***REMOVED***
	<-b.doneCh
***REMOVED***

func BenchmarkBroadcast(b *testing.B) ***REMOVED***
	benchmarks := []struct ***REMOVED***
		name        string
		usePrepared bool
		compression bool
	***REMOVED******REMOVED***
		***REMOVED***"NoCompression", false, false***REMOVED***,
		***REMOVED***"WithCompression", false, true***REMOVED***,
		***REMOVED***"NoCompressionPrepared", true, false***REMOVED***,
		***REMOVED***"WithCompressionPrepared", true, true***REMOVED***,
	***REMOVED***
	for _, bm := range benchmarks ***REMOVED***
		b.Run(bm.name, func(b *testing.B) ***REMOVED***
			bench := newBroadcastBench(bm.usePrepared, bm.compression)
			defer bench.close()
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				bench.runOnce()
			***REMOVED***
			b.ReportAllocs()
		***REMOVED***)
	***REMOVED***
***REMOVED***
