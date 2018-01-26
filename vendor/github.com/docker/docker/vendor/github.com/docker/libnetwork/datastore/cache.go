package datastore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/docker/libkv/store"
)

type kvMap map[string]KVObject

type cache struct ***REMOVED***
	sync.Mutex
	kmm map[string]kvMap
	ds  *datastore
***REMOVED***

func newCache(ds *datastore) *cache ***REMOVED***
	return &cache***REMOVED***kmm: make(map[string]kvMap), ds: ds***REMOVED***
***REMOVED***

func (c *cache) kmap(kvObject KVObject) (kvMap, error) ***REMOVED***
	var err error

	c.Lock()
	keyPrefix := Key(kvObject.KeyPrefix()...)
	kmap, ok := c.kmm[keyPrefix]
	c.Unlock()

	if ok ***REMOVED***
		return kmap, nil
	***REMOVED***

	kmap = kvMap***REMOVED******REMOVED***

	// Bail out right away if the kvObject does not implement KVConstructor
	ctor, ok := kvObject.(KVConstructor)
	if !ok ***REMOVED***
		return nil, errors.New("error while populating kmap, object does not implement KVConstructor interface")
	***REMOVED***

	kvList, err := c.ds.store.List(keyPrefix)
	if err != nil ***REMOVED***
		if err == store.ErrKeyNotFound ***REMOVED***
			// If the store doesn't have anything then there is nothing to
			// populate in the cache. Just bail out.
			goto out
		***REMOVED***

		return nil, fmt.Errorf("error while populating kmap: %v", err)
	***REMOVED***

	for _, kvPair := range kvList ***REMOVED***
		// Ignore empty kvPair values
		if len(kvPair.Value) == 0 ***REMOVED***
			continue
		***REMOVED***

		dstO := ctor.New()
		err = dstO.SetValue(kvPair.Value)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// Make sure the object has a correct view of the DB index in
		// case we need to modify it and update the DB.
		dstO.SetIndex(kvPair.LastIndex)

		kmap[Key(dstO.Key()...)] = dstO
	***REMOVED***

out:
	// There may multiple go routines racing to fill the
	// cache. The one which places the kmap in c.kmm first
	// wins. The others should just use what the first populated.
	c.Lock()
	kmapNew, ok := c.kmm[keyPrefix]
	if ok ***REMOVED***
		c.Unlock()
		return kmapNew, nil
	***REMOVED***

	c.kmm[keyPrefix] = kmap
	c.Unlock()

	return kmap, nil
***REMOVED***

func (c *cache) add(kvObject KVObject, atomic bool) error ***REMOVED***
	kmap, err := c.kmap(kvObject)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	// If atomic is true, cache needs to maintain its own index
	// for atomicity and the add needs to be atomic.
	if atomic ***REMOVED***
		if prev, ok := kmap[Key(kvObject.Key()...)]; ok ***REMOVED***
			if prev.Index() != kvObject.Index() ***REMOVED***
				c.Unlock()
				return ErrKeyModified
			***REMOVED***
		***REMOVED***

		// Increment index
		index := kvObject.Index()
		index++
		kvObject.SetIndex(index)
	***REMOVED***

	kmap[Key(kvObject.Key()...)] = kvObject
	c.Unlock()
	return nil
***REMOVED***

func (c *cache) del(kvObject KVObject, atomic bool) error ***REMOVED***
	kmap, err := c.kmap(kvObject)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	// If atomic is true, cache needs to maintain its own index
	// for atomicity and del needs to be atomic.
	if atomic ***REMOVED***
		if prev, ok := kmap[Key(kvObject.Key()...)]; ok ***REMOVED***
			if prev.Index() != kvObject.Index() ***REMOVED***
				c.Unlock()
				return ErrKeyModified
			***REMOVED***
		***REMOVED***
	***REMOVED***

	delete(kmap, Key(kvObject.Key()...))
	c.Unlock()
	return nil
***REMOVED***

func (c *cache) get(key string, kvObject KVObject) error ***REMOVED***
	kmap, err := c.kmap(kvObject)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	defer c.Unlock()

	o, ok := kmap[Key(kvObject.Key()...)]
	if !ok ***REMOVED***
		return ErrKeyNotFound
	***REMOVED***

	ctor, ok := o.(KVConstructor)
	if !ok ***REMOVED***
		return errors.New("kvobject does not implement KVConstructor interface. could not get object")
	***REMOVED***

	return ctor.CopyTo(kvObject)
***REMOVED***

func (c *cache) list(kvObject KVObject) ([]KVObject, error) ***REMOVED***
	kmap, err := c.kmap(kvObject)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c.Lock()
	defer c.Unlock()

	var kvol []KVObject
	for _, v := range kmap ***REMOVED***
		kvol = append(kvol, v)
	***REMOVED***

	return kvol, nil
***REMOVED***
