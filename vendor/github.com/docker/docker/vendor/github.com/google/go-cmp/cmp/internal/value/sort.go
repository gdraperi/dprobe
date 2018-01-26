// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package value

import (
	"fmt"
	"math"
	"reflect"
	"sort"
)

// SortKeys sorts a list of map keys, deduplicating keys if necessary.
// The type of each value must be comparable.
func SortKeys(vs []reflect.Value) []reflect.Value ***REMOVED***
	if len(vs) == 0 ***REMOVED***
		return vs
	***REMOVED***

	// Sort the map keys.
	sort.Sort(valueSorter(vs))

	// Deduplicate keys (fails for NaNs).
	vs2 := vs[:1]
	for _, v := range vs[1:] ***REMOVED***
		if v.Interface() != vs2[len(vs2)-1].Interface() ***REMOVED***
			vs2 = append(vs2, v)
		***REMOVED***
	***REMOVED***
	return vs2
***REMOVED***

// TODO: Use sort.Slice once Google AppEngine is on Go1.8 or above.
type valueSorter []reflect.Value

func (vs valueSorter) Len() int           ***REMOVED*** return len(vs) ***REMOVED***
func (vs valueSorter) Less(i, j int) bool ***REMOVED*** return isLess(vs[i], vs[j]) ***REMOVED***
func (vs valueSorter) Swap(i, j int)      ***REMOVED*** vs[i], vs[j] = vs[j], vs[i] ***REMOVED***

// isLess is a generic function for sorting arbitrary map keys.
// The inputs must be of the same type and must be comparable.
func isLess(x, y reflect.Value) bool ***REMOVED***
	switch x.Type().Kind() ***REMOVED***
	case reflect.Bool:
		return !x.Bool() && y.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return x.Int() < y.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return x.Uint() < y.Uint()
	case reflect.Float32, reflect.Float64:
		fx, fy := x.Float(), y.Float()
		return fx < fy || math.IsNaN(fx) && !math.IsNaN(fy)
	case reflect.Complex64, reflect.Complex128:
		cx, cy := x.Complex(), y.Complex()
		rx, ix, ry, iy := real(cx), imag(cx), real(cy), imag(cy)
		if rx == ry || (math.IsNaN(rx) && math.IsNaN(ry)) ***REMOVED***
			return ix < iy || math.IsNaN(ix) && !math.IsNaN(iy)
		***REMOVED***
		return rx < ry || math.IsNaN(rx) && !math.IsNaN(ry)
	case reflect.Ptr, reflect.UnsafePointer, reflect.Chan:
		return x.Pointer() < y.Pointer()
	case reflect.String:
		return x.String() < y.String()
	case reflect.Array:
		for i := 0; i < x.Len(); i++ ***REMOVED***
			if isLess(x.Index(i), y.Index(i)) ***REMOVED***
				return true
			***REMOVED***
			if isLess(y.Index(i), x.Index(i)) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return false
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ ***REMOVED***
			if isLess(x.Field(i), y.Field(i)) ***REMOVED***
				return true
			***REMOVED***
			if isLess(y.Field(i), x.Field(i)) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return false
	case reflect.Interface:
		vx, vy := x.Elem(), y.Elem()
		if !vx.IsValid() || !vy.IsValid() ***REMOVED***
			return !vx.IsValid() && vy.IsValid()
		***REMOVED***
		tx, ty := vx.Type(), vy.Type()
		if tx == ty ***REMOVED***
			return isLess(x.Elem(), y.Elem())
		***REMOVED***
		if tx.Kind() != ty.Kind() ***REMOVED***
			return vx.Kind() < vy.Kind()
		***REMOVED***
		if tx.String() != ty.String() ***REMOVED***
			return tx.String() < ty.String()
		***REMOVED***
		if tx.PkgPath() != ty.PkgPath() ***REMOVED***
			return tx.PkgPath() < ty.PkgPath()
		***REMOVED***
		// This can happen in rare situations, so we fallback to just comparing
		// the unique pointer for a reflect.Type. This guarantees deterministic
		// ordering within a program, but it is obviously not stable.
		return reflect.ValueOf(vx.Type()).Pointer() < reflect.ValueOf(vy.Type()).Pointer()
	default:
		// Must be Func, Map, or Slice; which are not comparable.
		panic(fmt.Sprintf("%T is not comparable", x.Type()))
	***REMOVED***
***REMOVED***
