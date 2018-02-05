// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cldrtree builds and generates a CLDR index file, including all
// inheritance.
//
package cldrtree

//go:generate go test -gen

// cldrtree stores CLDR data in a tree-like structure called Tree. In the CLDR
// data each branch in the tree is indicated by either an element name or an
// attribute value. A Tree does not distinguish between these two cases, but
// rather assumes that all branches can be accessed by an enum with a compact
// range of positive integer values starting from 0.
//
// Each Tree consists of three parts:
//    - a slice mapping compact language identifiers to an offset into a set of
//      indices,
//    - a set of indices, stored as a large blob of uint16 values that encode
//      the actual tree structure of data, and
//    - a set of buckets that each holds a collection of strings.
// each of which is explained in more detail below.
//
//
// Tree lookup
// A tree lookup is done by providing a locale and a "path", which is a
// sequence of enum values. The search starts with getting the index for the
// given locale and then incrementally jumping into the index using the path
// values. If an element cannot be found in the index, the search starts anew
// for the locale's parent locale. The path may change during lookup by means
// of aliasing, described below.
//
// Buckets
// Buckets hold the actual string data of the leaf values of the CLDR tree.
// This data is stored in buckets, rather than one large string, for multiple
// reasons:
//   - it allows representing leaf values more compactly, by storing all leaf
//     values in a single bucket and then needing only needing a uint16 to index
//     into this bucket for all leaf values,
//   - (TBD) allow multiple trees to share subsets of buckets, mostly to allow
//     linking in a smaller amount of data if only a subset of the buckets is
//     needed,
//   - to be nice to go fmt and the compiler.
//
// indices
// An index is a slice of uint16 for which the values are interpreted in one of
// two ways: as a node or a set of leaf values.
// A set of leaf values has the following form:
//      <max_size>, <bucket>, <offset>...
// max_size indicates the maximum enum value for which an offset is defined.
// An offset value of 0xFFFF (missingValue) also indicates an undefined value.
// If defined offset indicates the offset within the given bucket of the string.
// A node value has the following form:
//      <max_size>, <offset_or_alias>...
// max_size indicates the maximum value for which an offset is defined.
// A missing offset may also be indicated with 0. If the high bit (0x8000, or
// inheritMask) is not set, the offset points to the offset within the index
// for the current locale.
// An offset with high bit set is an alias. In this case the uint16 has the form
//       bits:
//         15: 1
//      14-12: negative offset into path relative to current position
//       0-11: new enum value for path element.
// On encountering an alias, the path is modified accordingly and the lookup is
// restarted for the given locale.

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

// TODO:
// - allow two Trees to share the same set of buckets.

// A Builder allows storing CLDR data in compact form.
type Builder struct ***REMOVED***
	table []string

	rootMeta    *metaData
	locales     []locale
	strToBucket map[string]stringInfo
	buckets     [][]byte
	enums       []*enum
	err         error

	// Stats
	size        int
	sizeAll     int
	bucketWaste int
***REMOVED***

const (
	maxBucketSize = 8 * 1024 // 8K
	maxStrlen     = 254      // allow 0xFF sentinel
)

func (b *Builder) setError(err error) ***REMOVED***
	if b.err == nil ***REMOVED***
		b.err = err
	***REMOVED***
***REMOVED***

func (b *Builder) addString(data string) stringInfo ***REMOVED***
	data = b.makeString(data)
	info, ok := b.strToBucket[data]
	if !ok ***REMOVED***
		b.size += len(data)
		x := len(b.buckets) - 1
		bucket := b.buckets[x]
		if len(bucket)+len(data) < maxBucketSize ***REMOVED***
			info.bucket = uint16(x)
			info.bucketPos = uint16(len(bucket))
			b.buckets[x] = append(bucket, data...)
		***REMOVED*** else ***REMOVED***
			info.bucket = uint16(len(b.buckets))
			info.bucketPos = 0
			b.buckets = append(b.buckets, []byte(data))
		***REMOVED***
		b.strToBucket[data] = info
	***REMOVED***
	return info
***REMOVED***

func (b *Builder) addStringToBucket(data string, bucket uint16) stringInfo ***REMOVED***
	data = b.makeString(data)
	info, ok := b.strToBucket[data]
	if !ok || info.bucket != bucket ***REMOVED***
		if ok ***REMOVED***
			b.bucketWaste += len(data)
		***REMOVED***
		b.size += len(data)
		bk := b.buckets[bucket]
		info.bucket = bucket
		info.bucketPos = uint16(len(bk))
		b.buckets[bucket] = append(bk, data...)
		b.strToBucket[data] = info
	***REMOVED***
	return info
***REMOVED***

func (b *Builder) makeString(data string) string ***REMOVED***
	if len(data) > maxStrlen ***REMOVED***
		b.setError(fmt.Errorf("string %q exceeds maximum length of %d", data, maxStrlen))
		data = data[:maxStrlen]
		for i := len(data) - 1; i > len(data)-4; i-- ***REMOVED***
			if utf8.RuneStart(data[i]) ***REMOVED***
				data = data[:i]
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	data = string([]byte***REMOVED***byte(len(data))***REMOVED***) + data
	b.sizeAll += len(data)
	return data
***REMOVED***

type stringInfo struct ***REMOVED***
	bufferPos uint32
	bucket    uint16
	bucketPos uint16
***REMOVED***

// New creates a new Builder.
func New(tableName string) *Builder ***REMOVED***
	b := &Builder***REMOVED***
		strToBucket: map[string]stringInfo***REMOVED******REMOVED***,
		buckets:     [][]byte***REMOVED***nil***REMOVED***, // initialize with first bucket.
	***REMOVED***
	b.rootMeta = &metaData***REMOVED***
		b:        b,
		typeInfo: &typeInfo***REMOVED******REMOVED***,
	***REMOVED***
	return b
***REMOVED***

// Gen writes all the tables and types for the collected data.
func (b *Builder) Gen(w *gen.CodeWriter) error ***REMOVED***
	t, err := build(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return generate(b, t, w)
***REMOVED***

// GenTestData generates tables useful for testing data generated with Gen.
func (b *Builder) GenTestData(w *gen.CodeWriter) error ***REMOVED***
	return generateTestData(b, w)
***REMOVED***

type locale struct ***REMOVED***
	tag  language.Tag
	root *Index
***REMOVED***

// Locale creates an index for the given locale.
func (b *Builder) Locale(t language.Tag) *Index ***REMOVED***
	index := &Index***REMOVED***
		meta: b.rootMeta,
	***REMOVED***
	b.locales = append(b.locales, locale***REMOVED***tag: t, root: index***REMOVED***)
	return index
***REMOVED***

// An Index holds a map of either leaf values or other indices.
type Index struct ***REMOVED***
	meta *metaData

	subIndex []*Index
	values   []keyValue
***REMOVED***

func (i *Index) setError(err error) ***REMOVED*** i.meta.b.setError(err) ***REMOVED***

type keyValue struct ***REMOVED***
	key   enumIndex
	value stringInfo
***REMOVED***

// Element is a CLDR XML element.
type Element interface ***REMOVED***
	GetCommon() *cldr.Common
***REMOVED***

// Index creates a subindex where the type and enum values are not shared
// with siblings by default. The name is derived from the elem. If elem is
// an alias reference, the alias will be resolved and linked. If elem is nil
// Index returns nil.
func (i *Index) Index(elem Element, opt ...Option) *Index ***REMOVED***
	if elem == nil || reflect.ValueOf(elem).IsNil() ***REMOVED***
		return nil
	***REMOVED***
	c := elem.GetCommon()
	o := &options***REMOVED***
		parent: i,
		name:   c.GetCommon().Element(),
	***REMOVED***
	o.fill(opt)
	o.setAlias(elem)
	return i.subIndexForKey(o)
***REMOVED***

// IndexWithName is like Section but derives the name from the given name.
func (i *Index) IndexWithName(name string, opt ...Option) *Index ***REMOVED***
	o := &options***REMOVED***parent: i, name: name***REMOVED***
	o.fill(opt)
	return i.subIndexForKey(o)
***REMOVED***

// IndexFromType creates a subindex the value of tye type attribute as key. It
// will also configure the Index to share the enumeration values with all
// sibling values. If elem is an alias, it will be resolved and linked.
func (i *Index) IndexFromType(elem Element, opts ...Option) *Index ***REMOVED***
	o := &options***REMOVED***
		parent: i,
		name:   elem.GetCommon().Type,
	***REMOVED***
	o.fill(opts)
	o.setAlias(elem)
	useSharedType()(o)
	return i.subIndexForKey(o)
***REMOVED***

// IndexFromAlt creates a subindex the value of tye alt attribute as key. It
// will also configure the Index to share the enumeration values with all
// sibling values. If elem is an alias, it will be resolved and linked.
func (i *Index) IndexFromAlt(elem Element, opts ...Option) *Index ***REMOVED***
	o := &options***REMOVED***
		parent: i,
		name:   elem.GetCommon().Alt,
	***REMOVED***
	o.fill(opts)
	o.setAlias(elem)
	useSharedType()(o)
	return i.subIndexForKey(o)
***REMOVED***

func (i *Index) subIndexForKey(opts *options) *Index ***REMOVED***
	key := opts.name
	if len(i.values) > 0 ***REMOVED***
		panic(fmt.Errorf("cldrtree: adding Index for %q when value already exists", key))
	***REMOVED***
	meta := i.meta.sub(key, opts)
	for _, x := range i.subIndex ***REMOVED***
		if x.meta == meta ***REMOVED***
			return x
		***REMOVED***
	***REMOVED***
	if alias := opts.alias; alias != nil ***REMOVED***
		if a := alias.GetCommon().Alias; a != nil ***REMOVED***
			if a.Source != "locale" ***REMOVED***
				i.setError(fmt.Errorf("cldrtree: non-locale alias not supported %v", a.Path))
			***REMOVED***
			if meta.inheritOffset < 0 ***REMOVED***
				i.setError(fmt.Errorf("cldrtree: alias was already set %v", a.Path))
			***REMOVED***
			path := a.Path
			for ; strings.HasPrefix(path, "../"); path = path[len("../"):] ***REMOVED***
				meta.inheritOffset--
			***REMOVED***
			m := aliasRe.FindStringSubmatch(path)
			if m == nil ***REMOVED***
				i.setError(fmt.Errorf("cldrtree: could not parse alias %q", a.Path))
			***REMOVED*** else ***REMOVED***
				key := m[4]
				if key == "" ***REMOVED***
					key = m[1]
				***REMOVED***
				meta.inheritIndex = key
			***REMOVED***
		***REMOVED***
	***REMOVED***
	x := &Index***REMOVED***meta: meta***REMOVED***
	i.subIndex = append(i.subIndex, x)
	return x
***REMOVED***

var aliasRe = regexp.MustCompile(`^([a-zA-Z]+)(\[@([a-zA-Z-]+)='([a-zA-Z-]+)'\])?`)

// SetValue sets the value, the data from a CLDR XML element, for the given key.
func (i *Index) SetValue(key string, value Element, opt ...Option) ***REMOVED***
	if len(i.subIndex) > 0 ***REMOVED***
		panic(fmt.Errorf("adding value for key %q when index already exists", key))
	***REMOVED***
	o := &options***REMOVED***parent: i***REMOVED***
	o.fill(opt)
	c := value.GetCommon()
	if c.Alias != nil ***REMOVED***
		i.setError(fmt.Errorf("cldrtree: alias not supported for SetValue %v", c.Alias.Path))
	***REMOVED***
	i.setValue(key, c.Data(), o)
***REMOVED***

func (i *Index) setValue(key, data string, o *options) ***REMOVED***
	index, _ := i.meta.typeInfo.lookupSubtype(key, o)
	kv := keyValue***REMOVED***key: index***REMOVED***
	if len(i.values) > 0 ***REMOVED***
		// Add string to the same bucket as the other values.
		bucket := i.values[0].value.bucket
		kv.value = i.meta.b.addStringToBucket(data, bucket)
	***REMOVED*** else ***REMOVED***
		kv.value = i.meta.b.addString(data)
	***REMOVED***
	i.values = append(i.values, kv)
***REMOVED***
