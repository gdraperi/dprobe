// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package korean

import (
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/enctest"
	"golang.org/x/text/transform"
)

func dec(e encoding.Encoding) (dir string, t transform.Transformer, err error) ***REMOVED***
	return "Decode", e.NewDecoder(), nil
***REMOVED***
func enc(e encoding.Encoding) (dir string, t transform.Transformer, err error) ***REMOVED***
	return "Encode", e.NewEncoder(), internal.ErrASCIIReplacement
***REMOVED***

func TestNonRepertoire(t *testing.T) ***REMOVED***
	// Pick n large enough to cause an overflow in the destination buffer of
	// transform.String.
	const n = 10000
	testCases := []struct ***REMOVED***
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	***REMOVED******REMOVED***
		***REMOVED***dec, EUCKR, "\xfe\xfe", "\ufffd"***REMOVED***,
		// ***REMOVED***dec, EUCKR, "א", "\ufffd"***REMOVED***, // TODO: why is this different?

		***REMOVED***enc, EUCKR, "א", ""***REMOVED***,
		***REMOVED***enc, EUCKR, "aא", "a"***REMOVED***,
		***REMOVED***enc, EUCKR, "\uac00א", "\xb0\xa1"***REMOVED***,
		// TODO: should we also handle Jamo?

		***REMOVED***dec, EUCKR, "\x80", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCKR, "\xff", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCKR, "\x81", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCKR, "\xb0\x40", "\ufffd@"***REMOVED***,
		***REMOVED***dec, EUCKR, "\xb0\xff", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCKR, "\xd0\x20", "\ufffd "***REMOVED***,
		***REMOVED***dec, EUCKR, "\xd0\xff", "\ufffd"***REMOVED***,

		***REMOVED***dec, EUCKR, strings.Repeat("\x81", n), strings.Repeat("걖", n/2)***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		dir, tr, wantErr := tc.init(tc.e)

		dst, _, err := transform.String(tr, tc.src)
		if err != wantErr ***REMOVED***
			t.Errorf("%s %v(%q): got %v; want %v", dir, tc.e, tc.src, err, wantErr)
		***REMOVED***
		if got := string(dst); got != tc.want ***REMOVED***
			t.Errorf("%s %v(%q):\ngot  %q\nwant %q", dir, tc.e, tc.src, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBasics(t *testing.T) ***REMOVED***
	// The encoded forms can be verified by the iconv program:
	// $ echo 月日は百代 | iconv -f UTF-8 -t SHIFT-JIS | xxd
	testCases := []struct ***REMOVED***
		e       encoding.Encoding
		encoded string
		utf8    string
	***REMOVED******REMOVED******REMOVED***
		// Korean tests.
		//
		// "A\uac02\uac35\uac56\ud401B\ud408\ud620\ud624C\u4f3d\u8a70D" is a
		// nonsense string that contains ASCII, Hangul and CJK ideographs.
		//
		// "세계야, 안녕" translates as "Hello, world".
		e:       EUCKR,
		encoded: "A\x81\x41\x81\x61\x81\x81\xc6\xfeB\xc7\xa1\xc7\xfe\xc8\xa1C\xca\xa1\xfd\xfeD",
		utf8:    "A\uac02\uac35\uac56\ud401B\ud408\ud620\ud624C\u4f3d\u8a70D",
	***REMOVED***, ***REMOVED***
		e:       EUCKR,
		encoded: "\xbc\xbc\xb0\xe8\xbe\xdf\x2c\x20\xbe\xc8\xb3\xe7",
		utf8:    "세계야, 안녕",
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		enctest.TestEncoding(t, tc.e, tc.encoded, tc.utf8, "", "")
	***REMOVED***
***REMOVED***

func TestFiles(t *testing.T) ***REMOVED*** enctest.TestFile(t, EUCKR) ***REMOVED***

func BenchmarkEncoding(b *testing.B) ***REMOVED*** enctest.Benchmark(b, EUCKR) ***REMOVED***
