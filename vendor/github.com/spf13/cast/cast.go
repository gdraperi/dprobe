// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cast provides easy and safe casting in Go.
package cast

import "time"

// ToBool casts an interface to a bool type.
func ToBool(i interface***REMOVED******REMOVED***) bool ***REMOVED***
	v, _ := ToBoolE(i)
	return v
***REMOVED***

// ToTime casts an interface to a time.Time type.
func ToTime(i interface***REMOVED******REMOVED***) time.Time ***REMOVED***
	v, _ := ToTimeE(i)
	return v
***REMOVED***

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface***REMOVED******REMOVED***) time.Duration ***REMOVED***
	v, _ := ToDurationE(i)
	return v
***REMOVED***

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i interface***REMOVED******REMOVED***) float64 ***REMOVED***
	v, _ := ToFloat64E(i)
	return v
***REMOVED***

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i interface***REMOVED******REMOVED***) float32 ***REMOVED***
	v, _ := ToFloat32E(i)
	return v
***REMOVED***

// ToInt64 casts an interface to an int64 type.
func ToInt64(i interface***REMOVED******REMOVED***) int64 ***REMOVED***
	v, _ := ToInt64E(i)
	return v
***REMOVED***

// ToInt32 casts an interface to an int32 type.
func ToInt32(i interface***REMOVED******REMOVED***) int32 ***REMOVED***
	v, _ := ToInt32E(i)
	return v
***REMOVED***

// ToInt16 casts an interface to an int16 type.
func ToInt16(i interface***REMOVED******REMOVED***) int16 ***REMOVED***
	v, _ := ToInt16E(i)
	return v
***REMOVED***

// ToInt8 casts an interface to an int8 type.
func ToInt8(i interface***REMOVED******REMOVED***) int8 ***REMOVED***
	v, _ := ToInt8E(i)
	return v
***REMOVED***

// ToInt casts an interface to an int type.
func ToInt(i interface***REMOVED******REMOVED***) int ***REMOVED***
	v, _ := ToIntE(i)
	return v
***REMOVED***

// ToUint casts an interface to a uint type.
func ToUint(i interface***REMOVED******REMOVED***) uint ***REMOVED***
	v, _ := ToUintE(i)
	return v
***REMOVED***

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface***REMOVED******REMOVED***) uint64 ***REMOVED***
	v, _ := ToUint64E(i)
	return v
***REMOVED***

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface***REMOVED******REMOVED***) uint32 ***REMOVED***
	v, _ := ToUint32E(i)
	return v
***REMOVED***

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface***REMOVED******REMOVED***) uint16 ***REMOVED***
	v, _ := ToUint16E(i)
	return v
***REMOVED***

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface***REMOVED******REMOVED***) uint8 ***REMOVED***
	v, _ := ToUint8E(i)
	return v
***REMOVED***

// ToString casts an interface to a string type.
func ToString(i interface***REMOVED******REMOVED***) string ***REMOVED***
	v, _ := ToStringE(i)
	return v
***REMOVED***

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapString(i interface***REMOVED******REMOVED***) map[string]string ***REMOVED***
	v, _ := ToStringMapStringE(i)
	return v
***REMOVED***

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i interface***REMOVED******REMOVED***) map[string][]string ***REMOVED***
	v, _ := ToStringMapStringSliceE(i)
	return v
***REMOVED***

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i interface***REMOVED******REMOVED***) map[string]bool ***REMOVED***
	v, _ := ToStringMapBoolE(i)
	return v
***REMOVED***

// ToStringMap casts an interface to a map[string]interface***REMOVED******REMOVED*** type.
func ToStringMap(i interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	v, _ := ToStringMapE(i)
	return v
***REMOVED***

// ToSlice casts an interface to a []interface***REMOVED******REMOVED*** type.
func ToSlice(i interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	v, _ := ToSliceE(i)
	return v
***REMOVED***

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i interface***REMOVED******REMOVED***) []bool ***REMOVED***
	v, _ := ToBoolSliceE(i)
	return v
***REMOVED***

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface***REMOVED******REMOVED***) []string ***REMOVED***
	v, _ := ToStringSliceE(i)
	return v
***REMOVED***

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i interface***REMOVED******REMOVED***) []int ***REMOVED***
	v, _ := ToIntSliceE(i)
	return v
***REMOVED***

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface***REMOVED******REMOVED***) []time.Duration ***REMOVED***
	v, _ := ToDurationSliceE(i)
	return v
***REMOVED***
