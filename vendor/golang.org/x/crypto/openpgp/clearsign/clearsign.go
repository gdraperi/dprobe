// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clearsign generates and processes OpenPGP, clear-signed data. See
// RFC 4880, section 7.
//
// Clearsigned messages are cryptographically signed, but the contents of the
// message are kept in plaintext so that it can be read without special tools.
package clearsign // import "golang.org/x/crypto/openpgp/clearsign"

import (
	"bufio"
	"bytes"
	"crypto"
	"hash"
	"io"
	"net/textproto"
	"strconv"

	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
)

// A Block represents a clearsigned message. A signature on a Block can
// be checked by passing Bytes into openpgp.CheckDetachedSignature.
type Block struct ***REMOVED***
	Headers          textproto.MIMEHeader // Optional message headers
	Plaintext        []byte               // The original message text
	Bytes            []byte               // The signed message
	ArmoredSignature *armor.Block         // The signature block
***REMOVED***

// start is the marker which denotes the beginning of a clearsigned message.
var start = []byte("\n-----BEGIN PGP SIGNED MESSAGE-----")

// dashEscape is prefixed to any lines that begin with a hyphen so that they
// can't be confused with endText.
var dashEscape = []byte("- ")

// endText is a marker which denotes the end of the message and the start of
// an armored signature.
var endText = []byte("-----BEGIN PGP SIGNATURE-----")

// end is a marker which denotes the end of the armored signature.
var end = []byte("\n-----END PGP SIGNATURE-----")

var crlf = []byte("\r\n")
var lf = byte('\n')

// getLine returns the first \r\n or \n delineated line from the given byte
// array. The line does not include the \r\n or \n. The remainder of the byte
// array (also not including the new line bytes) is also returned and this will
// always be smaller than the original argument.
func getLine(data []byte) (line, rest []byte) ***REMOVED***
	i := bytes.Index(data, []byte***REMOVED***'\n'***REMOVED***)
	var j int
	if i < 0 ***REMOVED***
		i = len(data)
		j = i
	***REMOVED*** else ***REMOVED***
		j = i + 1
		if i > 0 && data[i-1] == '\r' ***REMOVED***
			i--
		***REMOVED***
	***REMOVED***
	return data[0:i], data[j:]
***REMOVED***

// Decode finds the first clearsigned message in data and returns it, as well
// as the suffix of data which remains after the message.
func Decode(data []byte) (b *Block, rest []byte) ***REMOVED***
	// start begins with a newline. However, at the very beginning of
	// the byte array, we'll accept the start string without it.
	rest = data
	if bytes.HasPrefix(data, start[1:]) ***REMOVED***
		rest = rest[len(start)-1:]
	***REMOVED*** else if i := bytes.Index(data, start); i >= 0 ***REMOVED***
		rest = rest[i+len(start):]
	***REMOVED*** else ***REMOVED***
		return nil, data
	***REMOVED***

	// Consume the start line.
	_, rest = getLine(rest)

	var line []byte
	b = &Block***REMOVED***
		Headers: make(textproto.MIMEHeader),
	***REMOVED***

	// Next come a series of header lines.
	for ***REMOVED***
		// This loop terminates because getLine's second result is
		// always smaller than its argument.
		if len(rest) == 0 ***REMOVED***
			return nil, data
		***REMOVED***
		// An empty line marks the end of the headers.
		if line, rest = getLine(rest); len(line) == 0 ***REMOVED***
			break
		***REMOVED***

		i := bytes.Index(line, []byte***REMOVED***':'***REMOVED***)
		if i == -1 ***REMOVED***
			return nil, data
		***REMOVED***

		key, val := line[0:i], line[i+1:]
		key = bytes.TrimSpace(key)
		val = bytes.TrimSpace(val)
		b.Headers.Add(string(key), string(val))
	***REMOVED***

	firstLine := true
	for ***REMOVED***
		start := rest

		line, rest = getLine(rest)
		if len(line) == 0 && len(rest) == 0 ***REMOVED***
			// No armored data was found, so this isn't a complete message.
			return nil, data
		***REMOVED***
		if bytes.Equal(line, endText) ***REMOVED***
			// Back up to the start of the line because armor expects to see the
			// header line.
			rest = start
			break
		***REMOVED***

		// The final CRLF isn't included in the hash so we don't write it until
		// we've seen the next line.
		if firstLine ***REMOVED***
			firstLine = false
		***REMOVED*** else ***REMOVED***
			b.Bytes = append(b.Bytes, crlf...)
		***REMOVED***

		if bytes.HasPrefix(line, dashEscape) ***REMOVED***
			line = line[2:]
		***REMOVED***
		line = bytes.TrimRight(line, " \t")
		b.Bytes = append(b.Bytes, line...)

		b.Plaintext = append(b.Plaintext, line...)
		b.Plaintext = append(b.Plaintext, lf)
	***REMOVED***

	// We want to find the extent of the armored data (including any newlines at
	// the end).
	i := bytes.Index(rest, end)
	if i == -1 ***REMOVED***
		return nil, data
	***REMOVED***
	i += len(end)
	for i < len(rest) && (rest[i] == '\r' || rest[i] == '\n') ***REMOVED***
		i++
	***REMOVED***
	armored := rest[:i]
	rest = rest[i:]

	var err error
	b.ArmoredSignature, err = armor.Decode(bytes.NewBuffer(armored))
	if err != nil ***REMOVED***
		return nil, data
	***REMOVED***

	return b, rest
***REMOVED***

// A dashEscaper is an io.WriteCloser which processes the body of a clear-signed
// message. The clear-signed message is written to buffered and a hash, suitable
// for signing, is maintained in h.
//
// When closed, an armored signature is created and written to complete the
// message.
type dashEscaper struct ***REMOVED***
	buffered *bufio.Writer
	h        hash.Hash
	hashType crypto.Hash

	atBeginningOfLine bool
	isFirstLine       bool

	whitespace []byte
	byteBuf    []byte // a one byte buffer to save allocations

	privateKey *packet.PrivateKey
	config     *packet.Config
***REMOVED***

func (d *dashEscaper) Write(data []byte) (n int, err error) ***REMOVED***
	for _, b := range data ***REMOVED***
		d.byteBuf[0] = b

		if d.atBeginningOfLine ***REMOVED***
			// The final CRLF isn't included in the hash so we have to wait
			// until this point (the start of the next line) before writing it.
			if !d.isFirstLine ***REMOVED***
				d.h.Write(crlf)
			***REMOVED***
			d.isFirstLine = false
		***REMOVED***

		// Any whitespace at the end of the line has to be removed so we
		// buffer it until we find out whether there's more on this line.
		if b == ' ' || b == '\t' || b == '\r' ***REMOVED***
			d.whitespace = append(d.whitespace, b)
			d.atBeginningOfLine = false
			continue
		***REMOVED***

		if d.atBeginningOfLine ***REMOVED***
			// At the beginning of a line, hyphens have to be escaped.
			if b == '-' ***REMOVED***
				// The signature isn't calculated over the dash-escaped text so
				// the escape is only written to buffered.
				if _, err = d.buffered.Write(dashEscape); err != nil ***REMOVED***
					return
				***REMOVED***
				d.h.Write(d.byteBuf)
				d.atBeginningOfLine = false
			***REMOVED*** else if b == '\n' ***REMOVED***
				// Nothing to do because we delay writing CRLF to the hash.
			***REMOVED*** else ***REMOVED***
				d.h.Write(d.byteBuf)
				d.atBeginningOfLine = false
			***REMOVED***
			if err = d.buffered.WriteByte(b); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if b == '\n' ***REMOVED***
				// We got a raw \n. Drop any trailing whitespace and write a
				// CRLF.
				d.whitespace = d.whitespace[:0]
				// We delay writing CRLF to the hash until the start of the
				// next line.
				if err = d.buffered.WriteByte(b); err != nil ***REMOVED***
					return
				***REMOVED***
				d.atBeginningOfLine = true
			***REMOVED*** else ***REMOVED***
				// Any buffered whitespace wasn't at the end of the line so
				// we need to write it out.
				if len(d.whitespace) > 0 ***REMOVED***
					d.h.Write(d.whitespace)
					if _, err = d.buffered.Write(d.whitespace); err != nil ***REMOVED***
						return
					***REMOVED***
					d.whitespace = d.whitespace[:0]
				***REMOVED***
				d.h.Write(d.byteBuf)
				if err = d.buffered.WriteByte(b); err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	n = len(data)
	return
***REMOVED***

func (d *dashEscaper) Close() (err error) ***REMOVED***
	if !d.atBeginningOfLine ***REMOVED***
		if err = d.buffered.WriteByte(lf); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	sig := new(packet.Signature)
	sig.SigType = packet.SigTypeText
	sig.PubKeyAlgo = d.privateKey.PubKeyAlgo
	sig.Hash = d.hashType
	sig.CreationTime = d.config.Now()
	sig.IssuerKeyId = &d.privateKey.KeyId

	if err = sig.Sign(d.h, d.privateKey, d.config); err != nil ***REMOVED***
		return
	***REMOVED***

	out, err := armor.Encode(d.buffered, "PGP SIGNATURE", nil)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if err = sig.Serialize(out); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = out.Close(); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = d.buffered.Flush(); err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

// Encode returns a WriteCloser which will clear-sign a message with privateKey
// and write it to w. If config is nil, sensible defaults are used.
func Encode(w io.Writer, privateKey *packet.PrivateKey, config *packet.Config) (plaintext io.WriteCloser, err error) ***REMOVED***
	if privateKey.Encrypted ***REMOVED***
		return nil, errors.InvalidArgumentError("signing key is encrypted")
	***REMOVED***

	hashType := config.Hash()
	name := nameOfHash(hashType)
	if len(name) == 0 ***REMOVED***
		return nil, errors.UnsupportedError("unknown hash type: " + strconv.Itoa(int(hashType)))
	***REMOVED***

	if !hashType.Available() ***REMOVED***
		return nil, errors.UnsupportedError("unsupported hash type: " + strconv.Itoa(int(hashType)))
	***REMOVED***
	h := hashType.New()

	buffered := bufio.NewWriter(w)
	// start has a \n at the beginning that we don't want here.
	if _, err = buffered.Write(start[1:]); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = buffered.WriteByte(lf); err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err = buffered.WriteString("Hash: "); err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err = buffered.WriteString(name); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = buffered.WriteByte(lf); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = buffered.WriteByte(lf); err != nil ***REMOVED***
		return
	***REMOVED***

	plaintext = &dashEscaper***REMOVED***
		buffered: buffered,
		h:        h,
		hashType: hashType,

		atBeginningOfLine: true,
		isFirstLine:       true,

		byteBuf: make([]byte, 1),

		privateKey: privateKey,
		config:     config,
	***REMOVED***

	return
***REMOVED***

// nameOfHash returns the OpenPGP name for the given hash, or the empty string
// if the name isn't known. See RFC 4880, section 9.4.
func nameOfHash(h crypto.Hash) string ***REMOVED***
	switch h ***REMOVED***
	case crypto.MD5:
		return "MD5"
	case crypto.SHA1:
		return "SHA1"
	case crypto.RIPEMD160:
		return "RIPEMD160"
	case crypto.SHA224:
		return "SHA224"
	case crypto.SHA256:
		return "SHA256"
	case crypto.SHA384:
		return "SHA384"
	case crypto.SHA512:
		return "SHA512"
	***REMOVED***
	return ""
***REMOVED***
