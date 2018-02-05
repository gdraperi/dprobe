// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ianaindex

import (
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

var All = [][]encoding.Encoding***REMOVED***
	unicode.All,
	charmap.All,
	japanese.All,
	korean.All,
	simplifiedchinese.All,
	traditionalchinese.All,
***REMOVED***

// TestAllIANA tests whether an Encoding supported in x/text is defined by IANA but
// not supported by this package.
func TestAllIANA(t *testing.T) ***REMOVED***
	for _, ea := range All ***REMOVED***
		for _, e := range ea ***REMOVED***
			mib, _ := e.(identifier.Interface).ID()
			if x := findMIB(ianaToMIB, mib); x != -1 && encodings[x] == nil ***REMOVED***
				t.Errorf("supported MIB %v (%v) not in index", mib, e)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestNotSupported reports the encodings in IANA, but not by x/text.
func TestNotSupported(t *testing.T) ***REMOVED***
	mibs := map[identifier.MIB]bool***REMOVED******REMOVED***
	for _, ea := range All ***REMOVED***
		for _, e := range ea ***REMOVED***
			mib, _ := e.(identifier.Interface).ID()
			mibs[mib] = true
		***REMOVED***
	***REMOVED***

	// Many encodings in the IANA index will likely not be suppored by the
	// Go encodings. That is fine.
	// TODO: consider wheter we should add this test.
	// for code, mib := range ianaToMIB ***REMOVED***
	// 	t.Run(fmt.Sprint("IANA:", mib), func(t *testing.T) ***REMOVED***
	// 		if !mibs[mib] ***REMOVED***
	// 			t.Skipf("IANA encoding %s (MIB %v) not supported",
	// 				ianaNames[code], mib)
	// 		***REMOVED***
	// 	***REMOVED***)
	// ***REMOVED***
***REMOVED***

func TestEncoding(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		index     *Index
		name      string
		canonical string
		err       error
	***REMOVED******REMOVED***
		***REMOVED***MIME, "utf-8", "UTF-8", nil***REMOVED***,
		***REMOVED***MIME, "  utf-8  ", "UTF-8", nil***REMOVED***,
		***REMOVED***MIME, "  l5  ", "ISO-8859-9", nil***REMOVED***,
		***REMOVED***MIME, "latin5 ", "ISO-8859-9", nil***REMOVED***,
		***REMOVED***MIME, "LATIN5 ", "ISO-8859-9", nil***REMOVED***,
		***REMOVED***MIME, "latin 5", "", errInvalidName***REMOVED***,
		***REMOVED***MIME, "latin-5", "", errInvalidName***REMOVED***,

		***REMOVED***IANA, "utf-8", "UTF-8", nil***REMOVED***,
		***REMOVED***IANA, "  utf-8  ", "UTF-8", nil***REMOVED***,
		***REMOVED***IANA, "  l5  ", "ISO_8859-9:1989", nil***REMOVED***,
		***REMOVED***IANA, "latin5 ", "ISO_8859-9:1989", nil***REMOVED***,
		***REMOVED***IANA, "LATIN5 ", "ISO_8859-9:1989", nil***REMOVED***,
		***REMOVED***IANA, "latin 5", "", errInvalidName***REMOVED***,
		***REMOVED***IANA, "latin-5", "", errInvalidName***REMOVED***,

		***REMOVED***MIB, "utf-8", "UTF8", nil***REMOVED***,
		***REMOVED***MIB, "  utf-8  ", "UTF8", nil***REMOVED***,
		***REMOVED***MIB, "  l5  ", "ISOLatin5", nil***REMOVED***,
		***REMOVED***MIB, "latin5 ", "ISOLatin5", nil***REMOVED***,
		***REMOVED***MIB, "LATIN5 ", "ISOLatin5", nil***REMOVED***,
		***REMOVED***MIB, "latin 5", "", errInvalidName***REMOVED***,
		***REMOVED***MIB, "latin-5", "", errInvalidName***REMOVED***,
	***REMOVED***
	for i, tc := range testCases ***REMOVED***
		enc, err := tc.index.Encoding(tc.name)
		if err != tc.err ***REMOVED***
			t.Errorf("%d: error was %v; want %v", i, err, tc.err)
		***REMOVED***
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if got, err := tc.index.Name(enc); got != tc.canonical ***REMOVED***
			t.Errorf("%d: Name(Encoding(%q)) = %q; want %q (%v)", i, tc.name, got, tc.canonical, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTables(t *testing.T) ***REMOVED***
	for i, x := range []*Index***REMOVED***MIME, IANA***REMOVED*** ***REMOVED***
		for name, index := range x.alias ***REMOVED***
			got, err := x.Encoding(name)
			if err != nil ***REMOVED***
				t.Errorf("%d%s:err: unexpected error %v", i, name, err)
			***REMOVED***
			if want := x.enc[index]; got != want ***REMOVED***
				t.Errorf("%d%s:encoding: got %v; want %v", i, name, got, want)
			***REMOVED***
			if got != nil ***REMOVED***
				mib, _ := got.(identifier.Interface).ID()
				if i := findMIB(x.toMIB, mib); i != index ***REMOVED***
					t.Errorf("%d%s:mib: got %d; want %d", i, name, i, index)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type unsupported struct ***REMOVED***
	encoding.Encoding
***REMOVED***

func (unsupported) ID() (identifier.MIB, string) ***REMOVED*** return 9999, "" ***REMOVED***

func TestName(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		desc string
		enc  encoding.Encoding
		f    func(e encoding.Encoding) (string, error)
		name string
		err  error
	***REMOVED******REMOVED******REMOVED***
		"defined encoding",
		charmap.ISO8859_2,
		MIME.Name,
		"ISO-8859-2",
		nil,
	***REMOVED***, ***REMOVED***
		"defined Unicode encoding",
		unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
		IANA.Name,
		"UTF-16BE",
		nil,
	***REMOVED***, ***REMOVED***
		"another defined Unicode encoding",
		unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
		MIME.Name,
		"UTF-16",
		nil,
	***REMOVED***, ***REMOVED***
		"unknown Unicode encoding",
		unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM),
		MIME.Name,
		"",
		errUnknown,
	***REMOVED***, ***REMOVED***
		"undefined encoding",
		unsupported***REMOVED******REMOVED***,
		MIME.Name,
		"",
		errUnsupported,
	***REMOVED***, ***REMOVED***
		"undefined other encoding in HTML standard",
		charmap.CodePage437,
		IANA.Name,
		"IBM437",
		nil,
	***REMOVED***, ***REMOVED***
		"unknown encoding",
		encoding.Nop,
		IANA.Name,
		"",
		errUnknown,
	***REMOVED******REMOVED***
	for i, tc := range testCases ***REMOVED***
		name, err := tc.f(tc.enc)
		if name != tc.name || err != tc.err ***REMOVED***
			t.Errorf("%d:%s: got %q, %v; want %q, %v", i, tc.desc, name, err, tc.name, tc.err)
		***REMOVED***
	***REMOVED***
***REMOVED***
