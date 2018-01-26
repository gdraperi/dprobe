// Copyright (c) 2012, 2013 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a BSD-style license found in the LICENSE file.

package codec

// All non-std package dependencies live in this file,
// so porting to different environment is easy (just update functions).

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

var (
	raisePanicAfterRecover = false
	debugging              = true
)

func panicValToErr(panicVal interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	switch xerr := panicVal.(type) ***REMOVED***
	case error:
		*err = xerr
	case string:
		*err = errors.New(xerr)
	default:
		*err = fmt.Errorf("%v", panicVal)
	***REMOVED***
	if raisePanicAfterRecover ***REMOVED***
		panic(panicVal)
	***REMOVED***
	return
***REMOVED***

func isEmptyValueDeref(v reflect.Value, deref bool) bool ***REMOVED***
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
		if deref ***REMOVED***
			if v.IsNil() ***REMOVED***
				return true
			***REMOVED***
			return isEmptyValueDeref(v.Elem(), deref)
		***REMOVED*** else ***REMOVED***
			return v.IsNil()
		***REMOVED***
	case reflect.Struct:
		// return true if all fields are empty. else return false.

		// we cannot use equality check, because some fields may be maps/slices/etc
		// and consequently the structs are not comparable.
		// return v.Interface() == reflect.Zero(v.Type()).Interface()
		for i, n := 0, v.NumField(); i < n; i++ ***REMOVED***
			if !isEmptyValueDeref(v.Field(i), deref) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func isEmptyValue(v reflect.Value) bool ***REMOVED***
	return isEmptyValueDeref(v, true)
***REMOVED***

func debugf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if debugging ***REMOVED***
		if len(format) == 0 || format[len(format)-1] != '\n' ***REMOVED***
			format = format + "\n"
		***REMOVED***
		fmt.Printf(format, args...)
	***REMOVED***
***REMOVED***

func pruneSignExt(v []byte, pos bool) (n int) ***REMOVED***
	if len(v) < 2 ***REMOVED***
	***REMOVED*** else if pos && v[0] == 0 ***REMOVED***
		for ; v[n] == 0 && n+1 < len(v) && (v[n+1]&(1<<7) == 0); n++ ***REMOVED***
		***REMOVED***
	***REMOVED*** else if !pos && v[0] == 0xff ***REMOVED***
		for ; v[n] == 0xff && n+1 < len(v) && (v[n+1]&(1<<7) != 0); n++ ***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func implementsIntf(typ, iTyp reflect.Type) (success bool, indir int8) ***REMOVED***
	if typ == nil ***REMOVED***
		return
	***REMOVED***
	rt := typ
	// The type might be a pointer and we need to keep
	// dereferencing to the base type until we find an implementation.
	for ***REMOVED***
		if rt.Implements(iTyp) ***REMOVED***
			return true, indir
		***REMOVED***
		if p := rt; p.Kind() == reflect.Ptr ***REMOVED***
			indir++
			if indir >= math.MaxInt8 ***REMOVED*** // insane number of indirections
				return false, 0
			***REMOVED***
			rt = p.Elem()
			continue
		***REMOVED***
		break
	***REMOVED***
	// No luck yet, but if this is a base type (non-pointer), the pointer might satisfy.
	if typ.Kind() != reflect.Ptr ***REMOVED***
		// Not a pointer, but does the pointer work?
		if reflect.PtrTo(typ).Implements(iTyp) ***REMOVED***
			return true, -1
		***REMOVED***
	***REMOVED***
	return false, 0
***REMOVED***
