// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package s2k implements the various OpenPGP string-to-key transforms as
// specified in RFC 4800 section 3.7.1.
package s2k // import "golang.org/x/crypto/openpgp/s2k"

import (
	"crypto"
	"hash"
	"io"
	"strconv"

	"golang.org/x/crypto/openpgp/errors"
)

// Config collects configuration parameters for s2k key-stretching
// transformatioms. A nil *Config is valid and results in all default
// values. Currently, Config is used only by the Serialize function in
// this package.
type Config struct ***REMOVED***
	// Hash is the default hash function to be used. If
	// nil, SHA1 is used.
	Hash crypto.Hash
	// S2KCount is only used for symmetric encryption. It
	// determines the strength of the passphrase stretching when
	// the said passphrase is hashed to produce a key. S2KCount
	// should be between 1024 and 65011712, inclusive. If Config
	// is nil or S2KCount is 0, the value 65536 used. Not all
	// values in the above range can be represented. S2KCount will
	// be rounded up to the next representable value if it cannot
	// be encoded exactly. When set, it is strongly encrouraged to
	// use a value that is at least 65536. See RFC 4880 Section
	// 3.7.1.3.
	S2KCount int
***REMOVED***

func (c *Config) hash() crypto.Hash ***REMOVED***
	if c == nil || uint(c.Hash) == 0 ***REMOVED***
		// SHA1 is the historical default in this package.
		return crypto.SHA1
	***REMOVED***

	return c.Hash
***REMOVED***

func (c *Config) encodedCount() uint8 ***REMOVED***
	if c == nil || c.S2KCount == 0 ***REMOVED***
		return 96 // The common case. Correspoding to 65536
	***REMOVED***

	i := c.S2KCount
	switch ***REMOVED***
	// Behave like GPG. Should we make 65536 the lowest value used?
	case i < 1024:
		i = 1024
	case i > 65011712:
		i = 65011712
	***REMOVED***

	return encodeCount(i)
***REMOVED***

// encodeCount converts an iterative "count" in the range 1024 to
// 65011712, inclusive, to an encoded count. The return value is the
// octet that is actually stored in the GPG file. encodeCount panics
// if i is not in the above range (encodedCount above takes care to
// pass i in the correct range). See RFC 4880 Section 3.7.7.1.
func encodeCount(i int) uint8 ***REMOVED***
	if i < 1024 || i > 65011712 ***REMOVED***
		panic("count arg i outside the required range")
	***REMOVED***

	for encoded := 0; encoded < 256; encoded++ ***REMOVED***
		count := decodeCount(uint8(encoded))
		if count >= i ***REMOVED***
			return uint8(encoded)
		***REMOVED***
	***REMOVED***

	return 255
***REMOVED***

// decodeCount returns the s2k mode 3 iterative "count" corresponding to
// the encoded octet c.
func decodeCount(c uint8) int ***REMOVED***
	return (16 + int(c&15)) << (uint32(c>>4) + 6)
***REMOVED***

// Simple writes to out the result of computing the Simple S2K function (RFC
// 4880, section 3.7.1.1) using the given hash and input passphrase.
func Simple(out []byte, h hash.Hash, in []byte) ***REMOVED***
	Salted(out, h, in, nil)
***REMOVED***

var zero [1]byte

// Salted writes to out the result of computing the Salted S2K function (RFC
// 4880, section 3.7.1.2) using the given hash, input passphrase and salt.
func Salted(out []byte, h hash.Hash, in []byte, salt []byte) ***REMOVED***
	done := 0
	var digest []byte

	for i := 0; done < len(out); i++ ***REMOVED***
		h.Reset()
		for j := 0; j < i; j++ ***REMOVED***
			h.Write(zero[:])
		***REMOVED***
		h.Write(salt)
		h.Write(in)
		digest = h.Sum(digest[:0])
		n := copy(out[done:], digest)
		done += n
	***REMOVED***
***REMOVED***

// Iterated writes to out the result of computing the Iterated and Salted S2K
// function (RFC 4880, section 3.7.1.3) using the given hash, input passphrase,
// salt and iteration count.
func Iterated(out []byte, h hash.Hash, in []byte, salt []byte, count int) ***REMOVED***
	combined := make([]byte, len(in)+len(salt))
	copy(combined, salt)
	copy(combined[len(salt):], in)

	if count < len(combined) ***REMOVED***
		count = len(combined)
	***REMOVED***

	done := 0
	var digest []byte
	for i := 0; done < len(out); i++ ***REMOVED***
		h.Reset()
		for j := 0; j < i; j++ ***REMOVED***
			h.Write(zero[:])
		***REMOVED***
		written := 0
		for written < count ***REMOVED***
			if written+len(combined) > count ***REMOVED***
				todo := count - written
				h.Write(combined[:todo])
				written = count
			***REMOVED*** else ***REMOVED***
				h.Write(combined)
				written += len(combined)
			***REMOVED***
		***REMOVED***
		digest = h.Sum(digest[:0])
		n := copy(out[done:], digest)
		done += n
	***REMOVED***
***REMOVED***

// Parse reads a binary specification for a string-to-key transformation from r
// and returns a function which performs that transform.
func Parse(r io.Reader) (f func(out, in []byte), err error) ***REMOVED***
	var buf [9]byte

	_, err = io.ReadFull(r, buf[:2])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	hash, ok := HashIdToHash(buf[1])
	if !ok ***REMOVED***
		return nil, errors.UnsupportedError("hash for S2K function: " + strconv.Itoa(int(buf[1])))
	***REMOVED***
	if !hash.Available() ***REMOVED***
		return nil, errors.UnsupportedError("hash not available: " + strconv.Itoa(int(hash)))
	***REMOVED***
	h := hash.New()

	switch buf[0] ***REMOVED***
	case 0:
		f := func(out, in []byte) ***REMOVED***
			Simple(out, h, in)
		***REMOVED***
		return f, nil
	case 1:
		_, err = io.ReadFull(r, buf[:8])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		f := func(out, in []byte) ***REMOVED***
			Salted(out, h, in, buf[:8])
		***REMOVED***
		return f, nil
	case 3:
		_, err = io.ReadFull(r, buf[:9])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		count := decodeCount(buf[8])
		f := func(out, in []byte) ***REMOVED***
			Iterated(out, h, in, buf[:8], count)
		***REMOVED***
		return f, nil
	***REMOVED***

	return nil, errors.UnsupportedError("S2K function")
***REMOVED***

// Serialize salts and stretches the given passphrase and writes the
// resulting key into key. It also serializes an S2K descriptor to
// w. The key stretching can be configured with c, which may be
// nil. In that case, sensible defaults will be used.
func Serialize(w io.Writer, key []byte, rand io.Reader, passphrase []byte, c *Config) error ***REMOVED***
	var buf [11]byte
	buf[0] = 3 /* iterated and salted */
	buf[1], _ = HashToHashId(c.hash())
	salt := buf[2:10]
	if _, err := io.ReadFull(rand, salt); err != nil ***REMOVED***
		return err
	***REMOVED***
	encodedCount := c.encodedCount()
	count := decodeCount(encodedCount)
	buf[10] = encodedCount
	if _, err := w.Write(buf[:]); err != nil ***REMOVED***
		return err
	***REMOVED***

	Iterated(key, c.hash().New(), passphrase, salt, count)
	return nil
***REMOVED***

// hashToHashIdMapping contains pairs relating OpenPGP's hash identifier with
// Go's crypto.Hash type. See RFC 4880, section 9.4.
var hashToHashIdMapping = []struct ***REMOVED***
	id   byte
	hash crypto.Hash
	name string
***REMOVED******REMOVED***
	***REMOVED***1, crypto.MD5, "MD5"***REMOVED***,
	***REMOVED***2, crypto.SHA1, "SHA1"***REMOVED***,
	***REMOVED***3, crypto.RIPEMD160, "RIPEMD160"***REMOVED***,
	***REMOVED***8, crypto.SHA256, "SHA256"***REMOVED***,
	***REMOVED***9, crypto.SHA384, "SHA384"***REMOVED***,
	***REMOVED***10, crypto.SHA512, "SHA512"***REMOVED***,
	***REMOVED***11, crypto.SHA224, "SHA224"***REMOVED***,
***REMOVED***

// HashIdToHash returns a crypto.Hash which corresponds to the given OpenPGP
// hash id.
func HashIdToHash(id byte) (h crypto.Hash, ok bool) ***REMOVED***
	for _, m := range hashToHashIdMapping ***REMOVED***
		if m.id == id ***REMOVED***
			return m.hash, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

// HashIdToString returns the name of the hash function corresponding to the
// given OpenPGP hash id.
func HashIdToString(id byte) (name string, ok bool) ***REMOVED***
	for _, m := range hashToHashIdMapping ***REMOVED***
		if m.id == id ***REMOVED***
			return m.name, true
		***REMOVED***
	***REMOVED***

	return "", false
***REMOVED***

// HashIdToHash returns an OpenPGP hash id which corresponds the given Hash.
func HashToHashId(h crypto.Hash) (id byte, ok bool) ***REMOVED***
	for _, m := range hashToHashIdMapping ***REMOVED***
		if m.hash == h ***REMOVED***
			return m.id, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***
