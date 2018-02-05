// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stringset provides a way to represent a collection of strings
// compactly.
package stringset

import "sort"

// A Set holds a collection of strings that can be looked up by an index number.
type Set struct ***REMOVED***
	// These fields are exported to allow for code generation.

	Data  string
	Index []uint16
***REMOVED***

// Elem returns the string with index i. It panics if i is out of range.
func (s *Set) Elem(i int) string ***REMOVED***
	return s.Data[s.Index[i]:s.Index[i+1]]
***REMOVED***

// Len returns the number of strings in the set.
func (s *Set) Len() int ***REMOVED***
	return len(s.Index) - 1
***REMOVED***

// Search returns the index of the given string or -1 if it is not in the set.
// The Set must have been created with strings in sorted order.
func Search(s *Set, str string) int ***REMOVED***
	// TODO: optimize this if it gets used a lot.
	n := len(s.Index) - 1
	p := sort.Search(n, func(i int) bool ***REMOVED***
		return s.Elem(i) >= str
	***REMOVED***)
	if p == n || str != s.Elem(p) ***REMOVED***
		return -1
	***REMOVED***
	return p
***REMOVED***

// A Builder constructs Sets.
type Builder struct ***REMOVED***
	set   Set
	index map[string]int
***REMOVED***

// NewBuilder returns a new and initialized Builder.
func NewBuilder() *Builder ***REMOVED***
	return &Builder***REMOVED***
		set: Set***REMOVED***
			Index: []uint16***REMOVED***0***REMOVED***,
		***REMOVED***,
		index: map[string]int***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// Set creates the set created so far.
func (b *Builder) Set() Set ***REMOVED***
	return b.set
***REMOVED***

// Index returns the index for the given string, which must have been added
// before.
func (b *Builder) Index(s string) int ***REMOVED***
	return b.index[s]
***REMOVED***

// Add adds a string to the index. Strings that are added by a single Add will
// be stored together, unless they match an existing string.
func (b *Builder) Add(ss ...string) ***REMOVED***
	// First check if the string already exists.
	for _, s := range ss ***REMOVED***
		if _, ok := b.index[s]; ok ***REMOVED***
			continue
		***REMOVED***
		b.index[s] = len(b.set.Index) - 1
		b.set.Data += s
		x := len(b.set.Data)
		if x > 0xFFFF ***REMOVED***
			panic("Index too > 0xFFFF")
		***REMOVED***
		b.set.Index = append(b.set.Index, uint16(x))
	***REMOVED***
***REMOVED***
