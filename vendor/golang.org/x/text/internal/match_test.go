// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"strings"
	"testing"

	"golang.org/x/text/language"
)

func TestInheritanceMatcher(t *testing.T) ***REMOVED***
	for i, tt := range []struct ***REMOVED***
		haveTags string
		wantTags string
		match    string
		conf     language.Confidence
	***REMOVED******REMOVED***
		***REMOVED***"und,en,en-US", "en-US", "en-US", language.Exact***REMOVED***, // most specific match
		***REMOVED***"zh-Hant,zh", "zh-TW", "zh-Hant", language.High***REMOVED***,  // zh-TW implies Hant.
		***REMOVED***"und,zh", "zh-TW", "und", language.High***REMOVED***,          // zh-TW does not match zh.
		***REMOVED***"zh", "zh-TW", "und", language.No***REMOVED***,                // zh-TW does not match zh.
		***REMOVED***"iw,en,nl", "he", "he", language.Exact***REMOVED***,           // matches after canonicalization
		***REMOVED***"he,en,nl", "iw", "he", language.Exact***REMOVED***,           // matches after canonicalization
		// Prefer first match over more specific match for various reasons:
		// a) consistency of user interface is more important than an exact match,
		// b) _if_ und is specified, it should be considered a correct and useful match,
		// Note that a call to this Match will almost always be with a single tag.
		***REMOVED***"und,en,en-US", "he,en-US", "und", language.High***REMOVED***,
	***REMOVED*** ***REMOVED***
		have := parseTags(tt.haveTags)
		m := NewInheritanceMatcher(have)
		tag, index, conf := m.Match(parseTags(tt.wantTags)...)
		want := language.Raw.Make(tt.match)
		if tag != want ***REMOVED***
			t.Errorf("%d:tag: got %q; want %q", i, tag, want)
		***REMOVED***
		if conf != language.No ***REMOVED***
			if got, _ := language.All.Canonicalize(have[index]); got != want ***REMOVED***
				t.Errorf("%d:index: got %q; want %q ", i, got, want)
			***REMOVED***
		***REMOVED***
		if conf != tt.conf ***REMOVED***
			t.Errorf("%d:conf: got %v; want %v", i, conf, tt.conf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseTags(list string) (out []language.Tag) ***REMOVED***
	for _, s := range strings.Split(list, ",") ***REMOVED***
		out = append(out, language.Raw.Make(strings.TrimSpace(s)))
	***REMOVED***
	return out
***REMOVED***
