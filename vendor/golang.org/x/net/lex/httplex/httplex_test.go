// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httplex

import (
	"testing"
)

func isChar(c rune) bool ***REMOVED*** return c <= 127 ***REMOVED***

func isCtl(c rune) bool ***REMOVED*** return c <= 31 || c == 127 ***REMOVED***

func isSeparator(c rune) bool ***REMOVED***
	switch c ***REMOVED***
	case '(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?', '=', '***REMOVED***', '***REMOVED***', ' ', '\t':
		return true
	***REMOVED***
	return false
***REMOVED***

func TestIsToken(t *testing.T) ***REMOVED***
	for i := 0; i <= 130; i++ ***REMOVED***
		r := rune(i)
		expected := isChar(r) && !isCtl(r) && !isSeparator(r)
		if IsTokenRune(r) != expected ***REMOVED***
			t.Errorf("isToken(0x%x) = %v", r, !expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHeaderValuesContainsToken(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		vals  []string
		token string
		want  bool
	***REMOVED******REMOVED***
		***REMOVED***
			vals:  []string***REMOVED***"foo"***REMOVED***,
			token: "foo",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"bar", "foo"***REMOVED***,
			token: "foo",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"foo"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"foo"***REMOVED***,
			token: "bar",
			want:  false,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***" foo "***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"foo,bar"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"bar,foo,bar"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"bar , foo"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"foo ,bar "***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"bar, foo ,bar"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
		***REMOVED***
			vals:  []string***REMOVED***"bar , foo"***REMOVED***,
			token: "FOO",
			want:  true,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		got := HeaderValuesContainsToken(tt.vals, tt.token)
		if got != tt.want ***REMOVED***
			t.Errorf("headerValuesContainsToken(%q, %q) = %v; want %v", tt.vals, tt.token, got, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPunycodeHostPort(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in, want string
	***REMOVED******REMOVED***
		***REMOVED***"www.google.com", "www.google.com"***REMOVED***,
		***REMOVED***"гофер.рф", "xn--c1ae0ajs.xn--p1ai"***REMOVED***,
		***REMOVED***"bücher.de", "xn--bcher-kva.de"***REMOVED***,
		***REMOVED***"bücher.de:8080", "xn--bcher-kva.de:8080"***REMOVED***,
		***REMOVED***"[1::6]:8080", "[1::6]:8080"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		got, err := PunycodeHostPort(tt.in)
		if tt.want != got || err != nil ***REMOVED***
			t.Errorf("PunycodeHostPort(%q) = %q, %v, want %q, nil", tt.in, got, err, tt.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
