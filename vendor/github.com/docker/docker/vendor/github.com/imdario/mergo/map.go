// Copyright 2014 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func changeInitialCase(s string, mapper func(rune) rune) string ***REMOVED***
	if s == "" ***REMOVED***
		return s
	***REMOVED***
	r, n := utf8.DecodeRuneInString(s)
	return string(mapper(r)) + s[n:]
***REMOVED***

func isExported(field reflect.StructField) bool ***REMOVED***
	r, _ := utf8.DecodeRuneInString(field.Name)
	return r >= 'A' && r <= 'Z'
***REMOVED***

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMap(dst, src reflect.Value, visited map[uintptr]*visit, depth int, overwrite bool) (err error) ***REMOVED***
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
	zeroValue := reflect.Value***REMOVED******REMOVED***
	switch dst.Kind() ***REMOVED***
	case reflect.Map:
		dstMap := dst.Interface().(map[string]interface***REMOVED******REMOVED***)
		for i, n := 0, src.NumField(); i < n; i++ ***REMOVED***
			srcType := src.Type()
			field := srcType.Field(i)
			if !isExported(field) ***REMOVED***
				continue
			***REMOVED***
			fieldName := field.Name
			fieldName = changeInitialCase(fieldName, unicode.ToLower)
			if v, ok := dstMap[fieldName]; !ok || (isEmptyValue(reflect.ValueOf(v)) || overwrite) ***REMOVED***
				dstMap[fieldName] = src.Field(i).Interface()
			***REMOVED***
		***REMOVED***
	case reflect.Struct:
		srcMap := src.Interface().(map[string]interface***REMOVED******REMOVED***)
		for key := range srcMap ***REMOVED***
			srcValue := srcMap[key]
			fieldName := changeInitialCase(key, unicode.ToUpper)
			dstElement := dst.FieldByName(fieldName)
			if dstElement == zeroValue ***REMOVED***
				// We discard it because the field doesn't exist.
				continue
			***REMOVED***
			srcElement := reflect.ValueOf(srcValue)
			dstKind := dstElement.Kind()
			srcKind := srcElement.Kind()
			if srcKind == reflect.Ptr && dstKind != reflect.Ptr ***REMOVED***
				srcElement = srcElement.Elem()
				srcKind = reflect.TypeOf(srcElement.Interface()).Kind()
			***REMOVED*** else if dstKind == reflect.Ptr ***REMOVED***
				// Can this work? I guess it can't.
				if srcKind != reflect.Ptr && srcElement.CanAddr() ***REMOVED***
					srcPtr := srcElement.Addr()
					srcElement = reflect.ValueOf(srcPtr)
					srcKind = reflect.Ptr
				***REMOVED***
			***REMOVED***
			if !srcElement.IsValid() ***REMOVED***
				continue
			***REMOVED***
			if srcKind == dstKind ***REMOVED***
				if err = deepMerge(dstElement, srcElement, visited, depth+1, overwrite); err != nil ***REMOVED***
					return
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if srcKind == reflect.Map ***REMOVED***
					if err = deepMap(dstElement, srcElement, visited, depth+1, overwrite); err != nil ***REMOVED***
						return
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					return fmt.Errorf("type mismatch on %s field: found %v, expected %v", fieldName, srcKind, dstKind)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Map sets fields' values in dst from src.
// src can be a map with string keys or a struct. dst must be the opposite:
// if src is a map, dst must be a valid pointer to struct. If src is a struct,
// dst must be map[string]interface***REMOVED******REMOVED***.
// It won't merge unexported (private) fields and will do recursively
// any exported field.
// If dst is a map, keys will be src fields' names in lower camel case.
// Missing key in src that doesn't match a field in dst will be skipped. This
// doesn't apply if dst is a map.
// This is separated method from Merge because it is cleaner and it keeps sane
// semantics: merging equal types, mapping different (restricted) types.
func Map(dst, src interface***REMOVED******REMOVED***) error ***REMOVED***
	return _map(dst, src, false)
***REMOVED***

func MapWithOverwrite(dst, src interface***REMOVED******REMOVED***) error ***REMOVED***
	return _map(dst, src, true)
***REMOVED***

func _map(dst, src interface***REMOVED******REMOVED***, overwrite bool) error ***REMOVED***
	var (
		vDst, vSrc reflect.Value
		err        error
	)
	if vDst, vSrc, err = resolveValues(dst, src); err != nil ***REMOVED***
		return err
	***REMOVED***
	// To be friction-less, we redirect equal-type arguments
	// to deepMerge. Only because arguments can be anything.
	if vSrc.Kind() == vDst.Kind() ***REMOVED***
		return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, overwrite)
	***REMOVED***
	switch vSrc.Kind() ***REMOVED***
	case reflect.Struct:
		if vDst.Kind() != reflect.Map ***REMOVED***
			return ErrExpectedMapAsDestination
		***REMOVED***
	case reflect.Map:
		if vDst.Kind() != reflect.Struct ***REMOVED***
			return ErrExpectedStructAsDestination
		***REMOVED***
	default:
		return ErrNotSupported
	***REMOVED***
	return deepMap(vDst, vSrc, make(map[uintptr]*visit), 0, overwrite)
***REMOVED***
