// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package singleflight provides a duplicate function call suppression
// mechanism.
package singleflight // import "golang.org/x/sync/singleflight"

import "sync"

// call is an in-flight or completed singleflight.Do call
type call struct ***REMOVED***
	wg sync.WaitGroup

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val interface***REMOVED******REMOVED***
	err error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int
	chans []chan<- Result
***REMOVED***

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group struct ***REMOVED***
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
***REMOVED***

// Result holds the results of Do, so they can be passed
// on a channel.
type Result struct ***REMOVED***
	Val    interface***REMOVED******REMOVED***
	Err    error
	Shared bool
***REMOVED***

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (g *Group) Do(key string, fn func() (interface***REMOVED******REMOVED***, error)) (v interface***REMOVED******REMOVED***, err error, shared bool) ***REMOVED***
	g.mu.Lock()
	if g.m == nil ***REMOVED***
		g.m = make(map[string]*call)
	***REMOVED***
	if c, ok := g.m[key]; ok ***REMOVED***
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	***REMOVED***
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
***REMOVED***

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
func (g *Group) DoChan(key string, fn func() (interface***REMOVED******REMOVED***, error)) <-chan Result ***REMOVED***
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil ***REMOVED***
		g.m = make(map[string]*call)
	***REMOVED***
	if c, ok := g.m[key]; ok ***REMOVED***
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	***REMOVED***
	c := &call***REMOVED***chans: []chan<- Result***REMOVED***ch***REMOVED******REMOVED***
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
***REMOVED***

// doCall handles the single call for a key.
func (g *Group) doCall(c *call, key string, fn func() (interface***REMOVED******REMOVED***, error)) ***REMOVED***
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	for _, ch := range c.chans ***REMOVED***
		ch <- Result***REMOVED***c.val, c.err, c.dups > 0***REMOVED***
	***REMOVED***
	g.mu.Unlock()
***REMOVED***

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (g *Group) Forget(key string) ***REMOVED***
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
***REMOVED***
