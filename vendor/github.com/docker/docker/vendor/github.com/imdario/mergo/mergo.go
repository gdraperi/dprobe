// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"errors"
	"reflect"
)

// Errors reported by Mergo when it finds invalid arguments.
var (
	ErrNilArguments                = errors.New("src and dst must not be nil")
	ErrDifferentArgumentsTypes     = errors.New("src and dst must be of same type")
	ErrNotSupported                = errors.New("only structs and maps are supported")
	ErrExpectedMapAsDestination    = errors.New("dst was expected to be a map")
	ErrExpectedStructAsDestination = errors.New("dst was expected to be a struct")
)

// During deepMerge, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited are stored in a map indexed by 17 * a1 + a2;
type visit struct ***REMOVED***
	ptr  uintptr
	typ  reflect.Type
	next *visit
***REMOVED***

// From src/pkg/encoding/json.
func isEmptyValue(v reflect.Value) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	***REMOVED***
	return false
***REMOVED***

func resolveValues(dst, src interface***REMOVED******REMOVED***) (vDst, vSrc reflect.Value, err error) ***REMOVED***
	if dst == nil || src == nil ***REMOVED***
		err = ErrNilArguments
		return
	***REMOVED***
	vDst = reflect.ValueOf(dst).Elem()
	if vDst.Kind() != reflect.Struct && vDst.Kind() != reflect.Map ***REMOVED***
		err = ErrNotSupported
		return
	***REMOVED***
	vSrc = reflect.ValueOf(src)
	// We check if vSrc is a pointer to dereference it.
	if vSrc.Kind() == reflect.Ptr ***REMOVED***
		vSrc = vSrc.Elem()
	***REMOVED***
	return
***REMOVED***

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deeper(dst, src reflect.Value, visited map[uintptr]*visit, depth int) (err error) ***REMOVED***
	if dst.CanAddr() ***REMOVED***
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next ***REMOVED***
			if p.ptr == addr && p.typ == typ ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		// Remember, remember...
		visited[h] = &visit***REMOVED***addr, typ, seen***REMOVED***
	***REMOVED***
	return // TODO refactor
***REMOVED***
