package libkv

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/libkv/store"
)

// Initialize creates a new Store object, initializing the client
type Initialize func(addrs []string, options *store.Config) (store.Store, error)

var (
	// Backend initializers
	initializers = make(map[store.Backend]Initialize)

	supportedBackend = func() string ***REMOVED***
		keys := make([]string, 0, len(initializers))
		for k := range initializers ***REMOVED***
			keys = append(keys, string(k))
		***REMOVED***
		sort.Strings(keys)
		return strings.Join(keys, ", ")
	***REMOVED***()
)

// NewStore creates an instance of store
func NewStore(backend store.Backend, addrs []string, options *store.Config) (store.Store, error) ***REMOVED***
	if init, exists := initializers[backend]; exists ***REMOVED***
		return init(addrs, options)
	***REMOVED***

	return nil, fmt.Errorf("%s %s", store.ErrBackendNotSupported.Error(), supportedBackend)
***REMOVED***

// AddStore adds a new store backend to libkv
func AddStore(store store.Backend, init Initialize) ***REMOVED***
	initializers[store] = init
***REMOVED***
