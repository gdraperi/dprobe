// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldrtree

import (
	"golang.org/x/text/internal"
	"golang.org/x/text/language"
)

const (
	inheritOffsetShift        = 12
	inheritMask        uint16 = 0x8000
	inheritValueMask   uint16 = 0x0FFF

	missingValue uint16 = 0xFFFF
)

// Tree holds a tree of CLDR data.
type Tree struct ***REMOVED***
	Locales []uint32
	Indices []uint16
	Buckets []string
***REMOVED***

// Lookup looks up CLDR data for the given path. The lookup adheres to the alias
// and locale inheritance rules as defined in CLDR.
//
// Each subsequent element in path indicates which subtree to select data from.
// The last element of the path must select a leaf node. All other elements
// of the path select a subindex.
func (t *Tree) Lookup(tag int, path ...uint16) string ***REMOVED***
	return t.lookup(tag, false, path...)
***REMOVED***

// LookupFeature is like Lookup, but will first check whether a value of "other"
// as a fallback before traversing the inheritance chain.
func (t *Tree) LookupFeature(tag int, path ...uint16) string ***REMOVED***
	return t.lookup(tag, true, path...)
***REMOVED***

func (t *Tree) lookup(tag int, isFeature bool, path ...uint16) string ***REMOVED***
	origLang := tag
outer:
	for ***REMOVED***
		index := t.Indices[t.Locales[tag]:]

		k := uint16(0)
		for i := range path ***REMOVED***
			max := index[k]
			if i < len(path)-1 ***REMOVED***
				// index (non-leaf)
				if path[i] >= max ***REMOVED***
					break
				***REMOVED***
				k = index[k+1+path[i]]
				if k == 0 ***REMOVED***
					break
				***REMOVED***
				if v := k &^ inheritMask; k != v ***REMOVED***
					offset := v >> inheritOffsetShift
					value := v & inheritValueMask
					path[uint16(i)-offset] = value
					tag = origLang
					continue outer
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// leaf value
				offset := missingValue
				if path[i] < max ***REMOVED***
					offset = index[k+2+path[i]]
				***REMOVED***
				if offset == missingValue ***REMOVED***
					if !isFeature ***REMOVED***
						break
					***REMOVED***
					// "other" feature must exist
					offset = index[k+2]
				***REMOVED***
				data := t.Buckets[index[k+1]]
				n := uint16(data[offset])
				return data[offset+1 : offset+n+1]
			***REMOVED***
		***REMOVED***
		if tag == 0 ***REMOVED***
			break
		***REMOVED***
		tag = int(internal.Parent[tag])
	***REMOVED***
	return ""
***REMOVED***

func build(b *Builder) (*Tree, error) ***REMOVED***
	var t Tree

	t.Locales = make([]uint32, language.NumCompactTags)

	for _, loc := range b.locales ***REMOVED***
		tag, _ := language.CompactIndex(loc.tag)
		t.Locales[tag] = uint32(len(t.Indices))
		var x indexBuilder
		x.add(loc.root)
		t.Indices = append(t.Indices, x.index...)
	***REMOVED***
	// Set locales for which we don't have data to the parent's data.
	for i, v := range t.Locales ***REMOVED***
		p := uint16(i)
		for v == 0 && p != 0 ***REMOVED***
			p = internal.Parent[p]
			v = t.Locales[p]
		***REMOVED***
		t.Locales[i] = v
	***REMOVED***

	for _, b := range b.buckets ***REMOVED***
		t.Buckets = append(t.Buckets, string(b))
	***REMOVED***
	if b.err != nil ***REMOVED***
		return nil, b.err
	***REMOVED***
	return &t, nil
***REMOVED***

type indexBuilder struct ***REMOVED***
	index []uint16
***REMOVED***

func (b *indexBuilder) add(i *Index) uint16 ***REMOVED***
	offset := len(b.index)

	max := enumIndex(0)
	switch ***REMOVED***
	case len(i.values) > 0:
		for _, v := range i.values ***REMOVED***
			if v.key > max ***REMOVED***
				max = v.key
			***REMOVED***
		***REMOVED***
		b.index = append(b.index, make([]uint16, max+3)...)

		b.index[offset] = uint16(max) + 1

		b.index[offset+1] = i.values[0].value.bucket
		for i := offset + 2; i < len(b.index); i++ ***REMOVED***
			b.index[i] = missingValue
		***REMOVED***
		for _, v := range i.values ***REMOVED***
			b.index[offset+2+int(v.key)] = v.value.bucketPos
		***REMOVED***
		return uint16(offset)

	case len(i.subIndex) > 0:
		for _, s := range i.subIndex ***REMOVED***
			if s.meta.index > max ***REMOVED***
				max = s.meta.index
			***REMOVED***
		***REMOVED***
		b.index = append(b.index, make([]uint16, max+2)...)

		b.index[offset] = uint16(max) + 1

		for _, s := range i.subIndex ***REMOVED***
			x := b.add(s)
			b.index[offset+int(s.meta.index)+1] = x
		***REMOVED***
		return uint16(offset)

	case i.meta.inheritOffset < 0:
		v := uint16(-(i.meta.inheritOffset + 1)) << inheritOffsetShift
		p := i.meta
		for k := i.meta.inheritOffset; k < 0; k++ ***REMOVED***
			p = p.parent
		***REMOVED***
		v += uint16(p.typeInfo.enum.lookup(i.meta.inheritIndex))
		v |= inheritMask
		return v
	***REMOVED***

	return 0
***REMOVED***
