package daemon

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/libtrust"
)

// LoadOrCreateTrustKey attempts to load the libtrust key at the given path,
// otherwise generates a new one
// TODO: this should use more of libtrust.LoadOrCreateTrustKey which may need
// a refactor or this function to be moved into libtrust
func loadOrCreateTrustKey(trustKeyPath string) (libtrust.PrivateKey, error) ***REMOVED***
	err := system.MkdirAll(filepath.Dir(trustKeyPath), 0700, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	trustKey, err := libtrust.LoadKeyFile(trustKeyPath)
	if err == libtrust.ErrKeyFileDoesNotExist ***REMOVED***
		trustKey, err = libtrust.GenerateECP256PrivateKey()
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Error generating key: %s", err)
		***REMOVED***
		encodedKey, err := serializePrivateKey(trustKey, filepath.Ext(trustKeyPath))
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Error serializing key: %s", err)
		***REMOVED***
		if err := ioutils.AtomicWriteFile(trustKeyPath, encodedKey, os.FileMode(0600)); err != nil ***REMOVED***
			return nil, fmt.Errorf("Error saving key file: %s", err)
		***REMOVED***
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, fmt.Errorf("Error loading key file %s: %s", trustKeyPath, err)
	***REMOVED***
	return trustKey, nil
***REMOVED***

func serializePrivateKey(key libtrust.PrivateKey, ext string) (encoded []byte, err error) ***REMOVED***
	if ext == ".json" || ext == ".jwk" ***REMOVED***
		encoded, err = json.Marshal(key)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to encode private key JWK: %s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		pemBlock, err := key.PEMBlock()
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("unable to encode private key PEM: %s", err)
		***REMOVED***
		encoded = pem.EncodeToMemory(pemBlock)
	***REMOVED***
	return
***REMOVED***
