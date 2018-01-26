// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"reflect"
)

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int, overwrite bool) (err error) ***REMOVED***
	if !src.IsValid() ***REMOVED***
		return
	***REMOVED***
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
	switch dst.Kind() ***REMOVED***
	case reflect.Struct:
		for i, n := 0, dst.NumField(); i < n; i++ ***REMOVED***
			if err = deepMerge(dst.Field(i), src.Field(i), visited, depth+1, overwrite); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		for _, key := range src.MapKeys() ***REMOVED***
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() ***REMOVED***
				continue
			***REMOVED***
			dstElement := dst.MapIndex(key)
			switch srcElement.Kind() ***REMOVED***
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
				if srcElement.IsNil() ***REMOVED***
					continue
				***REMOVED***
				fallthrough
			default:
				switch reflect.TypeOf(srcElement.Interface()).Kind() ***REMOVED***
				case reflect.Struct:
					fallthrough
				case reflect.Ptr:
					fallthrough
				case reflect.Map:
					if err = deepMerge(dstElement, srcElement, visited, depth+1, overwrite); err != nil ***REMOVED***
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if !isEmptyValue(srcElement) && (overwrite || (!dstElement.IsValid() || isEmptyValue(dst))) ***REMOVED***
				if dst.IsNil() ***REMOVED***
					dst.Set(reflect.MakeMap(dst.Type()))
				***REMOVED***
				dst.SetMapIndex(key, srcElement)
			***REMOVED***
		***REMOVED***
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if src.IsNil() ***REMOVED***
			break
		***REMOVED*** else if dst.IsNil() ***REMOVED***
			if dst.CanSet() && (overwrite || isEmptyValue(dst)) ***REMOVED***
				dst.Set(src)
			***REMOVED***
		***REMOVED*** else if err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, overwrite); err != nil ***REMOVED***
			return
		***REMOVED***
	default:
		if dst.CanSet() && !isEmptyValue(src) && (overwrite || isEmptyValue(dst)) ***REMOVED***
			dst.Set(src)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Merge sets fields' values in dst from src if they have a zero
// value of their type.
// dst and src must be valid same-type structs and dst must be
// a pointer to struct.
// It won't merge unexported (private) fields and will do recursively
// any exported field.
func Merge(dst, src interface***REMOVED******REMOVED***) error ***REMOVED***
	return merge(dst, src, false)
***REMOVED***

func MergeWithOverwrite(dst, src interface***REMOVED******REMOVED***) error ***REMOVED***
	return merge(dst, src, true)
***REMOVED***

func merge(dst, src interface***REMOVED******REMOVED***, overwrite bool) error ***REMOVED***
	var (
		vDst, vSrc reflect.Value
		err        error
	)
	if vDst, vSrc, err = resolveValues(dst, src); err != nil ***REMOVED***
		return err
	***REMOVED***
	if vDst.Type() != vSrc.Type() ***REMOVED***
		return ErrDifferentArgumentsTypes
	***REMOVED***
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, overwrite)
***REMOVED***
