// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
)

// debugTransport if set, will print packet types as they go over the
// wire. No message decoding is done, to minimize the impact on timing.
const debugTransport = false

const (
	gcmCipherID    = "aes128-gcm@openssh.com"
	aes128cbcID    = "aes128-cbc"
	tripledescbcID = "3des-cbc"
)

// packetConn represents a transport that implements packet based
// operations.
type packetConn interface ***REMOVED***
	// Encrypt and send a packet of data to the remote peer.
	writePacket(packet []byte) error

	// Read a packet from the connection. The read is blocking,
	// i.e. if error is nil, then the returned byte slice is
	// always non-empty.
	readPacket() ([]byte, error)

	// Close closes the write-side of the connection.
	Close() error
***REMOVED***

// transport is the keyingTransport that implements the SSH packet
// protocol.
type transport struct ***REMOVED***
	reader connectionState
	writer connectionState

	bufReader *bufio.Reader
	bufWriter *bufio.Writer
	rand      io.Reader
	isClient  bool
	io.Closer
***REMOVED***

// packetCipher represents a combination of SSH encryption/MAC
// protocol.  A single instance should be used for one direction only.
type packetCipher interface ***REMOVED***
	// writePacket encrypts the packet and writes it to w. The
	// contents of the packet are generally scrambled.
	writePacket(seqnum uint32, w io.Writer, rand io.Reader, packet []byte) error

	// readPacket reads and decrypts a packet of data. The
	// returned packet may be overwritten by future calls of
	// readPacket.
	readPacket(seqnum uint32, r io.Reader) ([]byte, error)
***REMOVED***

// connectionState represents one side (read or write) of the
// connection. This is necessary because each direction has its own
// keys, and can even have its own algorithms
type connectionState struct ***REMOVED***
	packetCipher
	seqNum           uint32
	dir              direction
	pendingKeyChange chan packetCipher
***REMOVED***

// prepareKeyChange sets up key material for a keychange. The key changes in
// both directions are triggered by reading and writing a msgNewKey packet
// respectively.
func (t *transport) prepareKeyChange(algs *algorithms, kexResult *kexResult) error ***REMOVED***
	ciph, err := newPacketCipher(t.reader.dir, algs.r, kexResult)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.reader.pendingKeyChange <- ciph

	ciph, err = newPacketCipher(t.writer.dir, algs.w, kexResult)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.writer.pendingKeyChange <- ciph

	return nil
***REMOVED***

func (t *transport) printPacket(p []byte, write bool) ***REMOVED***
	if len(p) == 0 ***REMOVED***
		return
	***REMOVED***
	who := "server"
	if t.isClient ***REMOVED***
		who = "client"
	***REMOVED***
	what := "read"
	if write ***REMOVED***
		what = "write"
	***REMOVED***

	log.Println(what, who, p[0])
***REMOVED***

// Read and decrypt next packet.
func (t *transport) readPacket() (p []byte, err error) ***REMOVED***
	for ***REMOVED***
		p, err = t.reader.readPacket(t.bufReader)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if len(p) == 0 || (p[0] != msgIgnore && p[0] != msgDebug) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if debugTransport ***REMOVED***
		t.printPacket(p, false)
	***REMOVED***

	return p, err
***REMOVED***

func (s *connectionState) readPacket(r *bufio.Reader) ([]byte, error) ***REMOVED***
	packet, err := s.packetCipher.readPacket(s.seqNum, r)
	s.seqNum++
	if err == nil && len(packet) == 0 ***REMOVED***
		err = errors.New("ssh: zero length packet")
	***REMOVED***

	if len(packet) > 0 ***REMOVED***
		switch packet[0] ***REMOVED***
		case msgNewKeys:
			select ***REMOVED***
			case cipher := <-s.pendingKeyChange:
				s.packetCipher = cipher
			default:
				return nil, errors.New("ssh: got bogus newkeys message")
			***REMOVED***

		case msgDisconnect:
			// Transform a disconnect message into an
			// error. Since this is lowest level at which
			// we interpret message types, doing it here
			// ensures that we don't have to handle it
			// elsewhere.
			var msg disconnectMsg
			if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return nil, &msg
		***REMOVED***
	***REMOVED***

	// The packet may point to an internal buffer, so copy the
	// packet out here.
	fresh := make([]byte, len(packet))
	copy(fresh, packet)

	return fresh, err
***REMOVED***

func (t *transport) writePacket(packet []byte) error ***REMOVED***
	if debugTransport ***REMOVED***
		t.printPacket(packet, true)
	***REMOVED***
	return t.writer.writePacket(t.bufWriter, t.rand, packet)
***REMOVED***

func (s *connectionState) writePacket(w *bufio.Writer, rand io.Reader, packet []byte) error ***REMOVED***
	changeKeys := len(packet) > 0 && packet[0] == msgNewKeys

	err := s.packetCipher.writePacket(s.seqNum, w, rand, packet)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = w.Flush(); err != nil ***REMOVED***
		return err
	***REMOVED***
	s.seqNum++
	if changeKeys ***REMOVED***
		select ***REMOVED***
		case cipher := <-s.pendingKeyChange:
			s.packetCipher = cipher
		default:
			panic("ssh: no key material for msgNewKeys")
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func newTransport(rwc io.ReadWriteCloser, rand io.Reader, isClient bool) *transport ***REMOVED***
	t := &transport***REMOVED***
		bufReader: bufio.NewReader(rwc),
		bufWriter: bufio.NewWriter(rwc),
		rand:      rand,
		reader: connectionState***REMOVED***
			packetCipher:     &streamPacketCipher***REMOVED***cipher: noneCipher***REMOVED******REMOVED******REMOVED***,
			pendingKeyChange: make(chan packetCipher, 1),
		***REMOVED***,
		writer: connectionState***REMOVED***
			packetCipher:     &streamPacketCipher***REMOVED***cipher: noneCipher***REMOVED******REMOVED******REMOVED***,
			pendingKeyChange: make(chan packetCipher, 1),
		***REMOVED***,
		Closer: rwc,
	***REMOVED***
	t.isClient = isClient

	if isClient ***REMOVED***
		t.reader.dir = serverKeys
		t.writer.dir = clientKeys
	***REMOVED*** else ***REMOVED***
		t.reader.dir = clientKeys
		t.writer.dir = serverKeys
	***REMOVED***

	return t
***REMOVED***

type direction struct ***REMOVED***
	ivTag     []byte
	keyTag    []byte
	macKeyTag []byte
***REMOVED***

var (
	serverKeys = direction***REMOVED***[]byte***REMOVED***'B'***REMOVED***, []byte***REMOVED***'D'***REMOVED***, []byte***REMOVED***'F'***REMOVED******REMOVED***
	clientKeys = direction***REMOVED***[]byte***REMOVED***'A'***REMOVED***, []byte***REMOVED***'C'***REMOVED***, []byte***REMOVED***'E'***REMOVED******REMOVED***
)

// setupKeys sets the cipher and MAC keys from kex.K, kex.H and sessionId, as
// described in RFC 4253, section 6.4. direction should either be serverKeys
// (to setup server->client keys) or clientKeys (for client->server keys).
func newPacketCipher(d direction, algs directionAlgorithms, kex *kexResult) (packetCipher, error) ***REMOVED***
	cipherMode := cipherModes[algs.Cipher]
	macMode := macModes[algs.MAC]

	iv := make([]byte, cipherMode.ivSize)
	key := make([]byte, cipherMode.keySize)
	macKey := make([]byte, macMode.keySize)

	generateKeyMaterial(iv, d.ivTag, kex)
	generateKeyMaterial(key, d.keyTag, kex)
	generateKeyMaterial(macKey, d.macKeyTag, kex)

	return cipherModes[algs.Cipher].create(key, iv, macKey, algs)
***REMOVED***

// generateKeyMaterial fills out with key material generated from tag, K, H
// and sessionId, as specified in RFC 4253, section 7.2.
func generateKeyMaterial(out, tag []byte, r *kexResult) ***REMOVED***
	var digestsSoFar []byte

	h := r.Hash.New()
	for len(out) > 0 ***REMOVED***
		h.Reset()
		h.Write(r.K)
		h.Write(r.H)

		if len(digestsSoFar) == 0 ***REMOVED***
			h.Write(tag)
			h.Write(r.SessionID)
		***REMOVED*** else ***REMOVED***
			h.Write(digestsSoFar)
		***REMOVED***

		digest := h.Sum(nil)
		n := copy(out, digest)
		out = out[n:]
		if len(out) > 0 ***REMOVED***
			digestsSoFar = append(digestsSoFar, digest...)
		***REMOVED***
	***REMOVED***
***REMOVED***

const packageVersion = "SSH-2.0-Go"

// Sends and receives a version line.  The versionLine string should
// be US ASCII, start with "SSH-2.0-", and should not include a
// newline. exchangeVersions returns the other side's version line.
func exchangeVersions(rw io.ReadWriter, versionLine []byte) (them []byte, err error) ***REMOVED***
	// Contrary to the RFC, we do not ignore lines that don't
	// start with "SSH-2.0-" to make the library usable with
	// nonconforming servers.
	for _, c := range versionLine ***REMOVED***
		// The spec disallows non US-ASCII chars, and
		// specifically forbids null chars.
		if c < 32 ***REMOVED***
			return nil, errors.New("ssh: junk character in version line")
		***REMOVED***
	***REMOVED***
	if _, err = rw.Write(append(versionLine, '\r', '\n')); err != nil ***REMOVED***
		return
	***REMOVED***

	them, err = readVersion(rw)
	return them, err
***REMOVED***

// maxVersionStringBytes is the maximum number of bytes that we'll
// accept as a version string. RFC 4253 section 4.2 limits this at 255
// chars
const maxVersionStringBytes = 255

// Read version string as specified by RFC 4253, section 4.2.
func readVersion(r io.Reader) ([]byte, error) ***REMOVED***
	versionString := make([]byte, 0, 64)
	var ok bool
	var buf [1]byte

	for length := 0; length < maxVersionStringBytes; length++ ***REMOVED***
		_, err := io.ReadFull(r, buf[:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// The RFC says that the version should be terminated with \r\n
		// but several SSH servers actually only send a \n.
		if buf[0] == '\n' ***REMOVED***
			if !bytes.HasPrefix(versionString, []byte("SSH-")) ***REMOVED***
				// RFC 4253 says we need to ignore all version string lines
				// except the one containing the SSH version (provided that
				// all the lines do not exceed 255 bytes in total).
				versionString = versionString[:0]
				continue
			***REMOVED***
			ok = true
			break
		***REMOVED***

		// non ASCII chars are disallowed, but we are lenient,
		// since Go doesn't use null-terminated strings.

		// The RFC allows a comment after a space, however,
		// all of it (version and comments) goes into the
		// session hash.
		versionString = append(versionString, buf[0])
	***REMOVED***

	if !ok ***REMOVED***
		return nil, errors.New("ssh: overflow reading version string")
	***REMOVED***

	// There might be a '\r' on the end which we should remove.
	if len(versionString) > 0 && versionString[len(versionString)-1] == '\r' ***REMOVED***
		versionString = versionString[:len(versionString)-1]
	***REMOVED***
	return versionString, nil
***REMOVED***
