// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package assert provides helper functions for testing.
package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// skip defines the default call depth
const skip = 2

// Equal asserts that got and want are equal as defined by
// reflect.DeepEqual. The test fails with msg if they are not equal.
func Equal(t *testing.T, got, want interface***REMOVED******REMOVED***, msg ...string) ***REMOVED***
	if x := equal(2, got, want, msg...); x != "" ***REMOVED***
		fmt.Println(x)
		t.Fail()
	***REMOVED***
***REMOVED***

func equal(skip int, got, want interface***REMOVED******REMOVED***, msg ...string) string ***REMOVED***
	if !reflect.DeepEqual(got, want) ***REMOVED***
		return fail(skip, "got %v want %v %s", got, want, strings.Join(msg, " "))
	***REMOVED***
	return ""
***REMOVED***

// Panic asserts that function fn() panics.
// It assumes that recover() either returns a string or
// an error and fails if the message does not match
// the regular expression in 'matches'.
func Panic(t *testing.T, fn func(), matches string) ***REMOVED***
	if x := doesPanic(2, fn, matches); x != "" ***REMOVED***
		fmt.Println(x)
		t.Fail()
	***REMOVED***
***REMOVED***

func doesPanic(skip int, fn func(), expr string) (err string) ***REMOVED***
	defer func() ***REMOVED***
		r := recover()
		if r == nil ***REMOVED***
			err = fail(skip, "did not panic")
			return
		***REMOVED***
		var v string
		switch r.(type) ***REMOVED***
		case error:
			v = r.(error).Error()
		case string:
			v = r.(string)
		***REMOVED***
		err = matches(skip, v, expr)
	***REMOVED***()
	fn()
	return ""
***REMOVED***

// Matches asserts that a value matches a given regular expression.
func Matches(t *testing.T, value, expr string) ***REMOVED***
	if x := matches(2, value, expr); x != "" ***REMOVED***
		fmt.Println(x)
		t.Fail()
	***REMOVED***
***REMOVED***

func matches(skip int, value, expr string) string ***REMOVED***
	ok, err := regexp.MatchString(expr, value)
	if err != nil ***REMOVED***
		return fail(skip, "invalid pattern %q. %s", expr, err)
	***REMOVED***
	if !ok ***REMOVED***
		return fail(skip, "got %s which does not match %s", value, expr)
	***REMOVED***
	return ""
***REMOVED***

func fail(skip int, format string, args ...interface***REMOVED******REMOVED***) string ***REMOVED***
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("\t%s:%d: %s\n", filepath.Base(file), line, fmt.Sprintf(format, args...))
***REMOVED***
