package awsutil

import (
	"io"
	"reflect"
	"time"
)

// Copy deeply copies a src structure to dst. Useful for copying request and
// response structures.
//
// Can copy between structs of different type, but will only copy fields which
// are assignable, and exist in both structs. Fields which are not assignable,
// or do not exist in both structs are ignored.
func Copy(dst, src interface***REMOVED******REMOVED***) ***REMOVED***
	dstval := reflect.ValueOf(dst)
	if !dstval.IsValid() ***REMOVED***
		panic("Copy dst cannot be nil")
	***REMOVED***

	rcopy(dstval, reflect.ValueOf(src), true)
***REMOVED***

// CopyOf returns a copy of src while also allocating the memory for dst.
// src must be a pointer type or this operation will fail.
func CopyOf(src interface***REMOVED******REMOVED***) (dst interface***REMOVED******REMOVED***) ***REMOVED***
	dsti := reflect.New(reflect.TypeOf(src).Elem())
	dst = dsti.Interface()
	rcopy(dsti, reflect.ValueOf(src), true)
	return
***REMOVED***

// rcopy performs a recursive copy of values from the source to destination.
//
// root is used to skip certain aspects of the copy which are not valid
// for the root node of a object.
func rcopy(dst, src reflect.Value, root bool) ***REMOVED***
	if !src.IsValid() ***REMOVED***
		return
	***REMOVED***

	switch src.Kind() ***REMOVED***
	case reflect.Ptr:
		if _, ok := src.Interface().(io.Reader); ok ***REMOVED***
			if dst.Kind() == reflect.Ptr && dst.Elem().CanSet() ***REMOVED***
				dst.Elem().Set(src)
			***REMOVED*** else if dst.CanSet() ***REMOVED***
				dst.Set(src)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e := src.Type().Elem()
			if dst.CanSet() && !src.IsNil() ***REMOVED***
				if _, ok := src.Interface().(*time.Time); !ok ***REMOVED***
					dst.Set(reflect.New(e))
				***REMOVED*** else ***REMOVED***
					tempValue := reflect.New(e)
					tempValue.Elem().Set(src.Elem())
					// Sets time.Time's unexported values
					dst.Set(tempValue)
				***REMOVED***
			***REMOVED***
			if src.Elem().IsValid() ***REMOVED***
				// Keep the current root state since the depth hasn't changed
				rcopy(dst.Elem(), src.Elem(), root)
			***REMOVED***
		***REMOVED***
	case reflect.Struct:
		t := dst.Type()
		for i := 0; i < t.NumField(); i++ ***REMOVED***
			name := t.Field(i).Name
			srcVal := src.FieldByName(name)
			dstVal := dst.FieldByName(name)
			if srcVal.IsValid() && dstVal.CanSet() ***REMOVED***
				rcopy(dstVal, srcVal, false)
			***REMOVED***
		***REMOVED***
	case reflect.Slice:
		if src.IsNil() ***REMOVED***
			break
		***REMOVED***

		s := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		dst.Set(s)
		for i := 0; i < src.Len(); i++ ***REMOVED***
			rcopy(dst.Index(i), src.Index(i), false)
		***REMOVED***
	case reflect.Map:
		if src.IsNil() ***REMOVED***
			break
		***REMOVED***

		s := reflect.MakeMap(src.Type())
		dst.Set(s)
		for _, k := range src.MapKeys() ***REMOVED***
			v := src.MapIndex(k)
			v2 := reflect.New(v.Type()).Elem()
			rcopy(v2, v, false)
			dst.SetMapIndex(k, v2)
		***REMOVED***
	default:
		// Assign the value if possible. If its not assignable, the value would
		// need to be converted and the impact of that may be unexpected, or is
		// not compatible with the dst type.
		if src.Type().AssignableTo(dst.Type()) ***REMOVED***
			dst.Set(src)
		***REMOVED***
	***REMOVED***
***REMOVED***
