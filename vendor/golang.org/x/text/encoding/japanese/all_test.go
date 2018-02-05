// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package japanese

import (
	"fmt"
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
	// Pick n to cause the destination buffer in transform.String to overflow.
	const n = 100
	long := strings.Repeat(".", n)
	testCases := []struct ***REMOVED***
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	***REMOVED******REMOVED***
		***REMOVED***enc, EUCJP, "갂", ""***REMOVED***,
		***REMOVED***enc, EUCJP, "a갂", "a"***REMOVED***,
		***REMOVED***enc, EUCJP, "丌갂", "\x8f\xb0\xa4"***REMOVED***,

		***REMOVED***enc, ISO2022JP, "갂", ""***REMOVED***,
		***REMOVED***enc, ISO2022JP, "a갂", "a"***REMOVED***,
		***REMOVED***enc, ISO2022JP, "朗갂", "\x1b$BzF\x1b(B"***REMOVED***, // switch back to ASCII mode at end

		***REMOVED***enc, ShiftJIS, "갂", ""***REMOVED***,
		***REMOVED***enc, ShiftJIS, "a갂", "a"***REMOVED***,
		***REMOVED***enc, ShiftJIS, "\u2190갂", "\x81\xa9"***REMOVED***,

		// Continue correctly after errors
		***REMOVED***dec, EUCJP, "\x8e\xa0", "\ufffd\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8e\xe0", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8e\xff", "\ufffd\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8ea", "\ufffda"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa0", "\ufffd\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa1\xa0", "\ufffd\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa1a", "\ufffda"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa1a", "\ufffda"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa1a", "\ufffda"***REMOVED***,
		***REMOVED***dec, EUCJP, "\x8f\xa2\xa2", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\xfe", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\xfe\xfc", "\ufffd"***REMOVED***,
		***REMOVED***dec, EUCJP, "\xfe\xff", "\ufffd\ufffd"***REMOVED***,
		// Correct handling of end of source
		***REMOVED***dec, EUCJP, strings.Repeat("\x8e", n), strings.Repeat("\ufffd", n)***REMOVED***,
		***REMOVED***dec, EUCJP, strings.Repeat("\x8f", n), strings.Repeat("\ufffd", n)***REMOVED***,
		***REMOVED***dec, EUCJP, strings.Repeat("\x8f\xa0", n), strings.Repeat("\ufffd", 2*n)***REMOVED***,
		***REMOVED***dec, EUCJP, "a" + strings.Repeat("\x8f\xa1", n), "a" + strings.Repeat("\ufffd", n)***REMOVED***,
		***REMOVED***dec, EUCJP, "a" + strings.Repeat("\x8f\xa1\xff", n), "a" + strings.Repeat("\ufffd", 2*n)***REMOVED***,

		// Continue correctly after errors
		***REMOVED***dec, ShiftJIS, "\x80", "\u0080"***REMOVED***, // It's what the spec says.
		***REMOVED***dec, ShiftJIS, "\x81", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\x81\x7f", "\ufffd\u007f"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xe0", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xe0\x39", "\ufffd\u0039"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xe0\x9f", "燹"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xe0\xfd", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xef\xfc", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xfc\xfc", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xfc\xfd", "\ufffd"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xfdaa", "\ufffdaa"***REMOVED***,

		***REMOVED***dec, ShiftJIS, strings.Repeat("\x81\x81", n), strings.Repeat("＝", n)***REMOVED***,
		***REMOVED***dec, ShiftJIS, strings.Repeat("\xe0\xfd", n), strings.Repeat("\ufffd", n)***REMOVED***,
		***REMOVED***dec, ShiftJIS, "a" + strings.Repeat("\xe0\xfd", n), "a" + strings.Repeat("\ufffd", n)***REMOVED***,

		***REMOVED***dec, ISO2022JP, "\x1b$", "\ufffd$"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(", "\ufffd("***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b@", "\ufffd@"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1bZ", "\ufffdZ"***REMOVED***,
		// incomplete escapes
		***REMOVED***dec, ISO2022JP, "\x1b$", "\ufffd$"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$J.", "\ufffd$J."***REMOVED***,             // illegal
		***REMOVED***dec, ISO2022JP, "\x1b$B.", "\ufffd"***REMOVED***,                // JIS208
		***REMOVED***dec, ISO2022JP, "\x1b$(", "\ufffd$("***REMOVED***,               // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b$(..", "\ufffd$(.."***REMOVED***,           // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b$(" + long, "\ufffd$(" + long***REMOVED***, // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b$(D.", "\ufffd"***REMOVED***,               // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b$(D..", "\ufffd"***REMOVED***,              // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b$(D...", "\ufffd\ufffd"***REMOVED***,       // JIS212
		***REMOVED***dec, ISO2022JP, "\x1b(B.", "."***REMOVED***,                     // ascii
		***REMOVED***dec, ISO2022JP, "\x1b(B..", ".."***REMOVED***,                   // ascii
		***REMOVED***dec, ISO2022JP, "\x1b(J.", "."***REMOVED***,                     // roman
		***REMOVED***dec, ISO2022JP, "\x1b(J..", ".."***REMOVED***,                   // roman
		***REMOVED***dec, ISO2022JP, "\x1b(I\x20", "\ufffd"***REMOVED***,             // katakana
		***REMOVED***dec, ISO2022JP, "\x1b(I\x20\x20", "\ufffd\ufffd"***REMOVED***,   // katakana
		// recover to same state
		***REMOVED***dec, ISO2022JP, "\x1b(B\x1b.", "\ufffd."***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(I\x1b.", "\ufffdｮ"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(I\x1b$.", "\ufffd､ｮ"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(I\x1b(.", "\ufffdｨｮ"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$B\x7e\x7e", "\ufffd"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$@\x0a.", "\x0a."***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$B\x0a.", "\x0a."***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$(D\x0a.", "\x0a."***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b$(D\x7e\x7e", "\ufffd"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x80", "\ufffd"***REMOVED***,

		// TODO: according to https://encoding.spec.whatwg.org/#iso-2022-jp,
		// these should all be correct.
		// ***REMOVED***dec, ISO2022JP, "\x1b(B\x0E", "\ufffd"***REMOVED***,
		// ***REMOVED***dec, ISO2022JP, "\x1b(B\x0F", "\ufffd"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(B\x5C", "\u005C"***REMOVED***,
		***REMOVED***dec, ISO2022JP, "\x1b(B\x7E", "\u007E"***REMOVED***,
		// ***REMOVED***dec, ISO2022JP, "\x1b(J\x0E", "\ufffd"***REMOVED***,
		// ***REMOVED***dec, ISO2022JP, "\x1b(J\x0F", "\ufffd"***REMOVED***,
		// ***REMOVED***dec, ISO2022JP, "\x1b(J\x5C", "\u00A5"***REMOVED***,
		// ***REMOVED***dec, ISO2022JP, "\x1b(J\x7E", "\u203E"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		dir, tr, wantErr := tc.init(tc.e)
		t.Run(fmt.Sprintf("%s/%v/%q", dir, tc.e, tc.src), func(t *testing.T) ***REMOVED***
			dst := make([]byte, 100000)
			src := []byte(tc.src)
			for i := 0; i <= len(tc.src); i++ ***REMOVED***
				nDst, nSrc, err := tr.Transform(dst, src[:i], false)
				if err != nil && err != transform.ErrShortSrc && err != wantErr ***REMOVED***
					t.Fatalf("error on first call to Transform: %v", err)
				***REMOVED***
				n, _, err := tr.Transform(dst[nDst:], src[nSrc:], true)
				nDst += n
				if err != wantErr ***REMOVED***
					t.Fatalf("(%q|%q): got %v; want %v", tc.src[:i], tc.src[i:], err, wantErr)
				***REMOVED***
				if got := string(dst[:nDst]); got != tc.want ***REMOVED***
					t.Errorf("(%q|%q):\ngot  %q\nwant %q", tc.src[:i], tc.src[i:], got, tc.want)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestCorrect(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	***REMOVED******REMOVED***
		***REMOVED***dec, ShiftJIS, "\x9f\xfc", "滌"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xfb\xfc", "髙"***REMOVED***,
		***REMOVED***dec, ShiftJIS, "\xfa\xb1", "﨑"***REMOVED***,
		***REMOVED***enc, ShiftJIS, "滌", "\x9f\xfc"***REMOVED***,
		***REMOVED***enc, ShiftJIS, "﨑", "\xed\x95"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		dir, tr, _ := tc.init(tc.e)

		dst, _, err := transform.String(tr, tc.src)
		if err != nil ***REMOVED***
			t.Errorf("%s %v(%q): got %v; want %v", dir, tc.e, tc.src, err, nil)
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
		e         encoding.Encoding
		encPrefix string
		encSuffix string
		encoded   string
		utf8      string
	***REMOVED******REMOVED******REMOVED***
		// "A｡ｶﾟ 0208: etc 0212: etc" is a nonsense string that contains ASCII, half-width
		// kana, JIS X 0208 (including two near the kink in the Shift JIS second byte
		// encoding) and JIS X 0212 encodable codepoints.
		//
		// "月日は百代の過客にして、行かふ年も又旅人也。" is from the 17th century poem
		// "Oku no Hosomichi" and contains both hiragana and kanji.
		e: EUCJP,
		encoded: "A\x8e\xa1\x8e\xb6\x8e\xdf " +
			"0208: \xa1\xa1\xa1\xa2\xa1\xdf\xa1\xe0\xa1\xfd\xa1\xfe\xa2\xa1\xa2\xa2\xf4\xa6 " +
			"0212: \x8f\xa2\xaf\x8f\xed\xe3",
		utf8: "A｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199 " +
			"0212: \u02d8\u9fa5",
	***REMOVED***, ***REMOVED***
		e: EUCJP,
		encoded: "\xb7\xee\xc6\xfc\xa4\xcf\xc9\xb4\xc2\xe5\xa4\xce\xb2\xe1\xb5\xd2" +
			"\xa4\xcb\xa4\xb7\xa4\xc6\xa1\xa2\xb9\xd4\xa4\xab\xa4\xd5\xc7\xaf" +
			"\xa4\xe2\xcb\xf4\xce\xb9\xbf\xcd\xcc\xe9\xa1\xa3",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	***REMOVED***, ***REMOVED***
		e:         ISO2022JP,
		encSuffix: "\x1b\x28\x42",
		encoded: "\x1b\x28\x49\x21\x36\x5f\x1b\x28\x42 " +
			"0208: \x1b\x24\x42\x21\x21\x21\x22\x21\x5f\x21\x60\x21\x7d\x21\x7e\x22\x21\x22\x22\x74\x26",
		utf8: "｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199",
	***REMOVED***, ***REMOVED***
		e:         ISO2022JP,
		encPrefix: "\x1b\x24\x42",
		encSuffix: "\x1b\x28\x42",
		encoded: "\x37\x6e\x46\x7c\x24\x4f\x49\x34\x42\x65\x24\x4e\x32\x61\x35\x52" +
			"\x24\x4b\x24\x37\x24\x46\x21\x22\x39\x54\x24\x2b\x24\x55\x47\x2f" +
			"\x24\x62\x4b\x74\x4e\x39\x3f\x4d\x4c\x69\x21\x23",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	***REMOVED***, ***REMOVED***
		e: ShiftJIS,
		encoded: "A\xa1\xb6\xdf " +
			"0208: \x81\x40\x81\x41\x81\x7e\x81\x80\x81\x9d\x81\x9e\x81\x9f\x81\xa0\xea\xa4",
		utf8: "A｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199",
	***REMOVED***, ***REMOVED***
		e: ShiftJIS,
		encoded: "\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
			"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
			"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		enctest.TestEncoding(t, tc.e, tc.encoded, tc.utf8, tc.encPrefix, tc.encSuffix)
	***REMOVED***
***REMOVED***

func TestFiles(t *testing.T) ***REMOVED***
	enctest.TestFile(t, EUCJP)
	enctest.TestFile(t, ISO2022JP)
	enctest.TestFile(t, ShiftJIS)
***REMOVED***

func BenchmarkEncoding(b *testing.B) ***REMOVED***
	enctest.Benchmark(b, EUCJP)
	enctest.Benchmark(b, ISO2022JP)
	enctest.Benchmark(b, ShiftJIS)
***REMOVED***
