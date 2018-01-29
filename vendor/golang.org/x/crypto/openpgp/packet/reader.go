// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"golang.org/x/crypto/openpgp/errors"
	"io"
)

// Reader reads packets from an io.Reader and allows packets to be 'unread' so
// that they result from the next call to Next.
type Reader struct ***REMOVED***
	q       []Packet
	readers []io.Reader
***REMOVED***

// New io.Readers are pushed when a compressed or encrypted packet is processed
// and recursively treated as a new source of packets. However, a carefully
// crafted packet can trigger an infinite recursive sequence of packets. See
// http://mumble.net/~campbell/misc/pgp-quine
// https://web.nvd.nist.gov/view/vuln/detail?vulnId=CVE-2013-4402
// This constant limits the number of recursive packets that may be pushed.
const maxReaders = 32

// Next returns the most recently unread Packet, or reads another packet from
// the top-most io.Reader. Unknown packet types are skipped.
func (r *Reader) Next() (p Packet, err error) ***REMOVED***
	if len(r.q) > 0 ***REMOVED***
		p = r.q[len(r.q)-1]
		r.q = r.q[:len(r.q)-1]
		return
	***REMOVED***

	for len(r.readers) > 0 ***REMOVED***
		p, err = Read(r.readers[len(r.readers)-1])
		if err == nil ***REMOVED***
			return
		***REMOVED***
		if err == io.EOF ***REMOVED***
			r.readers = r.readers[:len(r.readers)-1]
			continue
		***REMOVED***
		if _, ok := err.(errors.UnknownPacketTypeError); !ok ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return nil, io.EOF
***REMOVED***

// Push causes the Reader to start reading from a new io.Reader. When an EOF
// error is seen from the new io.Reader, it is popped and the Reader continues
// to read from the next most recent io.Reader. Push returns a StructuralError
// if pushing the reader would exceed the maximum recursion level, otherwise it
// returns nil.
func (r *Reader) Push(reader io.Reader) (err error) ***REMOVED***
	if len(r.readers) >= maxReaders ***REMOVED***
		return errors.StructuralError("too many layers of packets")
	***REMOVED***
	r.readers = append(r.readers, reader)
	return nil
***REMOVED***

// Unread causes the given Packet to be returned from the next call to Next.
func (r *Reader) Unread(p Packet) ***REMOVED***
	r.q = append(r.q, p)
***REMOVED***

func NewReader(r io.Reader) *Reader ***REMOVED***
	return &Reader***REMOVED***
		q:       nil,
		readers: []io.Reader***REMOVED***r***REMOVED***,
	***REMOVED***
***REMOVED***
