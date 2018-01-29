// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto"
	"crypto/rand"
	"io"
	"time"
)

// Config collects a number of parameters along with sensible defaults.
// A nil *Config is valid and results in all default values.
type Config struct ***REMOVED***
	// Rand provides the source of entropy.
	// If nil, the crypto/rand Reader is used.
	Rand io.Reader
	// DefaultHash is the default hash function to be used.
	// If zero, SHA-256 is used.
	DefaultHash crypto.Hash
	// DefaultCipher is the cipher to be used.
	// If zero, AES-128 is used.
	DefaultCipher CipherFunction
	// Time returns the current time as the number of seconds since the
	// epoch. If Time is nil, time.Now is used.
	Time func() time.Time
	// DefaultCompressionAlgo is the compression algorithm to be
	// applied to the plaintext before encryption. If zero, no
	// compression is done.
	DefaultCompressionAlgo CompressionAlgo
	// CompressionConfig configures the compression settings.
	CompressionConfig *CompressionConfig
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
	// RSABits is the number of bits in new RSA keys made with NewEntity.
	// If zero, then 2048 bit keys are created.
	RSABits int
***REMOVED***

func (c *Config) Random() io.Reader ***REMOVED***
	if c == nil || c.Rand == nil ***REMOVED***
		return rand.Reader
	***REMOVED***
	return c.Rand
***REMOVED***

func (c *Config) Hash() crypto.Hash ***REMOVED***
	if c == nil || uint(c.DefaultHash) == 0 ***REMOVED***
		return crypto.SHA256
	***REMOVED***
	return c.DefaultHash
***REMOVED***

func (c *Config) Cipher() CipherFunction ***REMOVED***
	if c == nil || uint8(c.DefaultCipher) == 0 ***REMOVED***
		return CipherAES128
	***REMOVED***
	return c.DefaultCipher
***REMOVED***

func (c *Config) Now() time.Time ***REMOVED***
	if c == nil || c.Time == nil ***REMOVED***
		return time.Now()
	***REMOVED***
	return c.Time()
***REMOVED***

func (c *Config) Compression() CompressionAlgo ***REMOVED***
	if c == nil ***REMOVED***
		return CompressionNone
	***REMOVED***
	return c.DefaultCompressionAlgo
***REMOVED***

func (c *Config) PasswordHashIterations() int ***REMOVED***
	if c == nil || c.S2KCount == 0 ***REMOVED***
		return 0
	***REMOVED***
	return c.S2KCount
***REMOVED***
