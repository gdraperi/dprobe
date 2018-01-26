package datastore

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/types"
)

//DataStore exported
type DataStore interface ***REMOVED***
	// GetObject gets data from datastore and unmarshals to the specified object
	GetObject(key string, o KVObject) error
	// PutObject adds a new Record based on an object into the datastore
	PutObject(kvObject KVObject) error
	// PutObjectAtomic provides an atomic add and update operation for a Record
	PutObjectAtomic(kvObject KVObject) error
	// DeleteObject deletes a record
	DeleteObject(kvObject KVObject) error
	// DeleteObjectAtomic performs an atomic delete operation
	DeleteObjectAtomic(kvObject KVObject) error
	// DeleteTree deletes a record
	DeleteTree(kvObject KVObject) error
	// Watchable returns whether the store is watchable or not
	Watchable() bool
	// Watch for changes on a KVObject
	Watch(kvObject KVObject, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan KVObject, error)
	// RestartWatch retriggers stopped Watches
	RestartWatch()
	// Active returns if the store is active
	Active() bool
	// List returns of a list of KVObjects belonging to the parent
	// key. The caller must pass a KVObject of the same type as
	// the objects that need to be listed
	List(string, KVObject) ([]KVObject, error)
	// Map returns a Map of KVObjects
	Map(key string, kvObject KVObject) (map[string]KVObject, error)
	// Scope returns the scope of the store
	Scope() string
	// KVStore returns access to the KV Store
	KVStore() store.Store
	// Close closes the data store
	Close()
***REMOVED***

// ErrKeyModified is raised for an atomic update when the update is working on a stale state
var (
	ErrKeyModified = store.ErrKeyModified
	ErrKeyNotFound = store.ErrKeyNotFound
)

type datastore struct ***REMOVED***
	scope      string
	store      store.Store
	cache      *cache
	watchCh    chan struct***REMOVED******REMOVED***
	active     bool
	sequential bool
	sync.Mutex
***REMOVED***

// KVObject is Key/Value interface used by objects to be part of the DataStore
type KVObject interface ***REMOVED***
	// Key method lets an object provide the Key to be used in KV Store
	Key() []string
	// KeyPrefix method lets an object return immediate parent key that can be used for tree walk
	KeyPrefix() []string
	// Value method lets an object marshal its content to be stored in the KV store
	Value() []byte
	// SetValue is used by the datastore to set the object's value when loaded from the data store.
	SetValue([]byte) error
	// Index method returns the latest DB Index as seen by the object
	Index() uint64
	// SetIndex method allows the datastore to store the latest DB Index into the object
	SetIndex(uint64)
	// True if the object exists in the datastore, false if it hasn't been stored yet.
	// When SetIndex() is called, the object has been stored.
	Exists() bool
	// DataScope indicates the storage scope of the KV object
	DataScope() string
	// Skip provides a way for a KV Object to avoid persisting it in the KV Store
	Skip() bool
***REMOVED***

// KVConstructor interface defines methods which can construct a KVObject from another.
type KVConstructor interface ***REMOVED***
	// New returns a new object which is created based on the
	// source object
	New() KVObject
	// CopyTo deep copies the contents of the implementing object
	// to the passed destination object
	CopyTo(KVObject) error
***REMOVED***

// ScopeCfg represents Datastore configuration.
type ScopeCfg struct ***REMOVED***
	Client ScopeClientCfg
***REMOVED***

// ScopeClientCfg represents Datastore Client-only mode configuration
type ScopeClientCfg struct ***REMOVED***
	Provider string
	Address  string
	Config   *store.Config
***REMOVED***

const (
	// LocalScope indicates to store the KV object in local datastore such as boltdb
	LocalScope = "local"
	// GlobalScope indicates to store the KV object in global datastore such as consul/etcd/zookeeper
	GlobalScope = "global"
	// SwarmScope is not indicating a datastore location. It is defined here
	// along with the other two scopes just for consistency.
	SwarmScope    = "swarm"
	defaultPrefix = "/var/lib/docker/network/files"
)

const (
	// NetworkKeyPrefix is the prefix for network key in the kv store
	NetworkKeyPrefix = "network"
	// EndpointKeyPrefix is the prefix for endpoint key in the kv store
	EndpointKeyPrefix = "endpoint"
)

var (
	defaultScopes = makeDefaultScopes()
)

func makeDefaultScopes() map[string]*ScopeCfg ***REMOVED***
	def := make(map[string]*ScopeCfg)
	def[LocalScope] = &ScopeCfg***REMOVED***
		Client: ScopeClientCfg***REMOVED***
			Provider: string(store.BOLTDB),
			Address:  defaultPrefix + "/local-kv.db",
			Config: &store.Config***REMOVED***
				Bucket:            "libnetwork",
				ConnectionTimeout: time.Minute,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	return def
***REMOVED***

var defaultRootChain = []string***REMOVED***"docker", "network", "v1.0"***REMOVED***
var rootChain = defaultRootChain

// DefaultScopes returns a map of default scopes and its config for clients to use.
func DefaultScopes(dataDir string) map[string]*ScopeCfg ***REMOVED***
	if dataDir != "" ***REMOVED***
		defaultScopes[LocalScope].Client.Address = dataDir + "/network/files/local-kv.db"
		return defaultScopes
	***REMOVED***

	defaultScopes[LocalScope].Client.Address = defaultPrefix + "/local-kv.db"
	return defaultScopes
***REMOVED***

// IsValid checks if the scope config has valid configuration.
func (cfg *ScopeCfg) IsValid() bool ***REMOVED***
	if cfg == nil ||
		strings.TrimSpace(cfg.Client.Provider) == "" ||
		strings.TrimSpace(cfg.Client.Address) == "" ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

//Key provides convenient method to create a Key
func Key(key ...string) string ***REMOVED***
	keychain := append(rootChain, key...)
	str := strings.Join(keychain, "/")
	return str + "/"
***REMOVED***

//ParseKey provides convenient method to unpack the key to complement the Key function
func ParseKey(key string) ([]string, error) ***REMOVED***
	chain := strings.Split(strings.Trim(key, "/"), "/")

	// The key must atleast be equal to the rootChain in order to be considered as valid
	if len(chain) <= len(rootChain) || !reflect.DeepEqual(chain[0:len(rootChain)], rootChain) ***REMOVED***
		return nil, types.BadRequestErrorf("invalid Key : %s", key)
	***REMOVED***
	return chain[len(rootChain):], nil
***REMOVED***

// newClient used to connect to KV Store
func newClient(scope string, kv string, addr string, config *store.Config, cached bool) (DataStore, error) ***REMOVED***

	if cached && scope != LocalScope ***REMOVED***
		return nil, fmt.Errorf("caching supported only for scope %s", LocalScope)
	***REMOVED***
	sequential := false
	if scope == LocalScope ***REMOVED***
		sequential = true
	***REMOVED***

	if config == nil ***REMOVED***
		config = &store.Config***REMOVED******REMOVED***
	***REMOVED***

	var addrs []string

	if kv == string(store.BOLTDB) ***REMOVED***
		// Parse file path
		addrs = strings.Split(addr, ",")
	***REMOVED*** else ***REMOVED***
		// Parse URI
		parts := strings.SplitN(addr, "/", 2)
		addrs = strings.Split(parts[0], ",")

		// Add the custom prefix to the root chain
		if len(parts) == 2 ***REMOVED***
			rootChain = append([]string***REMOVED***parts[1]***REMOVED***, defaultRootChain...)
		***REMOVED***
	***REMOVED***

	store, err := libkv.NewStore(store.Backend(kv), addrs, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ds := &datastore***REMOVED***scope: scope, store: store, active: true, watchCh: make(chan struct***REMOVED******REMOVED***), sequential: sequential***REMOVED***
	if cached ***REMOVED***
		ds.cache = newCache(ds)
	***REMOVED***

	return ds, nil
***REMOVED***

// NewDataStore creates a new instance of LibKV data store
func NewDataStore(scope string, cfg *ScopeCfg) (DataStore, error) ***REMOVED***
	if cfg == nil || cfg.Client.Provider == "" || cfg.Client.Address == "" ***REMOVED***
		c, ok := defaultScopes[scope]
		if !ok || c.Client.Provider == "" || c.Client.Address == "" ***REMOVED***
			return nil, fmt.Errorf("unexpected scope %s without configuration passed", scope)
		***REMOVED***

		cfg = c
	***REMOVED***

	var cached bool
	if scope == LocalScope ***REMOVED***
		cached = true
	***REMOVED***

	return newClient(scope, cfg.Client.Provider, cfg.Client.Address, cfg.Client.Config, cached)
***REMOVED***

// NewDataStoreFromConfig creates a new instance of LibKV data store starting from the datastore config data
func NewDataStoreFromConfig(dsc discoverapi.DatastoreConfigData) (DataStore, error) ***REMOVED***
	var (
		ok    bool
		sCfgP *store.Config
	)

	sCfgP, ok = dsc.Config.(*store.Config)
	if !ok && dsc.Config != nil ***REMOVED***
		return nil, fmt.Errorf("cannot parse store configuration: %v", dsc.Config)
	***REMOVED***

	scopeCfg := &ScopeCfg***REMOVED***
		Client: ScopeClientCfg***REMOVED***
			Address:  dsc.Address,
			Provider: dsc.Provider,
			Config:   sCfgP,
		***REMOVED***,
	***REMOVED***

	ds, err := NewDataStore(dsc.Scope, scopeCfg)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to construct datastore client from datastore configuration %v: %v", dsc, err)
	***REMOVED***

	return ds, err
***REMOVED***

func (ds *datastore) Close() ***REMOVED***
	ds.store.Close()
***REMOVED***

func (ds *datastore) Scope() string ***REMOVED***
	return ds.scope
***REMOVED***

func (ds *datastore) Active() bool ***REMOVED***
	return ds.active
***REMOVED***

func (ds *datastore) Watchable() bool ***REMOVED***
	return ds.scope != LocalScope
***REMOVED***

func (ds *datastore) Watch(kvObject KVObject, stopCh <-chan struct***REMOVED******REMOVED***) (<-chan KVObject, error) ***REMOVED***
	sCh := make(chan struct***REMOVED******REMOVED***)

	ctor, ok := kvObject.(KVConstructor)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("error watching object type %T, object does not implement KVConstructor interface", kvObject)
	***REMOVED***

	kvpCh, err := ds.store.Watch(Key(kvObject.Key()...), sCh)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	kvoCh := make(chan KVObject)

	go func() ***REMOVED***
	retry_watch:
		var err error

		// Make sure to get a new instance of watch channel
		ds.Lock()
		watchCh := ds.watchCh
		ds.Unlock()

	loop:
		for ***REMOVED***
			select ***REMOVED***
			case <-stopCh:
				close(sCh)
				return
			case kvPair := <-kvpCh:
				// If the backend KV store gets reset libkv's go routine
				// for the watch can exit resulting in a nil value in
				// channel.
				if kvPair == nil ***REMOVED***
					ds.Lock()
					ds.active = false
					ds.Unlock()
					break loop
				***REMOVED***

				dstO := ctor.New()

				if err = dstO.SetValue(kvPair.Value); err != nil ***REMOVED***
					log.Printf("Could not unmarshal kvpair value = %s", string(kvPair.Value))
					break
				***REMOVED***

				dstO.SetIndex(kvPair.LastIndex)
				kvoCh <- dstO
			***REMOVED***
		***REMOVED***

		// Wait on watch channel for a re-trigger when datastore becomes active
		<-watchCh

		kvpCh, err = ds.store.Watch(Key(kvObject.Key()...), sCh)
		if err != nil ***REMOVED***
			log.Printf("Could not watch the key %s in store: %v", Key(kvObject.Key()...), err)
		***REMOVED***

		goto retry_watch
	***REMOVED***()

	return kvoCh, nil
***REMOVED***

func (ds *datastore) RestartWatch() ***REMOVED***
	ds.Lock()
	defer ds.Unlock()

	ds.active = true
	watchCh := ds.watchCh
	ds.watchCh = make(chan struct***REMOVED******REMOVED***)
	close(watchCh)
***REMOVED***

func (ds *datastore) KVStore() store.Store ***REMOVED***
	return ds.store
***REMOVED***

// PutObjectAtomic adds a new Record based on an object into the datastore
func (ds *datastore) PutObjectAtomic(kvObject KVObject) error ***REMOVED***
	var (
		previous *store.KVPair
		pair     *store.KVPair
		err      error
	)
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	if kvObject == nil ***REMOVED***
		return types.BadRequestErrorf("invalid KV Object : nil")
	***REMOVED***

	kvObjValue := kvObject.Value()

	if kvObjValue == nil ***REMOVED***
		return types.BadRequestErrorf("invalid KV Object with a nil Value for key %s", Key(kvObject.Key()...))
	***REMOVED***

	if kvObject.Skip() ***REMOVED***
		goto add_cache
	***REMOVED***

	if kvObject.Exists() ***REMOVED***
		previous = &store.KVPair***REMOVED***Key: Key(kvObject.Key()...), LastIndex: kvObject.Index()***REMOVED***
	***REMOVED*** else ***REMOVED***
		previous = nil
	***REMOVED***

	_, pair, err = ds.store.AtomicPut(Key(kvObject.Key()...), kvObjValue, previous, nil)
	if err != nil ***REMOVED***
		if err == store.ErrKeyExists ***REMOVED***
			return ErrKeyModified
		***REMOVED***
		return err
	***REMOVED***

	kvObject.SetIndex(pair.LastIndex)

add_cache:
	if ds.cache != nil ***REMOVED***
		// If persistent store is skipped, sequencing needs to
		// happen in cache.
		return ds.cache.add(kvObject, kvObject.Skip())
	***REMOVED***

	return nil
***REMOVED***

// PutObject adds a new Record based on an object into the datastore
func (ds *datastore) PutObject(kvObject KVObject) error ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	if kvObject == nil ***REMOVED***
		return types.BadRequestErrorf("invalid KV Object : nil")
	***REMOVED***

	if kvObject.Skip() ***REMOVED***
		goto add_cache
	***REMOVED***

	if err := ds.putObjectWithKey(kvObject, kvObject.Key()...); err != nil ***REMOVED***
		return err
	***REMOVED***

add_cache:
	if ds.cache != nil ***REMOVED***
		// If persistent store is skipped, sequencing needs to
		// happen in cache.
		return ds.cache.add(kvObject, kvObject.Skip())
	***REMOVED***

	return nil
***REMOVED***

func (ds *datastore) putObjectWithKey(kvObject KVObject, key ...string) error ***REMOVED***
	kvObjValue := kvObject.Value()

	if kvObjValue == nil ***REMOVED***
		return types.BadRequestErrorf("invalid KV Object with a nil Value for key %s", Key(kvObject.Key()...))
	***REMOVED***
	return ds.store.Put(Key(key...), kvObjValue, nil)
***REMOVED***

// GetObject returns a record matching the key
func (ds *datastore) GetObject(key string, o KVObject) error ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	if ds.cache != nil ***REMOVED***
		return ds.cache.get(key, o)
	***REMOVED***

	kvPair, err := ds.store.Get(key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := o.SetValue(kvPair.Value); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Make sure the object has a correct view of the DB index in
	// case we need to modify it and update the DB.
	o.SetIndex(kvPair.LastIndex)
	return nil
***REMOVED***

func (ds *datastore) ensureParent(parent string) error ***REMOVED***
	exists, err := ds.store.Exists(parent)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if exists ***REMOVED***
		return nil
	***REMOVED***
	return ds.store.Put(parent, []byte***REMOVED******REMOVED***, &store.WriteOptions***REMOVED***IsDir: true***REMOVED***)
***REMOVED***

func (ds *datastore) List(key string, kvObject KVObject) ([]KVObject, error) ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	if ds.cache != nil ***REMOVED***
		return ds.cache.list(kvObject)
	***REMOVED***

	var kvol []KVObject
	cb := func(key string, val KVObject) ***REMOVED***
		kvol = append(kvol, val)
	***REMOVED***
	err := ds.iterateKVPairsFromStore(key, kvObject, cb)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return kvol, nil
***REMOVED***

func (ds *datastore) iterateKVPairsFromStore(key string, kvObject KVObject, callback func(string, KVObject)) error ***REMOVED***
	// Bail out right away if the kvObject does not implement KVConstructor
	ctor, ok := kvObject.(KVConstructor)
	if !ok ***REMOVED***
		return fmt.Errorf("error listing objects, object does not implement KVConstructor interface")
	***REMOVED***

	// Make sure the parent key exists
	if err := ds.ensureParent(key); err != nil ***REMOVED***
		return err
	***REMOVED***

	kvList, err := ds.store.List(key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, kvPair := range kvList ***REMOVED***
		if len(kvPair.Value) == 0 ***REMOVED***
			continue
		***REMOVED***

		dstO := ctor.New()
		if err := dstO.SetValue(kvPair.Value); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Make sure the object has a correct view of the DB index in
		// case we need to modify it and update the DB.
		dstO.SetIndex(kvPair.LastIndex)
		callback(kvPair.Key, dstO)
	***REMOVED***

	return nil
***REMOVED***

func (ds *datastore) Map(key string, kvObject KVObject) (map[string]KVObject, error) ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	kvol := make(map[string]KVObject)
	cb := func(key string, val KVObject) ***REMOVED***
		// Trim the leading & trailing "/" to make it consistent across all stores
		kvol[strings.Trim(key, "/")] = val
	***REMOVED***
	err := ds.iterateKVPairsFromStore(key, kvObject, cb)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return kvol, nil
***REMOVED***

// DeleteObject unconditionally deletes a record from the store
func (ds *datastore) DeleteObject(kvObject KVObject) error ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	// cleaup the cache first
	if ds.cache != nil ***REMOVED***
		// If persistent store is skipped, sequencing needs to
		// happen in cache.
		ds.cache.del(kvObject, kvObject.Skip())
	***REMOVED***

	if kvObject.Skip() ***REMOVED***
		return nil
	***REMOVED***

	return ds.store.Delete(Key(kvObject.Key()...))
***REMOVED***

// DeleteObjectAtomic performs atomic delete on a record
func (ds *datastore) DeleteObjectAtomic(kvObject KVObject) error ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	if kvObject == nil ***REMOVED***
		return types.BadRequestErrorf("invalid KV Object : nil")
	***REMOVED***

	previous := &store.KVPair***REMOVED***Key: Key(kvObject.Key()...), LastIndex: kvObject.Index()***REMOVED***

	if kvObject.Skip() ***REMOVED***
		goto del_cache
	***REMOVED***

	if _, err := ds.store.AtomicDelete(Key(kvObject.Key()...), previous); err != nil ***REMOVED***
		if err == store.ErrKeyExists ***REMOVED***
			return ErrKeyModified
		***REMOVED***
		return err
	***REMOVED***

del_cache:
	// cleanup the cache only if AtomicDelete went through successfully
	if ds.cache != nil ***REMOVED***
		// If persistent store is skipped, sequencing needs to
		// happen in cache.
		return ds.cache.del(kvObject, kvObject.Skip())
	***REMOVED***

	return nil
***REMOVED***

// DeleteTree unconditionally deletes a record from the store
func (ds *datastore) DeleteTree(kvObject KVObject) error ***REMOVED***
	if ds.sequential ***REMOVED***
		ds.Lock()
		defer ds.Unlock()
	***REMOVED***

	// cleaup the cache first
	if ds.cache != nil ***REMOVED***
		// If persistent store is skipped, sequencing needs to
		// happen in cache.
		ds.cache.del(kvObject, kvObject.Skip())
	***REMOVED***

	if kvObject.Skip() ***REMOVED***
		return nil
	***REMOVED***

	return ds.store.DeleteTree(Key(kvObject.KeyPrefix()...))
***REMOVED***
