// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlindex

import (
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/language"
)

func TestGet(t *testing.T) ***REMOVED***
	for i, tc := range []struct ***REMOVED***
		name      string
		canonical string
		err       error
	***REMOVED******REMOVED***
		***REMOVED***"utf-8", "utf-8", nil***REMOVED***,
		***REMOVED***"  utf-8  ", "utf-8", nil***REMOVED***,
		***REMOVED***"  l5  ", "windows-1254", nil***REMOVED***,
		***REMOVED***"latin5 ", "windows-1254", nil***REMOVED***,
		***REMOVED***"latin 5", "", errInvalidName***REMOVED***,
		***REMOVED***"latin-5", "", errInvalidName***REMOVED***,
	***REMOVED*** ***REMOVED***
		enc, err := Get(tc.name)
		if err != tc.err ***REMOVED***
			t.Errorf("%d: error was %v; want %v", i, err, tc.err)
		***REMOVED***
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if got, err := Name(enc); got != tc.canonical ***REMOVED***
			t.Errorf("%d: Name(Get(%q)) = %q; want %q (%v)", i, tc.name, got, tc.canonical, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTables(t *testing.T) ***REMOVED***
	for name, index := range nameMap ***REMOVED***
		got, err := Get(name)
		if err != nil ***REMOVED***
			t.Errorf("%s:err: expected non-nil error", name)
		***REMOVED***
		if want := encodings[index]; got != want ***REMOVED***
			t.Errorf("%s:encoding: got %v; want %v", name, got, want)
		***REMOVED***
		mib, _ := got.(identifier.Interface).ID()
		if mibMap[mib] != index ***REMOVED***
			t.Errorf("%s:mibMab: got %d; want %d", name, mibMap[mib], index)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestName(t *testing.T) ***REMOVED***
	for i, tc := range []struct ***REMOVED***
		desc string
		enc  encoding.Encoding
		name string
		err  error
	***REMOVED******REMOVED******REMOVED***
		"defined encoding",
		charmap.ISO8859_2,
		"iso-8859-2",
		nil,
	***REMOVED***, ***REMOVED***
		"defined Unicode encoding",
		unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
		"utf-16be",
		nil,
	***REMOVED***, ***REMOVED***
		"undefined Unicode encoding in HTML standard",
		unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
		"",
		errUnsupported,
	***REMOVED***, ***REMOVED***
		"undefined other encoding in HTML standard",
		charmap.CodePage437,
		"",
		errUnsupported,
	***REMOVED***, ***REMOVED***
		"unknown encoding",
		encoding.Nop,
		"",
		errUnknown,
	***REMOVED******REMOVED*** ***REMOVED***
		name, err := Name(tc.enc)
		if name != tc.name || err != tc.err ***REMOVED***
			t.Errorf("%d:%s: got %q, %v; want %q, %v", i, tc.desc, name, err, tc.name, tc.err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLanguageDefault(t *testing.T) ***REMOVED***
	for _, tc := range []struct***REMOVED*** tag, want string ***REMOVED******REMOVED***
		***REMOVED***"und", "windows-1252"***REMOVED***, // The default value.
		***REMOVED***"ar", "windows-1256"***REMOVED***,
		***REMOVED***"ba", "windows-1251"***REMOVED***,
		***REMOVED***"be", "windows-1251"***REMOVED***,
		***REMOVED***"bg", "windows-1251"***REMOVED***,
		***REMOVED***"cs", "windows-1250"***REMOVED***,
		***REMOVED***"el", "iso-8859-7"***REMOVED***,
		***REMOVED***"et", "windows-1257"***REMOVED***,
		***REMOVED***"fa", "windows-1256"***REMOVED***,
		***REMOVED***"he", "windows-1255"***REMOVED***,
		***REMOVED***"hr", "windows-1250"***REMOVED***,
		***REMOVED***"hu", "iso-8859-2"***REMOVED***,
		***REMOVED***"ja", "shift_jis"***REMOVED***,
		***REMOVED***"kk", "windows-1251"***REMOVED***,
		***REMOVED***"ko", "euc-kr"***REMOVED***,
		***REMOVED***"ku", "windows-1254"***REMOVED***,
		***REMOVED***"ky", "windows-1251"***REMOVED***,
		***REMOVED***"lt", "windows-1257"***REMOVED***,
		***REMOVED***"lv", "windows-1257"***REMOVED***,
		***REMOVED***"mk", "windows-1251"***REMOVED***,
		***REMOVED***"pl", "iso-8859-2"***REMOVED***,
		***REMOVED***"ru", "windows-1251"***REMOVED***,
		***REMOVED***"sah", "windows-1251"***REMOVED***,
		***REMOVED***"sk", "windows-1250"***REMOVED***,
		***REMOVED***"sl", "iso-8859-2"***REMOVED***,
		***REMOVED***"sr", "windows-1251"***REMOVED***,
		***REMOVED***"tg", "windows-1251"***REMOVED***,
		***REMOVED***"th", "windows-874"***REMOVED***,
		***REMOVED***"tr", "windows-1254"***REMOVED***,
		***REMOVED***"tt", "windows-1251"***REMOVED***,
		***REMOVED***"uk", "windows-1251"***REMOVED***,
		***REMOVED***"vi", "windows-1258"***REMOVED***,
		***REMOVED***"zh-hans", "gb18030"***REMOVED***,
		***REMOVED***"zh-hant", "big5"***REMOVED***,
		// Variants and close approximates of the above.
		***REMOVED***"ar_EG", "windows-1256"***REMOVED***,
		***REMOVED***"bs", "windows-1250"***REMOVED***, // Bosnian Latin maps to Croatian.
		// Use default fallback in case of miss.
		***REMOVED***"nl", "windows-1252"***REMOVED***,
	***REMOVED*** ***REMOVED***
		if got := LanguageDefault(language.MustParse(tc.tag)); got != tc.want ***REMOVED***
			t.Errorf("LanguageDefault(%s) = %s; want %s", tc.tag, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***
