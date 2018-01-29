// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/internal/chacha20"
	"golang.org/x/crypto/poly1305"
)

const (
	packetSizeMultiple = 16 // TODO(huin) this should be determined by the cipher.

	// RFC 4253 section 6.1 defines a minimum packet size of 32768 that implementations
	// MUST be able to process (plus a few more kilobytes for padding and mac). The RFC
	// indicates implementations SHOULD be able to handle larger packet sizes, but then
	// waffles on about reasonable limits.
	//
	// OpenSSH caps their maxPacket at 256kB so we choose to do
	// the same. maxPacket is also used to ensure that uint32
	// length fields do not overflow, so it should remain well
	// below 4G.
	maxPacket = 256 * 1024
)

// noneCipher implements cipher.Stream and provides no encryption. It is used
// by the transport before the first key-exchange.
type noneCipher struct***REMOVED******REMOVED***

func (c noneCipher) XORKeyStream(dst, src []byte) ***REMOVED***
	copy(dst, src)
***REMOVED***

func newAESCTR(key, iv []byte) (cipher.Stream, error) ***REMOVED***
	c, err := aes.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return cipher.NewCTR(c, iv), nil
***REMOVED***

func newRC4(key, iv []byte) (cipher.Stream, error) ***REMOVED***
	return rc4.NewCipher(key)
***REMOVED***

type cipherMode struct ***REMOVED***
	keySize int
	ivSize  int
	create  func(key, iv []byte, macKey []byte, algs directionAlgorithms) (packetCipher, error)
***REMOVED***

func streamCipherMode(skip int, createFunc func(key, iv []byte) (cipher.Stream, error)) func(key, iv []byte, macKey []byte, algs directionAlgorithms) (packetCipher, error) ***REMOVED***
	return func(key, iv, macKey []byte, algs directionAlgorithms) (packetCipher, error) ***REMOVED***
		stream, err := createFunc(key, iv)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var streamDump []byte
		if skip > 0 ***REMOVED***
			streamDump = make([]byte, 512)
		***REMOVED***

		for remainingToDump := skip; remainingToDump > 0; ***REMOVED***
			dumpThisTime := remainingToDump
			if dumpThisTime > len(streamDump) ***REMOVED***
				dumpThisTime = len(streamDump)
			***REMOVED***
			stream.XORKeyStream(streamDump[:dumpThisTime], streamDump[:dumpThisTime])
			remainingToDump -= dumpThisTime
		***REMOVED***

		mac := macModes[algs.MAC].new(macKey)
		return &streamPacketCipher***REMOVED***
			mac:       mac,
			etm:       macModes[algs.MAC].etm,
			macResult: make([]byte, mac.Size()),
			cipher:    stream,
		***REMOVED***, nil
	***REMOVED***
***REMOVED***

// cipherModes documents properties of supported ciphers. Ciphers not included
// are not supported and will not be negotiated, even if explicitly requested in
// ClientConfig.Crypto.Ciphers.
var cipherModes = map[string]*cipherMode***REMOVED***
	// Ciphers from RFC4344, which introduced many CTR-based ciphers. Algorithms
	// are defined in the order specified in the RFC.
	"aes128-ctr": ***REMOVED***16, aes.BlockSize, streamCipherMode(0, newAESCTR)***REMOVED***,
	"aes192-ctr": ***REMOVED***24, aes.BlockSize, streamCipherMode(0, newAESCTR)***REMOVED***,
	"aes256-ctr": ***REMOVED***32, aes.BlockSize, streamCipherMode(0, newAESCTR)***REMOVED***,

	// Ciphers from RFC4345, which introduces security-improved arcfour ciphers.
	// They are defined in the order specified in the RFC.
	"arcfour128": ***REMOVED***16, 0, streamCipherMode(1536, newRC4)***REMOVED***,
	"arcfour256": ***REMOVED***32, 0, streamCipherMode(1536, newRC4)***REMOVED***,

	// Cipher defined in RFC 4253, which describes SSH Transport Layer Protocol.
	// Note that this cipher is not safe, as stated in RFC 4253: "Arcfour (and
	// RC4) has problems with weak keys, and should be used with caution."
	// RFC4345 introduces improved versions of Arcfour.
	"arcfour": ***REMOVED***16, 0, streamCipherMode(0, newRC4)***REMOVED***,

	// AEAD ciphers
	gcmCipherID:        ***REMOVED***16, 12, newGCMCipher***REMOVED***,
	chacha20Poly1305ID: ***REMOVED***64, 0, newChaCha20Cipher***REMOVED***,

	// CBC mode is insecure and so is not included in the default config.
	// (See http://www.isg.rhul.ac.uk/~kp/SandPfinal.pdf). If absolutely
	// needed, it's possible to specify a custom Config to enable it.
	// You should expect that an active attacker can recover plaintext if
	// you do.
	aes128cbcID: ***REMOVED***16, aes.BlockSize, newAESCBCCipher***REMOVED***,

	// 3des-cbc is insecure and is not included in the default
	// config.
	tripledescbcID: ***REMOVED***24, des.BlockSize, newTripleDESCBCCipher***REMOVED***,
***REMOVED***

// prefixLen is the length of the packet prefix that contains the packet length
// and number of padding bytes.
const prefixLen = 5

// streamPacketCipher is a packetCipher using a stream cipher.
type streamPacketCipher struct ***REMOVED***
	mac    hash.Hash
	cipher cipher.Stream
	etm    bool

	// The following members are to avoid per-packet allocations.
	prefix      [prefixLen]byte
	seqNumBytes [4]byte
	padding     [2 * packetSizeMultiple]byte
	packetData  []byte
	macResult   []byte
***REMOVED***

// readPacket reads and decrypt a single packet from the reader argument.
func (s *streamPacketCipher) readPacket(seqNum uint32, r io.Reader) ([]byte, error) ***REMOVED***
	if _, err := io.ReadFull(r, s.prefix[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var encryptedPaddingLength [1]byte
	if s.mac != nil && s.etm ***REMOVED***
		copy(encryptedPaddingLength[:], s.prefix[4:5])
		s.cipher.XORKeyStream(s.prefix[4:5], s.prefix[4:5])
	***REMOVED*** else ***REMOVED***
		s.cipher.XORKeyStream(s.prefix[:], s.prefix[:])
	***REMOVED***

	length := binary.BigEndian.Uint32(s.prefix[0:4])
	paddingLength := uint32(s.prefix[4])

	var macSize uint32
	if s.mac != nil ***REMOVED***
		s.mac.Reset()
		binary.BigEndian.PutUint32(s.seqNumBytes[:], seqNum)
		s.mac.Write(s.seqNumBytes[:])
		if s.etm ***REMOVED***
			s.mac.Write(s.prefix[:4])
			s.mac.Write(encryptedPaddingLength[:])
		***REMOVED*** else ***REMOVED***
			s.mac.Write(s.prefix[:])
		***REMOVED***
		macSize = uint32(s.mac.Size())
	***REMOVED***

	if length <= paddingLength+1 ***REMOVED***
		return nil, errors.New("ssh: invalid packet length, packet too small")
	***REMOVED***

	if length > maxPacket ***REMOVED***
		return nil, errors.New("ssh: invalid packet length, packet too large")
	***REMOVED***

	// the maxPacket check above ensures that length-1+macSize
	// does not overflow.
	if uint32(cap(s.packetData)) < length-1+macSize ***REMOVED***
		s.packetData = make([]byte, length-1+macSize)
	***REMOVED*** else ***REMOVED***
		s.packetData = s.packetData[:length-1+macSize]
	***REMOVED***

	if _, err := io.ReadFull(r, s.packetData); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mac := s.packetData[length-1:]
	data := s.packetData[:length-1]

	if s.mac != nil && s.etm ***REMOVED***
		s.mac.Write(data)
	***REMOVED***

	s.cipher.XORKeyStream(data, data)

	if s.mac != nil ***REMOVED***
		if !s.etm ***REMOVED***
			s.mac.Write(data)
		***REMOVED***
		s.macResult = s.mac.Sum(s.macResult[:0])
		if subtle.ConstantTimeCompare(s.macResult, mac) != 1 ***REMOVED***
			return nil, errors.New("ssh: MAC failure")
		***REMOVED***
	***REMOVED***

	return s.packetData[:length-paddingLength-1], nil
***REMOVED***

// writePacket encrypts and sends a packet of data to the writer argument
func (s *streamPacketCipher) writePacket(seqNum uint32, w io.Writer, rand io.Reader, packet []byte) error ***REMOVED***
	if len(packet) > maxPacket ***REMOVED***
		return errors.New("ssh: packet too large")
	***REMOVED***

	aadlen := 0
	if s.mac != nil && s.etm ***REMOVED***
		// packet length is not encrypted for EtM modes
		aadlen = 4
	***REMOVED***

	paddingLength := packetSizeMultiple - (prefixLen+len(packet)-aadlen)%packetSizeMultiple
	if paddingLength < 4 ***REMOVED***
		paddingLength += packetSizeMultiple
	***REMOVED***

	length := len(packet) + 1 + paddingLength
	binary.BigEndian.PutUint32(s.prefix[:], uint32(length))
	s.prefix[4] = byte(paddingLength)
	padding := s.padding[:paddingLength]
	if _, err := io.ReadFull(rand, padding); err != nil ***REMOVED***
		return err
	***REMOVED***

	if s.mac != nil ***REMOVED***
		s.mac.Reset()
		binary.BigEndian.PutUint32(s.seqNumBytes[:], seqNum)
		s.mac.Write(s.seqNumBytes[:])

		if s.etm ***REMOVED***
			// For EtM algorithms, the packet length must stay unencrypted,
			// but the following data (padding length) must be encrypted
			s.cipher.XORKeyStream(s.prefix[4:5], s.prefix[4:5])
		***REMOVED***

		s.mac.Write(s.prefix[:])

		if !s.etm ***REMOVED***
			// For non-EtM algorithms, the algorithm is applied on unencrypted data
			s.mac.Write(packet)
			s.mac.Write(padding)
		***REMOVED***
	***REMOVED***

	if !(s.mac != nil && s.etm) ***REMOVED***
		// For EtM algorithms, the padding length has already been encrypted
		// and the packet length must remain unencrypted
		s.cipher.XORKeyStream(s.prefix[:], s.prefix[:])
	***REMOVED***

	s.cipher.XORKeyStream(packet, packet)
	s.cipher.XORKeyStream(padding, padding)

	if s.mac != nil && s.etm ***REMOVED***
		// For EtM algorithms, packet and padding must be encrypted
		s.mac.Write(packet)
		s.mac.Write(padding)
	***REMOVED***

	if _, err := w.Write(s.prefix[:]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.Write(packet); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.Write(padding); err != nil ***REMOVED***
		return err
	***REMOVED***

	if s.mac != nil ***REMOVED***
		s.macResult = s.mac.Sum(s.macResult[:0])
		if _, err := w.Write(s.macResult); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type gcmCipher struct ***REMOVED***
	aead   cipher.AEAD
	prefix [4]byte
	iv     []byte
	buf    []byte
***REMOVED***

func newGCMCipher(key, iv, unusedMacKey []byte, unusedAlgs directionAlgorithms) (packetCipher, error) ***REMOVED***
	c, err := aes.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	aead, err := cipher.NewGCM(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &gcmCipher***REMOVED***
		aead: aead,
		iv:   iv,
	***REMOVED***, nil
***REMOVED***

const gcmTagSize = 16

func (c *gcmCipher) writePacket(seqNum uint32, w io.Writer, rand io.Reader, packet []byte) error ***REMOVED***
	// Pad out to multiple of 16 bytes. This is different from the
	// stream cipher because that encrypts the length too.
	padding := byte(packetSizeMultiple - (1+len(packet))%packetSizeMultiple)
	if padding < 4 ***REMOVED***
		padding += packetSizeMultiple
	***REMOVED***

	length := uint32(len(packet) + int(padding) + 1)
	binary.BigEndian.PutUint32(c.prefix[:], length)
	if _, err := w.Write(c.prefix[:]); err != nil ***REMOVED***
		return err
	***REMOVED***

	if cap(c.buf) < int(length) ***REMOVED***
		c.buf = make([]byte, length)
	***REMOVED*** else ***REMOVED***
		c.buf = c.buf[:length]
	***REMOVED***

	c.buf[0] = padding
	copy(c.buf[1:], packet)
	if _, err := io.ReadFull(rand, c.buf[1+len(packet):]); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.buf = c.aead.Seal(c.buf[:0], c.iv, c.buf, c.prefix[:])
	if _, err := w.Write(c.buf); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.incIV()

	return nil
***REMOVED***

func (c *gcmCipher) incIV() ***REMOVED***
	for i := 4 + 7; i >= 4; i-- ***REMOVED***
		c.iv[i]++
		if c.iv[i] != 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *gcmCipher) readPacket(seqNum uint32, r io.Reader) ([]byte, error) ***REMOVED***
	if _, err := io.ReadFull(r, c.prefix[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	length := binary.BigEndian.Uint32(c.prefix[:])
	if length > maxPacket ***REMOVED***
		return nil, errors.New("ssh: max packet length exceeded")
	***REMOVED***

	if cap(c.buf) < int(length+gcmTagSize) ***REMOVED***
		c.buf = make([]byte, length+gcmTagSize)
	***REMOVED*** else ***REMOVED***
		c.buf = c.buf[:length+gcmTagSize]
	***REMOVED***

	if _, err := io.ReadFull(r, c.buf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	plain, err := c.aead.Open(c.buf[:0], c.iv, c.buf, c.prefix[:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.incIV()

	padding := plain[0]
	if padding < 4 ***REMOVED***
		// padding is a byte, so it automatically satisfies
		// the maximum size, which is 255.
		return nil, fmt.Errorf("ssh: illegal padding %d", padding)
	***REMOVED***

	if int(padding+1) >= len(plain) ***REMOVED***
		return nil, fmt.Errorf("ssh: padding %d too large", padding)
	***REMOVED***
	plain = plain[1 : length-uint32(padding)]
	return plain, nil
***REMOVED***

// cbcCipher implements aes128-cbc cipher defined in RFC 4253 section 6.1
type cbcCipher struct ***REMOVED***
	mac       hash.Hash
	macSize   uint32
	decrypter cipher.BlockMode
	encrypter cipher.BlockMode

	// The following members are to avoid per-packet allocations.
	seqNumBytes [4]byte
	packetData  []byte
	macResult   []byte

	// Amount of data we should still read to hide which
	// verification error triggered.
	oracleCamouflage uint32
***REMOVED***

func newCBCCipher(c cipher.Block, key, iv, macKey []byte, algs directionAlgorithms) (packetCipher, error) ***REMOVED***
	cbc := &cbcCipher***REMOVED***
		mac:        macModes[algs.MAC].new(macKey),
		decrypter:  cipher.NewCBCDecrypter(c, iv),
		encrypter:  cipher.NewCBCEncrypter(c, iv),
		packetData: make([]byte, 1024),
	***REMOVED***
	if cbc.mac != nil ***REMOVED***
		cbc.macSize = uint32(cbc.mac.Size())
	***REMOVED***

	return cbc, nil
***REMOVED***

func newAESCBCCipher(key, iv, macKey []byte, algs directionAlgorithms) (packetCipher, error) ***REMOVED***
	c, err := aes.NewCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cbc, err := newCBCCipher(c, key, iv, macKey, algs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return cbc, nil
***REMOVED***

func newTripleDESCBCCipher(key, iv, macKey []byte, algs directionAlgorithms) (packetCipher, error) ***REMOVED***
	c, err := des.NewTripleDESCipher(key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cbc, err := newCBCCipher(c, key, iv, macKey, algs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return cbc, nil
***REMOVED***

func maxUInt32(a, b int) uint32 ***REMOVED***
	if a > b ***REMOVED***
		return uint32(a)
	***REMOVED***
	return uint32(b)
***REMOVED***

const (
	cbcMinPacketSizeMultiple = 8
	cbcMinPacketSize         = 16
	cbcMinPaddingSize        = 4
)

// cbcError represents a verification error that may leak information.
type cbcError string

func (e cbcError) Error() string ***REMOVED*** return string(e) ***REMOVED***

func (c *cbcCipher) readPacket(seqNum uint32, r io.Reader) ([]byte, error) ***REMOVED***
	p, err := c.readPacketLeaky(seqNum, r)
	if err != nil ***REMOVED***
		if _, ok := err.(cbcError); ok ***REMOVED***
			// Verification error: read a fixed amount of
			// data, to make distinguishing between
			// failing MAC and failing length check more
			// difficult.
			io.CopyN(ioutil.Discard, r, int64(c.oracleCamouflage))
		***REMOVED***
	***REMOVED***
	return p, err
***REMOVED***

func (c *cbcCipher) readPacketLeaky(seqNum uint32, r io.Reader) ([]byte, error) ***REMOVED***
	blockSize := c.decrypter.BlockSize()

	// Read the header, which will include some of the subsequent data in the
	// case of block ciphers - this is copied back to the payload later.
	// How many bytes of payload/padding will be read with this first read.
	firstBlockLength := uint32((prefixLen + blockSize - 1) / blockSize * blockSize)
	firstBlock := c.packetData[:firstBlockLength]
	if _, err := io.ReadFull(r, firstBlock); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.oracleCamouflage = maxPacket + 4 + c.macSize - firstBlockLength

	c.decrypter.CryptBlocks(firstBlock, firstBlock)
	length := binary.BigEndian.Uint32(firstBlock[:4])
	if length > maxPacket ***REMOVED***
		return nil, cbcError("ssh: packet too large")
	***REMOVED***
	if length+4 < maxUInt32(cbcMinPacketSize, blockSize) ***REMOVED***
		// The minimum size of a packet is 16 (or the cipher block size, whichever
		// is larger) bytes.
		return nil, cbcError("ssh: packet too small")
	***REMOVED***
	// The length of the packet (including the length field but not the MAC) must
	// be a multiple of the block size or 8, whichever is larger.
	if (length+4)%maxUInt32(cbcMinPacketSizeMultiple, blockSize) != 0 ***REMOVED***
		return nil, cbcError("ssh: invalid packet length multiple")
	***REMOVED***

	paddingLength := uint32(firstBlock[4])
	if paddingLength < cbcMinPaddingSize || length <= paddingLength+1 ***REMOVED***
		return nil, cbcError("ssh: invalid packet length")
	***REMOVED***

	// Positions within the c.packetData buffer:
	macStart := 4 + length
	paddingStart := macStart - paddingLength

	// Entire packet size, starting before length, ending at end of mac.
	entirePacketSize := macStart + c.macSize

	// Ensure c.packetData is large enough for the entire packet data.
	if uint32(cap(c.packetData)) < entirePacketSize ***REMOVED***
		// Still need to upsize and copy, but this should be rare at runtime, only
		// on upsizing the packetData buffer.
		c.packetData = make([]byte, entirePacketSize)
		copy(c.packetData, firstBlock)
	***REMOVED*** else ***REMOVED***
		c.packetData = c.packetData[:entirePacketSize]
	***REMOVED***

	n, err := io.ReadFull(r, c.packetData[firstBlockLength:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.oracleCamouflage -= uint32(n)

	remainingCrypted := c.packetData[firstBlockLength:macStart]
	c.decrypter.CryptBlocks(remainingCrypted, remainingCrypted)

	mac := c.packetData[macStart:]
	if c.mac != nil ***REMOVED***
		c.mac.Reset()
		binary.BigEndian.PutUint32(c.seqNumBytes[:], seqNum)
		c.mac.Write(c.seqNumBytes[:])
		c.mac.Write(c.packetData[:macStart])
		c.macResult = c.mac.Sum(c.macResult[:0])
		if subtle.ConstantTimeCompare(c.macResult, mac) != 1 ***REMOVED***
			return nil, cbcError("ssh: MAC failure")
		***REMOVED***
	***REMOVED***

	return c.packetData[prefixLen:paddingStart], nil
***REMOVED***

func (c *cbcCipher) writePacket(seqNum uint32, w io.Writer, rand io.Reader, packet []byte) error ***REMOVED***
	effectiveBlockSize := maxUInt32(cbcMinPacketSizeMultiple, c.encrypter.BlockSize())

	// Length of encrypted portion of the packet (header, payload, padding).
	// Enforce minimum padding and packet size.
	encLength := maxUInt32(prefixLen+len(packet)+cbcMinPaddingSize, cbcMinPaddingSize)
	// Enforce block size.
	encLength = (encLength + effectiveBlockSize - 1) / effectiveBlockSize * effectiveBlockSize

	length := encLength - 4
	paddingLength := int(length) - (1 + len(packet))

	// Overall buffer contains: header, payload, padding, mac.
	// Space for the MAC is reserved in the capacity but not the slice length.
	bufferSize := encLength + c.macSize
	if uint32(cap(c.packetData)) < bufferSize ***REMOVED***
		c.packetData = make([]byte, encLength, bufferSize)
	***REMOVED*** else ***REMOVED***
		c.packetData = c.packetData[:encLength]
	***REMOVED***

	p := c.packetData

	// Packet header.
	binary.BigEndian.PutUint32(p, length)
	p = p[4:]
	p[0] = byte(paddingLength)

	// Payload.
	p = p[1:]
	copy(p, packet)

	// Padding.
	p = p[len(packet):]
	if _, err := io.ReadFull(rand, p); err != nil ***REMOVED***
		return err
	***REMOVED***

	if c.mac != nil ***REMOVED***
		c.mac.Reset()
		binary.BigEndian.PutUint32(c.seqNumBytes[:], seqNum)
		c.mac.Write(c.seqNumBytes[:])
		c.mac.Write(c.packetData)
		// The MAC is now appended into the capacity reserved for it earlier.
		c.packetData = c.mac.Sum(c.packetData)
	***REMOVED***

	c.encrypter.CryptBlocks(c.packetData[:encLength], c.packetData[:encLength])

	if _, err := w.Write(c.packetData); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

const chacha20Poly1305ID = "chacha20-poly1305@openssh.com"

// chacha20Poly1305Cipher implements the chacha20-poly1305@openssh.com
// AEAD, which is described here:
//
//   https://tools.ietf.org/html/draft-josefsson-ssh-chacha20-poly1305-openssh-00
//
// the methods here also implement padding, which RFC4253 Section 6
// also requires of stream ciphers.
type chacha20Poly1305Cipher struct ***REMOVED***
	lengthKey  [32]byte
	contentKey [32]byte
	buf        []byte
***REMOVED***

func newChaCha20Cipher(key, unusedIV, unusedMACKey []byte, unusedAlgs directionAlgorithms) (packetCipher, error) ***REMOVED***
	if len(key) != 64 ***REMOVED***
		panic(len(key))
	***REMOVED***

	c := &chacha20Poly1305Cipher***REMOVED***
		buf: make([]byte, 256),
	***REMOVED***

	copy(c.contentKey[:], key[:32])
	copy(c.lengthKey[:], key[32:])
	return c, nil
***REMOVED***

// The Poly1305 key is obtained by encrypting 32 0-bytes.
var chacha20PolyKeyInput [32]byte

func (c *chacha20Poly1305Cipher) readPacket(seqNum uint32, r io.Reader) ([]byte, error) ***REMOVED***
	var counter [16]byte
	binary.BigEndian.PutUint64(counter[8:], uint64(seqNum))

	var polyKey [32]byte
	chacha20.XORKeyStream(polyKey[:], chacha20PolyKeyInput[:], &counter, &c.contentKey)

	encryptedLength := c.buf[:4]
	if _, err := io.ReadFull(r, encryptedLength); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var lenBytes [4]byte
	chacha20.XORKeyStream(lenBytes[:], encryptedLength, &counter, &c.lengthKey)

	length := binary.BigEndian.Uint32(lenBytes[:])
	if length > maxPacket ***REMOVED***
		return nil, errors.New("ssh: invalid packet length, packet too large")
	***REMOVED***

	contentEnd := 4 + length
	packetEnd := contentEnd + poly1305.TagSize
	if uint32(cap(c.buf)) < packetEnd ***REMOVED***
		c.buf = make([]byte, packetEnd)
		copy(c.buf[:], encryptedLength)
	***REMOVED*** else ***REMOVED***
		c.buf = c.buf[:packetEnd]
	***REMOVED***

	if _, err := io.ReadFull(r, c.buf[4:packetEnd]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var mac [poly1305.TagSize]byte
	copy(mac[:], c.buf[contentEnd:packetEnd])
	if !poly1305.Verify(&mac, c.buf[:contentEnd], &polyKey) ***REMOVED***
		return nil, errors.New("ssh: MAC failure")
	***REMOVED***

	counter[0] = 1

	plain := c.buf[4:contentEnd]
	chacha20.XORKeyStream(plain, plain, &counter, &c.contentKey)

	padding := plain[0]
	if padding < 4 ***REMOVED***
		// padding is a byte, so it automatically satisfies
		// the maximum size, which is 255.
		return nil, fmt.Errorf("ssh: illegal padding %d", padding)
	***REMOVED***

	if int(padding)+1 >= len(plain) ***REMOVED***
		return nil, fmt.Errorf("ssh: padding %d too large", padding)
	***REMOVED***

	plain = plain[1 : len(plain)-int(padding)]

	return plain, nil
***REMOVED***

func (c *chacha20Poly1305Cipher) writePacket(seqNum uint32, w io.Writer, rand io.Reader, payload []byte) error ***REMOVED***
	var counter [16]byte
	binary.BigEndian.PutUint64(counter[8:], uint64(seqNum))

	var polyKey [32]byte
	chacha20.XORKeyStream(polyKey[:], chacha20PolyKeyInput[:], &counter, &c.contentKey)

	// There is no blocksize, so fall back to multiple of 8 byte
	// padding, as described in RFC 4253, Sec 6.
	const packetSizeMultiple = 8

	padding := packetSizeMultiple - (1+len(payload))%packetSizeMultiple
	if padding < 4 ***REMOVED***
		padding += packetSizeMultiple
	***REMOVED***

	// size (4 bytes), padding (1), payload, padding, tag.
	totalLength := 4 + 1 + len(payload) + padding + poly1305.TagSize
	if cap(c.buf) < totalLength ***REMOVED***
		c.buf = make([]byte, totalLength)
	***REMOVED*** else ***REMOVED***
		c.buf = c.buf[:totalLength]
	***REMOVED***

	binary.BigEndian.PutUint32(c.buf, uint32(1+len(payload)+padding))
	chacha20.XORKeyStream(c.buf, c.buf[:4], &counter, &c.lengthKey)
	c.buf[4] = byte(padding)
	copy(c.buf[5:], payload)
	packetEnd := 5 + len(payload) + padding
	if _, err := io.ReadFull(rand, c.buf[5+len(payload):packetEnd]); err != nil ***REMOVED***
		return err
	***REMOVED***

	counter[0] = 1
	chacha20.XORKeyStream(c.buf[4:], c.buf[4:packetEnd], &counter, &c.contentKey)

	var mac [poly1305.TagSize]byte
	poly1305.Sum(&mac, c.buf[:packetEnd], &polyKey)

	copy(c.buf[packetEnd:], mac[:])

	if _, err := w.Write(c.buf); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
