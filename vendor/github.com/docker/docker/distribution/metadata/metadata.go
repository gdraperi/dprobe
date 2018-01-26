package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
)

// Store implements a K/V store for mapping distribution-related IDs
// to on-disk layer IDs and image IDs. The namespace identifies the type of
// mapping (i.e. "v1ids" or "artifacts"). MetadataStore is goroutine-safe.
type Store interface ***REMOVED***
	// Get retrieves data by namespace and key.
	Get(namespace string, key string) ([]byte, error)
	// Set writes data indexed by namespace and key.
	Set(namespace, key string, value []byte) error
	// Delete removes data indexed by namespace and key.
	Delete(namespace, key string) error
***REMOVED***

// FSMetadataStore uses the filesystem to associate metadata with layer and
// image IDs.
type FSMetadataStore struct ***REMOVED***
	sync.RWMutex
	basePath string
***REMOVED***

// NewFSMetadataStore creates a new filesystem-based metadata store.
func NewFSMetadataStore(basePath string) (*FSMetadataStore, error) ***REMOVED***
	if err := os.MkdirAll(basePath, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &FSMetadataStore***REMOVED***
		basePath: basePath,
	***REMOVED***, nil
***REMOVED***

func (store *FSMetadataStore) path(namespace, key string) string ***REMOVED***
	return filepath.Join(store.basePath, namespace, key)
***REMOVED***

// Get retrieves data by namespace and key. The data is read from a file named
// after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Get(namespace string, key string) ([]byte, error) ***REMOVED***
	store.RLock()
	defer store.RUnlock()

	return ioutil.ReadFile(store.path(namespace, key))
***REMOVED***

// Set writes data indexed by namespace and key. The data is written to a file
// named after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Set(namespace, key string, value []byte) error ***REMOVED***
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutils.AtomicWriteFile(path, value, 0644)
***REMOVED***

// Delete removes data indexed by namespace and key. The data file named after
// the key, stored in the namespace's directory is deleted.
func (store *FSMetadataStore) Delete(namespace, key string) error ***REMOVED***
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	return os.Remove(path)
***REMOVED***
