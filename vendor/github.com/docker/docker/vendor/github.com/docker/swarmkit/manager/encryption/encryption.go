package encryption

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// This package defines the interfaces and encryption package

const humanReadablePrefix = "SWMKEY-1-"

// ErrCannotDecrypt is the type of error returned when some data cannot be decryptd as plaintext
type ErrCannotDecrypt struct ***REMOVED***
	msg string
***REMOVED***

func (e ErrCannotDecrypt) Error() string ***REMOVED***
	return e.msg
***REMOVED***

// A Decrypter can decrypt an encrypted record
type Decrypter interface ***REMOVED***
	Decrypt(api.MaybeEncryptedRecord) ([]byte, error)
***REMOVED***

// A Encrypter can encrypt some bytes into an encrypted record
type Encrypter interface ***REMOVED***
	Encrypt(data []byte) (*api.MaybeEncryptedRecord, error)
***REMOVED***

type noopCrypter struct***REMOVED******REMOVED***

func (n noopCrypter) Decrypt(e api.MaybeEncryptedRecord) ([]byte, error) ***REMOVED***
	if e.Algorithm != n.Algorithm() ***REMOVED***
		return nil, fmt.Errorf("record is encrypted")
	***REMOVED***
	return e.Data, nil
***REMOVED***

func (n noopCrypter) Encrypt(data []byte) (*api.MaybeEncryptedRecord, error) ***REMOVED***
	return &api.MaybeEncryptedRecord***REMOVED***
		Algorithm: n.Algorithm(),
		Data:      data,
	***REMOVED***, nil
***REMOVED***

func (n noopCrypter) Algorithm() api.MaybeEncryptedRecord_Algorithm ***REMOVED***
	return api.MaybeEncryptedRecord_NotEncrypted
***REMOVED***

// NoopCrypter is just a pass-through crypter - it does not actually encrypt or
// decrypt any data
var NoopCrypter = noopCrypter***REMOVED******REMOVED***

// Decrypt turns a slice of bytes serialized as an MaybeEncryptedRecord into a slice of plaintext bytes
func Decrypt(encryptd []byte, decrypter Decrypter) ([]byte, error) ***REMOVED***
	if decrypter == nil ***REMOVED***
		return nil, ErrCannotDecrypt***REMOVED***msg: "no decrypter specified"***REMOVED***
	***REMOVED***
	r := api.MaybeEncryptedRecord***REMOVED******REMOVED***
	if err := proto.Unmarshal(encryptd, &r); err != nil ***REMOVED***
		// nope, this wasn't marshalled as a MaybeEncryptedRecord
		return nil, ErrCannotDecrypt***REMOVED***msg: "unable to unmarshal as MaybeEncryptedRecord"***REMOVED***
	***REMOVED***
	plaintext, err := decrypter.Decrypt(r)
	if err != nil ***REMOVED***
		return nil, ErrCannotDecrypt***REMOVED***msg: err.Error()***REMOVED***
	***REMOVED***
	return plaintext, nil
***REMOVED***

// Encrypt turns a slice of bytes into a serialized MaybeEncryptedRecord slice of bytes
func Encrypt(plaintext []byte, encrypter Encrypter) ([]byte, error) ***REMOVED***
	if encrypter == nil ***REMOVED***
		return nil, fmt.Errorf("no encrypter specified")
	***REMOVED***

	encryptedRecord, err := encrypter.Encrypt(plaintext)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "unable to encrypt data")
	***REMOVED***

	data, err := proto.Marshal(encryptedRecord)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "unable to marshal as MaybeEncryptedRecord")
	***REMOVED***

	return data, nil
***REMOVED***

// Defaults returns a default encrypter and decrypter
func Defaults(key []byte) (Encrypter, Decrypter) ***REMOVED***
	n := NewNACLSecretbox(key)
	return n, n
***REMOVED***

// GenerateSecretKey generates a secret key that can be used for encrypting data
// using this package
func GenerateSecretKey() []byte ***REMOVED***
	secretData := make([]byte, naclSecretboxKeySize)
	if _, err := io.ReadFull(cryptorand.Reader, secretData); err != nil ***REMOVED***
		// panic if we can't read random data
		panic(errors.Wrap(err, "failed to read random bytes"))
	***REMOVED***
	return secretData
***REMOVED***

// HumanReadableKey displays a secret key in a human readable way
func HumanReadableKey(key []byte) string ***REMOVED***
	// base64-encode the key
	return humanReadablePrefix + base64.RawStdEncoding.EncodeToString(key)
***REMOVED***

// ParseHumanReadableKey returns a key as bytes from recognized serializations of
// said keys
func ParseHumanReadableKey(key string) ([]byte, error) ***REMOVED***
	if !strings.HasPrefix(key, humanReadablePrefix) ***REMOVED***
		return nil, fmt.Errorf("invalid key string")
	***REMOVED***
	keyBytes, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(key, humanReadablePrefix))
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("invalid key string")
	***REMOVED***
	return keyBytes, nil
***REMOVED***
