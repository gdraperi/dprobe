// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/language"
)

type (
	renamedBool       bool
	renamedInt        int
	renamedInt8       int8
	renamedInt16      int16
	renamedInt32      int32
	renamedInt64      int64
	renamedUint       uint
	renamedUint8      uint8
	renamedUint16     uint16
	renamedUint32     uint32
	renamedUint64     uint64
	renamedUintptr    uintptr
	renamedString     string
	renamedBytes      []byte
	renamedFloat32    float32
	renamedFloat64    float64
	renamedComplex64  complex64
	renamedComplex128 complex128
)

func TestFmtInterface(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	var i1 interface***REMOVED******REMOVED***
	i1 = "abc"
	s := p.Sprintf("%s", i1)
	if s != "abc" ***REMOVED***
		t.Errorf(`Sprintf("%%s", empty("abc")) = %q want %q`, s, "abc")
	***REMOVED***
***REMOVED***

var (
	NaN    = math.NaN()
	posInf = math.Inf(1)
	negInf = math.Inf(-1)

	intVar = 0

	array  = [5]int***REMOVED***1, 2, 3, 4, 5***REMOVED***
	iarray = [4]interface***REMOVED******REMOVED******REMOVED***1, "hello", 2.5, nil***REMOVED***
	slice  = array[:]
	islice = iarray[:]
)

type A struct ***REMOVED***
	i int
	j uint
	s string
	x []int
***REMOVED***

type I int

func (i I) String() string ***REMOVED***
	p := NewPrinter(language.Und)
	return p.Sprintf("<%d>", int(i))
***REMOVED***

type B struct ***REMOVED***
	I I
	j int
***REMOVED***

type C struct ***REMOVED***
	i int
	B
***REMOVED***

type F int

func (f F) Format(s fmt.State, c rune) ***REMOVED***
	p := NewPrinter(language.Und)
	p.Fprintf(s, "<%c=F(%d)>", c, int(f))
***REMOVED***

type G int

func (g G) GoString() string ***REMOVED***
	p := NewPrinter(language.Und)
	return p.Sprintf("GoString(%d)", int(g))
***REMOVED***

type S struct ***REMOVED***
	F F // a struct field that Formats
	G G // a struct field that GoStrings
***REMOVED***

type SI struct ***REMOVED***
	I interface***REMOVED******REMOVED***
***REMOVED***

// P is a type with a String method with pointer receiver for testing %p.
type P int

var pValue P

func (p *P) String() string ***REMOVED***
	return "String(p)"
***REMOVED***

var barray = [5]renamedUint8***REMOVED***1, 2, 3, 4, 5***REMOVED***
var bslice = barray[:]

type byteStringer byte

func (byteStringer) String() string ***REMOVED***
	return "X"
***REMOVED***

var byteStringerSlice = []byteStringer***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***

type byteFormatter byte

func (byteFormatter) Format(f fmt.State, _ rune) ***REMOVED***
	p := NewPrinter(language.Und)
	p.Fprint(f, "X")
***REMOVED***

var byteFormatterSlice = []byteFormatter***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***

var fmtTests = []struct ***REMOVED***
	fmt string
	val interface***REMOVED******REMOVED***
	out string
***REMOVED******REMOVED***
	// The behavior of the following tests differs from that of the fmt package.

	// Unlike with the fmt package, it is okay to have extra arguments for
	// strings without format parameters. This is because it is impossible to
	// distinguish between reordered or ordered format strings in this case.
	// (For reordered format strings it is okay to not use arguments.)
	***REMOVED***"", nil, ""***REMOVED***,
	***REMOVED***"", 2, ""***REMOVED***,
	***REMOVED***"no args", "hello", "no args"***REMOVED***,

	***REMOVED***"%017091901790959340919092959340919017929593813360", 0, "%!(NOVERB)"***REMOVED***,
	***REMOVED***"%184467440737095516170v", 0, "%!(NOVERB)"***REMOVED***,
	// Extra argument errors should format without flags set.
	***REMOVED***"%010.2", "12345", "%!(NOVERB)"***REMOVED***,

	// Some key other differences, asides from localized values:
	// - NaN values should not use affixes; so no signs (CLDR requirement)
	// - Infinity uses patterns, so signs may be different (CLDR requirement)
	// - The # flag is used to disable localization.

	// All following tests are analogous to those of the fmt package, but with
	// localized numbers when appropriate.
	***REMOVED***"%d", 12345, "12,345"***REMOVED***,
	***REMOVED***"%v", 12345, "12,345"***REMOVED***,
	***REMOVED***"%t", true, "true"***REMOVED***,

	// basic string
	***REMOVED***"%s", "abc", "abc"***REMOVED***,
	***REMOVED***"%q", "abc", `"abc"`***REMOVED***,
	***REMOVED***"%x", "abc", "616263"***REMOVED***,
	***REMOVED***"%x", "\xff\xf0\x0f\xff", "fff00fff"***REMOVED***,
	***REMOVED***"%X", "\xff\xf0\x0f\xff", "FFF00FFF"***REMOVED***,
	***REMOVED***"%x", "", ""***REMOVED***,
	***REMOVED***"% x", "", ""***REMOVED***,
	***REMOVED***"%#x", "", ""***REMOVED***,
	***REMOVED***"%# x", "", ""***REMOVED***,
	***REMOVED***"%x", "xyz", "78797a"***REMOVED***,
	***REMOVED***"%X", "xyz", "78797A"***REMOVED***,
	***REMOVED***"% x", "xyz", "78 79 7a"***REMOVED***,
	***REMOVED***"% X", "xyz", "78 79 7A"***REMOVED***,
	***REMOVED***"%#x", "xyz", "0x78797a"***REMOVED***,
	***REMOVED***"%#X", "xyz", "0X78797A"***REMOVED***,
	***REMOVED***"%# x", "xyz", "0x78 0x79 0x7a"***REMOVED***,
	***REMOVED***"%# X", "xyz", "0X78 0X79 0X7A"***REMOVED***,

	// basic bytes
	***REMOVED***"%s", []byte("abc"), "abc"***REMOVED***,
	***REMOVED***"%s", [3]byte***REMOVED***'a', 'b', 'c'***REMOVED***, "abc"***REMOVED***,
	***REMOVED***"%s", &[3]byte***REMOVED***'a', 'b', 'c'***REMOVED***, "&abc"***REMOVED***,
	***REMOVED***"%q", []byte("abc"), `"abc"`***REMOVED***,
	***REMOVED***"%x", []byte("abc"), "616263"***REMOVED***,
	***REMOVED***"%x", []byte("\xff\xf0\x0f\xff"), "fff00fff"***REMOVED***,
	***REMOVED***"%X", []byte("\xff\xf0\x0f\xff"), "FFF00FFF"***REMOVED***,
	***REMOVED***"%x", []byte(""), ""***REMOVED***,
	***REMOVED***"% x", []byte(""), ""***REMOVED***,
	***REMOVED***"%#x", []byte(""), ""***REMOVED***,
	***REMOVED***"%# x", []byte(""), ""***REMOVED***,
	***REMOVED***"%x", []byte("xyz"), "78797a"***REMOVED***,
	***REMOVED***"%X", []byte("xyz"), "78797A"***REMOVED***,
	***REMOVED***"% x", []byte("xyz"), "78 79 7a"***REMOVED***,
	***REMOVED***"% X", []byte("xyz"), "78 79 7A"***REMOVED***,
	***REMOVED***"%#x", []byte("xyz"), "0x78797a"***REMOVED***,
	***REMOVED***"%#X", []byte("xyz"), "0X78797A"***REMOVED***,
	***REMOVED***"%# x", []byte("xyz"), "0x78 0x79 0x7a"***REMOVED***,
	***REMOVED***"%# X", []byte("xyz"), "0X78 0X79 0X7A"***REMOVED***,

	// escaped strings
	***REMOVED***"%q", "", `""`***REMOVED***,
	***REMOVED***"%#q", "", "``"***REMOVED***,
	***REMOVED***"%q", "\"", `"\""`***REMOVED***,
	***REMOVED***"%#q", "\"", "`\"`"***REMOVED***,
	***REMOVED***"%q", "`", `"` + "`" + `"`***REMOVED***,
	***REMOVED***"%#q", "`", `"` + "`" + `"`***REMOVED***,
	***REMOVED***"%q", "\n", `"\n"`***REMOVED***,
	***REMOVED***"%#q", "\n", `"\n"`***REMOVED***,
	***REMOVED***"%q", `\n`, `"\\n"`***REMOVED***,
	***REMOVED***"%#q", `\n`, "`\\n`"***REMOVED***,
	***REMOVED***"%q", "abc", `"abc"`***REMOVED***,
	***REMOVED***"%#q", "abc", "`abc`"***REMOVED***,
	***REMOVED***"%q", "日本語", `"日本語"`***REMOVED***,
	***REMOVED***"%+q", "日本語", `"\u65e5\u672c\u8a9e"`***REMOVED***,
	***REMOVED***"%#q", "日本語", "`日本語`"***REMOVED***,
	***REMOVED***"%#+q", "日本語", "`日本語`"***REMOVED***,
	***REMOVED***"%q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`***REMOVED***,
	***REMOVED***"%+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`***REMOVED***,
	***REMOVED***"%#q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`***REMOVED***,
	***REMOVED***"%#+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`***REMOVED***,
	***REMOVED***"%q", "☺", `"☺"`***REMOVED***,
	***REMOVED***"% q", "☺", `"☺"`***REMOVED***, // The space modifier should have no effect.
	***REMOVED***"%+q", "☺", `"\u263a"`***REMOVED***,
	***REMOVED***"%#q", "☺", "`☺`"***REMOVED***,
	***REMOVED***"%#+q", "☺", "`☺`"***REMOVED***,
	***REMOVED***"%10q", "⌘", `       "⌘"`***REMOVED***,
	***REMOVED***"%+10q", "⌘", `  "\u2318"`***REMOVED***,
	***REMOVED***"%-10q", "⌘", `"⌘"       `***REMOVED***,
	***REMOVED***"%+-10q", "⌘", `"\u2318"  `***REMOVED***,
	***REMOVED***"%010q", "⌘", `0000000"⌘"`***REMOVED***,
	***REMOVED***"%+010q", "⌘", `00"\u2318"`***REMOVED***,
	***REMOVED***"%-010q", "⌘", `"⌘"       `***REMOVED***, // 0 has no effect when - is present.
	***REMOVED***"%+-010q", "⌘", `"\u2318"  `***REMOVED***,
	***REMOVED***"%#8q", "\n", `    "\n"`***REMOVED***,
	***REMOVED***"%#+8q", "\r", `    "\r"`***REMOVED***,
	***REMOVED***"%#-8q", "\t", "`	`     "***REMOVED***,
	***REMOVED***"%#+-8q", "\b", `"\b"    `***REMOVED***,
	***REMOVED***"%q", "abc\xffdef", `"abc\xffdef"`***REMOVED***,
	***REMOVED***"%+q", "abc\xffdef", `"abc\xffdef"`***REMOVED***,
	***REMOVED***"%#q", "abc\xffdef", `"abc\xffdef"`***REMOVED***,
	***REMOVED***"%#+q", "abc\xffdef", `"abc\xffdef"`***REMOVED***,
	// Runes that are not printable.
	***REMOVED***"%q", "\U0010ffff", `"\U0010ffff"`***REMOVED***,
	***REMOVED***"%+q", "\U0010ffff", `"\U0010ffff"`***REMOVED***,
	***REMOVED***"%#q", "\U0010ffff", "`􏿿`"***REMOVED***,
	***REMOVED***"%#+q", "\U0010ffff", "`􏿿`"***REMOVED***,
	// Runes that are not valid.
	***REMOVED***"%q", string(0x110000), `"�"`***REMOVED***,
	***REMOVED***"%+q", string(0x110000), `"\ufffd"`***REMOVED***,
	***REMOVED***"%#q", string(0x110000), "`�`"***REMOVED***,
	***REMOVED***"%#+q", string(0x110000), "`�`"***REMOVED***,

	// characters
	***REMOVED***"%c", uint('x'), "x"***REMOVED***,
	***REMOVED***"%c", 0xe4, "ä"***REMOVED***,
	***REMOVED***"%c", 0x672c, "本"***REMOVED***,
	***REMOVED***"%c", '日', "日"***REMOVED***,
	***REMOVED***"%.0c", '⌘', "⌘"***REMOVED***, // Specifying precision should have no effect.
	***REMOVED***"%3c", '⌘', "  ⌘"***REMOVED***,
	***REMOVED***"%-3c", '⌘', "⌘  "***REMOVED***,
	// Runes that are not printable.
	***REMOVED***"%c", '\U00000e00', "\u0e00"***REMOVED***,
	***REMOVED***"%c", '\U0010ffff', "\U0010ffff"***REMOVED***,
	// Runes that are not valid.
	***REMOVED***"%c", -1, "�"***REMOVED***,
	***REMOVED***"%c", 0xDC80, "�"***REMOVED***,
	***REMOVED***"%c", rune(0x110000), "�"***REMOVED***,
	***REMOVED***"%c", int64(0xFFFFFFFFF), "�"***REMOVED***,
	***REMOVED***"%c", uint64(0xFFFFFFFFF), "�"***REMOVED***,

	// escaped characters
	***REMOVED***"%q", uint(0), `'\x00'`***REMOVED***,
	***REMOVED***"%+q", uint(0), `'\x00'`***REMOVED***,
	***REMOVED***"%q", '"', `'"'`***REMOVED***,
	***REMOVED***"%+q", '"', `'"'`***REMOVED***,
	***REMOVED***"%q", '\'', `'\''`***REMOVED***,
	***REMOVED***"%+q", '\'', `'\''`***REMOVED***,
	***REMOVED***"%q", '`', "'`'"***REMOVED***,
	***REMOVED***"%+q", '`', "'`'"***REMOVED***,
	***REMOVED***"%q", 'x', `'x'`***REMOVED***,
	***REMOVED***"%+q", 'x', `'x'`***REMOVED***,
	***REMOVED***"%q", 'ÿ', `'ÿ'`***REMOVED***,
	***REMOVED***"%+q", 'ÿ', `'\u00ff'`***REMOVED***,
	***REMOVED***"%q", '\n', `'\n'`***REMOVED***,
	***REMOVED***"%+q", '\n', `'\n'`***REMOVED***,
	***REMOVED***"%q", '☺', `'☺'`***REMOVED***,
	***REMOVED***"%+q", '☺', `'\u263a'`***REMOVED***,
	***REMOVED***"% q", '☺', `'☺'`***REMOVED***,  // The space modifier should have no effect.
	***REMOVED***"%.0q", '☺', `'☺'`***REMOVED***, // Specifying precision should have no effect.
	***REMOVED***"%10q", '⌘', `       '⌘'`***REMOVED***,
	***REMOVED***"%+10q", '⌘', `  '\u2318'`***REMOVED***,
	***REMOVED***"%-10q", '⌘', `'⌘'       `***REMOVED***,
	***REMOVED***"%+-10q", '⌘', `'\u2318'  `***REMOVED***,
	***REMOVED***"%010q", '⌘', `0000000'⌘'`***REMOVED***,
	***REMOVED***"%+010q", '⌘', `00'\u2318'`***REMOVED***,
	***REMOVED***"%-010q", '⌘', `'⌘'       `***REMOVED***, // 0 has no effect when - is present.
	***REMOVED***"%+-010q", '⌘', `'\u2318'  `***REMOVED***,
	// Runes that are not printable.
	***REMOVED***"%q", '\U00000e00', `'\u0e00'`***REMOVED***,
	***REMOVED***"%q", '\U0010ffff', `'\U0010ffff'`***REMOVED***,
	// Runes that are not valid.
	***REMOVED***"%q", int32(-1), "%!q(int32=-1)"***REMOVED***,
	***REMOVED***"%q", 0xDC80, `'�'`***REMOVED***,
	***REMOVED***"%q", rune(0x110000), "%!q(int32=1,114,112)"***REMOVED***,
	***REMOVED***"%q", int64(0xFFFFFFFFF), "%!q(int64=68,719,476,735)"***REMOVED***,
	***REMOVED***"%q", uint64(0xFFFFFFFFF), "%!q(uint64=68,719,476,735)"***REMOVED***,

	// width
	***REMOVED***"%5s", "abc", "  abc"***REMOVED***,
	***REMOVED***"%2s", "\u263a", " ☺"***REMOVED***,
	***REMOVED***"%-5s", "abc", "abc  "***REMOVED***,
	***REMOVED***"%-8q", "abc", `"abc"   `***REMOVED***,
	***REMOVED***"%05s", "abc", "00abc"***REMOVED***,
	***REMOVED***"%08q", "abc", `000"abc"`***REMOVED***,
	***REMOVED***"%5s", "abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxyz"***REMOVED***,
	***REMOVED***"%.5s", "abcdefghijklmnopqrstuvwxyz", "abcde"***REMOVED***,
	***REMOVED***"%.0s", "日本語日本語", ""***REMOVED***,
	***REMOVED***"%.5s", "日本語日本語", "日本語日本"***REMOVED***,
	***REMOVED***"%.10s", "日本語日本語", "日本語日本語"***REMOVED***,
	***REMOVED***"%.5s", []byte("日本語日本語"), "日本語日本"***REMOVED***,
	***REMOVED***"%.5q", "abcdefghijklmnopqrstuvwxyz", `"abcde"`***REMOVED***,
	***REMOVED***"%.5x", "abcdefghijklmnopqrstuvwxyz", "6162636465"***REMOVED***,
	***REMOVED***"%.5q", []byte("abcdefghijklmnopqrstuvwxyz"), `"abcde"`***REMOVED***,
	***REMOVED***"%.5x", []byte("abcdefghijklmnopqrstuvwxyz"), "6162636465"***REMOVED***,
	***REMOVED***"%.3q", "日本語日本語", `"日本語"`***REMOVED***,
	***REMOVED***"%.3q", []byte("日本語日本語"), `"日本語"`***REMOVED***,
	***REMOVED***"%.1q", "日本語", `"日"`***REMOVED***,
	***REMOVED***"%.1q", []byte("日本語"), `"日"`***REMOVED***,
	***REMOVED***"%.1x", "日本語", "e6"***REMOVED***,
	***REMOVED***"%.1X", []byte("日本語"), "E6"***REMOVED***,
	***REMOVED***"%10.1q", "日本語日本語", `       "日"`***REMOVED***,
	***REMOVED***"%10v", nil, "     <nil>"***REMOVED***,
	***REMOVED***"%-10v", nil, "<nil>     "***REMOVED***,

	// integers
	***REMOVED***"%d", uint(12345), "12,345"***REMOVED***,
	***REMOVED***"%d", int(-12345), "-12,345"***REMOVED***,
	***REMOVED***"%d", ^uint8(0), "255"***REMOVED***,
	***REMOVED***"%d", ^uint16(0), "65,535"***REMOVED***,
	***REMOVED***"%d", ^uint32(0), "4,294,967,295"***REMOVED***,
	***REMOVED***"%d", ^uint64(0), "18,446,744,073,709,551,615"***REMOVED***,
	***REMOVED***"%d", int8(-1 << 7), "-128"***REMOVED***,
	***REMOVED***"%d", int16(-1 << 15), "-32,768"***REMOVED***,
	***REMOVED***"%d", int32(-1 << 31), "-2,147,483,648"***REMOVED***,
	***REMOVED***"%d", int64(-1 << 63), "-9,223,372,036,854,775,808"***REMOVED***,
	***REMOVED***"%.d", 0, ""***REMOVED***,
	***REMOVED***"%.0d", 0, ""***REMOVED***,
	***REMOVED***"%6.0d", 0, "      "***REMOVED***,
	***REMOVED***"%06.0d", 0, "      "***REMOVED***,
	***REMOVED***"% d", 12345, " 12,345"***REMOVED***,
	***REMOVED***"%+d", 12345, "+12,345"***REMOVED***,
	***REMOVED***"%+d", -12345, "-12,345"***REMOVED***,
	***REMOVED***"%b", 7, "111"***REMOVED***,
	***REMOVED***"%b", -6, "-110"***REMOVED***,
	***REMOVED***"%b", ^uint32(0), "11111111111111111111111111111111"***REMOVED***,
	***REMOVED***"%b", ^uint64(0), "1111111111111111111111111111111111111111111111111111111111111111"***REMOVED***,
	***REMOVED***"%b", int64(-1 << 63), zeroFill("-1", 63, "")***REMOVED***,
	***REMOVED***"%o", 01234, "1234"***REMOVED***,
	***REMOVED***"%#o", 01234, "01234"***REMOVED***,
	***REMOVED***"%o", ^uint32(0), "37777777777"***REMOVED***,
	***REMOVED***"%o", ^uint64(0), "1777777777777777777777"***REMOVED***,
	***REMOVED***"%#X", 0, "0X0"***REMOVED***,
	***REMOVED***"%x", 0x12abcdef, "12abcdef"***REMOVED***,
	***REMOVED***"%X", 0x12abcdef, "12ABCDEF"***REMOVED***,
	***REMOVED***"%x", ^uint32(0), "ffffffff"***REMOVED***,
	***REMOVED***"%X", ^uint64(0), "FFFFFFFFFFFFFFFF"***REMOVED***,
	***REMOVED***"%.20b", 7, "00000000000000000111"***REMOVED***,
	***REMOVED***"%10d", 12345, "    12,345"***REMOVED***,
	***REMOVED***"%10d", -12345, "   -12,345"***REMOVED***,
	***REMOVED***"%+10d", 12345, "   +12,345"***REMOVED***,
	***REMOVED***"%010d", 12345, "0,000,012,345"***REMOVED***,
	***REMOVED***"%010d", -12345, "-0,000,012,345"***REMOVED***,
	***REMOVED***"%20.8d", 1234, "          00,001,234"***REMOVED***,
	***REMOVED***"%20.8d", -1234, "         -00,001,234"***REMOVED***,
	***REMOVED***"%020.8d", 1234, "          00,001,234"***REMOVED***,
	***REMOVED***"%020.8d", -1234, "         -00,001,234"***REMOVED***,
	***REMOVED***"%-20.8d", 1234, "00,001,234          "***REMOVED***,
	***REMOVED***"%-20.8d", -1234, "-00,001,234         "***REMOVED***,
	***REMOVED***"%-#20.8x", 0x1234abc, "0x01234abc          "***REMOVED***,
	***REMOVED***"%-#20.8X", 0x1234abc, "0X01234ABC          "***REMOVED***,
	***REMOVED***"%-#20.8o", 01234, "00001234            "***REMOVED***,

	// Test correct f.intbuf overflow checks.
	***REMOVED***"%068d", 1, "00," + strings.Repeat("000,", 21) + "001"***REMOVED***,
	***REMOVED***"%068d", -1, "-00," + strings.Repeat("000,", 21) + "001"***REMOVED***,
	***REMOVED***"%#.68x", 42, zeroFill("0x", 68, "2a")***REMOVED***,
	***REMOVED***"%.68d", -42, "-00," + strings.Repeat("000,", 21) + "042"***REMOVED***,
	***REMOVED***"%+.68d", 42, "+00," + strings.Repeat("000,", 21) + "042"***REMOVED***,
	***REMOVED***"% .68d", 42, " 00," + strings.Repeat("000,", 21) + "042"***REMOVED***,
	***REMOVED***"% +.68d", 42, "+00," + strings.Repeat("000,", 21) + "042"***REMOVED***,

	// unicode format
	***REMOVED***"%U", 0, "U+0000"***REMOVED***,
	***REMOVED***"%U", -1, "U+FFFFFFFFFFFFFFFF"***REMOVED***,
	***REMOVED***"%U", '\n', `U+000A`***REMOVED***,
	***REMOVED***"%#U", '\n', `U+000A`***REMOVED***,
	***REMOVED***"%+U", 'x', `U+0078`***REMOVED***,       // Plus flag should have no effect.
	***REMOVED***"%# U", 'x', `U+0078 'x'`***REMOVED***,  // Space flag should have no effect.
	***REMOVED***"%#.2U", 'x', `U+0078 'x'`***REMOVED***, // Precisions below 4 should print 4 digits.
	***REMOVED***"%U", '\u263a', `U+263A`***REMOVED***,
	***REMOVED***"%#U", '\u263a', `U+263A '☺'`***REMOVED***,
	***REMOVED***"%U", '\U0001D6C2', `U+1D6C2`***REMOVED***,
	***REMOVED***"%#U", '\U0001D6C2', `U+1D6C2 '𝛂'`***REMOVED***,
	***REMOVED***"%#14.6U", '⌘', "  U+002318 '⌘'"***REMOVED***,
	***REMOVED***"%#-14.6U", '⌘', "U+002318 '⌘'  "***REMOVED***,
	***REMOVED***"%#014.6U", '⌘', "  U+002318 '⌘'"***REMOVED***,
	***REMOVED***"%#-014.6U", '⌘', "U+002318 '⌘'  "***REMOVED***,
	***REMOVED***"%.68U", uint(42), zeroFill("U+", 68, "2A")***REMOVED***,
	***REMOVED***"%#.68U", '日', zeroFill("U+", 68, "65E5") + " '日'"***REMOVED***,

	// floats
	***REMOVED***"%+.3e", 0.0, "+0.000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%+.3e", 1.0, "+1.000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%+.3f", -1.0, "-1.000"***REMOVED***,
	***REMOVED***"%+.3F", -1.0, "-1.000"***REMOVED***,
	***REMOVED***"%+.3F", float32(-1.0), "-1.000"***REMOVED***,
	***REMOVED***"%+07.2f", 1.0, "+001.00"***REMOVED***,
	***REMOVED***"%+07.2f", -1.0, "-001.00"***REMOVED***,
	***REMOVED***"%-07.2f", 1.0, "1.00   "***REMOVED***,
	***REMOVED***"%-07.2f", -1.0, "-1.00  "***REMOVED***,
	***REMOVED***"%+-07.2f", 1.0, "+1.00  "***REMOVED***,
	***REMOVED***"%+-07.2f", -1.0, "-1.00  "***REMOVED***,
	***REMOVED***"%-+07.2f", 1.0, "+1.00  "***REMOVED***,
	***REMOVED***"%-+07.2f", -1.0, "-1.00  "***REMOVED***,
	***REMOVED***"%+10.2f", +1.0, "     +1.00"***REMOVED***,
	***REMOVED***"%+10.2f", -1.0, "     -1.00"***REMOVED***,
	***REMOVED***"% .3E", -1.0, "-1.000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"% .3e", 1.0, " 1.000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%+.3g", 0.0, "+0"***REMOVED***,
	***REMOVED***"%+.3g", 1.0, "+1"***REMOVED***,
	***REMOVED***"%+.3g", -1.0, "-1"***REMOVED***,
	***REMOVED***"% .3g", -1.0, "-1"***REMOVED***,
	***REMOVED***"% .3g", 1.0, " 1"***REMOVED***,
	***REMOVED***"%b", float32(1.0), "8388608p-23"***REMOVED***,
	***REMOVED***"%b", 1.0, "4503599627370496p-52"***REMOVED***,
	// Test sharp flag used with floats.
	***REMOVED***"%#g", 1e-323, "1.00000e-323"***REMOVED***,
	***REMOVED***"%#g", -1.0, "-1.00000"***REMOVED***,
	***REMOVED***"%#g", 1.1, "1.10000"***REMOVED***,
	***REMOVED***"%#g", 123456.0, "123456."***REMOVED***,
	***REMOVED***"%#g", 1234567.0, "1.234567e+06"***REMOVED***,
	***REMOVED***"%#g", 1230000.0, "1.23000e+06"***REMOVED***,
	***REMOVED***"%#g", 1000000.0, "1.00000e+06"***REMOVED***,
	***REMOVED***"%#.0f", 1.0, "1."***REMOVED***,
	***REMOVED***"%#.0e", 1.0, "1.e+00"***REMOVED***,
	***REMOVED***"%#.0g", 1.0, "1."***REMOVED***,
	***REMOVED***"%#.0g", 1100000.0, "1.e+06"***REMOVED***,
	***REMOVED***"%#.4f", 1.0, "1.0000"***REMOVED***,
	***REMOVED***"%#.4e", 1.0, "1.0000e+00"***REMOVED***,
	***REMOVED***"%#.4g", 1.0, "1.000"***REMOVED***,
	***REMOVED***"%#.4g", 100000.0, "1.000e+05"***REMOVED***,
	***REMOVED***"%#.0f", 123.0, "123."***REMOVED***,
	***REMOVED***"%#.0e", 123.0, "1.e+02"***REMOVED***,
	***REMOVED***"%#.0g", 123.0, "1.e+02"***REMOVED***,
	***REMOVED***"%#.4f", 123.0, "123.0000"***REMOVED***,
	***REMOVED***"%#.4e", 123.0, "1.2300e+02"***REMOVED***,
	***REMOVED***"%#.4g", 123.0, "123.0"***REMOVED***,
	***REMOVED***"%#.4g", 123000.0, "1.230e+05"***REMOVED***,
	***REMOVED***"%#9.4g", 1.0, "    1.000"***REMOVED***,
	// The sharp flag has no effect for binary float format.
	***REMOVED***"%#b", 1.0, "4503599627370496p-52"***REMOVED***,
	// Precision has no effect for binary float format.
	***REMOVED***"%.4b", float32(1.0), "8388608p-23"***REMOVED***,
	***REMOVED***"%.4b", -1.0, "-4503599627370496p-52"***REMOVED***,
	// Test correct f.intbuf boundary checks.
	***REMOVED***"%.68f", 1.0, zeroFill("1.", 68, "")***REMOVED***,
	***REMOVED***"%.68f", -1.0, zeroFill("-1.", 68, "")***REMOVED***,
	// float infinites and NaNs
	***REMOVED***"%f", posInf, "∞"***REMOVED***,
	***REMOVED***"%.1f", negInf, "-∞"***REMOVED***,
	***REMOVED***"% f", NaN, "NaN"***REMOVED***,
	***REMOVED***"%20f", posInf, "                   ∞"***REMOVED***,
	***REMOVED***"% 20F", posInf, "                   ∞"***REMOVED***,
	***REMOVED***"% 20e", negInf, "                  -∞"***REMOVED***,
	***REMOVED***"%+20E", negInf, "                  -∞"***REMOVED***,
	***REMOVED***"% +20g", negInf, "                  -∞"***REMOVED***,
	***REMOVED***"%+-20G", posInf, "+∞                  "***REMOVED***,
	***REMOVED***"%20e", NaN, "                 NaN"***REMOVED***,
	***REMOVED***"% +20E", NaN, "                 NaN"***REMOVED***,
	***REMOVED***"% -20g", NaN, "NaN                 "***REMOVED***,
	***REMOVED***"%+-20G", NaN, "NaN                 "***REMOVED***,
	// Zero padding does not apply to infinities and NaN.
	***REMOVED***"%+020e", posInf, "                  +∞"***REMOVED***,
	***REMOVED***"%-020f", negInf, "-∞                  "***REMOVED***,
	***REMOVED***"%-020E", NaN, "NaN                 "***REMOVED***,

	// complex values
	***REMOVED***"%.f", 0i, "(0+0i)"***REMOVED***,
	***REMOVED***"% .f", 0i, "( 0+0i)"***REMOVED***,
	***REMOVED***"%+.f", 0i, "(+0+0i)"***REMOVED***,
	***REMOVED***"% +.f", 0i, "(+0+0i)"***REMOVED***,
	***REMOVED***"%+.3e", 0i, "(+0.000\u202f×\u202f10⁰⁰+0.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%+.3f", 0i, "(+0.000+0.000i)"***REMOVED***,
	***REMOVED***"%+.3g", 0i, "(+0+0i)"***REMOVED***,
	***REMOVED***"%+.3e", 1 + 2i, "(+1.000\u202f×\u202f10⁰⁰+2.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%+.3f", 1 + 2i, "(+1.000+2.000i)"***REMOVED***,
	***REMOVED***"%+.3g", 1 + 2i, "(+1+2i)"***REMOVED***,
	***REMOVED***"%.3e", 0i, "(0.000\u202f×\u202f10⁰⁰+0.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%.3f", 0i, "(0.000+0.000i)"***REMOVED***,
	***REMOVED***"%.3F", 0i, "(0.000+0.000i)"***REMOVED***,
	***REMOVED***"%.3F", complex64(0i), "(0.000+0.000i)"***REMOVED***,
	***REMOVED***"%.3g", 0i, "(0+0i)"***REMOVED***,
	***REMOVED***"%.3e", 1 + 2i, "(1.000\u202f×\u202f10⁰⁰+2.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%.3f", 1 + 2i, "(1.000+2.000i)"***REMOVED***,
	***REMOVED***"%.3g", 1 + 2i, "(1+2i)"***REMOVED***,
	***REMOVED***"%.3e", -1 - 2i, "(-1.000\u202f×\u202f10⁰⁰-2.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%.3f", -1 - 2i, "(-1.000-2.000i)"***REMOVED***,
	***REMOVED***"%.3g", -1 - 2i, "(-1-2i)"***REMOVED***,
	***REMOVED***"% .3E", -1 - 2i, "(-1.000\u202f×\u202f10⁰⁰-2.000\u202f×\u202f10⁰⁰i)"***REMOVED***,
	***REMOVED***"%+.3g", 1 + 2i, "(+1+2i)"***REMOVED***,
	***REMOVED***"%+.3g", complex64(1 + 2i), "(+1+2i)"***REMOVED***,
	***REMOVED***"%#g", 1 + 2i, "(1.00000+2.00000i)"***REMOVED***,
	***REMOVED***"%#g", 123456 + 789012i, "(123456.+789012.i)"***REMOVED***,
	***REMOVED***"%#g", 1e-10i, "(0.00000+1.00000e-10i)"***REMOVED***,
	***REMOVED***"%#g", -1e10 - 1.11e100i, "(-1.00000e+10-1.11000e+100i)"***REMOVED***,
	***REMOVED***"%#.0f", 1.23 + 1.0i, "(1.+1.i)"***REMOVED***,
	***REMOVED***"%#.0e", 1.23 + 1.0i, "(1.e+00+1.e+00i)"***REMOVED***,
	***REMOVED***"%#.0g", 1.23 + 1.0i, "(1.+1.i)"***REMOVED***,
	***REMOVED***"%#.0g", 0 + 100000i, "(0.+1.e+05i)"***REMOVED***,
	***REMOVED***"%#.0g", 1230000 + 0i, "(1.e+06+0.i)"***REMOVED***,
	***REMOVED***"%#.4f", 1 + 1.23i, "(1.0000+1.2300i)"***REMOVED***,
	***REMOVED***"%#.4e", 123 + 1i, "(1.2300e+02+1.0000e+00i)"***REMOVED***,
	***REMOVED***"%#.4g", 123 + 1.23i, "(123.0+1.230i)"***REMOVED***,
	***REMOVED***"%#12.5g", 0 + 100000i, "(      0.0000 +1.0000e+05i)"***REMOVED***,
	***REMOVED***"%#12.5g", 1230000 - 0i, "(  1.2300e+06     +0.0000i)"***REMOVED***,
	***REMOVED***"%b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"***REMOVED***,
	***REMOVED***"%b", complex64(1 + 2i), "(8388608p-23+8388608p-22i)"***REMOVED***,
	// The sharp flag has no effect for binary complex format.
	***REMOVED***"%#b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"***REMOVED***,
	// Precision has no effect for binary complex format.
	***REMOVED***"%.4b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"***REMOVED***,
	***REMOVED***"%.4b", complex64(1 + 2i), "(8388608p-23+8388608p-22i)"***REMOVED***,
	// complex infinites and NaNs
	***REMOVED***"%f", complex(posInf, posInf), "(∞+∞i)"***REMOVED***,
	***REMOVED***"%f", complex(negInf, negInf), "(-∞-∞i)"***REMOVED***,
	***REMOVED***"%f", complex(NaN, NaN), "(NaN+NaNi)"***REMOVED***,
	***REMOVED***"%.1f", complex(posInf, posInf), "(∞+∞i)"***REMOVED***,
	***REMOVED***"% f", complex(posInf, posInf), "( ∞+∞i)"***REMOVED***,
	***REMOVED***"% f", complex(negInf, negInf), "(-∞-∞i)"***REMOVED***,
	***REMOVED***"% f", complex(NaN, NaN), "(NaN+NaNi)"***REMOVED***,
	***REMOVED***"%8e", complex(posInf, posInf), "(       ∞      +∞i)"***REMOVED***,
	***REMOVED***"% 8E", complex(posInf, posInf), "(       ∞      +∞i)"***REMOVED***,
	***REMOVED***"%+8f", complex(negInf, negInf), "(      -∞      -∞i)"***REMOVED***,
	***REMOVED***"% +8g", complex(negInf, negInf), "(      -∞      -∞i)"***REMOVED***, // TODO(g)
	***REMOVED***"% -8G", complex(NaN, NaN), "(NaN     +NaN    i)"***REMOVED***,
	***REMOVED***"%+-8b", complex(NaN, NaN), "(+NaN    +NaN    i)"***REMOVED***,
	// Zero padding does not apply to infinities and NaN.
	***REMOVED***"%08f", complex(posInf, posInf), "(       ∞      +∞i)"***REMOVED***,
	***REMOVED***"%-08g", complex(negInf, negInf), "(-∞      -∞      i)"***REMOVED***,
	***REMOVED***"%-08G", complex(NaN, NaN), "(NaN     +NaN    i)"***REMOVED***,

	// old test/fmt_test.go
	***REMOVED***"%e", 1.0, "1.000000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%e", 1234.5678e3, "1.234568\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%e", 1234.5678e-8, "1.234568\u202f×\u202f10⁻⁰⁵"***REMOVED***,
	***REMOVED***"%e", -7.0, "-7.000000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%e", -1e-9, "-1.000000\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%f", 1234.5678e3, "1,234,567.800000"***REMOVED***,
	***REMOVED***"%f", 1234.5678e-8, "0.000012"***REMOVED***,
	***REMOVED***"%f", -7.0, "-7.000000"***REMOVED***,
	***REMOVED***"%f", -1e-9, "-0.000000"***REMOVED***,
	***REMOVED***"%g", 1234.5678e3, "1.2345678\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%g", float32(1234.5678e3), "1.2345678\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%g", 1234.5678e-8, "1.2345678\u202f×\u202f10⁻⁰⁵"***REMOVED***,
	***REMOVED***"%g", -7.0, "-7"***REMOVED***,
	***REMOVED***"%g", -1e-9, "-1\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%g", float32(-1e-9), "-1\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%E", 1.0, "1.000000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%E", 1234.5678e3, "1.234568\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%E", 1234.5678e-8, "1.234568\u202f×\u202f10⁻⁰⁵"***REMOVED***,
	***REMOVED***"%E", -7.0, "-7.000000\u202f×\u202f10⁰⁰"***REMOVED***,
	***REMOVED***"%E", -1e-9, "-1.000000\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%G", 1234.5678e3, "1.2345678\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%G", float32(1234.5678e3), "1.2345678\u202f×\u202f10⁰⁶"***REMOVED***,
	***REMOVED***"%G", 1234.5678e-8, "1.2345678\u202f×\u202f10⁻⁰⁵"***REMOVED***,
	***REMOVED***"%G", -7.0, "-7"***REMOVED***,
	***REMOVED***"%G", -1e-9, "-1\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%G", float32(-1e-9), "-1\u202f×\u202f10⁻⁰⁹"***REMOVED***,
	***REMOVED***"%20.5s", "qwertyuiop", "               qwert"***REMOVED***,
	***REMOVED***"%.5s", "qwertyuiop", "qwert"***REMOVED***,
	***REMOVED***"%-20.5s", "qwertyuiop", "qwert               "***REMOVED***,
	***REMOVED***"%20c", 'x', "                   x"***REMOVED***,
	***REMOVED***"%-20c", 'x', "x                   "***REMOVED***,
	***REMOVED***"%20.6e", 1.2345e3, "     1.234500\u202f×\u202f10⁰³"***REMOVED***,
	***REMOVED***"%20.6e", 1.2345e-3, "    1.234500\u202f×\u202f10⁻⁰³"***REMOVED***,
	***REMOVED***"%20e", 1.2345e3, "     1.234500\u202f×\u202f10⁰³"***REMOVED***,
	***REMOVED***"%20e", 1.2345e-3, "    1.234500\u202f×\u202f10⁻⁰³"***REMOVED***,
	***REMOVED***"%20.8e", 1.2345e3, "   1.23450000\u202f×\u202f10⁰³"***REMOVED***,
	***REMOVED***"%20f", 1.23456789e3, "        1,234.567890"***REMOVED***,
	***REMOVED***"%20f", 1.23456789e-3, "            0.001235"***REMOVED***,
	***REMOVED***"%20f", 12345678901.23456789, "12,345,678,901.234568"***REMOVED***,
	***REMOVED***"%-20f", 1.23456789e3, "1,234.567890        "***REMOVED***,
	***REMOVED***"%20.8f", 1.23456789e3, "      1,234.56789000"***REMOVED***,
	***REMOVED***"%20.8f", 1.23456789e-3, "          0.00123457"***REMOVED***,
	***REMOVED***"%g", 1.23456789e3, "1,234.56789"***REMOVED***,
	***REMOVED***"%g", 1.23456789e-3, "0.00123456789"***REMOVED***,
	***REMOVED***"%g", 1.23456789e20, "1.23456789\u202f×\u202f10²⁰"***REMOVED***,

	// arrays
	***REMOVED***"%v", array, "[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", iarray, "[1 hello 2.5 <nil>]"***REMOVED***,
	***REMOVED***"%v", barray, "[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", &array, "&[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", &iarray, "&[1 hello 2.5 <nil>]"***REMOVED***,
	***REMOVED***"%v", &barray, "&[1 2 3 4 5]"***REMOVED***,

	// slices
	***REMOVED***"%v", slice, "[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", islice, "[1 hello 2.5 <nil>]"***REMOVED***,
	***REMOVED***"%v", bslice, "[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", &slice, "&[1 2 3 4 5]"***REMOVED***,
	***REMOVED***"%v", &islice, "&[1 hello 2.5 <nil>]"***REMOVED***,
	***REMOVED***"%v", &bslice, "&[1 2 3 4 5]"***REMOVED***,

	// byte arrays and slices with %b,%c,%d,%o,%U and %v
	***REMOVED***"%b", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[1000001 1000010 1000011]"***REMOVED***,
	***REMOVED***"%c", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[A B C]"***REMOVED***,
	***REMOVED***"%d", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[65 66 67]"***REMOVED***,
	***REMOVED***"%o", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[101 102 103]"***REMOVED***,
	***REMOVED***"%U", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[U+0041 U+0042 U+0043]"***REMOVED***,
	***REMOVED***"%v", [3]byte***REMOVED***65, 66, 67***REMOVED***, "[65 66 67]"***REMOVED***,
	***REMOVED***"%v", [1]byte***REMOVED***123***REMOVED***, "[123]"***REMOVED***,
	***REMOVED***"%012v", []byte***REMOVED******REMOVED***, "[]"***REMOVED***,
	***REMOVED***"%#012v", []byte***REMOVED******REMOVED***, "[]byte***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%6v", []byte***REMOVED***1, 11, 111***REMOVED***, "[     1     11    111]"***REMOVED***,
	***REMOVED***"%06v", []byte***REMOVED***1, 11, 111***REMOVED***, "[000001 000011 000111]"***REMOVED***,
	***REMOVED***"%-6v", []byte***REMOVED***1, 11, 111***REMOVED***, "[1      11     111   ]"***REMOVED***,
	***REMOVED***"%-06v", []byte***REMOVED***1, 11, 111***REMOVED***, "[1      11     111   ]"***REMOVED***,
	***REMOVED***"%#v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED***0x1, 0xb, 0x6f***REMOVED***"***REMOVED***,
	***REMOVED***"%#6v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED***   0x1,    0xb,   0x6f***REMOVED***"***REMOVED***,
	***REMOVED***"%#06v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED***0x000001, 0x00000b, 0x00006f***REMOVED***"***REMOVED***,
	***REMOVED***"%#-6v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED***0x1   , 0xb   , 0x6f  ***REMOVED***"***REMOVED***,
	***REMOVED***"%#-06v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED***0x1   , 0xb   , 0x6f  ***REMOVED***"***REMOVED***,
	// f.space should and f.plus should not have an effect with %v.
	***REMOVED***"% v", []byte***REMOVED***1, 11, 111***REMOVED***, "[ 1  11  111]"***REMOVED***,
	***REMOVED***"%+v", [3]byte***REMOVED***1, 11, 111***REMOVED***, "[1 11 111]"***REMOVED***,
	***REMOVED***"%# -6v", []byte***REMOVED***1, 11, 111***REMOVED***, "[]byte***REMOVED*** 0x1  ,  0xb  ,  0x6f ***REMOVED***"***REMOVED***,
	***REMOVED***"%#+-6v", [3]byte***REMOVED***1, 11, 111***REMOVED***, "[3]uint8***REMOVED***0x1   , 0xb   , 0x6f  ***REMOVED***"***REMOVED***,
	// f.space and f.plus should have an effect with %d.
	***REMOVED***"% d", []byte***REMOVED***1, 11, 111***REMOVED***, "[ 1  11  111]"***REMOVED***,
	***REMOVED***"%+d", [3]byte***REMOVED***1, 11, 111***REMOVED***, "[+1 +11 +111]"***REMOVED***,
	***REMOVED***"%# -6d", []byte***REMOVED***1, 11, 111***REMOVED***, "[ 1      11     111  ]"***REMOVED***,
	***REMOVED***"%#+-6d", [3]byte***REMOVED***1, 11, 111***REMOVED***, "[+1     +11    +111  ]"***REMOVED***,

	// floates with %v
	***REMOVED***"%v", 1.2345678, "1.2345678"***REMOVED***,
	***REMOVED***"%v", float32(1.2345678), "1.2345678"***REMOVED***,

	// complexes with %v
	***REMOVED***"%v", 1 + 2i, "(1+2i)"***REMOVED***,
	***REMOVED***"%v", complex64(1 + 2i), "(1+2i)"***REMOVED***,

	// structs
	***REMOVED***"%v", A***REMOVED***1, 2, "a", []int***REMOVED***1, 2***REMOVED******REMOVED***, `***REMOVED***1 2 a [1 2]***REMOVED***`***REMOVED***,
	***REMOVED***"%+v", A***REMOVED***1, 2, "a", []int***REMOVED***1, 2***REMOVED******REMOVED***, `***REMOVED***i:1 j:2 s:a x:[1 2]***REMOVED***`***REMOVED***,

	// +v on structs with Stringable items
	***REMOVED***"%+v", B***REMOVED***1, 2***REMOVED***, `***REMOVED***I:<1> j:2***REMOVED***`***REMOVED***,
	***REMOVED***"%+v", C***REMOVED***1, B***REMOVED***2, 3***REMOVED******REMOVED***, `***REMOVED***i:1 B:***REMOVED***I:<2> j:3***REMOVED******REMOVED***`***REMOVED***,

	// other formats on Stringable items
	***REMOVED***"%s", I(23), `<23>`***REMOVED***,
	***REMOVED***"%q", I(23), `"<23>"`***REMOVED***,
	***REMOVED***"%x", I(23), `3c32333e`***REMOVED***,
	***REMOVED***"%#x", I(23), `0x3c32333e`***REMOVED***,
	***REMOVED***"%# x", I(23), `0x3c 0x32 0x33 0x3e`***REMOVED***,
	// Stringer applies only to string formats.
	***REMOVED***"%d", I(23), `23`***REMOVED***,
	// Stringer applies to the extracted value.
	***REMOVED***"%s", reflect.ValueOf(I(23)), `<23>`***REMOVED***,

	// go syntax
	***REMOVED***"%#v", A***REMOVED***1, 2, "a", []int***REMOVED***1, 2***REMOVED******REMOVED***, `message.A***REMOVED***i:1, j:0x2, s:"a", x:[]int***REMOVED***1, 2***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***"%#v", new(byte), "(*uint8)(0xPTR)"***REMOVED***,
	***REMOVED***"%#v", TestFmtInterface, "(func(*testing.T))(0xPTR)"***REMOVED***,
	***REMOVED***"%#v", make(chan int), "(chan int)(0xPTR)"***REMOVED***,
	***REMOVED***"%#v", uint64(1<<64 - 1), "0xffffffffffffffff"***REMOVED***,
	***REMOVED***"%#v", 1000000000, "1000000000"***REMOVED***,
	***REMOVED***"%#v", map[string]int***REMOVED***"a": 1***REMOVED***, `map[string]int***REMOVED***"a":1***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", map[string]B***REMOVED***"a": ***REMOVED***1, 2***REMOVED******REMOVED***, `map[string]message.B***REMOVED***"a":message.B***REMOVED***I:1, j:2***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***"%#v", []string***REMOVED***"a", "b"***REMOVED***, `[]string***REMOVED***"a", "b"***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", SI***REMOVED******REMOVED***, `message.SI***REMOVED***I:interface ***REMOVED******REMOVED***(nil)***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", []int(nil), `[]int(nil)`***REMOVED***,
	***REMOVED***"%#v", []int***REMOVED******REMOVED***, `[]int***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***"%#v", array, `[5]int***REMOVED***1, 2, 3, 4, 5***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", &array, `&[5]int***REMOVED***1, 2, 3, 4, 5***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", iarray, `[4]interface ***REMOVED******REMOVED******REMOVED***1, "hello", 2.5, interface ***REMOVED******REMOVED***(nil)***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", &iarray, `&[4]interface ***REMOVED******REMOVED******REMOVED***1, "hello", 2.5, interface ***REMOVED******REMOVED***(nil)***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", map[int]byte(nil), `map[int]uint8(nil)`***REMOVED***,
	***REMOVED***"%#v", map[int]byte***REMOVED******REMOVED***, `map[int]uint8***REMOVED******REMOVED***`***REMOVED***,
	***REMOVED***"%#v", "foo", `"foo"`***REMOVED***,
	***REMOVED***"%#v", barray, `[5]message.renamedUint8***REMOVED***0x1, 0x2, 0x3, 0x4, 0x5***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", bslice, `[]message.renamedUint8***REMOVED***0x1, 0x2, 0x3, 0x4, 0x5***REMOVED***`***REMOVED***,
	***REMOVED***"%#v", []int32(nil), "[]int32(nil)"***REMOVED***,
	***REMOVED***"%#v", 1.2345678, "1.2345678"***REMOVED***,
	***REMOVED***"%#v", float32(1.2345678), "1.2345678"***REMOVED***,
	// Only print []byte and []uint8 as type []byte if they appear at the top level.
	***REMOVED***"%#v", []byte(nil), "[]byte(nil)"***REMOVED***,
	***REMOVED***"%#v", []uint8(nil), "[]byte(nil)"***REMOVED***,
	***REMOVED***"%#v", []byte***REMOVED******REMOVED***, "[]byte***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", []uint8***REMOVED******REMOVED***, "[]byte***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", reflect.ValueOf([]byte***REMOVED******REMOVED***), "[]uint8***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", reflect.ValueOf([]uint8***REMOVED******REMOVED***), "[]uint8***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", &[]byte***REMOVED******REMOVED***, "&[]uint8***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", &[]byte***REMOVED******REMOVED***, "&[]uint8***REMOVED******REMOVED***"***REMOVED***,
	***REMOVED***"%#v", [3]byte***REMOVED******REMOVED***, "[3]uint8***REMOVED***0x0, 0x0, 0x0***REMOVED***"***REMOVED***,
	***REMOVED***"%#v", [3]uint8***REMOVED******REMOVED***, "[3]uint8***REMOVED***0x0, 0x0, 0x0***REMOVED***"***REMOVED***,

	// slices with other formats
	***REMOVED***"%#x", []int***REMOVED***1, 2, 15***REMOVED***, `[0x1 0x2 0xf]`***REMOVED***,
	***REMOVED***"%x", []int***REMOVED***1, 2, 15***REMOVED***, `[1 2 f]`***REMOVED***,
	***REMOVED***"%d", []int***REMOVED***1, 2, 15***REMOVED***, `[1 2 15]`***REMOVED***,
	***REMOVED***"%d", []byte***REMOVED***1, 2, 15***REMOVED***, `[1 2 15]`***REMOVED***,
	***REMOVED***"%q", []string***REMOVED***"a", "b"***REMOVED***, `["a" "b"]`***REMOVED***,
	***REMOVED***"% 02x", []byte***REMOVED***1***REMOVED***, "01"***REMOVED***,
	***REMOVED***"% 02x", []byte***REMOVED***1, 2, 3***REMOVED***, "01 02 03"***REMOVED***,

	// Padding with byte slices.
	***REMOVED***"%2x", []byte***REMOVED******REMOVED***, "  "***REMOVED***,
	***REMOVED***"%#2x", []byte***REMOVED******REMOVED***, "  "***REMOVED***,
	***REMOVED***"% 02x", []byte***REMOVED******REMOVED***, "00"***REMOVED***,
	***REMOVED***"%# 02x", []byte***REMOVED******REMOVED***, "00"***REMOVED***,
	***REMOVED***"%-2x", []byte***REMOVED******REMOVED***, "  "***REMOVED***,
	***REMOVED***"%-02x", []byte***REMOVED******REMOVED***, "  "***REMOVED***,
	***REMOVED***"%8x", []byte***REMOVED***0xab***REMOVED***, "      ab"***REMOVED***,
	***REMOVED***"% 8x", []byte***REMOVED***0xab***REMOVED***, "      ab"***REMOVED***,
	***REMOVED***"%#8x", []byte***REMOVED***0xab***REMOVED***, "    0xab"***REMOVED***,
	***REMOVED***"%# 8x", []byte***REMOVED***0xab***REMOVED***, "    0xab"***REMOVED***,
	***REMOVED***"%08x", []byte***REMOVED***0xab***REMOVED***, "000000ab"***REMOVED***,
	***REMOVED***"% 08x", []byte***REMOVED***0xab***REMOVED***, "000000ab"***REMOVED***,
	***REMOVED***"%#08x", []byte***REMOVED***0xab***REMOVED***, "00000xab"***REMOVED***,
	***REMOVED***"%# 08x", []byte***REMOVED***0xab***REMOVED***, "00000xab"***REMOVED***,
	***REMOVED***"%10x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "      abcd"***REMOVED***,
	***REMOVED***"% 10x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "     ab cd"***REMOVED***,
	***REMOVED***"%#10x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "    0xabcd"***REMOVED***,
	***REMOVED***"%# 10x", []byte***REMOVED***0xab, 0xcd***REMOVED***, " 0xab 0xcd"***REMOVED***,
	***REMOVED***"%010x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "000000abcd"***REMOVED***,
	***REMOVED***"% 010x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "00000ab cd"***REMOVED***,
	***REMOVED***"%#010x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "00000xabcd"***REMOVED***,
	***REMOVED***"%# 010x", []byte***REMOVED***0xab, 0xcd***REMOVED***, "00xab 0xcd"***REMOVED***,
	***REMOVED***"%-10X", []byte***REMOVED***0xab***REMOVED***, "AB        "***REMOVED***,
	***REMOVED***"% -010X", []byte***REMOVED***0xab***REMOVED***, "AB        "***REMOVED***,
	***REMOVED***"%#-10X", []byte***REMOVED***0xab, 0xcd***REMOVED***, "0XABCD    "***REMOVED***,
	***REMOVED***"%# -010X", []byte***REMOVED***0xab, 0xcd***REMOVED***, "0XAB 0XCD "***REMOVED***,
	// Same for strings
	***REMOVED***"%2x", "", "  "***REMOVED***,
	***REMOVED***"%#2x", "", "  "***REMOVED***,
	***REMOVED***"% 02x", "", "00"***REMOVED***,
	***REMOVED***"%# 02x", "", "00"***REMOVED***,
	***REMOVED***"%-2x", "", "  "***REMOVED***,
	***REMOVED***"%-02x", "", "  "***REMOVED***,
	***REMOVED***"%8x", "\xab", "      ab"***REMOVED***,
	***REMOVED***"% 8x", "\xab", "      ab"***REMOVED***,
	***REMOVED***"%#8x", "\xab", "    0xab"***REMOVED***,
	***REMOVED***"%# 8x", "\xab", "    0xab"***REMOVED***,
	***REMOVED***"%08x", "\xab", "000000ab"***REMOVED***,
	***REMOVED***"% 08x", "\xab", "000000ab"***REMOVED***,
	***REMOVED***"%#08x", "\xab", "00000xab"***REMOVED***,
	***REMOVED***"%# 08x", "\xab", "00000xab"***REMOVED***,
	***REMOVED***"%10x", "\xab\xcd", "      abcd"***REMOVED***,
	***REMOVED***"% 10x", "\xab\xcd", "     ab cd"***REMOVED***,
	***REMOVED***"%#10x", "\xab\xcd", "    0xabcd"***REMOVED***,
	***REMOVED***"%# 10x", "\xab\xcd", " 0xab 0xcd"***REMOVED***,
	***REMOVED***"%010x", "\xab\xcd", "000000abcd"***REMOVED***,
	***REMOVED***"% 010x", "\xab\xcd", "00000ab cd"***REMOVED***,
	***REMOVED***"%#010x", "\xab\xcd", "00000xabcd"***REMOVED***,
	***REMOVED***"%# 010x", "\xab\xcd", "00xab 0xcd"***REMOVED***,
	***REMOVED***"%-10X", "\xab", "AB        "***REMOVED***,
	***REMOVED***"% -010X", "\xab", "AB        "***REMOVED***,
	***REMOVED***"%#-10X", "\xab\xcd", "0XABCD    "***REMOVED***,
	***REMOVED***"%# -010X", "\xab\xcd", "0XAB 0XCD "***REMOVED***,

	// renamings
	***REMOVED***"%v", renamedBool(true), "true"***REMOVED***,
	***REMOVED***"%d", renamedBool(true), "%!d(message.renamedBool=true)"***REMOVED***,
	***REMOVED***"%o", renamedInt(8), "10"***REMOVED***,
	***REMOVED***"%d", renamedInt8(-9), "-9"***REMOVED***,
	***REMOVED***"%v", renamedInt16(10), "10"***REMOVED***,
	***REMOVED***"%v", renamedInt32(-11), "-11"***REMOVED***,
	***REMOVED***"%X", renamedInt64(255), "FF"***REMOVED***,
	***REMOVED***"%v", renamedUint(13), "13"***REMOVED***,
	***REMOVED***"%o", renamedUint8(14), "16"***REMOVED***,
	***REMOVED***"%X", renamedUint16(15), "F"***REMOVED***,
	***REMOVED***"%d", renamedUint32(16), "16"***REMOVED***,
	***REMOVED***"%X", renamedUint64(17), "11"***REMOVED***,
	***REMOVED***"%o", renamedUintptr(18), "22"***REMOVED***,
	***REMOVED***"%x", renamedString("thing"), "7468696e67"***REMOVED***,
	***REMOVED***"%d", renamedBytes([]byte***REMOVED***1, 2, 15***REMOVED***), `[1 2 15]`***REMOVED***,
	***REMOVED***"%q", renamedBytes([]byte("hello")), `"hello"`***REMOVED***,
	***REMOVED***"%x", []renamedUint8***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***, "68656c6c6f"***REMOVED***,
	***REMOVED***"%X", []renamedUint8***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***, "68656C6C6F"***REMOVED***,
	***REMOVED***"%s", []renamedUint8***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***, "hello"***REMOVED***,
	***REMOVED***"%q", []renamedUint8***REMOVED***'h', 'e', 'l', 'l', 'o'***REMOVED***, `"hello"`***REMOVED***,
	***REMOVED***"%v", renamedFloat32(22), "22"***REMOVED***,
	***REMOVED***"%v", renamedFloat64(33), "33"***REMOVED***,
	***REMOVED***"%v", renamedComplex64(3 + 4i), "(3+4i)"***REMOVED***,
	***REMOVED***"%v", renamedComplex128(4 - 3i), "(4-3i)"***REMOVED***,

	// Formatter
	***REMOVED***"%x", F(1), "<x=F(1)>"***REMOVED***,
	***REMOVED***"%x", G(2), "2"***REMOVED***,
	***REMOVED***"%+v", S***REMOVED***F(4), G(5)***REMOVED***, "***REMOVED***F:<v=F(4)> G:5***REMOVED***"***REMOVED***,

	// GoStringer
	***REMOVED***"%#v", G(6), "GoString(6)"***REMOVED***,
	***REMOVED***"%#v", S***REMOVED***F(7), G(8)***REMOVED***, "message.S***REMOVED***F:<v=F(7)>, G:GoString(8)***REMOVED***"***REMOVED***,

	// %T
	***REMOVED***"%T", byte(0), "uint8"***REMOVED***,
	***REMOVED***"%T", reflect.ValueOf(nil), "reflect.Value"***REMOVED***,
	***REMOVED***"%T", (4 - 3i), "complex128"***REMOVED***,
	***REMOVED***"%T", renamedComplex128(4 - 3i), "message.renamedComplex128"***REMOVED***,
	***REMOVED***"%T", intVar, "int"***REMOVED***,
	***REMOVED***"%6T", &intVar, "  *int"***REMOVED***,
	***REMOVED***"%10T", nil, "     <nil>"***REMOVED***,
	***REMOVED***"%-10T", nil, "<nil>     "***REMOVED***,

	// %p with pointers
	***REMOVED***"%p", (*int)(nil), "0x0"***REMOVED***,
	***REMOVED***"%#p", (*int)(nil), "0"***REMOVED***,
	***REMOVED***"%p", &intVar, "0xPTR"***REMOVED***,
	***REMOVED***"%#p", &intVar, "PTR"***REMOVED***,
	***REMOVED***"%p", &array, "0xPTR"***REMOVED***,
	***REMOVED***"%p", &slice, "0xPTR"***REMOVED***,
	***REMOVED***"%8.2p", (*int)(nil), "    0x00"***REMOVED***,
	***REMOVED***"%-20.16p", &intVar, "0xPTR  "***REMOVED***,
	// %p on non-pointers
	***REMOVED***"%p", make(chan int), "0xPTR"***REMOVED***,
	***REMOVED***"%p", make(map[int]int), "0xPTR"***REMOVED***,
	***REMOVED***"%p", func() ***REMOVED******REMOVED***, "0xPTR"***REMOVED***,
	***REMOVED***"%p", 27, "%!p(int=27)"***REMOVED***,  // not a pointer at all
	***REMOVED***"%p", nil, "%!p(<nil>)"***REMOVED***,  // nil on its own has no type ...
	***REMOVED***"%#p", nil, "%!p(<nil>)"***REMOVED***, // ... and hence is not a pointer type.
	// pointers with specified base
	***REMOVED***"%b", &intVar, "PTR_b"***REMOVED***,
	***REMOVED***"%d", &intVar, "PTR_d"***REMOVED***,
	***REMOVED***"%o", &intVar, "PTR_o"***REMOVED***,
	***REMOVED***"%x", &intVar, "PTR_x"***REMOVED***,
	***REMOVED***"%X", &intVar, "PTR_X"***REMOVED***,
	// %v on pointers
	***REMOVED***"%v", nil, "<nil>"***REMOVED***,
	***REMOVED***"%#v", nil, "<nil>"***REMOVED***,
	***REMOVED***"%v", (*int)(nil), "<nil>"***REMOVED***,
	***REMOVED***"%#v", (*int)(nil), "(*int)(nil)"***REMOVED***,
	***REMOVED***"%v", &intVar, "0xPTR"***REMOVED***,
	***REMOVED***"%#v", &intVar, "(*int)(0xPTR)"***REMOVED***,
	***REMOVED***"%8.2v", (*int)(nil), "   <nil>"***REMOVED***,
	***REMOVED***"%-20.16v", &intVar, "0xPTR  "***REMOVED***,
	// string method on pointer
	***REMOVED***"%s", &pValue, "String(p)"***REMOVED***, // String method...
	***REMOVED***"%p", &pValue, "0xPTR"***REMOVED***,     // ... is not called with %p.

	// %d on Stringer should give integer if possible
	***REMOVED***"%s", time.Time***REMOVED******REMOVED***.Month(), "January"***REMOVED***,
	***REMOVED***"%d", time.Time***REMOVED******REMOVED***.Month(), "1"***REMOVED***,

	// erroneous things
	***REMOVED***"%s %", "hello", "hello %!(NOVERB)"***REMOVED***,
	***REMOVED***"%s %.2", "hello", "hello %!(NOVERB)"***REMOVED***,

	// The "<nil>" show up because maps are printed by
	// first obtaining a list of keys and then looking up
	// each key. Since NaNs can be map keys but cannot
	// be fetched directly, the lookup fails and returns a
	// zero reflect.Value, which formats as <nil>.
	// This test is just to check that it shows the two NaNs at all.
	***REMOVED***"%v", map[float64]int***REMOVED***NaN: 1, NaN: 2***REMOVED***, "map[NaN:<nil> NaN:<nil>]"***REMOVED***,

	// Comparison of padding rules with C printf.
	/*
		C program:
		#include <stdio.h>

		char *format[] = ***REMOVED***
			"[%.2f]",
			"[% .2f]",
			"[%+.2f]",
			"[%7.2f]",
			"[% 7.2f]",
			"[%+7.2f]",
			"[% +7.2f]",
			"[%07.2f]",
			"[% 07.2f]",
			"[%+07.2f]",
			"[% +07.2f]"
		***REMOVED***;

		int main(void) ***REMOVED***
			int i;
			for(i = 0; i < 11; i++) ***REMOVED***
				printf("%s: ", format[i]);
				printf(format[i], 1.0);
				printf(" ");
				printf(format[i], -1.0);
				printf("\n");
			***REMOVED***
		***REMOVED***

		Output:
			[%.2f]: [1.00] [-1.00]
			[% .2f]: [ 1.00] [-1.00]
			[%+.2f]: [+1.00] [-1.00]
			[%7.2f]: [   1.00] [  -1.00]
			[% 7.2f]: [   1.00] [  -1.00]
			[%+7.2f]: [  +1.00] [  -1.00]
			[% +7.2f]: [  +1.00] [  -1.00]
			[%07.2f]: [0001.00] [-001.00]
			[% 07.2f]: [ 001.00] [-001.00]
			[%+07.2f]: [+001.00] [-001.00]
			[% +07.2f]: [+001.00] [-001.00]

	*/
	***REMOVED***"%.2f", 1.0, "1.00"***REMOVED***,
	***REMOVED***"%.2f", -1.0, "-1.00"***REMOVED***,
	***REMOVED***"% .2f", 1.0, " 1.00"***REMOVED***,
	***REMOVED***"% .2f", -1.0, "-1.00"***REMOVED***,
	***REMOVED***"%+.2f", 1.0, "+1.00"***REMOVED***,
	***REMOVED***"%+.2f", -1.0, "-1.00"***REMOVED***,
	***REMOVED***"%7.2f", 1.0, "   1.00"***REMOVED***,
	***REMOVED***"%7.2f", -1.0, "  -1.00"***REMOVED***,
	***REMOVED***"% 7.2f", 1.0, "   1.00"***REMOVED***,
	***REMOVED***"% 7.2f", -1.0, "  -1.00"***REMOVED***,
	***REMOVED***"%+7.2f", 1.0, "  +1.00"***REMOVED***,
	***REMOVED***"%+7.2f", -1.0, "  -1.00"***REMOVED***,
	***REMOVED***"% +7.2f", 1.0, "  +1.00"***REMOVED***,
	***REMOVED***"% +7.2f", -1.0, "  -1.00"***REMOVED***,
	// Padding with 0's indicates minimum number of integer digits minus the
	// period, if present, and minus the sign if it is fixed.
	// TODO: consider making this number the number of significant digits.
	***REMOVED***"%07.2f", 1.0, "0,001.00"***REMOVED***,
	***REMOVED***"%07.2f", -1.0, "-0,001.00"***REMOVED***,
	***REMOVED***"% 07.2f", 1.0, " 001.00"***REMOVED***,
	***REMOVED***"% 07.2f", -1.0, "-001.00"***REMOVED***,
	***REMOVED***"%+07.2f", 1.0, "+001.00"***REMOVED***,
	***REMOVED***"%+07.2f", -1.0, "-001.00"***REMOVED***,
	***REMOVED***"% +07.2f", 1.0, "+001.00"***REMOVED***,
	***REMOVED***"% +07.2f", -1.0, "-001.00"***REMOVED***,

	// Complex numbers: exhaustively tested in TestComplexFormatting.
	***REMOVED***"%7.2f", 1 + 2i, "(   1.00  +2.00i)"***REMOVED***,
	***REMOVED***"%+07.2f", -1 - 2i, "(-001.00-002.00i)"***REMOVED***,

	// Use spaces instead of zero if padding to the right.
	***REMOVED***"%0-5s", "abc", "abc  "***REMOVED***,
	***REMOVED***"%-05.1f", 1.0, "1.0  "***REMOVED***,

	// float and complex formatting should not change the padding width
	// for other elements. See issue 14642.
	***REMOVED***"%06v", []interface***REMOVED******REMOVED******REMOVED***+10.0, 10***REMOVED***, "[000,010 000,010]"***REMOVED***,
	***REMOVED***"%06v", []interface***REMOVED******REMOVED******REMOVED***-10.0, 10***REMOVED***, "[-000,010 000,010]"***REMOVED***,
	***REMOVED***"%06v", []interface***REMOVED******REMOVED******REMOVED***+10.0 + 10i, 10***REMOVED***, "[(000,010+00,010i) 000,010]"***REMOVED***,
	***REMOVED***"%06v", []interface***REMOVED******REMOVED******REMOVED***-10.0 + 10i, 10***REMOVED***, "[(-000,010+00,010i) 000,010]"***REMOVED***,

	// integer formatting should not alter padding for other elements.
	***REMOVED***"%03.6v", []interface***REMOVED******REMOVED******REMOVED***1, 2.0, "x"***REMOVED***, "[000,001 002 00x]"***REMOVED***,
	***REMOVED***"%03.0v", []interface***REMOVED******REMOVED******REMOVED***0, 2.0, "x"***REMOVED***, "[    002 000]"***REMOVED***,

	// Complex fmt used to leave the plus flag set for future entries in the array
	// causing +2+0i and +3+0i instead of 2+0i and 3+0i.
	***REMOVED***"%v", []complex64***REMOVED***1, 2, 3***REMOVED***, "[(1+0i) (2+0i) (3+0i)]"***REMOVED***,
	***REMOVED***"%v", []complex128***REMOVED***1, 2, 3***REMOVED***, "[(1+0i) (2+0i) (3+0i)]"***REMOVED***,

	// Incomplete format specification caused crash.
	***REMOVED***"%.", 3, "%!.(int=3)"***REMOVED***,

	// Padding for complex numbers. Has been bad, then fixed, then bad again.
	***REMOVED***"%+10.2f", +104.66 + 440.51i, "(   +104.66   +440.51i)"***REMOVED***,
	***REMOVED***"%+10.2f", -104.66 + 440.51i, "(   -104.66   +440.51i)"***REMOVED***,
	***REMOVED***"%+10.2f", +104.66 - 440.51i, "(   +104.66   -440.51i)"***REMOVED***,
	***REMOVED***"%+10.2f", -104.66 - 440.51i, "(   -104.66   -440.51i)"***REMOVED***,
	***REMOVED***"%010.2f", +104.66 + 440.51i, "(0,000,104.66+000,440.51i)"***REMOVED***,
	***REMOVED***"%+010.2f", +104.66 + 440.51i, "(+000,104.66+000,440.51i)"***REMOVED***,
	***REMOVED***"%+010.2f", -104.66 + 440.51i, "(-000,104.66+000,440.51i)"***REMOVED***,
	***REMOVED***"%+010.2f", +104.66 - 440.51i, "(+000,104.66-000,440.51i)"***REMOVED***,
	***REMOVED***"%+010.2f", -104.66 - 440.51i, "(-000,104.66-000,440.51i)"***REMOVED***,

	// []T where type T is a byte with a Stringer method.
	***REMOVED***"%v", byteStringerSlice, "[X X X X X]"***REMOVED***,
	***REMOVED***"%s", byteStringerSlice, "hello"***REMOVED***,
	***REMOVED***"%q", byteStringerSlice, "\"hello\""***REMOVED***,
	***REMOVED***"%x", byteStringerSlice, "68656c6c6f"***REMOVED***,
	***REMOVED***"%X", byteStringerSlice, "68656C6C6F"***REMOVED***,
	***REMOVED***"%#v", byteStringerSlice, "[]message.byteStringer***REMOVED***0x68, 0x65, 0x6c, 0x6c, 0x6f***REMOVED***"***REMOVED***,

	// And the same for Formatter.
	***REMOVED***"%v", byteFormatterSlice, "[X X X X X]"***REMOVED***,
	***REMOVED***"%s", byteFormatterSlice, "hello"***REMOVED***,
	***REMOVED***"%q", byteFormatterSlice, "\"hello\""***REMOVED***,
	***REMOVED***"%x", byteFormatterSlice, "68656c6c6f"***REMOVED***,
	***REMOVED***"%X", byteFormatterSlice, "68656C6C6F"***REMOVED***,
	// This next case seems wrong, but the docs say the Formatter wins here.
	***REMOVED***"%#v", byteFormatterSlice, "[]message.byteFormatter***REMOVED***X, X, X, X, X***REMOVED***"***REMOVED***,

	// reflect.Value handled specially in Go 1.5, making it possible to
	// see inside non-exported fields (which cannot be accessed with Interface()).
	// Issue 8965.
	***REMOVED***"%v", reflect.ValueOf(A***REMOVED******REMOVED***).Field(0).String(), "<int Value>"***REMOVED***, // Equivalent to the old way.
	***REMOVED***"%v", reflect.ValueOf(A***REMOVED******REMOVED***).Field(0), "0"***REMOVED***,                    // Sees inside the field.

	// verbs apply to the extracted value too.
	***REMOVED***"%s", reflect.ValueOf("hello"), "hello"***REMOVED***,
	***REMOVED***"%q", reflect.ValueOf("hello"), `"hello"`***REMOVED***,
	***REMOVED***"%#04x", reflect.ValueOf(256), "0x0100"***REMOVED***,

	// invalid reflect.Value doesn't crash.
	***REMOVED***"%v", reflect.Value***REMOVED******REMOVED***, "<invalid reflect.Value>"***REMOVED***,
	***REMOVED***"%v", &reflect.Value***REMOVED******REMOVED***, "<invalid Value>"***REMOVED***,
	***REMOVED***"%v", SI***REMOVED***reflect.Value***REMOVED******REMOVED******REMOVED***, "***REMOVED***<invalid Value>***REMOVED***"***REMOVED***,

	// Tests to check that not supported verbs generate an error string.
	***REMOVED***"%☠", nil, "%!☠(<nil>)"***REMOVED***,
	***REMOVED***"%☠", interface***REMOVED******REMOVED***(nil), "%!☠(<nil>)"***REMOVED***,
	***REMOVED***"%☠", int(0), "%!☠(int=0)"***REMOVED***,
	***REMOVED***"%☠", uint(0), "%!☠(uint=0)"***REMOVED***,
	***REMOVED***"%☠", []byte***REMOVED***0, 1***REMOVED***, "[%!☠(uint8=0) %!☠(uint8=1)]"***REMOVED***,
	***REMOVED***"%☠", []uint8***REMOVED***0, 1***REMOVED***, "[%!☠(uint8=0) %!☠(uint8=1)]"***REMOVED***,
	***REMOVED***"%☠", [1]byte***REMOVED***0***REMOVED***, "[%!☠(uint8=0)]"***REMOVED***,
	***REMOVED***"%☠", [1]uint8***REMOVED***0***REMOVED***, "[%!☠(uint8=0)]"***REMOVED***,
	***REMOVED***"%☠", "hello", "%!☠(string=hello)"***REMOVED***,
	***REMOVED***"%☠", 1.2345678, "%!☠(float64=1.2345678)"***REMOVED***,
	***REMOVED***"%☠", float32(1.2345678), "%!☠(float32=1.2345678)"***REMOVED***,
	***REMOVED***"%☠", 1.2345678 + 1.2345678i, "%!☠(complex128=(1.2345678+1.2345678i))"***REMOVED***,
	***REMOVED***"%☠", complex64(1.2345678 + 1.2345678i), "%!☠(complex64=(1.2345678+1.2345678i))"***REMOVED***,
	***REMOVED***"%☠", &intVar, "%!☠(*int=0xPTR)"***REMOVED***,
	***REMOVED***"%☠", make(chan int), "%!☠(chan int=0xPTR)"***REMOVED***,
	***REMOVED***"%☠", func() ***REMOVED******REMOVED***, "%!☠(func()=0xPTR)"***REMOVED***,
	***REMOVED***"%☠", reflect.ValueOf(renamedInt(0)), "%!☠(message.renamedInt=0)"***REMOVED***,
	***REMOVED***"%☠", SI***REMOVED***renamedInt(0)***REMOVED***, "***REMOVED***%!☠(message.renamedInt=0)***REMOVED***"***REMOVED***,
	***REMOVED***"%☠", &[]interface***REMOVED******REMOVED******REMOVED***I(1), G(2)***REMOVED***, "&[%!☠(message.I=1) %!☠(message.G=2)]"***REMOVED***,
	***REMOVED***"%☠", SI***REMOVED***&[]interface***REMOVED******REMOVED******REMOVED***I(1), G(2)***REMOVED******REMOVED***, "***REMOVED***%!☠(*[]interface ***REMOVED******REMOVED***=&[1 2])***REMOVED***"***REMOVED***,
	***REMOVED***"%☠", reflect.Value***REMOVED******REMOVED***, "<invalid reflect.Value>"***REMOVED***,
	***REMOVED***"%☠", map[float64]int***REMOVED***NaN: 1***REMOVED***, "map[%!☠(float64=NaN):%!☠(<nil>)]"***REMOVED***,
***REMOVED***

// zeroFill generates zero-filled strings of the specified width. The length
// of the suffix (but not the prefix) is compensated for in the width calculation.
func zeroFill(prefix string, width int, suffix string) string ***REMOVED***
	return prefix + strings.Repeat("0", width-len(suffix)) + suffix
***REMOVED***

func TestSprintf(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	for _, tt := range fmtTests ***REMOVED***
		t.Run(fmt.Sprint(tt.fmt, "/", tt.val), func(t *testing.T) ***REMOVED***
			s := p.Sprintf(tt.fmt, tt.val)
			i := strings.Index(tt.out, "PTR")
			if i >= 0 && i < len(s) ***REMOVED***
				var pattern, chars string
				switch ***REMOVED***
				case strings.HasPrefix(tt.out[i:], "PTR_b"):
					pattern = "PTR_b"
					chars = "01"
				case strings.HasPrefix(tt.out[i:], "PTR_o"):
					pattern = "PTR_o"
					chars = "01234567"
				case strings.HasPrefix(tt.out[i:], "PTR_d"):
					pattern = "PTR_d"
					chars = "0123456789"
				case strings.HasPrefix(tt.out[i:], "PTR_x"):
					pattern = "PTR_x"
					chars = "0123456789abcdef"
				case strings.HasPrefix(tt.out[i:], "PTR_X"):
					pattern = "PTR_X"
					chars = "0123456789ABCDEF"
				default:
					pattern = "PTR"
					chars = "0123456789abcdefABCDEF"
				***REMOVED***
				p := s[:i] + pattern
				for j := i; j < len(s); j++ ***REMOVED***
					if !strings.ContainsRune(chars, rune(s[j])) ***REMOVED***
						p += s[j:]
						break
					***REMOVED***
				***REMOVED***
				s = p
			***REMOVED***
			if s != tt.out ***REMOVED***
				if _, ok := tt.val.(string); ok ***REMOVED***
					// Don't requote the already-quoted strings.
					// It's too confusing to read the errors.
					t.Errorf("Sprintf(%q, %q) = <%s> want <%s>", tt.fmt, tt.val, s, tt.out)
				***REMOVED*** else ***REMOVED***
					t.Errorf("Sprintf(%q, %v) = %q want %q", tt.fmt, tt.val, s, tt.out)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

var f float64

// TestComplexFormatting checks that a complex always formats to the same
// thing as if done by hand with two singleton prints.
func TestComplexFormatting(t *testing.T) ***REMOVED***
	var yesNo = []bool***REMOVED***true, false***REMOVED***
	var values = []float64***REMOVED***1, 0, -1, posInf, negInf, NaN***REMOVED***
	p := NewPrinter(language.Und)
	for _, plus := range yesNo ***REMOVED***
		for _, zero := range yesNo ***REMOVED***
			for _, space := range yesNo ***REMOVED***
				for _, char := range "fFeEgG" ***REMOVED***
					realFmt := "%"
					if zero ***REMOVED***
						realFmt += "0"
					***REMOVED***
					if space ***REMOVED***
						realFmt += " "
					***REMOVED***
					if plus ***REMOVED***
						realFmt += "+"
					***REMOVED***
					realFmt += "10.2"
					realFmt += string(char)
					// Imaginary part always has a sign, so force + and ignore space.
					imagFmt := "%"
					if zero ***REMOVED***
						imagFmt += "0"
					***REMOVED***
					imagFmt += "+"
					imagFmt += "10.2"
					imagFmt += string(char)
					for _, realValue := range values ***REMOVED***
						for _, imagValue := range values ***REMOVED***
							one := p.Sprintf(realFmt, complex(realValue, imagValue))
							two := p.Sprintf("("+realFmt+imagFmt+"i)", realValue, imagValue)
							if math.IsNaN(imagValue) ***REMOVED***
								p := len(two) - len("NaNi)") - 1
								if two[p] == ' ' ***REMOVED***
									two = two[:p] + "+" + two[p+1:]
								***REMOVED*** else ***REMOVED***
									two = two[:p+1] + "+" + two[p+1:]
								***REMOVED***
							***REMOVED***
							if one != two ***REMOVED***
								t.Error(f, one, two)
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type SE []interface***REMOVED******REMOVED*** // slice of empty; notational compactness.

var reorderTests = []struct ***REMOVED***
	format string
	args   SE
	out    string
***REMOVED******REMOVED***
	***REMOVED***"%[1]d", SE***REMOVED***1***REMOVED***, "1"***REMOVED***,
	***REMOVED***"%[2]d", SE***REMOVED***2, 1***REMOVED***, "1"***REMOVED***,
	***REMOVED***"%[2]d %[1]d", SE***REMOVED***1, 2***REMOVED***, "2 1"***REMOVED***,
	***REMOVED***"%[2]*[1]d", SE***REMOVED***2, 5***REMOVED***, "    2"***REMOVED***,
	***REMOVED***"%6.2f", SE***REMOVED***12.0***REMOVED***, " 12.00"***REMOVED***, // Explicit version of next line.
	***REMOVED***"%[3]*.[2]*[1]f", SE***REMOVED***12.0, 2, 6***REMOVED***, " 12.00"***REMOVED***,
	***REMOVED***"%[1]*.[2]*[3]f", SE***REMOVED***6, 2, 12.0***REMOVED***, " 12.00"***REMOVED***,
	***REMOVED***"%10f", SE***REMOVED***12.0***REMOVED***, " 12.000000"***REMOVED***,
	***REMOVED***"%[1]*[3]f", SE***REMOVED***10, 99, 12.0***REMOVED***, " 12.000000"***REMOVED***,
	***REMOVED***"%.6f", SE***REMOVED***12.0***REMOVED***, "12.000000"***REMOVED***, // Explicit version of next line.
	***REMOVED***"%.[1]*[3]f", SE***REMOVED***6, 99, 12.0***REMOVED***, "12.000000"***REMOVED***,
	***REMOVED***"%6.f", SE***REMOVED***12.0***REMOVED***, "    12"***REMOVED***, //  // Explicit version of next line; empty precision means zero.
	***REMOVED***"%[1]*.[3]f", SE***REMOVED***6, 3, 12.0***REMOVED***, "    12"***REMOVED***,
	// An actual use! Print the same arguments twice.
	***REMOVED***"%d %d %d %#[1]o %#o %#o", SE***REMOVED***11, 12, 13***REMOVED***, "11 12 13 013 014 015"***REMOVED***,

	// Erroneous cases.
	***REMOVED***"%[d", SE***REMOVED***2, 1***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%]d", SE***REMOVED***2, 1***REMOVED***, "%!](int=2)d%!(EXTRA int=1)"***REMOVED***,
	***REMOVED***"%[]d", SE***REMOVED***2, 1***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%[-3]d", SE***REMOVED***2, 1***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%[99]d", SE***REMOVED***2, 1***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%[3]", SE***REMOVED***2, 1***REMOVED***, "%!(NOVERB)"***REMOVED***,
	***REMOVED***"%[1].2d", SE***REMOVED***5, 6***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%[1]2d", SE***REMOVED***2, 1***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%3.[2]d", SE***REMOVED***7***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%.[2]d", SE***REMOVED***7***REMOVED***, "%!d(BADINDEX)"***REMOVED***,
	***REMOVED***"%d %d %d %#[1]o %#o %#o %#o", SE***REMOVED***11, 12, 13***REMOVED***, "11 12 13 013 014 015 %!o(MISSING)"***REMOVED***,
	***REMOVED***"%[5]d %[2]d %d", SE***REMOVED***1, 2, 3***REMOVED***, "%!d(BADINDEX) 2 3"***REMOVED***,
	***REMOVED***"%d %[3]d %d", SE***REMOVED***1, 2***REMOVED***, "1 %!d(BADINDEX) 2"***REMOVED***, // Erroneous index does not affect sequence.
	***REMOVED***"%.[]", SE***REMOVED******REMOVED***, "%!](BADINDEX)"***REMOVED***,                // Issue 10675
	***REMOVED***"%.-3d", SE***REMOVED***42***REMOVED***, "%!-(int=42)3d"***REMOVED***,             // TODO: Should this set return better error messages?
	// The following messages are interpreted as if there is no substitution,
	// in which case it is okay to have extra arguments. This is different
	// semantics from the fmt package.
	***REMOVED***"%2147483648d", SE***REMOVED***42***REMOVED***, "%!(NOVERB)"***REMOVED***,
	***REMOVED***"%-2147483648d", SE***REMOVED***42***REMOVED***, "%!(NOVERB)"***REMOVED***,
	***REMOVED***"%.2147483648d", SE***REMOVED***42***REMOVED***, "%!(NOVERB)"***REMOVED***,
***REMOVED***

func TestReorder(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	for _, tc := range reorderTests ***REMOVED***
		t.Run(fmt.Sprint(tc.format, "/", tc.args), func(t *testing.T) ***REMOVED***
			s := p.Sprintf(tc.format, tc.args...)
			if s != tc.out ***REMOVED***
				t.Errorf("Sprintf(%q, %v) = %q want %q", tc.format, tc.args, s, tc.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkSprintfPadding(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%16f", 1.0)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfEmpty(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfString(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%s", "hello")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfTruncateString(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%.3s", "日本語日本語日本語")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfQuoteString(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%q", "日本語日本語日本語")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfInt(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%d", 5)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfIntInt(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%d %d", 5, 6)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfPrefixedInt(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("This is some meaningless prefix text that needs to be scanned %d", 6)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfFloat(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%g", 5.23184)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfComplex(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%f", 5.23184+5.23184i)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfBoolean(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%t", true)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfHexString(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("% #x", "0123456789abcdef")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfHexBytes(b *testing.B) ***REMOVED***
	data := []byte("0123456789abcdef")
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("% #x", data)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfBytes(b *testing.B) ***REMOVED***
	data := []byte("0123456789abcdef")
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%v", data)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfStringer(b *testing.B) ***REMOVED***
	stringer := I(12345)
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%v", stringer)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkSprintfStructure(b *testing.B) ***REMOVED***
	s := &[]interface***REMOVED******REMOVED******REMOVED***SI***REMOVED***12345***REMOVED***, map[int]string***REMOVED***0: "hello"***REMOVED******REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			p.Sprintf("%#v", s)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkManyArgs(b *testing.B) ***REMOVED***
	b.RunParallel(func(pb *testing.PB) ***REMOVED***
		var buf bytes.Buffer
		p := NewPrinter(language.English)
		for pb.Next() ***REMOVED***
			buf.Reset()
			p.Fprintf(&buf, "%2d/%2d/%2d %d:%d:%d %s %s\n", 3, 4, 5, 11, 12, 13, "hello", "world")
		***REMOVED***
	***REMOVED***)
***REMOVED***

func BenchmarkFprintInt(b *testing.B) ***REMOVED***
	var buf bytes.Buffer
	p := NewPrinter(language.English)
	for i := 0; i < b.N; i++ ***REMOVED***
		buf.Reset()
		p.Fprint(&buf, 123456)
	***REMOVED***
***REMOVED***

func BenchmarkFprintfBytes(b *testing.B) ***REMOVED***
	data := []byte(string("0123456789"))
	var buf bytes.Buffer
	p := NewPrinter(language.English)
	for i := 0; i < b.N; i++ ***REMOVED***
		buf.Reset()
		p.Fprintf(&buf, "%s", data)
	***REMOVED***
***REMOVED***

func BenchmarkFprintIntNoAlloc(b *testing.B) ***REMOVED***
	var x interface***REMOVED******REMOVED*** = 123456
	var buf bytes.Buffer
	p := NewPrinter(language.English)
	for i := 0; i < b.N; i++ ***REMOVED***
		buf.Reset()
		p.Fprint(&buf, x)
	***REMOVED***
***REMOVED***

var mallocBuf bytes.Buffer
var mallocPointer *int // A pointer so we know the interface value won't allocate.

var mallocTest = []struct ***REMOVED***
	count int
	desc  string
	fn    func(p *Printer)
***REMOVED******REMOVED***
	***REMOVED***0, `Sprintf("")`, func(p *Printer) ***REMOVED*** p.Sprintf("") ***REMOVED******REMOVED***,
	***REMOVED***1, `Sprintf("xxx")`, func(p *Printer) ***REMOVED*** p.Sprintf("xxx") ***REMOVED******REMOVED***,
	***REMOVED***2, `Sprintf("%x")`, func(p *Printer) ***REMOVED*** p.Sprintf("%x", 7) ***REMOVED******REMOVED***,
	***REMOVED***2, `Sprintf("%s")`, func(p *Printer) ***REMOVED*** p.Sprintf("%s", "hello") ***REMOVED******REMOVED***,
	***REMOVED***3, `Sprintf("%x %x")`, func(p *Printer) ***REMOVED*** p.Sprintf("%x %x", 7, 112) ***REMOVED******REMOVED***,
	***REMOVED***2, `Sprintf("%g")`, func(p *Printer) ***REMOVED*** p.Sprintf("%g", float32(3.14159)) ***REMOVED******REMOVED***, // TODO: Can this be 1?
	***REMOVED***1, `Fprintf(buf, "%s")`, func(p *Printer) ***REMOVED*** mallocBuf.Reset(); p.Fprintf(&mallocBuf, "%s", "hello") ***REMOVED******REMOVED***,
	// If the interface value doesn't need to allocate, amortized allocation overhead should be zero.
	***REMOVED***0, `Fprintf(buf, "%x %x %x")`, func(p *Printer) ***REMOVED***
		mallocBuf.Reset()
		p.Fprintf(&mallocBuf, "%x %x %x", mallocPointer, mallocPointer, mallocPointer)
	***REMOVED******REMOVED***,
***REMOVED***

var _ bytes.Buffer

func TestCountMallocs(t *testing.T) ***REMOVED***
	switch ***REMOVED***
	case testing.Short():
		t.Skip("skipping malloc count in short mode")
	case runtime.GOMAXPROCS(0) > 1:
		t.Skip("skipping; GOMAXPROCS>1")
		// TODO: detect race detecter enabled.
		// case race.Enabled:
		// 	t.Skip("skipping malloc count under race detector")
	***REMOVED***
	p := NewPrinter(language.English)
	for _, mt := range mallocTest ***REMOVED***
		mallocs := testing.AllocsPerRun(100, func() ***REMOVED*** mt.fn(p) ***REMOVED***)
		if got, max := mallocs, float64(mt.count); got > max ***REMOVED***
			t.Errorf("%s: got %v allocs, want <=%v", mt.desc, got, max)
		***REMOVED***
	***REMOVED***
***REMOVED***

type flagPrinter struct***REMOVED******REMOVED***

func (flagPrinter) Format(f fmt.State, c rune) ***REMOVED***
	s := "%"
	for i := 0; i < 128; i++ ***REMOVED***
		if f.Flag(i) ***REMOVED***
			s += string(i)
		***REMOVED***
	***REMOVED***
	if w, ok := f.Width(); ok ***REMOVED***
		s += fmt.Sprintf("%d", w)
	***REMOVED***
	if p, ok := f.Precision(); ok ***REMOVED***
		s += fmt.Sprintf(".%d", p)
	***REMOVED***
	s += string(c)
	io.WriteString(f, "["+s+"]")
***REMOVED***

var flagtests = []struct ***REMOVED***
	in  string
	out string
***REMOVED******REMOVED***
	***REMOVED***"%a", "[%a]"***REMOVED***,
	***REMOVED***"%-a", "[%-a]"***REMOVED***,
	***REMOVED***"%+a", "[%+a]"***REMOVED***,
	***REMOVED***"%#a", "[%#a]"***REMOVED***,
	***REMOVED***"% a", "[% a]"***REMOVED***,
	***REMOVED***"%0a", "[%0a]"***REMOVED***,
	***REMOVED***"%1.2a", "[%1.2a]"***REMOVED***,
	***REMOVED***"%-1.2a", "[%-1.2a]"***REMOVED***,
	***REMOVED***"%+1.2a", "[%+1.2a]"***REMOVED***,
	***REMOVED***"%-+1.2a", "[%+-1.2a]"***REMOVED***,
	***REMOVED***"%-+1.2abc", "[%+-1.2a]bc"***REMOVED***,
	***REMOVED***"%-1.2abc", "[%-1.2a]bc"***REMOVED***,
***REMOVED***

func TestFlagParser(t *testing.T) ***REMOVED***
	var flagprinter flagPrinter
	for _, tt := range flagtests ***REMOVED***
		s := NewPrinter(language.Und).Sprintf(tt.in, &flagprinter)
		if s != tt.out ***REMOVED***
			t.Errorf("Sprintf(%q, &flagprinter) => %q, want %q", tt.in, s, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStructPrinter(t *testing.T) ***REMOVED***
	type T struct ***REMOVED***
		a string
		b string
		c int
	***REMOVED***
	var s T
	s.a = "abc"
	s.b = "def"
	s.c = 123
	var tests = []struct ***REMOVED***
		fmt string
		out string
	***REMOVED******REMOVED***
		***REMOVED***"%v", "***REMOVED***abc def 123***REMOVED***"***REMOVED***,
		***REMOVED***"%+v", "***REMOVED***a:abc b:def c:123***REMOVED***"***REMOVED***,
		***REMOVED***"%#v", `message.T***REMOVED***a:"abc", b:"def", c:123***REMOVED***`***REMOVED***,
	***REMOVED***
	p := NewPrinter(language.Und)
	for _, tt := range tests ***REMOVED***
		out := p.Sprintf(tt.fmt, s)
		if out != tt.out ***REMOVED***
			t.Errorf("Sprintf(%q, s) = %#q, want %#q", tt.fmt, out, tt.out)
		***REMOVED***
		// The same but with a pointer.
		out = p.Sprintf(tt.fmt, &s)
		if out != "&"+tt.out ***REMOVED***
			t.Errorf("Sprintf(%q, &s) = %#q, want %#q", tt.fmt, out, "&"+tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSlicePrinter(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	slice := []int***REMOVED******REMOVED***
	s := p.Sprint(slice)
	if s != "[]" ***REMOVED***
		t.Errorf("empty slice printed as %q not %q", s, "[]")
	***REMOVED***
	slice = []int***REMOVED***1, 2, 3***REMOVED***
	s = p.Sprint(slice)
	if s != "[1 2 3]" ***REMOVED***
		t.Errorf("slice: got %q expected %q", s, "[1 2 3]")
	***REMOVED***
	s = p.Sprint(&slice)
	if s != "&[1 2 3]" ***REMOVED***
		t.Errorf("&slice: got %q expected %q", s, "&[1 2 3]")
	***REMOVED***
***REMOVED***

// presentInMap checks map printing using substrings so we don't depend on the
// print order.
func presentInMap(s string, a []string, t *testing.T) ***REMOVED***
	for i := 0; i < len(a); i++ ***REMOVED***
		loc := strings.Index(s, a[i])
		if loc < 0 ***REMOVED***
			t.Errorf("map print: expected to find %q in %q", a[i], s)
		***REMOVED***
		// make sure the match ends here
		loc += len(a[i])
		if loc >= len(s) || (s[loc] != ' ' && s[loc] != ']') ***REMOVED***
			t.Errorf("map print: %q not properly terminated in %q", a[i], s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMapPrinter(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	m0 := make(map[int]string)
	s := p.Sprint(m0)
	if s != "map[]" ***REMOVED***
		t.Errorf("empty map printed as %q not %q", s, "map[]")
	***REMOVED***
	m1 := map[int]string***REMOVED***1: "one", 2: "two", 3: "three"***REMOVED***
	a := []string***REMOVED***"1:one", "2:two", "3:three"***REMOVED***
	presentInMap(p.Sprintf("%v", m1), a, t)
	presentInMap(p.Sprint(m1), a, t)
	// Pointer to map prints the same but with initial &.
	if !strings.HasPrefix(p.Sprint(&m1), "&") ***REMOVED***
		t.Errorf("no initial & for address of map")
	***REMOVED***
	presentInMap(p.Sprintf("%v", &m1), a, t)
	presentInMap(p.Sprint(&m1), a, t)
***REMOVED***

func TestEmptyMap(t *testing.T) ***REMOVED***
	const emptyMapStr = "map[]"
	var m map[string]int
	p := NewPrinter(language.Und)
	s := p.Sprint(m)
	if s != emptyMapStr ***REMOVED***
		t.Errorf("nil map printed as %q not %q", s, emptyMapStr)
	***REMOVED***
	m = make(map[string]int)
	s = p.Sprint(m)
	if s != emptyMapStr ***REMOVED***
		t.Errorf("empty map printed as %q not %q", s, emptyMapStr)
	***REMOVED***
***REMOVED***

// TestBlank checks that Sprint (and hence Print, Fprint) puts spaces in the
// right places, that is, between arg pairs in which neither is a string.
func TestBlank(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	got := p.Sprint("<", 1, ">:", 1, 2, 3, "!")
	expect := "<1>:1 2 3!"
	if got != expect ***REMOVED***
		t.Errorf("got %q expected %q", got, expect)
	***REMOVED***
***REMOVED***

// TestBlankln checks that Sprintln (and hence Println, Fprintln) puts spaces in
// the right places, that is, between all arg pairs.
func TestBlankln(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	got := p.Sprintln("<", 1, ">:", 1, 2, 3, "!")
	expect := "< 1 >: 1 2 3 !\n"
	if got != expect ***REMOVED***
		t.Errorf("got %q expected %q", got, expect)
	***REMOVED***
***REMOVED***

// TestFormatterPrintln checks Formatter with Sprint, Sprintln, Sprintf.
func TestFormatterPrintln(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	f := F(1)
	expect := "<v=F(1)>\n"
	s := p.Sprint(f, "\n")
	if s != expect ***REMOVED***
		t.Errorf("Sprint wrong with Formatter: expected %q got %q", expect, s)
	***REMOVED***
	s = p.Sprintln(f)
	if s != expect ***REMOVED***
		t.Errorf("Sprintln wrong with Formatter: expected %q got %q", expect, s)
	***REMOVED***
	s = p.Sprintf("%v\n", f)
	if s != expect ***REMOVED***
		t.Errorf("Sprintf wrong with Formatter: expected %q got %q", expect, s)
	***REMOVED***
***REMOVED***

func args(a ...interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED*** return a ***REMOVED***

var startests = []struct ***REMOVED***
	fmt string
	in  []interface***REMOVED******REMOVED***
	out string
***REMOVED******REMOVED***
	***REMOVED***"%*d", args(4, 42), "  42"***REMOVED***,
	***REMOVED***"%-*d", args(4, 42), "42  "***REMOVED***,
	***REMOVED***"%*d", args(-4, 42), "42  "***REMOVED***,
	***REMOVED***"%-*d", args(-4, 42), "42  "***REMOVED***,
	***REMOVED***"%.*d", args(4, 42), "0,042"***REMOVED***,
	***REMOVED***"%*.*d", args(8, 4, 42), "   0,042"***REMOVED***,
	***REMOVED***"%0*d", args(4, 42), "0,042"***REMOVED***,
	// Some non-int types for width. (Issue 10732).
	***REMOVED***"%0*d", args(uint(4), 42), "0,042"***REMOVED***,
	***REMOVED***"%0*d", args(uint64(4), 42), "0,042"***REMOVED***,
	***REMOVED***"%0*d", args('\x04', 42), "0,042"***REMOVED***,
	***REMOVED***"%0*d", args(uintptr(4), 42), "0,042"***REMOVED***,

	// erroneous
	***REMOVED***"%*d", args(nil, 42), "%!(BADWIDTH)42"***REMOVED***,
	***REMOVED***"%*d", args(int(1e7), 42), "%!(BADWIDTH)42"***REMOVED***,
	***REMOVED***"%*d", args(int(-1e7), 42), "%!(BADWIDTH)42"***REMOVED***,
	***REMOVED***"%.*d", args(nil, 42), "%!(BADPREC)42"***REMOVED***,
	***REMOVED***"%.*d", args(-1, 42), "%!(BADPREC)42"***REMOVED***,
	***REMOVED***"%.*d", args(int(1e7), 42), "%!(BADPREC)42"***REMOVED***,
	***REMOVED***"%.*d", args(uint(1e7), 42), "%!(BADPREC)42"***REMOVED***,
	***REMOVED***"%.*d", args(uint64(1<<63), 42), "%!(BADPREC)42"***REMOVED***,   // Huge negative (-inf).
	***REMOVED***"%.*d", args(uint64(1<<64-1), 42), "%!(BADPREC)42"***REMOVED***, // Small negative (-1).
	***REMOVED***"%*d", args(5, "foo"), "%!d(string=  foo)"***REMOVED***,
	***REMOVED***"%*% %d", args(20, 5), "% 5"***REMOVED***,
	***REMOVED***"%*", args(4), "%!(NOVERB)"***REMOVED***,
***REMOVED***

func TestWidthAndPrecision(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	for i, tt := range startests ***REMOVED***
		t.Run(fmt.Sprint(tt.fmt, tt.in), func(t *testing.T) ***REMOVED***
			s := p.Sprintf(tt.fmt, tt.in...)
			if s != tt.out ***REMOVED***
				t.Errorf("#%d: %q: got %q expected %q", i, tt.fmt, s, tt.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// PanicS is a type that panics in String.
type PanicS struct ***REMOVED***
	message interface***REMOVED******REMOVED***
***REMOVED***

// Value receiver.
func (p PanicS) String() string ***REMOVED***
	panic(p.message)
***REMOVED***

// PanicGo is a type that panics in GoString.
type PanicGo struct ***REMOVED***
	message interface***REMOVED******REMOVED***
***REMOVED***

// Value receiver.
func (p PanicGo) GoString() string ***REMOVED***
	panic(p.message)
***REMOVED***

// PanicF is a type that panics in Format.
type PanicF struct ***REMOVED***
	message interface***REMOVED******REMOVED***
***REMOVED***

// Value receiver.
func (p PanicF) Format(f fmt.State, c rune) ***REMOVED***
	panic(p.message)
***REMOVED***

var panictests = []struct ***REMOVED***
	desc string
	fmt  string
	in   interface***REMOVED******REMOVED***
	out  string
***REMOVED******REMOVED***
	// String
	***REMOVED***"String", "%s", (*PanicS)(nil), "<nil>"***REMOVED***, // nil pointer special case
	***REMOVED***"String", "%s", PanicS***REMOVED***io.ErrUnexpectedEOF***REMOVED***, "%!s(PANIC=unexpected EOF)"***REMOVED***,
	***REMOVED***"String", "%s", PanicS***REMOVED***3***REMOVED***, "%!s(PANIC=3)"***REMOVED***,
	// GoString
	***REMOVED***"GoString", "%#v", (*PanicGo)(nil), "<nil>"***REMOVED***, // nil pointer special case
	***REMOVED***"GoString", "%#v", PanicGo***REMOVED***io.ErrUnexpectedEOF***REMOVED***, "%!v(PANIC=unexpected EOF)"***REMOVED***,
	***REMOVED***"GoString", "%#v", PanicGo***REMOVED***3***REMOVED***, "%!v(PANIC=3)"***REMOVED***,
	// Issue 18282. catchPanic should not clear fmtFlags permanently.
	***REMOVED***"Issue 18282", "%#v", []interface***REMOVED******REMOVED******REMOVED***PanicGo***REMOVED***3***REMOVED***, PanicGo***REMOVED***3***REMOVED******REMOVED***, "[]interface ***REMOVED******REMOVED******REMOVED***%!v(PANIC=3), %!v(PANIC=3)***REMOVED***"***REMOVED***,
	// Format
	***REMOVED***"Format", "%s", (*PanicF)(nil), "<nil>"***REMOVED***, // nil pointer special case
	***REMOVED***"Format", "%s", PanicF***REMOVED***io.ErrUnexpectedEOF***REMOVED***, "%!s(PANIC=unexpected EOF)"***REMOVED***,
	***REMOVED***"Format", "%s", PanicF***REMOVED***3***REMOVED***, "%!s(PANIC=3)"***REMOVED***,
***REMOVED***

func TestPanics(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	for i, tt := range panictests ***REMOVED***
		t.Run(fmt.Sprint(tt.desc, "/", tt.fmt, "/", tt.in), func(t *testing.T) ***REMOVED***
			s := p.Sprintf(tt.fmt, tt.in)
			if s != tt.out ***REMOVED***
				t.Errorf("%d: %q: got %q expected %q", i, tt.fmt, s, tt.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// recurCount tests that erroneous String routine doesn't cause fatal recursion.
var recurCount = 0

type Recur struct ***REMOVED***
	i      int
	failed *bool
***REMOVED***

func (r *Recur) String() string ***REMOVED***
	p := NewPrinter(language.Und)
	if recurCount++; recurCount > 10 ***REMOVED***
		*r.failed = true
		return "FAIL"
	***REMOVED***
	// This will call badVerb. Before the fix, that would cause us to recur into
	// this routine to print %!p(value). Now we don't call the user's method
	// during an error.
	return p.Sprintf("recur@%p value: %d", r, r.i)
***REMOVED***

func TestBadVerbRecursion(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	failed := false
	r := &Recur***REMOVED***3, &failed***REMOVED***
	p.Sprintf("recur@%p value: %d\n", &r, r.i)
	if failed ***REMOVED***
		t.Error("fail with pointer")
	***REMOVED***
	failed = false
	r = &Recur***REMOVED***4, &failed***REMOVED***
	p.Sprintf("recur@%p, value: %d\n", r, r.i)
	if failed ***REMOVED***
		t.Error("fail with value")
	***REMOVED***
***REMOVED***

func TestNilDoesNotBecomeTyped(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	type A struct***REMOVED******REMOVED***
	type B struct***REMOVED******REMOVED***
	var a *A = nil
	var b B = B***REMOVED******REMOVED***

	// indirect the Sprintf call through this noVetWarn variable to avoid
	// "go test" failing vet checks in Go 1.10+.
	noVetWarn := p.Sprintf
	got := noVetWarn("%s %s %s %s %s", nil, a, nil, b, nil)

	const expect = "%!s(<nil>) %!s(*message.A=<nil>) %!s(<nil>) ***REMOVED******REMOVED*** %!s(<nil>)"
	if got != expect ***REMOVED***
		t.Errorf("expected:\n\t%q\ngot:\n\t%q", expect, got)
	***REMOVED***
***REMOVED***

var formatterFlagTests = []struct ***REMOVED***
	in  string
	val interface***REMOVED******REMOVED***
	out string
***REMOVED******REMOVED***
	// scalar values with the (unused by fmt) 'a' verb.
	***REMOVED***"%a", flagPrinter***REMOVED******REMOVED***, "[%a]"***REMOVED***,
	***REMOVED***"%-a", flagPrinter***REMOVED******REMOVED***, "[%-a]"***REMOVED***,
	***REMOVED***"%+a", flagPrinter***REMOVED******REMOVED***, "[%+a]"***REMOVED***,
	***REMOVED***"%#a", flagPrinter***REMOVED******REMOVED***, "[%#a]"***REMOVED***,
	***REMOVED***"% a", flagPrinter***REMOVED******REMOVED***, "[% a]"***REMOVED***,
	***REMOVED***"%0a", flagPrinter***REMOVED******REMOVED***, "[%0a]"***REMOVED***,
	***REMOVED***"%1.2a", flagPrinter***REMOVED******REMOVED***, "[%1.2a]"***REMOVED***,
	***REMOVED***"%-1.2a", flagPrinter***REMOVED******REMOVED***, "[%-1.2a]"***REMOVED***,
	***REMOVED***"%+1.2a", flagPrinter***REMOVED******REMOVED***, "[%+1.2a]"***REMOVED***,
	***REMOVED***"%-+1.2a", flagPrinter***REMOVED******REMOVED***, "[%+-1.2a]"***REMOVED***,
	***REMOVED***"%-+1.2abc", flagPrinter***REMOVED******REMOVED***, "[%+-1.2a]bc"***REMOVED***,
	***REMOVED***"%-1.2abc", flagPrinter***REMOVED******REMOVED***, "[%-1.2a]bc"***REMOVED***,

	// composite values with the 'a' verb
	***REMOVED***"%a", [1]flagPrinter***REMOVED******REMOVED***, "[[%a]]"***REMOVED***,
	***REMOVED***"%-a", [1]flagPrinter***REMOVED******REMOVED***, "[[%-a]]"***REMOVED***,
	***REMOVED***"%+a", [1]flagPrinter***REMOVED******REMOVED***, "[[%+a]]"***REMOVED***,
	***REMOVED***"%#a", [1]flagPrinter***REMOVED******REMOVED***, "[[%#a]]"***REMOVED***,
	***REMOVED***"% a", [1]flagPrinter***REMOVED******REMOVED***, "[[% a]]"***REMOVED***,
	***REMOVED***"%0a", [1]flagPrinter***REMOVED******REMOVED***, "[[%0a]]"***REMOVED***,
	***REMOVED***"%1.2a", [1]flagPrinter***REMOVED******REMOVED***, "[[%1.2a]]"***REMOVED***,
	***REMOVED***"%-1.2a", [1]flagPrinter***REMOVED******REMOVED***, "[[%-1.2a]]"***REMOVED***,
	***REMOVED***"%+1.2a", [1]flagPrinter***REMOVED******REMOVED***, "[[%+1.2a]]"***REMOVED***,
	***REMOVED***"%-+1.2a", [1]flagPrinter***REMOVED******REMOVED***, "[[%+-1.2a]]"***REMOVED***,
	***REMOVED***"%-+1.2abc", [1]flagPrinter***REMOVED******REMOVED***, "[[%+-1.2a]]bc"***REMOVED***,
	***REMOVED***"%-1.2abc", [1]flagPrinter***REMOVED******REMOVED***, "[[%-1.2a]]bc"***REMOVED***,

	// simple values with the 'v' verb
	***REMOVED***"%v", flagPrinter***REMOVED******REMOVED***, "[%v]"***REMOVED***,
	***REMOVED***"%-v", flagPrinter***REMOVED******REMOVED***, "[%-v]"***REMOVED***,
	***REMOVED***"%+v", flagPrinter***REMOVED******REMOVED***, "[%+v]"***REMOVED***,
	***REMOVED***"%#v", flagPrinter***REMOVED******REMOVED***, "[%#v]"***REMOVED***,
	***REMOVED***"% v", flagPrinter***REMOVED******REMOVED***, "[% v]"***REMOVED***,
	***REMOVED***"%0v", flagPrinter***REMOVED******REMOVED***, "[%0v]"***REMOVED***,
	***REMOVED***"%1.2v", flagPrinter***REMOVED******REMOVED***, "[%1.2v]"***REMOVED***,
	***REMOVED***"%-1.2v", flagPrinter***REMOVED******REMOVED***, "[%-1.2v]"***REMOVED***,
	***REMOVED***"%+1.2v", flagPrinter***REMOVED******REMOVED***, "[%+1.2v]"***REMOVED***,
	***REMOVED***"%-+1.2v", flagPrinter***REMOVED******REMOVED***, "[%+-1.2v]"***REMOVED***,
	***REMOVED***"%-+1.2vbc", flagPrinter***REMOVED******REMOVED***, "[%+-1.2v]bc"***REMOVED***,
	***REMOVED***"%-1.2vbc", flagPrinter***REMOVED******REMOVED***, "[%-1.2v]bc"***REMOVED***,

	// composite values with the 'v' verb.
	***REMOVED***"%v", [1]flagPrinter***REMOVED******REMOVED***, "[[%v]]"***REMOVED***,
	***REMOVED***"%-v", [1]flagPrinter***REMOVED******REMOVED***, "[[%-v]]"***REMOVED***,
	***REMOVED***"%+v", [1]flagPrinter***REMOVED******REMOVED***, "[[%+v]]"***REMOVED***,
	***REMOVED***"%#v", [1]flagPrinter***REMOVED******REMOVED***, "[1]message.flagPrinter***REMOVED***[%#v]***REMOVED***"***REMOVED***,
	***REMOVED***"% v", [1]flagPrinter***REMOVED******REMOVED***, "[[% v]]"***REMOVED***,
	***REMOVED***"%0v", [1]flagPrinter***REMOVED******REMOVED***, "[[%0v]]"***REMOVED***,
	***REMOVED***"%1.2v", [1]flagPrinter***REMOVED******REMOVED***, "[[%1.2v]]"***REMOVED***,
	***REMOVED***"%-1.2v", [1]flagPrinter***REMOVED******REMOVED***, "[[%-1.2v]]"***REMOVED***,
	***REMOVED***"%+1.2v", [1]flagPrinter***REMOVED******REMOVED***, "[[%+1.2v]]"***REMOVED***,
	***REMOVED***"%-+1.2v", [1]flagPrinter***REMOVED******REMOVED***, "[[%+-1.2v]]"***REMOVED***,
	***REMOVED***"%-+1.2vbc", [1]flagPrinter***REMOVED******REMOVED***, "[[%+-1.2v]]bc"***REMOVED***,
	***REMOVED***"%-1.2vbc", [1]flagPrinter***REMOVED******REMOVED***, "[[%-1.2v]]bc"***REMOVED***,
***REMOVED***

func TestFormatterFlags(t *testing.T) ***REMOVED***
	p := NewPrinter(language.Und)
	for _, tt := range formatterFlagTests ***REMOVED***
		s := p.Sprintf(tt.in, tt.val)
		if s != tt.out ***REMOVED***
			t.Errorf("Sprintf(%q, %T) = %q, want %q", tt.in, tt.val, s, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***
