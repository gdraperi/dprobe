// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing
// algorithm. See http://www.usenix.org/event/usenix99/provos/provos.pdf
package bcrypt // import "golang.org/x/crypto/bcrypt"

// The code is a port of Provos and Mazières's C implementation.
import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"strconv"

	"golang.org/x/crypto/blowfish"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

// The error returned from CompareHashAndPassword when a password and hash do
// not match.
var ErrMismatchedHashAndPassword = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")

// The error returned from CompareHashAndPassword when a hash is too short to
// be a bcrypt hash.
var ErrHashTooShort = errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password")

// The error returned from CompareHashAndPassword when a hash was created with
// a bcrypt algorithm newer than this implementation.
type HashVersionTooNewError byte

func (hv HashVersionTooNewError) Error() string ***REMOVED***
	return fmt.Sprintf("crypto/bcrypt: bcrypt algorithm version '%c' requested is newer than current version '%c'", byte(hv), majorVersion)
***REMOVED***

// The error returned from CompareHashAndPassword when a hash starts with something other than '$'
type InvalidHashPrefixError byte

func (ih InvalidHashPrefixError) Error() string ***REMOVED***
	return fmt.Sprintf("crypto/bcrypt: bcrypt hashes must start with '$', but hashedSecret started with '%c'", byte(ih))
***REMOVED***

type InvalidCostError int

func (ic InvalidCostError) Error() string ***REMOVED***
	return fmt.Sprintf("crypto/bcrypt: cost %d is outside allowed range (%d,%d)", int(ic), int(MinCost), int(MaxCost))
***REMOVED***

const (
	majorVersion       = '2'
	minorVersion       = 'a'
	maxSaltSize        = 16
	maxCryptedHashSize = 23
	encodedSaltSize    = 22
	encodedHashSize    = 31
	minHashSize        = 59
)

// magicCipherData is an IV for the 64 Blowfish encryption calls in
// bcrypt(). It's the string "OrpheanBeholderScryDoubt" in big-endian bytes.
var magicCipherData = []byte***REMOVED***
	0x4f, 0x72, 0x70, 0x68,
	0x65, 0x61, 0x6e, 0x42,
	0x65, 0x68, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x53,
	0x63, 0x72, 0x79, 0x44,
	0x6f, 0x75, 0x62, 0x74,
***REMOVED***

type hashed struct ***REMOVED***
	hash  []byte
	salt  []byte
	cost  int // allowed range is MinCost to MaxCost
	major byte
	minor byte
***REMOVED***

// GenerateFromPassword returns the bcrypt hash of the password at the given
// cost. If the cost given is less than MinCost, the cost will be set to
// DefaultCost, instead. Use CompareHashAndPassword, as defined in this package,
// to compare the returned hashed password with its cleartext version.
func GenerateFromPassword(password []byte, cost int) ([]byte, error) ***REMOVED***
	p, err := newFromPassword(password, cost)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p.Hash(), nil
***REMOVED***

// CompareHashAndPassword compares a bcrypt hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password []byte) error ***REMOVED***
	p, err := newFromHash(hashedPassword)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	otherHash, err := bcrypt(password, p.cost, p.salt)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	otherP := &hashed***REMOVED***otherHash, p.salt, p.cost, p.major, p.minor***REMOVED***
	if subtle.ConstantTimeCompare(p.Hash(), otherP.Hash()) == 1 ***REMOVED***
		return nil
	***REMOVED***

	return ErrMismatchedHashAndPassword
***REMOVED***

// Cost returns the hashing cost used to create the given hashed
// password. When, in the future, the hashing cost of a password system needs
// to be increased in order to adjust for greater computational power, this
// function allows one to establish which passwords need to be updated.
func Cost(hashedPassword []byte) (int, error) ***REMOVED***
	p, err := newFromHash(hashedPassword)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return p.cost, nil
***REMOVED***

func newFromPassword(password []byte, cost int) (*hashed, error) ***REMOVED***
	if cost < MinCost ***REMOVED***
		cost = DefaultCost
	***REMOVED***
	p := new(hashed)
	p.major = majorVersion
	p.minor = minorVersion

	err := checkCost(cost)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p.cost = cost

	unencodedSalt := make([]byte, maxSaltSize)
	_, err = io.ReadFull(rand.Reader, unencodedSalt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p.salt = base64Encode(unencodedSalt)
	hash, err := bcrypt(password, p.cost, p.salt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p.hash = hash
	return p, err
***REMOVED***

func newFromHash(hashedSecret []byte) (*hashed, error) ***REMOVED***
	if len(hashedSecret) < minHashSize ***REMOVED***
		return nil, ErrHashTooShort
	***REMOVED***
	p := new(hashed)
	n, err := p.decodeVersion(hashedSecret)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hashedSecret = hashedSecret[n:]
	n, err = p.decodeCost(hashedSecret)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	hashedSecret = hashedSecret[n:]

	// The "+2" is here because we'll have to append at most 2 '=' to the salt
	// when base64 decoding it in expensiveBlowfishSetup().
	p.salt = make([]byte, encodedSaltSize, encodedSaltSize+2)
	copy(p.salt, hashedSecret[:encodedSaltSize])

	hashedSecret = hashedSecret[encodedSaltSize:]
	p.hash = make([]byte, len(hashedSecret))
	copy(p.hash, hashedSecret)

	return p, nil
***REMOVED***

func bcrypt(password []byte, cost int, salt []byte) ([]byte, error) ***REMOVED***
	cipherData := make([]byte, len(magicCipherData))
	copy(cipherData, magicCipherData)

	c, err := expensiveBlowfishSetup(password, uint32(cost), salt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for i := 0; i < 24; i += 8 ***REMOVED***
		for j := 0; j < 64; j++ ***REMOVED***
			c.Encrypt(cipherData[i:i+8], cipherData[i:i+8])
		***REMOVED***
	***REMOVED***

	// Bug compatibility with C bcrypt implementations. We only encode 23 of
	// the 24 bytes encrypted.
	hsh := base64Encode(cipherData[:maxCryptedHashSize])
	return hsh, nil
***REMOVED***

func expensiveBlowfishSetup(key []byte, cost uint32, salt []byte) (*blowfish.Cipher, error) ***REMOVED***
	csalt, err := base64Decode(salt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Bug compatibility with C bcrypt implementations. They use the trailing
	// NULL in the key string during expansion.
	// We copy the key to prevent changing the underlying array.
	ckey := append(key[:len(key):len(key)], 0)

	c, err := blowfish.NewSaltedCipher(ckey, csalt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var i, rounds uint64
	rounds = 1 << cost
	for i = 0; i < rounds; i++ ***REMOVED***
		blowfish.ExpandKey(ckey, c)
		blowfish.ExpandKey(csalt, c)
	***REMOVED***

	return c, nil
***REMOVED***

func (p *hashed) Hash() []byte ***REMOVED***
	arr := make([]byte, 60)
	arr[0] = '$'
	arr[1] = p.major
	n := 2
	if p.minor != 0 ***REMOVED***
		arr[2] = p.minor
		n = 3
	***REMOVED***
	arr[n] = '$'
	n++
	copy(arr[n:], []byte(fmt.Sprintf("%02d", p.cost)))
	n += 2
	arr[n] = '$'
	n++
	copy(arr[n:], p.salt)
	n += encodedSaltSize
	copy(arr[n:], p.hash)
	n += encodedHashSize
	return arr[:n]
***REMOVED***

func (p *hashed) decodeVersion(sbytes []byte) (int, error) ***REMOVED***
	if sbytes[0] != '$' ***REMOVED***
		return -1, InvalidHashPrefixError(sbytes[0])
	***REMOVED***
	if sbytes[1] > majorVersion ***REMOVED***
		return -1, HashVersionTooNewError(sbytes[1])
	***REMOVED***
	p.major = sbytes[1]
	n := 3
	if sbytes[2] != '$' ***REMOVED***
		p.minor = sbytes[2]
		n++
	***REMOVED***
	return n, nil
***REMOVED***

// sbytes should begin where decodeVersion left off.
func (p *hashed) decodeCost(sbytes []byte) (int, error) ***REMOVED***
	cost, err := strconv.Atoi(string(sbytes[0:2]))
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	err = checkCost(cost)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	p.cost = cost
	return 3, nil
***REMOVED***

func (p *hashed) String() string ***REMOVED***
	return fmt.Sprintf("&***REMOVED***hash: %#v, salt: %#v, cost: %d, major: %c, minor: %c***REMOVED***", string(p.hash), p.salt, p.cost, p.major, p.minor)
***REMOVED***

func checkCost(cost int) error ***REMOVED***
	if cost < MinCost || cost > MaxCost ***REMOVED***
		return InvalidCostError(cost)
	***REMOVED***
	return nil
***REMOVED***
