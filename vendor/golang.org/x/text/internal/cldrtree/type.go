// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldrtree

import (
	"log"
	"strconv"
)

// enumIndex is the numerical value of an enum value.
type enumIndex int

// An enum is a collection of enum values.
type enum struct ***REMOVED***
	name   string // the Go type of the enum
	rename func(string) string
	keyMap map[string]enumIndex
	keys   []string
***REMOVED***

// lookup returns the index for the enum corresponding to the string. If s
// currently does not exist it will add the entry.
func (e *enum) lookup(s string) enumIndex ***REMOVED***
	if e.rename != nil ***REMOVED***
		s = e.rename(s)
	***REMOVED***
	x, ok := e.keyMap[s]
	if !ok ***REMOVED***
		if e.keyMap == nil ***REMOVED***
			e.keyMap = map[string]enumIndex***REMOVED******REMOVED***
		***REMOVED***
		u, err := strconv.ParseUint(s, 10, 32)
		if err == nil ***REMOVED***
			for len(e.keys) <= int(u) ***REMOVED***
				x := enumIndex(len(e.keys))
				s := strconv.Itoa(int(x))
				e.keyMap[s] = x
				e.keys = append(e.keys, s)
			***REMOVED***
			if e.keyMap[s] != enumIndex(u) ***REMOVED***
				// TODO: handle more gracefully.
				log.Fatalf("cldrtree: mix of integer and non-integer for %q %v", s, e.keys)
			***REMOVED***
			return enumIndex(u)
		***REMOVED***
		x = enumIndex(len(e.keys))
		e.keyMap[s] = x
		e.keys = append(e.keys, s)
	***REMOVED***
	return x
***REMOVED***

// A typeInfo indicates the set of possible enum values and a mapping from
// these values to subtypes.
type typeInfo struct ***REMOVED***
	enum        *enum
	entries     map[enumIndex]*typeInfo
	keyTypeInfo *typeInfo
	shareKeys   bool
***REMOVED***

func (t *typeInfo) sharedKeys() bool ***REMOVED***
	return t.shareKeys
***REMOVED***

func (t *typeInfo) lookupSubtype(s string, opts *options) (x enumIndex, sub *typeInfo) ***REMOVED***
	if t.enum == nil ***REMOVED***
		if t.enum = opts.sharedEnums; t.enum == nil ***REMOVED***
			t.enum = &enum***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	if opts.sharedEnums != nil && t.enum != opts.sharedEnums ***REMOVED***
		panic("incompatible enums defined")
	***REMOVED***
	x = t.enum.lookup(s)
	if t.entries == nil ***REMOVED***
		t.entries = map[enumIndex]*typeInfo***REMOVED******REMOVED***
	***REMOVED***
	sub, ok := t.entries[x]
	if !ok ***REMOVED***
		sub = opts.sharedType
		if sub == nil ***REMOVED***
			sub = &typeInfo***REMOVED******REMOVED***
		***REMOVED***
		t.entries[x] = sub
	***REMOVED***
	t.shareKeys = opts.sharedType != nil // For analysis purposes.
	return x, sub
***REMOVED***

// metaData includes information about subtypes, possibly sharing commonality
// with sibling branches, and information about inheritance, which may differ
// per branch.
type metaData struct ***REMOVED***
	b *Builder

	parent *metaData

	index    enumIndex // index into the parent's subtype index
	key      string
	elem     string // XML element corresponding to this type.
	typeInfo *typeInfo

	lookup map[enumIndex]*metaData
	subs   []*metaData

	inheritOffset int    // always negative when applicable
	inheritIndex  string // new value for field indicated by inheritOffset
	// inheritType   *metaData
***REMOVED***

func (m *metaData) sub(key string, opts *options) *metaData ***REMOVED***
	if m.lookup == nil ***REMOVED***
		m.lookup = map[enumIndex]*metaData***REMOVED******REMOVED***
	***REMOVED***
	enum, info := m.typeInfo.lookupSubtype(key, opts)
	sub := m.lookup[enum]
	if sub == nil ***REMOVED***
		sub = &metaData***REMOVED***
			b:      m.b,
			parent: m,

			index:    enum,
			key:      key,
			typeInfo: info,
		***REMOVED***
		m.lookup[enum] = sub
		m.subs = append(m.subs, sub)
	***REMOVED***
	return sub
***REMOVED***

func (m *metaData) validate() ***REMOVED***
	for _, s := range m.subs ***REMOVED***
		s.validate()
	***REMOVED***
***REMOVED***
