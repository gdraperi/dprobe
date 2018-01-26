package ca

import (
	cryptorand "crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"crypto/tls"

	"github.com/docker/swarmkit/ioutils"
	"github.com/pkg/errors"
)

const (
	// keyPerms are the permissions used to write the TLS keys
	keyPerms = 0600
	// certPerms are the permissions used to write TLS certificates
	certPerms = 0644
	// versionHeader is the TLS PEM key header that contains the KEK version
	versionHeader = "kek-version"
)

// PEMKeyHeaders is something that needs to know about PEM headers when reading
// or writing TLS keys.
type PEMKeyHeaders interface ***REMOVED***
	// UnmarshalHeaders loads the headers map given the current KEK
	UnmarshalHeaders(map[string]string, KEKData) (PEMKeyHeaders, error)
	// MarshalHeaders returns a header map given the current KEK
	MarshalHeaders(KEKData) (map[string]string, error)
	// UpdateKEK may get a new PEMKeyHeaders if the KEK changes
	UpdateKEK(KEKData, KEKData) PEMKeyHeaders
***REMOVED***

// KeyReader reads a TLS cert and key from disk
type KeyReader interface ***REMOVED***
	Read() ([]byte, []byte, error)
	Target() string
***REMOVED***

// KeyWriter writes a TLS key and cert to disk
type KeyWriter interface ***REMOVED***
	Write([]byte, []byte, *KEKData) error
	ViewAndUpdateHeaders(func(PEMKeyHeaders) (PEMKeyHeaders, error)) error
	ViewAndRotateKEK(func(KEKData, PEMKeyHeaders) (KEKData, PEMKeyHeaders, error)) error
	GetCurrentState() (PEMKeyHeaders, KEKData)
	Target() string
***REMOVED***

// KEKData provides an optional update to the kek when writing.  The structure
// is needed so that we can tell the difference between "do not encrypt anymore"
// and there is "no update".
type KEKData struct ***REMOVED***
	KEK     []byte
	Version uint64
***REMOVED***

// ErrInvalidKEK means that we cannot decrypt the TLS key for some reason
type ErrInvalidKEK struct ***REMOVED***
	Wrapped error
***REMOVED***

func (e ErrInvalidKEK) Error() string ***REMOVED***
	return e.Wrapped.Error()
***REMOVED***

// KeyReadWriter is an object that knows how to read and write TLS keys and certs to disk,
// optionally encrypted and optionally updating PEM headers.
type KeyReadWriter struct ***REMOVED***
	mu         sync.Mutex
	kekData    KEKData
	paths      CertPaths
	headersObj PEMKeyHeaders
***REMOVED***

// NewKeyReadWriter creates a new KeyReadWriter
func NewKeyReadWriter(paths CertPaths, kek []byte, headersObj PEMKeyHeaders) *KeyReadWriter ***REMOVED***
	return &KeyReadWriter***REMOVED***
		kekData:    KEKData***REMOVED***KEK: kek***REMOVED***,
		paths:      paths,
		headersObj: headersObj,
	***REMOVED***
***REMOVED***

// Migrate checks to see if a temporary key file exists.  Older versions of
// swarmkit wrote temporary keys instead of temporary certificates, so
// migrate that temporary key if it exists.  We want to write temporary certificates,
// instead of temporary keys, because we may need to periodically re-encrypt the
// keys and modify the headers, and it's easier to have a single canonical key
// location than two possible key locations.
func (k *KeyReadWriter) Migrate() error ***REMOVED***
	tmpPaths := k.genTempPaths()
	keyBytes, err := ioutil.ReadFile(tmpPaths.Key)
	if err != nil ***REMOVED***
		return nil // no key?  no migration
	***REMOVED***

	// it does exist - no need to decrypt, because previous versions of swarmkit
	// which supported this temporary key did not support encrypting TLS keys
	cert, err := ioutil.ReadFile(k.paths.Cert)
	if err != nil ***REMOVED***
		return os.RemoveAll(tmpPaths.Key) // no cert?  no migration
	***REMOVED***

	// nope, this does not match the cert
	if _, err = tls.X509KeyPair(cert, keyBytes); err != nil ***REMOVED***
		return os.RemoveAll(tmpPaths.Key)
	***REMOVED***

	return os.Rename(tmpPaths.Key, k.paths.Key)
***REMOVED***

// Read will read a TLS cert and key from the given paths
func (k *KeyReadWriter) Read() ([]byte, []byte, error) ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()
	keyBlock, err := k.readKey()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	if version, ok := keyBlock.Headers[versionHeader]; ok ***REMOVED***
		if versionInt, err := strconv.ParseUint(version, 10, 64); err == nil ***REMOVED***
			k.kekData.Version = versionInt
		***REMOVED***
	***REMOVED***
	delete(keyBlock.Headers, versionHeader)

	if k.headersObj != nil ***REMOVED***
		newHeaders, err := k.headersObj.UnmarshalHeaders(keyBlock.Headers, k.kekData)
		if err != nil ***REMOVED***
			return nil, nil, errors.Wrap(err, "unable to read TLS key headers")
		***REMOVED***
		k.headersObj = newHeaders
	***REMOVED***

	keyBytes := pem.EncodeToMemory(keyBlock)
	cert, err := ioutil.ReadFile(k.paths.Cert)
	// The cert is written to a temporary file first, then the key, and then
	// the cert gets renamed - so, if interrupted, it's possible to end up with
	// a cert that only exists in the temporary location.
	switch ***REMOVED***
	case err == nil:
		_, err = tls.X509KeyPair(cert, keyBytes)
	case os.IsNotExist(err): //continue to try temp location
		break
	default:
		return nil, nil, err
	***REMOVED***

	// either the cert doesn't exist, or it doesn't match the key - try the temp file, if it exists
	if err != nil ***REMOVED***
		var tempErr error
		tmpPaths := k.genTempPaths()
		cert, tempErr = ioutil.ReadFile(tmpPaths.Cert)
		if tempErr != nil ***REMOVED***
			return nil, nil, err // return the original error
		***REMOVED***
		if _, tempErr := tls.X509KeyPair(cert, keyBytes); tempErr != nil ***REMOVED***
			os.RemoveAll(tmpPaths.Cert) // nope, it doesn't match either - remove and return the original error
			return nil, nil, err
		***REMOVED***
		os.Rename(tmpPaths.Cert, k.paths.Cert) // try to move the temp cert back to the regular location

	***REMOVED***

	return cert, keyBytes, nil
***REMOVED***

// ViewAndRotateKEK re-encrypts the key with a new KEK
func (k *KeyReadWriter) ViewAndRotateKEK(cb func(KEKData, PEMKeyHeaders) (KEKData, PEMKeyHeaders, error)) error ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()

	updatedKEK, updatedHeaderObj, err := cb(k.kekData, k.headersObj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	keyBlock, err := k.readKey()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return k.writeKey(keyBlock, updatedKEK, updatedHeaderObj)
***REMOVED***

// ViewAndUpdateHeaders updates the header manager, and updates any headers on the existing key
func (k *KeyReadWriter) ViewAndUpdateHeaders(cb func(PEMKeyHeaders) (PEMKeyHeaders, error)) error ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()

	pkh, err := cb(k.headersObj)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	keyBlock, err := k.readKeyblock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	headers := make(map[string]string)
	if pkh != nil ***REMOVED***
		var err error
		headers, err = pkh.MarshalHeaders(k.kekData)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// we WANT any original encryption headers
	for key, value := range keyBlock.Headers ***REMOVED***
		normalizedKey := strings.TrimSpace(strings.ToLower(key))
		if normalizedKey == "proc-type" || normalizedKey == "dek-info" ***REMOVED***
			headers[key] = value
		***REMOVED***
	***REMOVED***
	headers[versionHeader] = strconv.FormatUint(k.kekData.Version, 10)
	keyBlock.Headers = headers

	if err = ioutils.AtomicWriteFile(k.paths.Key, pem.EncodeToMemory(keyBlock), keyPerms); err != nil ***REMOVED***
		return err
	***REMOVED***
	k.headersObj = pkh
	return nil
***REMOVED***

// GetCurrentState returns the current KEK data, including version
func (k *KeyReadWriter) GetCurrentState() (PEMKeyHeaders, KEKData) ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.headersObj, k.kekData
***REMOVED***

// Write attempts write a cert and key to text.  This can also optionally update
// the KEK while writing, if an updated KEK is provided.  If the pointer to the
// update KEK is nil, then we don't update. If the updated KEK itself is nil,
// then we update the KEK to be nil (data should be unencrypted).
func (k *KeyReadWriter) Write(certBytes, plaintextKeyBytes []byte, kekData *KEKData) error ***REMOVED***
	k.mu.Lock()
	defer k.mu.Unlock()

	// current assumption is that the cert and key will be in the same directory
	if err := os.MkdirAll(filepath.Dir(k.paths.Key), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Ensure that we will have a keypair on disk at all times by writing the cert to a
	// temp path first.  This is because we want to have only a single copy of the key
	// for rotation and header modification.
	tmpPaths := k.genTempPaths()
	if err := ioutils.AtomicWriteFile(tmpPaths.Cert, certBytes, certPerms); err != nil ***REMOVED***
		return err
	***REMOVED***

	keyBlock, _ := pem.Decode(plaintextKeyBytes)
	if keyBlock == nil ***REMOVED***
		return errors.New("invalid PEM-encoded private key")
	***REMOVED***

	if kekData == nil ***REMOVED***
		kekData = &k.kekData
	***REMOVED***
	pkh := k.headersObj
	if k.headersObj != nil ***REMOVED***
		pkh = k.headersObj.UpdateKEK(k.kekData, *kekData)
	***REMOVED***

	if err := k.writeKey(keyBlock, *kekData, pkh); err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.Rename(tmpPaths.Cert, k.paths.Cert)
***REMOVED***

func (k *KeyReadWriter) genTempPaths() CertPaths ***REMOVED***
	return CertPaths***REMOVED***
		Key:  filepath.Join(filepath.Dir(k.paths.Key), "."+filepath.Base(k.paths.Key)),
		Cert: filepath.Join(filepath.Dir(k.paths.Cert), "."+filepath.Base(k.paths.Cert)),
	***REMOVED***
***REMOVED***

// Target returns a string representation of this KeyReadWriter, namely where
// it is writing to
func (k *KeyReadWriter) Target() string ***REMOVED***
	return k.paths.Cert
***REMOVED***

func (k *KeyReadWriter) readKeyblock() (*pem.Block, error) ***REMOVED***
	key, err := ioutil.ReadFile(k.paths.Key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Decode the PEM private key
	keyBlock, _ := pem.Decode(key)
	if keyBlock == nil ***REMOVED***
		return nil, errors.New("invalid PEM-encoded private key")
	***REMOVED***

	return keyBlock, nil
***REMOVED***

// readKey returns the decrypted key pem bytes, and enforces the KEK if applicable
// (writes it back with the correct encryption if it is not correctly encrypted)
func (k *KeyReadWriter) readKey() (*pem.Block, error) ***REMOVED***
	keyBlock, err := k.readKeyblock()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !x509.IsEncryptedPEMBlock(keyBlock) ***REMOVED***
		return keyBlock, nil
	***REMOVED***

	// If it's encrypted, we can't read without a passphrase (we're assuming
	// empty passphrases are invalid)
	if k.kekData.KEK == nil ***REMOVED***
		return nil, ErrInvalidKEK***REMOVED***Wrapped: x509.IncorrectPasswordError***REMOVED***
	***REMOVED***

	derBytes, err := x509.DecryptPEMBlock(keyBlock, k.kekData.KEK)
	if err != nil ***REMOVED***
		return nil, ErrInvalidKEK***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***
	// remove encryption PEM headers
	headers := make(map[string]string)
	mergePEMHeaders(headers, keyBlock.Headers)

	return &pem.Block***REMOVED***
		Type:    keyBlock.Type, // the key type doesn't change
		Bytes:   derBytes,
		Headers: headers,
	***REMOVED***, nil
***REMOVED***

// writeKey takes an unencrypted keyblock and, if the kek is not nil, encrypts it before
// writing it to disk.  If the kek is nil, writes it to disk unencrypted.
func (k *KeyReadWriter) writeKey(keyBlock *pem.Block, kekData KEKData, pkh PEMKeyHeaders) error ***REMOVED***
	if kekData.KEK != nil ***REMOVED***
		encryptedPEMBlock, err := x509.EncryptPEMBlock(cryptorand.Reader,
			keyBlock.Type,
			keyBlock.Bytes,
			kekData.KEK,
			x509.PEMCipherAES256)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if encryptedPEMBlock.Headers == nil ***REMOVED***
			return errors.New("unable to encrypt key - invalid PEM file produced")
		***REMOVED***
		keyBlock = encryptedPEMBlock
	***REMOVED***

	if pkh != nil ***REMOVED***
		headers, err := pkh.MarshalHeaders(kekData)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		mergePEMHeaders(keyBlock.Headers, headers)
	***REMOVED***
	keyBlock.Headers[versionHeader] = strconv.FormatUint(kekData.Version, 10)

	if err := ioutils.AtomicWriteFile(k.paths.Key, pem.EncodeToMemory(keyBlock), keyPerms); err != nil ***REMOVED***
		return err
	***REMOVED***
	k.kekData = kekData
	k.headersObj = pkh
	return nil
***REMOVED***

// merges one set of PEM headers onto another, excepting for key encryption value
// "proc-type" and "dek-info"
func mergePEMHeaders(original, newSet map[string]string) ***REMOVED***
	for key, value := range newSet ***REMOVED***
		normalizedKey := strings.TrimSpace(strings.ToLower(key))
		if normalizedKey != "proc-type" && normalizedKey != "dek-info" ***REMOVED***
			original[key] = value
		***REMOVED***
	***REMOVED***
***REMOVED***
