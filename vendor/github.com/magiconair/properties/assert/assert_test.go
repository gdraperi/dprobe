// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assert

import "testing"

func TestEqualEquals(t *testing.T) ***REMOVED***
	if got, want := equal(2, "a", "a"), ""; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestEqualFails(t *testing.T) ***REMOVED***
	if got, want := equal(2, "a", "b"), "\tassert_test.go:16: got a want b \n"; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestPanicPanics(t *testing.T) ***REMOVED***
	if got, want := doesPanic(2, func() ***REMOVED*** panic("foo") ***REMOVED***, ""), ""; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestPanicPanicsAndMatches(t *testing.T) ***REMOVED***
	if got, want := doesPanic(2, func() ***REMOVED*** panic("foo") ***REMOVED***, "foo"), ""; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestPanicPanicsAndDoesNotMatch(t *testing.T) ***REMOVED***
	if got, want := doesPanic(2, func() ***REMOVED*** panic("foo") ***REMOVED***, "bar"), "\tassert.go:62: got foo which does not match bar\n"; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestPanicPanicsAndDoesNotPanic(t *testing.T) ***REMOVED***
	if got, want := doesPanic(2, func() ***REMOVED******REMOVED***, "bar"), "\tassert.go:65: did not panic\n"; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestMatchesMatches(t *testing.T) ***REMOVED***
	if got, want := matches(2, "aaa", "a"), ""; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***

func TestMatchesDoesNotMatch(t *testing.T) ***REMOVED***
	if got, want := matches(2, "aaa", "b"), "\tassert_test.go:52: got aaa which does not match b\n"; got != want ***REMOVED***
		t.Fatalf("got %q want %q", got, want)
	***REMOVED***
***REMOVED***
