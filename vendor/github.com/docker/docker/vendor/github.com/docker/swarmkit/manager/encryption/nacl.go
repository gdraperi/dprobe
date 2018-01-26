package encryption

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"

	"github.com/docker/swarmkit/api"

	"golang.org/x/crypto/nacl/secretbox"
)

const naclSecretboxKeySize = 32
const naclSecretboxNonceSize = 24

// This provides the default implementation of an encrypter and decrypter, as well
// as the default KDF function.

// NACLSecretbox is an implementation of an encrypter/decrypter.  Encrypting
// generates random Nonces.
type NACLSecretbox struct ***REMOVED***
	key [naclSecretboxKeySize]byte
***REMOVED***

// NewNACLSecretbox returns a new NACL secretbox encrypter/decrypter with the given key
func NewNACLSecretbox(key []byte) NACLSecretbox ***REMOVED***
	secretbox := NACLSecretbox***REMOVED******REMOVED***
	copy(secretbox.key[:], key)
	return secretbox
***REMOVED***

// Algorithm returns the type of algorithm this is (NACL Secretbox using XSalsa20 and Poly1305)
func (n NACLSecretbox) Algorithm() api.MaybeEncryptedRecord_Algorithm ***REMOVED***
	return api.MaybeEncryptedRecord_NACLSecretboxSalsa20Poly1305
***REMOVED***

// Encrypt encrypts some bytes and returns an encrypted record
func (n NACLSecretbox) Encrypt(data []byte) (*api.MaybeEncryptedRecord, error) ***REMOVED***
	var nonce [24]byte
	if _, err := io.ReadFull(cryptorand.Reader, nonce[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Seal's first argument is an "out", the data that the new encrypted message should be
	// appended to.  Since we don't want to append anything, we pass nil.
	encrypted := secretbox.Seal(nil, data, &nonce, &n.key)
	return &api.MaybeEncryptedRecord***REMOVED***
		Algorithm: n.Algorithm(),
		Data:      encrypted,
		Nonce:     nonce[:],
	***REMOVED***, nil
***REMOVED***

// Decrypt decrypts a MaybeEncryptedRecord and returns some bytes
func (n NACLSecretbox) Decrypt(record api.MaybeEncryptedRecord) ([]byte, error) ***REMOVED***
	if record.Algorithm != n.Algorithm() ***REMOVED***
		return nil, fmt.Errorf("not a NACL secretbox record")
	***REMOVED***
	if len(record.Nonce) != naclSecretboxNonceSize ***REMOVED***
		return nil, fmt.Errorf("invalid nonce size for NACL secretbox: require 24, got %d", len(record.Nonce))
	***REMOVED***

	var decryptNonce [naclSecretboxNonceSize]byte
	copy(decryptNonce[:], record.Nonce[:naclSecretboxNonceSize])

	// Open's first argument is an "out", the data that the decrypted message should be
	// appended to.  Since we don't want to append anything, we pass nil.
	decrypted, ok := secretbox.Open(nil, record.Data, &decryptNonce, &n.key)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("decryption error using NACL secretbox")
	***REMOVED***
	return decrypted, nil
***REMOVED***
