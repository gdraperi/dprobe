package libtrust

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	// ErrKeyFileDoesNotExist indicates that the private key file does not exist.
	ErrKeyFileDoesNotExist = errors.New("key file does not exist")
)

func readKeyFileBytes(filename string) ([]byte, error) ***REMOVED***
	data, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			err = ErrKeyFileDoesNotExist
		***REMOVED*** else ***REMOVED***
			err = fmt.Errorf("unable to read key file %s: %s", filename, err)
		***REMOVED***

		return nil, err
	***REMOVED***

	return data, nil
***REMOVED***

/*
	Loading and Saving of Public and Private Keys in either PEM or JWK format.
*/

// LoadKeyFile opens the given filename and attempts to read a Private Key
// encoded in either PEM or JWK format (if .json or .jwk file extension).
func LoadKeyFile(filename string) (PrivateKey, error) ***REMOVED***
	contents, err := readKeyFileBytes(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var key PrivateKey

	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		key, err = UnmarshalPrivateKeyJWK(contents)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to decode private key JWK: %s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		key, err = UnmarshalPrivateKeyPEM(contents)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to decode private key PEM: %s", err)
		***REMOVED***
	***REMOVED***

	return key, nil
***REMOVED***

// LoadPublicKeyFile opens the given filename and attempts to read a Public Key
// encoded in either PEM or JWK format (if .json or .jwk file extension).
func LoadPublicKeyFile(filename string) (PublicKey, error) ***REMOVED***
	contents, err := readKeyFileBytes(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var key PublicKey

	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		key, err = UnmarshalPublicKeyJWK(contents)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to decode public key JWK: %s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		key, err = UnmarshalPublicKeyPEM(contents)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to decode public key PEM: %s", err)
		***REMOVED***
	***REMOVED***

	return key, nil
***REMOVED***

// SaveKey saves the given key to a file using the provided filename.
// This process will overwrite any existing file at the provided location.
func SaveKey(filename string, key PrivateKey) error ***REMOVED***
	var encodedKey []byte
	var err error

	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		// Encode in JSON Web Key format.
		encodedKey, err = json.MarshalIndent(key, "", "    ")
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to encode private key JWK: %s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Encode in PEM format.
		pemBlock, err := key.PEMBlock()
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to encode private key PEM: %s", err)
		***REMOVED***
		encodedKey = pem.EncodeToMemory(pemBlock)
	***REMOVED***

	err = ioutil.WriteFile(filename, encodedKey, os.FileMode(0600))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to write private key file %s: %s", filename, err)
	***REMOVED***

	return nil
***REMOVED***

// SavePublicKey saves the given public key to the file.
func SavePublicKey(filename string, key PublicKey) error ***REMOVED***
	var encodedKey []byte
	var err error

	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		// Encode in JSON Web Key format.
		encodedKey, err = json.MarshalIndent(key, "", "    ")
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to encode public key JWK: %s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Encode in PEM format.
		pemBlock, err := key.PEMBlock()
		if err != nil ***REMOVED***
			return fmt.Errorf("unable to encode public key PEM: %s", err)
		***REMOVED***
		encodedKey = pem.EncodeToMemory(pemBlock)
	***REMOVED***

	err = ioutil.WriteFile(filename, encodedKey, os.FileMode(0644))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to write public key file %s: %s", filename, err)
	***REMOVED***

	return nil
***REMOVED***

// Public Key Set files

type jwkSet struct ***REMOVED***
	Keys []json.RawMessage `json:"keys"`
***REMOVED***

// LoadKeySetFile loads a key set
func LoadKeySetFile(filename string) ([]PublicKey, error) ***REMOVED***
	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		return loadJSONKeySetFile(filename)
	***REMOVED***

	// Must be a PEM format file
	return loadPEMKeySetFile(filename)
***REMOVED***

func loadJSONKeySetRaw(data []byte) ([]json.RawMessage, error) ***REMOVED***
	if len(data) == 0 ***REMOVED***
		// This is okay, just return an empty slice.
		return []json.RawMessage***REMOVED******REMOVED***, nil
	***REMOVED***

	keySet := jwkSet***REMOVED******REMOVED***

	err := json.Unmarshal(data, &keySet)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unable to decode JSON Web Key Set: %s", err)
	***REMOVED***

	return keySet.Keys, nil
***REMOVED***

func loadJSONKeySetFile(filename string) ([]PublicKey, error) ***REMOVED***
	contents, err := readKeyFileBytes(filename)
	if err != nil && err != ErrKeyFileDoesNotExist ***REMOVED***
		return nil, err
	***REMOVED***

	return UnmarshalPublicKeyJWKSet(contents)
***REMOVED***

func loadPEMKeySetFile(filename string) ([]PublicKey, error) ***REMOVED***
	data, err := readKeyFileBytes(filename)
	if err != nil && err != ErrKeyFileDoesNotExist ***REMOVED***
		return nil, err
	***REMOVED***

	return UnmarshalPublicKeyPEMBundle(data)
***REMOVED***

// AddKeySetFile adds a key to a key set
func AddKeySetFile(filename string, key PublicKey) error ***REMOVED***
	if strings.HasSuffix(filename, ".json") || strings.HasSuffix(filename, ".jwk") ***REMOVED***
		return addKeySetJSONFile(filename, key)
	***REMOVED***

	// Must be a PEM format file
	return addKeySetPEMFile(filename, key)
***REMOVED***

func addKeySetJSONFile(filename string, key PublicKey) error ***REMOVED***
	encodedKey, err := json.Marshal(key)
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to encode trusted client key: %s", err)
	***REMOVED***

	contents, err := readKeyFileBytes(filename)
	if err != nil && err != ErrKeyFileDoesNotExist ***REMOVED***
		return err
	***REMOVED***

	rawEntries, err := loadJSONKeySetRaw(contents)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	rawEntries = append(rawEntries, json.RawMessage(encodedKey))
	entriesWrapper := jwkSet***REMOVED***Keys: rawEntries***REMOVED***

	encodedEntries, err := json.MarshalIndent(entriesWrapper, "", "    ")
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to encode trusted client keys: %s", err)
	***REMOVED***

	err = ioutil.WriteFile(filename, encodedEntries, os.FileMode(0644))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to write trusted client keys file %s: %s", filename, err)
	***REMOVED***

	return nil
***REMOVED***

func addKeySetPEMFile(filename string, key PublicKey) error ***REMOVED***
	// Encode to PEM, open file for appending, write PEM.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0644))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to open trusted client keys file %s: %s", filename, err)
	***REMOVED***
	defer file.Close()

	pemBlock, err := key.PEMBlock()
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to encoded trusted key: %s", err)
	***REMOVED***

	_, err = file.Write(pem.EncodeToMemory(pemBlock))
	if err != nil ***REMOVED***
		return fmt.Errorf("unable to write trusted keys file: %s", err)
	***REMOVED***

	return nil
***REMOVED***
