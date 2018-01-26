package aws

import "time"

// String returns a pointer to the string value passed in.
func String(v string) *string ***REMOVED***
	return &v
***REMOVED***

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return ""
***REMOVED***

// StringSlice converts a slice of string values into a slice of
// string pointers
func StringSlice(src []string) []*string ***REMOVED***
	dst := make([]*string, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// StringValueSlice converts a slice of string pointers into a slice of
// string values
func StringValueSlice(src []*string) []string ***REMOVED***
	dst := make([]string, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// StringMap converts a string map of string values into a string
// map of string pointers
func StringMap(src map[string]string) map[string]*string ***REMOVED***
	dst := make(map[string]*string)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// StringValueMap converts a string map of string pointers into a string
// map of string values
func StringValueMap(src map[string]*string) map[string]string ***REMOVED***
	dst := make(map[string]string)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool ***REMOVED***
	return &v
***REMOVED***

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return false
***REMOVED***

// BoolSlice converts a slice of bool values into a slice of
// bool pointers
func BoolSlice(src []bool) []*bool ***REMOVED***
	dst := make([]*bool, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// BoolValueSlice converts a slice of bool pointers into a slice of
// bool values
func BoolValueSlice(src []*bool) []bool ***REMOVED***
	dst := make([]bool, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// BoolMap converts a string map of bool values into a string
// map of bool pointers
func BoolMap(src map[string]bool) map[string]*bool ***REMOVED***
	dst := make(map[string]*bool)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// BoolValueMap converts a string map of bool pointers into a string
// map of bool values
func BoolValueMap(src map[string]*bool) map[string]bool ***REMOVED***
	dst := make(map[string]bool)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Int returns a pointer to the int value passed in.
func Int(v int) *int ***REMOVED***
	return &v
***REMOVED***

// IntValue returns the value of the int pointer passed in or
// 0 if the pointer is nil.
func IntValue(v *int) int ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return 0
***REMOVED***

// IntSlice converts a slice of int values into a slice of
// int pointers
func IntSlice(src []int) []*int ***REMOVED***
	dst := make([]*int, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// IntValueSlice converts a slice of int pointers into a slice of
// int values
func IntValueSlice(src []*int) []int ***REMOVED***
	dst := make([]int, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// IntMap converts a string map of int values into a string
// map of int pointers
func IntMap(src map[string]int) map[string]*int ***REMOVED***
	dst := make(map[string]*int)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// IntValueMap converts a string map of int pointers into a string
// map of int values
func IntValueMap(src map[string]*int) map[string]int ***REMOVED***
	dst := make(map[string]int)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 ***REMOVED***
	return &v
***REMOVED***

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return 0
***REMOVED***

// Int64Slice converts a slice of int64 values into a slice of
// int64 pointers
func Int64Slice(src []int64) []*int64 ***REMOVED***
	dst := make([]*int64, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// Int64ValueSlice converts a slice of int64 pointers into a slice of
// int64 values
func Int64ValueSlice(src []*int64) []int64 ***REMOVED***
	dst := make([]int64, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Int64Map converts a string map of int64 values into a string
// map of int64 pointers
func Int64Map(src map[string]int64) map[string]*int64 ***REMOVED***
	dst := make(map[string]*int64)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// Int64ValueMap converts a string map of int64 pointers into a string
// map of int64 values
func Int64ValueMap(src map[string]*int64) map[string]int64 ***REMOVED***
	dst := make(map[string]int64)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Float64 returns a pointer to the float64 value passed in.
func Float64(v float64) *float64 ***REMOVED***
	return &v
***REMOVED***

// Float64Value returns the value of the float64 pointer passed in or
// 0 if the pointer is nil.
func Float64Value(v *float64) float64 ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return 0
***REMOVED***

// Float64Slice converts a slice of float64 values into a slice of
// float64 pointers
func Float64Slice(src []float64) []*float64 ***REMOVED***
	dst := make([]*float64, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// Float64ValueSlice converts a slice of float64 pointers into a slice of
// float64 values
func Float64ValueSlice(src []*float64) []float64 ***REMOVED***
	dst := make([]float64, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Float64Map converts a string map of float64 values into a string
// map of float64 pointers
func Float64Map(src map[string]float64) map[string]*float64 ***REMOVED***
	dst := make(map[string]*float64)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// Float64ValueMap converts a string map of float64 pointers into a string
// map of float64 values
func Float64ValueMap(src map[string]*float64) map[string]float64 ***REMOVED***
	dst := make(map[string]float64)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Time returns a pointer to the time.Time value passed in.
func Time(v time.Time) *time.Time ***REMOVED***
	return &v
***REMOVED***

// TimeValue returns the value of the time.Time pointer passed in or
// time.Time***REMOVED******REMOVED*** if the pointer is nil.
func TimeValue(v *time.Time) time.Time ***REMOVED***
	if v != nil ***REMOVED***
		return *v
	***REMOVED***
	return time.Time***REMOVED******REMOVED***
***REMOVED***

// SecondsTimeValue converts an int64 pointer to a time.Time value
// representing seconds since Epoch or time.Time***REMOVED******REMOVED*** if the pointer is nil.
func SecondsTimeValue(v *int64) time.Time ***REMOVED***
	if v != nil ***REMOVED***
		return time.Unix((*v / 1000), 0)
	***REMOVED***
	return time.Time***REMOVED******REMOVED***
***REMOVED***

// MillisecondsTimeValue converts an int64 pointer to a time.Time value
// representing milliseconds sinch Epoch or time.Time***REMOVED******REMOVED*** if the pointer is nil.
func MillisecondsTimeValue(v *int64) time.Time ***REMOVED***
	if v != nil ***REMOVED***
		return time.Unix(0, (*v * 1000000))
	***REMOVED***
	return time.Time***REMOVED******REMOVED***
***REMOVED***

// TimeUnixMilli returns a Unix timestamp in milliseconds from "January 1, 1970 UTC".
// The result is undefined if the Unix time cannot be represented by an int64.
// Which includes calling TimeUnixMilli on a zero Time is undefined.
//
// This utility is useful for service API's such as CloudWatch Logs which require
// their unix time values to be in milliseconds.
//
// See Go stdlib https://golang.org/pkg/time/#Time.UnixNano for more information.
func TimeUnixMilli(t time.Time) int64 ***REMOVED***
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
***REMOVED***

// TimeSlice converts a slice of time.Time values into a slice of
// time.Time pointers
func TimeSlice(src []time.Time) []*time.Time ***REMOVED***
	dst := make([]*time.Time, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		dst[i] = &(src[i])
	***REMOVED***
	return dst
***REMOVED***

// TimeValueSlice converts a slice of time.Time pointers into a slice of
// time.Time values
func TimeValueSlice(src []*time.Time) []time.Time ***REMOVED***
	dst := make([]time.Time, len(src))
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] != nil ***REMOVED***
			dst[i] = *(src[i])
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// TimeMap converts a string map of time.Time values into a string
// map of time.Time pointers
func TimeMap(src map[string]time.Time) map[string]*time.Time ***REMOVED***
	dst := make(map[string]*time.Time)
	for k, val := range src ***REMOVED***
		v := val
		dst[k] = &v
	***REMOVED***
	return dst
***REMOVED***

// TimeValueMap converts a string map of time.Time pointers into a string
// map of time.Time values
func TimeValueMap(src map[string]*time.Time) map[string]time.Time ***REMOVED***
	dst := make(map[string]time.Time)
	for k, val := range src ***REMOVED***
		if val != nil ***REMOVED***
			dst[k] = *val
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***
