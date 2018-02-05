// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldrtree

import (
	"reflect"

	"golang.org/x/text/unicode/cldr"
)

// An Option configures an Index.
type Option func(*options)

type options struct ***REMOVED***
	parent *Index

	name  string
	alias *cldr.Common

	sharedType  *typeInfo
	sharedEnums *enum
***REMOVED***

func (o *options) fill(opt []Option) ***REMOVED***
	for _, f := range opt ***REMOVED***
		f(o)
	***REMOVED***
***REMOVED***

// aliasOpt sets an alias from the given node, if the node defines one.
func (o *options) setAlias(n Element) ***REMOVED***
	if n != nil && !reflect.ValueOf(n).IsNil() ***REMOVED***
		o.alias = n.GetCommon()
	***REMOVED***
***REMOVED***

// Enum defines a enumeration type. The resulting option may be passed for the
// construction of multiple Indexes, which they will share the same enum values.
// Calling Gen on a Builder will generate the Enum for the given name. The
// optional values fix the values for the given identifier to the argument
// position (starting at 0). Other values may still be added and will be
// assigned to subsequent values.
func Enum(name string, value ...string) Option ***REMOVED***
	return EnumFunc(name, nil, value...)
***REMOVED***

// EnumFunc is like Enum but also takes a function that allows rewriting keys.
func EnumFunc(name string, rename func(string) string, value ...string) Option ***REMOVED***
	enum := &enum***REMOVED***name: name, rename: rename, keyMap: map[string]enumIndex***REMOVED******REMOVED******REMOVED***
	for _, e := range value ***REMOVED***
		enum.lookup(e)
	***REMOVED***
	return func(o *options) ***REMOVED***
		found := false
		for _, e := range o.parent.meta.b.enums ***REMOVED***
			if e.name == enum.name ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			o.parent.meta.b.enums = append(o.parent.meta.b.enums, enum)
		***REMOVED***
		o.sharedEnums = enum
	***REMOVED***
***REMOVED***

// SharedType returns an option which causes all Indexes to which this option is
// passed to have the same type.
func SharedType() Option ***REMOVED***
	info := &typeInfo***REMOVED******REMOVED***
	return func(o *options) ***REMOVED*** o.sharedType = info ***REMOVED***
***REMOVED***

func useSharedType() Option ***REMOVED***
	return func(o *options) ***REMOVED***
		sub := o.parent.meta.typeInfo.keyTypeInfo
		if sub == nil ***REMOVED***
			sub = &typeInfo***REMOVED******REMOVED***
			o.parent.meta.typeInfo.keyTypeInfo = sub
		***REMOVED***
		o.sharedType = sub
	***REMOVED***
***REMOVED***
