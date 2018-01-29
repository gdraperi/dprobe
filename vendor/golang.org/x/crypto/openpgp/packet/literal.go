// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"encoding/binary"
	"io"
)

// LiteralData represents an encrypted file. See RFC 4880, section 5.9.
type LiteralData struct ***REMOVED***
	IsBinary bool
	FileName string
	Time     uint32 // Unix epoch time. Either creation time or modification time. 0 means undefined.
	Body     io.Reader
***REMOVED***

// ForEyesOnly returns whether the contents of the LiteralData have been marked
// as especially sensitive.
func (l *LiteralData) ForEyesOnly() bool ***REMOVED***
	return l.FileName == "_CONSOLE"
***REMOVED***

func (l *LiteralData) parse(r io.Reader) (err error) ***REMOVED***
	var buf [256]byte

	_, err = readFull(r, buf[:2])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	l.IsBinary = buf[0] == 'b'
	fileNameLen := int(buf[1])

	_, err = readFull(r, buf[:fileNameLen])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	l.FileName = string(buf[:fileNameLen])

	_, err = readFull(r, buf[:4])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	l.Time = binary.BigEndian.Uint32(buf[:4])
	l.Body = r
	return
***REMOVED***

// SerializeLiteral serializes a literal data packet to w and returns a
// WriteCloser to which the data itself can be written and which MUST be closed
// on completion. The fileName is truncated to 255 bytes.
func SerializeLiteral(w io.WriteCloser, isBinary bool, fileName string, time uint32) (plaintext io.WriteCloser, err error) ***REMOVED***
	var buf [4]byte
	buf[0] = 't'
	if isBinary ***REMOVED***
		buf[0] = 'b'
	***REMOVED***
	if len(fileName) > 255 ***REMOVED***
		fileName = fileName[:255]
	***REMOVED***
	buf[1] = byte(len(fileName))

	inner, err := serializeStreamHeader(w, packetTypeLiteralData)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = inner.Write(buf[:2])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = inner.Write([]byte(fileName))
	if err != nil ***REMOVED***
		return
	***REMOVED***
	binary.BigEndian.PutUint32(buf[:], time)
	_, err = inner.Write(buf[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	plaintext = inner
	return
***REMOVED***
