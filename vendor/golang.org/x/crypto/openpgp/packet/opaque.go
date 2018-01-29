// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/openpgp/errors"
)

// OpaquePacket represents an OpenPGP packet as raw, unparsed data. This is
// useful for splitting and storing the original packet contents separately,
// handling unsupported packet types or accessing parts of the packet not yet
// implemented by this package.
type OpaquePacket struct ***REMOVED***
	// Packet type
	Tag uint8
	// Reason why the packet was parsed opaquely
	Reason error
	// Binary contents of the packet data
	Contents []byte
***REMOVED***

func (op *OpaquePacket) parse(r io.Reader) (err error) ***REMOVED***
	op.Contents, err = ioutil.ReadAll(r)
	return
***REMOVED***

// Serialize marshals the packet to a writer in its original form, including
// the packet header.
func (op *OpaquePacket) Serialize(w io.Writer) (err error) ***REMOVED***
	err = serializeHeader(w, packetType(op.Tag), len(op.Contents))
	if err == nil ***REMOVED***
		_, err = w.Write(op.Contents)
	***REMOVED***
	return
***REMOVED***

// Parse attempts to parse the opaque contents into a structure supported by
// this package. If the packet is not known then the result will be another
// OpaquePacket.
func (op *OpaquePacket) Parse() (p Packet, err error) ***REMOVED***
	hdr := bytes.NewBuffer(nil)
	err = serializeHeader(hdr, packetType(op.Tag), len(op.Contents))
	if err != nil ***REMOVED***
		op.Reason = err
		return op, err
	***REMOVED***
	p, err = Read(io.MultiReader(hdr, bytes.NewBuffer(op.Contents)))
	if err != nil ***REMOVED***
		op.Reason = err
		p = op
	***REMOVED***
	return
***REMOVED***

// OpaqueReader reads OpaquePackets from an io.Reader.
type OpaqueReader struct ***REMOVED***
	r io.Reader
***REMOVED***

func NewOpaqueReader(r io.Reader) *OpaqueReader ***REMOVED***
	return &OpaqueReader***REMOVED***r: r***REMOVED***
***REMOVED***

// Read the next OpaquePacket.
func (or *OpaqueReader) Next() (op *OpaquePacket, err error) ***REMOVED***
	tag, _, contents, err := readHeader(or.r)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	op = &OpaquePacket***REMOVED***Tag: uint8(tag), Reason: err***REMOVED***
	err = op.parse(contents)
	if err != nil ***REMOVED***
		consumeAll(contents)
	***REMOVED***
	return
***REMOVED***

// OpaqueSubpacket represents an unparsed OpenPGP subpacket,
// as found in signature and user attribute packets.
type OpaqueSubpacket struct ***REMOVED***
	SubType  uint8
	Contents []byte
***REMOVED***

// OpaqueSubpackets extracts opaque, unparsed OpenPGP subpackets from
// their byte representation.
func OpaqueSubpackets(contents []byte) (result []*OpaqueSubpacket, err error) ***REMOVED***
	var (
		subHeaderLen int
		subPacket    *OpaqueSubpacket
	)
	for len(contents) > 0 ***REMOVED***
		subHeaderLen, subPacket, err = nextSubpacket(contents)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		result = append(result, subPacket)
		contents = contents[subHeaderLen+len(subPacket.Contents):]
	***REMOVED***
	return
***REMOVED***

func nextSubpacket(contents []byte) (subHeaderLen int, subPacket *OpaqueSubpacket, err error) ***REMOVED***
	// RFC 4880, section 5.2.3.1
	var subLen uint32
	if len(contents) < 1 ***REMOVED***
		goto Truncated
	***REMOVED***
	subPacket = &OpaqueSubpacket***REMOVED******REMOVED***
	switch ***REMOVED***
	case contents[0] < 192:
		subHeaderLen = 2 // 1 length byte, 1 subtype byte
		if len(contents) < subHeaderLen ***REMOVED***
			goto Truncated
		***REMOVED***
		subLen = uint32(contents[0])
		contents = contents[1:]
	case contents[0] < 255:
		subHeaderLen = 3 // 2 length bytes, 1 subtype
		if len(contents) < subHeaderLen ***REMOVED***
			goto Truncated
		***REMOVED***
		subLen = uint32(contents[0]-192)<<8 + uint32(contents[1]) + 192
		contents = contents[2:]
	default:
		subHeaderLen = 6 // 5 length bytes, 1 subtype
		if len(contents) < subHeaderLen ***REMOVED***
			goto Truncated
		***REMOVED***
		subLen = uint32(contents[1])<<24 |
			uint32(contents[2])<<16 |
			uint32(contents[3])<<8 |
			uint32(contents[4])
		contents = contents[5:]
	***REMOVED***
	if subLen > uint32(len(contents)) || subLen == 0 ***REMOVED***
		goto Truncated
	***REMOVED***
	subPacket.SubType = contents[0]
	subPacket.Contents = contents[1:subLen]
	return
Truncated:
	err = errors.StructuralError("subpacket truncated")
	return
***REMOVED***

func (osp *OpaqueSubpacket) Serialize(w io.Writer) (err error) ***REMOVED***
	buf := make([]byte, 6)
	n := serializeSubpacketLength(buf, len(osp.Contents)+1)
	buf[n] = osp.SubType
	if _, err = w.Write(buf[:n+1]); err != nil ***REMOVED***
		return
	***REMOVED***
	_, err = w.Write(osp.Contents)
	return
***REMOVED***
