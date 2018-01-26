// Copyright 2012 SocialCode. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Reader struct ***REMOVED***
	mu   sync.Mutex
	conn net.Conn
***REMOVED***

func NewReader(addr string) (*Reader, error) ***REMOVED***
	var err error
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("ResolveUDPAddr('%s'): %s", addr, err)
	***REMOVED***

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("ListenUDP: %s", err)
	***REMOVED***

	r := new(Reader)
	r.conn = conn
	return r, nil
***REMOVED***

func (r *Reader) Addr() string ***REMOVED***
	return r.conn.LocalAddr().String()
***REMOVED***

// FIXME: this will discard data if p isn't big enough to hold the
// full message.
func (r *Reader) Read(p []byte) (int, error) ***REMOVED***
	msg, err := r.ReadMessage()
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	var data string

	if msg.Full == "" ***REMOVED***
		data = msg.Short
	***REMOVED*** else ***REMOVED***
		data = msg.Full
	***REMOVED***

	return strings.NewReader(data).Read(p)
***REMOVED***

func (r *Reader) ReadMessage() (*Message, error) ***REMOVED***
	cBuf := make([]byte, ChunkSize)
	var (
		err        error
		n, length  int
		cid, ocid  []byte
		seq, total uint8
		cHead      []byte
		cReader    io.Reader
		chunks     [][]byte
	)

	for got := 0; got < 128 && (total == 0 || got < int(total)); got++ ***REMOVED***
		if n, err = r.conn.Read(cBuf); err != nil ***REMOVED***
			return nil, fmt.Errorf("Read: %s", err)
		***REMOVED***
		cHead, cBuf = cBuf[:2], cBuf[:n]

		if bytes.Equal(cHead, magicChunked) ***REMOVED***
			//fmt.Printf("chunked %v\n", cBuf[:14])
			cid, seq, total = cBuf[2:2+8], cBuf[2+8], cBuf[2+8+1]
			if ocid != nil && !bytes.Equal(cid, ocid) ***REMOVED***
				return nil, fmt.Errorf("out-of-band message %v (awaited %v)", cid, ocid)
			***REMOVED*** else if ocid == nil ***REMOVED***
				ocid = cid
				chunks = make([][]byte, total)
			***REMOVED***
			n = len(cBuf) - chunkedHeaderLen
			//fmt.Printf("setting chunks[%d]: %d\n", seq, n)
			chunks[seq] = append(make([]byte, 0, n), cBuf[chunkedHeaderLen:]...)
			length += n
		***REMOVED*** else ***REMOVED*** //not chunked
			if total > 0 ***REMOVED***
				return nil, fmt.Errorf("out-of-band message (not chunked)")
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	//fmt.Printf("\nchunks: %v\n", chunks)

	if length > 0 ***REMOVED***
		if cap(cBuf) < length ***REMOVED***
			cBuf = append(cBuf, make([]byte, 0, length-cap(cBuf))...)
		***REMOVED***
		cBuf = cBuf[:0]
		for i := range chunks ***REMOVED***
			//fmt.Printf("appending %d %v\n", i, chunks[i])
			cBuf = append(cBuf, chunks[i]...)
		***REMOVED***
		cHead = cBuf[:2]
	***REMOVED***

	// the data we get from the wire is compressed
	if bytes.Equal(cHead, magicGzip) ***REMOVED***
		cReader, err = gzip.NewReader(bytes.NewReader(cBuf))
	***REMOVED*** else if cHead[0] == magicZlib[0] &&
		(int(cHead[0])*256+int(cHead[1]))%31 == 0 ***REMOVED***
		// zlib is slightly more complicated, but correct
		cReader, err = zlib.NewReader(bytes.NewReader(cBuf))
	***REMOVED*** else ***REMOVED***
		// compliance with https://github.com/Graylog2/graylog2-server
		// treating all messages as uncompressed if  they are not gzip, zlib or
		// chunked
		cReader = bytes.NewReader(cBuf)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, fmt.Errorf("NewReader: %s", err)
	***REMOVED***

	msg := new(Message)
	if err := json.NewDecoder(cReader).Decode(&msg); err != nil ***REMOVED***
		return nil, fmt.Errorf("json.Unmarshal: %s", err)
	***REMOVED***

	return msg, nil
***REMOVED***
