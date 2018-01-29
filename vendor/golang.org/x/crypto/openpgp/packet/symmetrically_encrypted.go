// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto/cipher"
	"crypto/sha1"
	"crypto/subtle"
	"golang.org/x/crypto/openpgp/errors"
	"hash"
	"io"
	"strconv"
)

// SymmetricallyEncrypted represents a symmetrically encrypted byte string. The
// encrypted contents will consist of more OpenPGP packets. See RFC 4880,
// sections 5.7 and 5.13.
type SymmetricallyEncrypted struct ***REMOVED***
	MDC      bool // true iff this is a type 18 packet and thus has an embedded MAC.
	contents io.Reader
	prefix   []byte
***REMOVED***

const symmetricallyEncryptedVersion = 1

func (se *SymmetricallyEncrypted) parse(r io.Reader) error ***REMOVED***
	if se.MDC ***REMOVED***
		// See RFC 4880, section 5.13.
		var buf [1]byte
		_, err := readFull(r, buf[:])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if buf[0] != symmetricallyEncryptedVersion ***REMOVED***
			return errors.UnsupportedError("unknown SymmetricallyEncrypted version")
		***REMOVED***
	***REMOVED***
	se.contents = r
	return nil
***REMOVED***

// Decrypt returns a ReadCloser, from which the decrypted contents of the
// packet can be read. An incorrect key can, with high probability, be detected
// immediately and this will result in a KeyIncorrect error being returned.
func (se *SymmetricallyEncrypted) Decrypt(c CipherFunction, key []byte) (io.ReadCloser, error) ***REMOVED***
	keySize := c.KeySize()
	if keySize == 0 ***REMOVED***
		return nil, errors.UnsupportedError("unknown cipher: " + strconv.Itoa(int(c)))
	***REMOVED***
	if len(key) != keySize ***REMOVED***
		return nil, errors.InvalidArgumentError("SymmetricallyEncrypted: incorrect key length")
	***REMOVED***

	if se.prefix == nil ***REMOVED***
		se.prefix = make([]byte, c.blockSize()+2)
		_, err := readFull(se.contents, se.prefix)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else if len(se.prefix) != c.blockSize()+2 ***REMOVED***
		return nil, errors.InvalidArgumentError("can't try ciphers with different block lengths")
	***REMOVED***

	ocfbResync := OCFBResync
	if se.MDC ***REMOVED***
		// MDC packets use a different form of OCFB mode.
		ocfbResync = OCFBNoResync
	***REMOVED***

	s := NewOCFBDecrypter(c.new(key), se.prefix, ocfbResync)
	if s == nil ***REMOVED***
		return nil, errors.ErrKeyIncorrect
	***REMOVED***

	plaintext := cipher.StreamReader***REMOVED***S: s, R: se.contents***REMOVED***

	if se.MDC ***REMOVED***
		// MDC packets have an embedded hash that we need to check.
		h := sha1.New()
		h.Write(se.prefix)
		return &seMDCReader***REMOVED***in: plaintext, h: h***REMOVED***, nil
	***REMOVED***

	// Otherwise, we just need to wrap plaintext so that it's a valid ReadCloser.
	return seReader***REMOVED***plaintext***REMOVED***, nil
***REMOVED***

// seReader wraps an io.Reader with a no-op Close method.
type seReader struct ***REMOVED***
	in io.Reader
***REMOVED***

func (ser seReader) Read(buf []byte) (int, error) ***REMOVED***
	return ser.in.Read(buf)
***REMOVED***

func (ser seReader) Close() error ***REMOVED***
	return nil
***REMOVED***

const mdcTrailerSize = 1 /* tag byte */ + 1 /* length byte */ + sha1.Size

// An seMDCReader wraps an io.Reader, maintains a running hash and keeps hold
// of the most recent 22 bytes (mdcTrailerSize). Upon EOF, those bytes form an
// MDC packet containing a hash of the previous contents which is checked
// against the running hash. See RFC 4880, section 5.13.
type seMDCReader struct ***REMOVED***
	in          io.Reader
	h           hash.Hash
	trailer     [mdcTrailerSize]byte
	scratch     [mdcTrailerSize]byte
	trailerUsed int
	error       bool
	eof         bool
***REMOVED***

func (ser *seMDCReader) Read(buf []byte) (n int, err error) ***REMOVED***
	if ser.error ***REMOVED***
		err = io.ErrUnexpectedEOF
		return
	***REMOVED***
	if ser.eof ***REMOVED***
		err = io.EOF
		return
	***REMOVED***

	// If we haven't yet filled the trailer buffer then we must do that
	// first.
	for ser.trailerUsed < mdcTrailerSize ***REMOVED***
		n, err = ser.in.Read(ser.trailer[ser.trailerUsed:])
		ser.trailerUsed += n
		if err == io.EOF ***REMOVED***
			if ser.trailerUsed != mdcTrailerSize ***REMOVED***
				n = 0
				err = io.ErrUnexpectedEOF
				ser.error = true
				return
			***REMOVED***
			ser.eof = true
			n = 0
			return
		***REMOVED***

		if err != nil ***REMOVED***
			n = 0
			return
		***REMOVED***
	***REMOVED***

	// If it's a short read then we read into a temporary buffer and shift
	// the data into the caller's buffer.
	if len(buf) <= mdcTrailerSize ***REMOVED***
		n, err = readFull(ser.in, ser.scratch[:len(buf)])
		copy(buf, ser.trailer[:n])
		ser.h.Write(buf[:n])
		copy(ser.trailer[:], ser.trailer[n:])
		copy(ser.trailer[mdcTrailerSize-n:], ser.scratch[:])
		if n < len(buf) ***REMOVED***
			ser.eof = true
			err = io.EOF
		***REMOVED***
		return
	***REMOVED***

	n, err = ser.in.Read(buf[mdcTrailerSize:])
	copy(buf, ser.trailer[:])
	ser.h.Write(buf[:n])
	copy(ser.trailer[:], buf[n:])

	if err == io.EOF ***REMOVED***
		ser.eof = true
	***REMOVED***
	return
***REMOVED***

// This is a new-format packet tag byte for a type 19 (MDC) packet.
const mdcPacketTagByte = byte(0x80) | 0x40 | 19

func (ser *seMDCReader) Close() error ***REMOVED***
	if ser.error ***REMOVED***
		return errors.SignatureError("error during reading")
	***REMOVED***

	for !ser.eof ***REMOVED***
		// We haven't seen EOF so we need to read to the end
		var buf [1024]byte
		_, err := ser.Read(buf[:])
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return errors.SignatureError("error during reading")
		***REMOVED***
	***REMOVED***

	if ser.trailer[0] != mdcPacketTagByte || ser.trailer[1] != sha1.Size ***REMOVED***
		return errors.SignatureError("MDC packet not found")
	***REMOVED***
	ser.h.Write(ser.trailer[:2])

	final := ser.h.Sum(nil)
	if subtle.ConstantTimeCompare(final, ser.trailer[2:]) != 1 ***REMOVED***
		return errors.SignatureError("hash mismatch")
	***REMOVED***
	return nil
***REMOVED***

// An seMDCWriter writes through to an io.WriteCloser while maintains a running
// hash of the data written. On close, it emits an MDC packet containing the
// running hash.
type seMDCWriter struct ***REMOVED***
	w io.WriteCloser
	h hash.Hash
***REMOVED***

func (w *seMDCWriter) Write(buf []byte) (n int, err error) ***REMOVED***
	w.h.Write(buf)
	return w.w.Write(buf)
***REMOVED***

func (w *seMDCWriter) Close() (err error) ***REMOVED***
	var buf [mdcTrailerSize]byte

	buf[0] = mdcPacketTagByte
	buf[1] = sha1.Size
	w.h.Write(buf[:2])
	digest := w.h.Sum(nil)
	copy(buf[2:], digest)

	_, err = w.w.Write(buf[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return w.w.Close()
***REMOVED***

// noOpCloser is like an ioutil.NopCloser, but for an io.Writer.
type noOpCloser struct ***REMOVED***
	w io.Writer
***REMOVED***

func (c noOpCloser) Write(data []byte) (n int, err error) ***REMOVED***
	return c.w.Write(data)
***REMOVED***

func (c noOpCloser) Close() error ***REMOVED***
	return nil
***REMOVED***

// SerializeSymmetricallyEncrypted serializes a symmetrically encrypted packet
// to w and returns a WriteCloser to which the to-be-encrypted packets can be
// written.
// If config is nil, sensible defaults will be used.
func SerializeSymmetricallyEncrypted(w io.Writer, c CipherFunction, key []byte, config *Config) (contents io.WriteCloser, err error) ***REMOVED***
	if c.KeySize() != len(key) ***REMOVED***
		return nil, errors.InvalidArgumentError("SymmetricallyEncrypted.Serialize: bad key length")
	***REMOVED***
	writeCloser := noOpCloser***REMOVED***w***REMOVED***
	ciphertext, err := serializeStreamHeader(writeCloser, packetTypeSymmetricallyEncryptedMDC)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = ciphertext.Write([]byte***REMOVED***symmetricallyEncryptedVersion***REMOVED***)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	block := c.new(key)
	blockSize := block.BlockSize()
	iv := make([]byte, blockSize)
	_, err = config.Random().Read(iv)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	s, prefix := NewOCFBEncrypter(block, iv, OCFBNoResync)
	_, err = ciphertext.Write(prefix)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	plaintext := cipher.StreamWriter***REMOVED***S: s, W: ciphertext***REMOVED***

	h := sha1.New()
	h.Write(iv)
	h.Write(iv[blockSize-2:])
	contents = &seMDCWriter***REMOVED***w: plaintext, h: h***REMOVED***
	return
***REMOVED***
