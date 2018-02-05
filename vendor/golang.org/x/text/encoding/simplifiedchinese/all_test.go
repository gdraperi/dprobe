// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simplifiedchinese

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
	// Pick n large enough to overflow the destination buffer of transform.String.
	const n = 10000
	testCases := []struct ***REMOVED***
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	***REMOVED******REMOVED***
		***REMOVED***dec, GBK, "a\xfe\xfeb", "a\ufffdb"***REMOVED***,
		***REMOVED***dec, HZGB2312, "~***REMOVED***z~", "\ufffd"***REMOVED***,

		***REMOVED***enc, GBK, "갂", ""***REMOVED***,
		***REMOVED***enc, GBK, "a갂", "a"***REMOVED***,
		***REMOVED***enc, GBK, "\u4e02갂", "\x81@"***REMOVED***,

		***REMOVED***enc, HZGB2312, "갂", ""***REMOVED***,
		***REMOVED***enc, HZGB2312, "a갂", "a"***REMOVED***,
		***REMOVED***enc, HZGB2312, "\u6cf5갂", "~***REMOVED***1C~***REMOVED***"***REMOVED***,

		***REMOVED***dec, GB18030, "\x80", "€"***REMOVED***,
		***REMOVED***dec, GB18030, "\x81", "\ufffd"***REMOVED***,
		***REMOVED***dec, GB18030, "\x81\x20", "\ufffd "***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\xfe", "\ufffd"***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\xff", "\ufffd\ufffd"***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\x30", "\ufffd0"***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\x30\x30 ", "\ufffd00 "***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\x30\xff ", "\ufffd0\ufffd "***REMOVED***,
		***REMOVED***dec, GB18030, "\xfe\x30\x81\x21", "\ufffd0\ufffd!"***REMOVED***,

		***REMOVED***dec, GB18030, strings.Repeat("\xfe\x30", n), strings.Repeat("\ufffd0", n)***REMOVED***,

		***REMOVED***dec, HZGB2312, "~/", "\ufffd"***REMOVED***,
		***REMOVED***dec, HZGB2312, "~***REMOVED***a\x80", "\ufffd"***REMOVED***,
		***REMOVED***dec, HZGB2312, "~***REMOVED***a\x80", "\ufffd"***REMOVED***,
		***REMOVED***dec, HZGB2312, "~***REMOVED***" + strings.Repeat("z~", n), strings.Repeat("\ufffd", n)***REMOVED***,
		***REMOVED***dec, HZGB2312, "~***REMOVED***" + strings.Repeat("\xfe\x30", n), strings.Repeat("\ufffd", n*2)***REMOVED***,
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
		e         encoding.Encoding
		encPrefix string
		encoded   string
		utf8      string
	***REMOVED******REMOVED******REMOVED***
		// "\u0081\u00de\u00df\u00e0\u00e1\u00e2\u00e3\uffff\U00010000" is a
		// nonsense string that contains GB18030 encodable codepoints of which
		// only U+00E0 and U+00E1 are GBK encodable.
		//
		// "A\u3000\u554a\u4e02\u4e90\u72dc\u7349\u02ca\u2588Z€" is a nonsense
		// string that contains ASCII and GBK encodable codepoints from Levels
		// 1-5 as well as the Euro sign.
		//
		// "A\u43f0\u4c32\U00027267\u3000\U0002910d\u79d4Z€" is a nonsense string
		// that contains ASCII and Big5 encodable codepoints from the Basic
		// Multilingual Plane and the Supplementary Ideographic Plane as well as
		// the Euro sign.
		//
		// "花间一壶酒，独酌无相亲。" (simplified) and
		// "花間一壺酒，獨酌無相親。" (traditional)
		// are from the 8th century poem "Yuè Xià Dú Zhuó".
		e: GB18030,
		encoded: "\x81\x30\x81\x31\x81\x30\x89\x37\x81\x30\x89\x38\xa8\xa4\xa8\xa2" +
			"\x81\x30\x89\x39\x81\x30\x8a\x30\x84\x31\xa4\x39\x90\x30\x81\x30",
		utf8: "\u0081\u00de\u00df\u00e0\u00e1\u00e2\u00e3\uffff\U00010000",
	***REMOVED***, ***REMOVED***
		e: GB18030,
		encoded: "\xbb\xa8\xbc\xe4\xd2\xbb\xba\xf8\xbe\xc6\xa3\xac\xb6\xc0\xd7\xc3" +
			"\xce\xde\xcf\xe0\xc7\xd7\xa1\xa3",
		utf8: "花间一壶酒，独酌无相亲。",
	***REMOVED***, ***REMOVED***
		e:       GBK,
		encoded: "A\xa1\xa1\xb0\xa1\x81\x40\x81\x80\xaa\x40\xaa\x80\xa8\x40\xa8\x80Z\x80",
		utf8:    "A\u3000\u554a\u4e02\u4e90\u72dc\u7349\u02ca\u2588Z€",
	***REMOVED***, ***REMOVED***
		e: GBK,
		encoded: "\xbb\xa8\xbc\xe4\xd2\xbb\xba\xf8\xbe\xc6\xa3\xac\xb6\xc0\xd7\xc3" +
			"\xce\xde\xcf\xe0\xc7\xd7\xa1\xa3",
		utf8: "花间一壶酒，独酌无相亲。",
	***REMOVED***, ***REMOVED***
		e:       HZGB2312,
		encoded: "A~***REMOVED***\x21\x21~~\x30\x21~***REMOVED***Z~~",
		utf8:    "A\u3000~\u554aZ~",
	***REMOVED***, ***REMOVED***
		e:         HZGB2312,
		encPrefix: "~***REMOVED***",
		encoded:   ";(<dR;:x>F#,6@WCN^O`GW!#",
		utf8:      "花间一壶酒，独酌无相亲。",
	***REMOVED******REMOVED***

	for _, tc := range testCases ***REMOVED***
		enctest.TestEncoding(t, tc.e, tc.encoded, tc.utf8, tc.encPrefix, "")
	***REMOVED***
***REMOVED***

func TestFiles(t *testing.T) ***REMOVED***
	enctest.TestFile(t, GB18030)
	enctest.TestFile(t, GBK)
	enctest.TestFile(t, HZGB2312)
***REMOVED***

func BenchmarkEncoding(b *testing.B) ***REMOVED***
	enctest.Benchmark(b, GB18030)
	enctest.Benchmark(b, GBK)
	enctest.Benchmark(b, HZGB2312)
***REMOVED***
