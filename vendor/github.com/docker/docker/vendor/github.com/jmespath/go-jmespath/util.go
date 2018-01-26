package jmespath

import (
	"errors"
	"reflect"
)

// IsFalse determines if an object is false based on the JMESPath spec.
// JMESPath defines false values to be any of:
// - An empty string array, or hash.
// - The boolean value false.
// - nil
func isFalse(value interface***REMOVED******REMOVED***) bool ***REMOVED***
	switch v := value.(type) ***REMOVED***
	case bool:
		return !v
	case []interface***REMOVED******REMOVED***:
		return len(v) == 0
	case map[string]interface***REMOVED******REMOVED***:
		return len(v) == 0
	case string:
		return len(v) == 0
	case nil:
		return true
	***REMOVED***
	// Try the reflection cases before returning false.
	rv := reflect.ValueOf(value)
	switch rv.Kind() ***REMOVED***
	case reflect.Struct:
		// A struct type will never be false, even if
		// all of its values are the zero type.
		return false
	case reflect.Slice, reflect.Map:
		return rv.Len() == 0
	case reflect.Ptr:
		if rv.IsNil() ***REMOVED***
			return true
		***REMOVED***
		// If it's a pointer type, we'll try to deref the pointer
		// and evaluate the pointer value for isFalse.
		element := rv.Elem()
		return isFalse(element.Interface())
	***REMOVED***
	return false
***REMOVED***

// ObjsEqual is a generic object equality check.
// It will take two arbitrary objects and recursively determine
// if they are equal.
func objsEqual(left interface***REMOVED******REMOVED***, right interface***REMOVED******REMOVED***) bool ***REMOVED***
	return reflect.DeepEqual(left, right)
***REMOVED***

// SliceParam refers to a single part of a slice.
// A slice consists of a start, a stop, and a step, similar to
// python slices.
type sliceParam struct ***REMOVED***
	N         int
	Specified bool
***REMOVED***

// Slice supports [start:stop:step] style slicing that's supported in JMESPath.
func slice(slice []interface***REMOVED******REMOVED***, parts []sliceParam) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	computed, err := computeSliceParams(len(slice), parts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	start, stop, step := computed[0], computed[1], computed[2]
	result := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	if step > 0 ***REMOVED***
		for i := start; i < stop; i += step ***REMOVED***
			result = append(result, slice[i])
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := start; i > stop; i += step ***REMOVED***
			result = append(result, slice[i])
		***REMOVED***
	***REMOVED***
	return result, nil
***REMOVED***

func computeSliceParams(length int, parts []sliceParam) ([]int, error) ***REMOVED***
	var start, stop, step int
	if !parts[2].Specified ***REMOVED***
		step = 1
	***REMOVED*** else if parts[2].N == 0 ***REMOVED***
		return nil, errors.New("Invalid slice, step cannot be 0")
	***REMOVED*** else ***REMOVED***
		step = parts[2].N
	***REMOVED***
	var stepValueNegative bool
	if step < 0 ***REMOVED***
		stepValueNegative = true
	***REMOVED*** else ***REMOVED***
		stepValueNegative = false
	***REMOVED***

	if !parts[0].Specified ***REMOVED***
		if stepValueNegative ***REMOVED***
			start = length - 1
		***REMOVED*** else ***REMOVED***
			start = 0
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		start = capSlice(length, parts[0].N, step)
	***REMOVED***

	if !parts[1].Specified ***REMOVED***
		if stepValueNegative ***REMOVED***
			stop = -1
		***REMOVED*** else ***REMOVED***
			stop = length
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		stop = capSlice(length, parts[1].N, step)
	***REMOVED***
	return []int***REMOVED***start, stop, step***REMOVED***, nil
***REMOVED***

func capSlice(length int, actual int, step int) int ***REMOVED***
	if actual < 0 ***REMOVED***
		actual += length
		if actual < 0 ***REMOVED***
			if step < 0 ***REMOVED***
				actual = -1
			***REMOVED*** else ***REMOVED***
				actual = 0
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if actual >= length ***REMOVED***
		if step < 0 ***REMOVED***
			actual = length - 1
		***REMOVED*** else ***REMOVED***
			actual = length
		***REMOVED***
	***REMOVED***
	return actual
***REMOVED***

// ToArrayNum converts an empty interface type to a slice of float64.
// If any element in the array cannot be converted, then nil is returned
// along with a second value of false.
func toArrayNum(data interface***REMOVED******REMOVED***) ([]float64, bool) ***REMOVED***
	// Is there a better way to do this with reflect?
	if d, ok := data.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		result := make([]float64, len(d))
		for i, el := range d ***REMOVED***
			item, ok := el.(float64)
			if !ok ***REMOVED***
				return nil, false
			***REMOVED***
			result[i] = item
		***REMOVED***
		return result, true
	***REMOVED***
	return nil, false
***REMOVED***

// ToArrayStr converts an empty interface type to a slice of strings.
// If any element in the array cannot be converted, then nil is returned
// along with a second value of false.  If the input data could be entirely
// converted, then the converted data, along with a second value of true,
// will be returned.
func toArrayStr(data interface***REMOVED******REMOVED***) ([]string, bool) ***REMOVED***
	// Is there a better way to do this with reflect?
	if d, ok := data.([]interface***REMOVED******REMOVED***); ok ***REMOVED***
		result := make([]string, len(d))
		for i, el := range d ***REMOVED***
			item, ok := el.(string)
			if !ok ***REMOVED***
				return nil, false
			***REMOVED***
			result[i] = item
		***REMOVED***
		return result, true
	***REMOVED***
	return nil, false
***REMOVED***

func isSliceType(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	if v == nil ***REMOVED***
		return false
	***REMOVED***
	return reflect.TypeOf(v).Kind() == reflect.Slice
***REMOVED***
