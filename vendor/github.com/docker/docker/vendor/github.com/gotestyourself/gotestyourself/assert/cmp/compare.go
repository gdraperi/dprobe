/*Package cmp provides Comparisons for Assert and Check*/
package cmp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/pmezard/go-difflib/difflib"
)

// Compare two complex values using https://godoc.org/github.com/google/go-cmp/cmp
// and succeeds if the values are equal.
//
// The comparison can be customized using comparison Options.
func Compare(x, y interface***REMOVED******REMOVED***, opts ...cmp.Option) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		diff := cmp.Diff(x, y, opts...)
		return diff == "", "\n" + diff
	***REMOVED***
***REMOVED***

// Equal succeeds if x == y.
func Equal(x, y interface***REMOVED******REMOVED***) func() (success bool, message string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		return x == y, fmt.Sprintf("%v (%T) != %v (%T)", x, x, y, y)
	***REMOVED***
***REMOVED***

// Len succeeds if the sequence has the expected length.
func Len(seq interface***REMOVED******REMOVED***, expected int) func() (bool, string) ***REMOVED***
	return func() (success bool, message string) ***REMOVED***
		defer func() ***REMOVED***
			if e := recover(); e != nil ***REMOVED***
				success = false
				message = fmt.Sprintf("type %T does not have a length", seq)
			***REMOVED***
		***REMOVED***()
		value := reflect.ValueOf(seq)
		length := value.Len()
		if length == expected ***REMOVED***
			return true, ""
		***REMOVED***
		msg := fmt.Sprintf("expected %s (length %d) to have length %d", seq, length, expected)
		return false, msg
	***REMOVED***
***REMOVED***

// NilError succeeds if the last argument is a nil error.
func NilError(arg interface***REMOVED******REMOVED***, args ...interface***REMOVED******REMOVED***) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		msgFunc := func(value reflect.Value) string ***REMOVED***
			return fmt.Sprintf("error is not nil: %s", value.Interface().(error).Error())
		***REMOVED***
		if len(args) == 0 ***REMOVED***
			return isNil(arg, msgFunc)()
		***REMOVED***
		return isNil(args[len(args)-1], msgFunc)()
	***REMOVED***
***REMOVED***

// Contains succeeds if item is in collection. Collection may be a string, map,
// slice, or array.
//
// If collection is a string, item must also be a string, and is compared using
// strings.Contains().
// If collection is a Map, contains will succeed if item is a key in the map.
// If collection is a slice or array, item is compared to each item in the
// sequence using reflect.DeepEqual().
func Contains(collection interface***REMOVED******REMOVED***, item interface***REMOVED******REMOVED***) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		colValue := reflect.ValueOf(collection)
		if !colValue.IsValid() ***REMOVED***
			return false, fmt.Sprintf("nil does not contain items")
		***REMOVED***
		msg := fmt.Sprintf("%v does not contain %v", collection, item)

		itemValue := reflect.ValueOf(item)
		switch colValue.Type().Kind() ***REMOVED***
		case reflect.String:
			if itemValue.Type().Kind() != reflect.String ***REMOVED***
				return false, "string may only contain strings"
			***REMOVED***
			success := strings.Contains(colValue.String(), itemValue.String())
			return success, fmt.Sprintf("string %q does not contain %q", collection, item)

		case reflect.Map:
			if itemValue.Type() != colValue.Type().Key() ***REMOVED***
				return false, fmt.Sprintf(
					"%v can not contain a %v key", colValue.Type(), itemValue.Type())
			***REMOVED***
			index := colValue.MapIndex(itemValue)
			return index.IsValid(), msg

		case reflect.Slice, reflect.Array:
			for i := 0; i < colValue.Len(); i++ ***REMOVED***
				if reflect.DeepEqual(colValue.Index(i).Interface(), item) ***REMOVED***
					return true, ""
				***REMOVED***
			***REMOVED***
			return false, msg
		default:
			return false, fmt.Sprintf("type %T does not contain items", collection)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Panics succeeds if f() panics.
func Panics(f func()) func() (bool, string) ***REMOVED***
	return func() (success bool, message string) ***REMOVED***
		defer func() ***REMOVED***
			if err := recover(); err != nil ***REMOVED***
				success = true
			***REMOVED***
		***REMOVED***()
		f()
		return false, "did not panic"
	***REMOVED***
***REMOVED***

// EqualMultiLine succeeds if the two strings are equal. If they are not equal
// the failure message will be the difference between the two strings.
func EqualMultiLine(x, y string) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		if x == y ***REMOVED***
			return true, ""
		***REMOVED***

		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff***REMOVED***
			A:        difflib.SplitLines(x),
			B:        difflib.SplitLines(y),
			FromFile: "left",
			ToFile:   "right",
			Context:  3,
		***REMOVED***)
		if err != nil ***REMOVED***
			return false, fmt.Sprintf("failed to produce diff: %s", err)
		***REMOVED***
		return false, "\n" + diff
	***REMOVED***
***REMOVED***

// Error succeeds if err is a non-nil error, and the error message equals the
// expected message.
func Error(err error, message string) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		switch ***REMOVED***
		case err == nil:
			return false, "expected an error, got nil"
		case err.Error() != message:
			return false, fmt.Sprintf(
				"expected error message %q, got %q", message, err.Error())
		***REMOVED***
		return true, ""
	***REMOVED***
***REMOVED***

// ErrorContains succeeds if err is a non-nil error, and the error message contains
// the expected substring.
func ErrorContains(err error, substring string) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		switch ***REMOVED***
		case err == nil:
			return false, "expected an error, got nil"
		case !strings.Contains(err.Error(), substring):
			return false, fmt.Sprintf(
				"expected error message to contain %q, got %q", substring, err.Error())
		***REMOVED***
		return true, ""
	***REMOVED***
***REMOVED***

// Nil succeeds if obj is a nil interface, pointer, or function.
//
// Use NilError() for comparing errors. Use Len(obj, 0) for comparing slices,
// maps, and channels.
func Nil(obj interface***REMOVED******REMOVED***) func() (bool, string) ***REMOVED***
	msgFunc := func(value reflect.Value) string ***REMOVED***
		return fmt.Sprintf("%v (type %s) is not nil", reflect.Indirect(value), value.Type())
	***REMOVED***
	return isNil(obj, msgFunc)
***REMOVED***

func isNil(obj interface***REMOVED******REMOVED***, msgFunc func(reflect.Value) string) func() (bool, string) ***REMOVED***
	return func() (bool, string) ***REMOVED***
		if obj == nil ***REMOVED***
			return true, ""
		***REMOVED***
		value := reflect.ValueOf(obj)
		kind := value.Type().Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice ***REMOVED***
			if value.IsNil() ***REMOVED***
				return true, ""
			***REMOVED***
			return false, msgFunc(value)
		***REMOVED***

		return false, fmt.Sprintf("%v (type %s) can not be nil", value, value.Type())
	***REMOVED***
***REMOVED***
