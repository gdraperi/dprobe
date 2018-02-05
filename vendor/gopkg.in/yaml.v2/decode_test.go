package yaml_test

import (
	"errors"
	. "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"
	"math"
	"net"
	"reflect"
	"strings"
	"time"
)

var unmarshalIntTest = 123

var unmarshalTests = []struct ***REMOVED***
	data  string
	value interface***REMOVED******REMOVED***
***REMOVED******REMOVED***
	***REMOVED***
		"",
		&struct***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"***REMOVED******REMOVED***", &struct***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: hi",
		map[string]string***REMOVED***"v": "hi"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: hi", map[string]interface***REMOVED******REMOVED******REMOVED***"v": "hi"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: true",
		map[string]string***REMOVED***"v": "true"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: true",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": true***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 10",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 10***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 0b10",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 2***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 0xA",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 10***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 4294967296",
		map[string]int64***REMOVED***"v": 4294967296***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 0.1",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 0.1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: .1",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 0.1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: .Inf",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": math.Inf(+1)***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: -.Inf",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": math.Inf(-1)***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: -10",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": -10***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: -.1",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": -0.1***REMOVED***,
	***REMOVED***,

	// Simple values.
	***REMOVED***
		"123",
		&unmarshalIntTest,
	***REMOVED***,

	// Floats from spec
	***REMOVED***
		"canonical: 6.8523e+5",
		map[string]interface***REMOVED******REMOVED******REMOVED***"canonical": 6.8523e+5***REMOVED***,
	***REMOVED***, ***REMOVED***
		"expo: 685.230_15e+03",
		map[string]interface***REMOVED******REMOVED******REMOVED***"expo": 685.23015e+03***REMOVED***,
	***REMOVED***, ***REMOVED***
		"fixed: 685_230.15",
		map[string]interface***REMOVED******REMOVED******REMOVED***"fixed": 685230.15***REMOVED***,
	***REMOVED***, ***REMOVED***
		"neginf: -.inf",
		map[string]interface***REMOVED******REMOVED******REMOVED***"neginf": math.Inf(-1)***REMOVED***,
	***REMOVED***, ***REMOVED***
		"fixed: 685_230.15",
		map[string]float64***REMOVED***"fixed": 685230.15***REMOVED***,
	***REMOVED***,
	//***REMOVED***"sexa: 190:20:30.15", map[string]interface***REMOVED******REMOVED******REMOVED***"sexa": 0***REMOVED******REMOVED***, // Unsupported
	//***REMOVED***"notanum: .NaN", map[string]interface***REMOVED******REMOVED******REMOVED***"notanum": math.NaN()***REMOVED******REMOVED***, // Equality of NaN fails.

	// Bools from spec
	***REMOVED***
		"canonical: y",
		map[string]interface***REMOVED******REMOVED******REMOVED***"canonical": true***REMOVED***,
	***REMOVED***, ***REMOVED***
		"answer: NO",
		map[string]interface***REMOVED******REMOVED******REMOVED***"answer": false***REMOVED***,
	***REMOVED***, ***REMOVED***
		"logical: True",
		map[string]interface***REMOVED******REMOVED******REMOVED***"logical": true***REMOVED***,
	***REMOVED***, ***REMOVED***
		"option: on",
		map[string]interface***REMOVED******REMOVED******REMOVED***"option": true***REMOVED***,
	***REMOVED***, ***REMOVED***
		"option: on",
		map[string]bool***REMOVED***"option": true***REMOVED***,
	***REMOVED***,
	// Ints from spec
	***REMOVED***
		"canonical: 685230",
		map[string]interface***REMOVED******REMOVED******REMOVED***"canonical": 685230***REMOVED***,
	***REMOVED***, ***REMOVED***
		"decimal: +685_230",
		map[string]interface***REMOVED******REMOVED******REMOVED***"decimal": 685230***REMOVED***,
	***REMOVED***, ***REMOVED***
		"octal: 02472256",
		map[string]interface***REMOVED******REMOVED******REMOVED***"octal": 685230***REMOVED***,
	***REMOVED***, ***REMOVED***
		"hexa: 0x_0A_74_AE",
		map[string]interface***REMOVED******REMOVED******REMOVED***"hexa": 685230***REMOVED***,
	***REMOVED***, ***REMOVED***
		"bin: 0b1010_0111_0100_1010_1110",
		map[string]interface***REMOVED******REMOVED******REMOVED***"bin": 685230***REMOVED***,
	***REMOVED***, ***REMOVED***
		"bin: -0b101010",
		map[string]interface***REMOVED******REMOVED******REMOVED***"bin": -42***REMOVED***,
	***REMOVED***, ***REMOVED***
		"decimal: +685_230",
		map[string]int***REMOVED***"decimal": 685230***REMOVED***,
	***REMOVED***,

	//***REMOVED***"sexa: 190:20:30", map[string]interface***REMOVED******REMOVED******REMOVED***"sexa": 0***REMOVED******REMOVED***, // Unsupported

	// Nulls from spec
	***REMOVED***
		"empty:",
		map[string]interface***REMOVED******REMOVED******REMOVED***"empty": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"canonical: ~",
		map[string]interface***REMOVED******REMOVED******REMOVED***"canonical": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"english: null",
		map[string]interface***REMOVED******REMOVED******REMOVED***"english": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"~: null key",
		map[interface***REMOVED******REMOVED***]string***REMOVED***nil: "null key"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"empty:",
		map[string]*bool***REMOVED***"empty": nil***REMOVED***,
	***REMOVED***,

	// Flow sequence
	***REMOVED***
		"seq: [A,B]",
		map[string]interface***REMOVED******REMOVED******REMOVED***"seq": []interface***REMOVED******REMOVED******REMOVED***"A", "B"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq: [A,B,C,]",
		map[string][]string***REMOVED***"seq": []string***REMOVED***"A", "B", "C"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq: [A,1,C]",
		map[string][]string***REMOVED***"seq": []string***REMOVED***"A", "1", "C"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq: [A,1,C]",
		map[string][]int***REMOVED***"seq": []int***REMOVED***1***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq: [A,1,C]",
		map[string]interface***REMOVED******REMOVED******REMOVED***"seq": []interface***REMOVED******REMOVED******REMOVED***"A", 1, "C"***REMOVED******REMOVED***,
	***REMOVED***,
	// Block sequence
	***REMOVED***
		"seq:\n - A\n - B",
		map[string]interface***REMOVED******REMOVED******REMOVED***"seq": []interface***REMOVED******REMOVED******REMOVED***"A", "B"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq:\n - A\n - B\n - C",
		map[string][]string***REMOVED***"seq": []string***REMOVED***"A", "B", "C"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq:\n - A\n - 1\n - C",
		map[string][]string***REMOVED***"seq": []string***REMOVED***"A", "1", "C"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq:\n - A\n - 1\n - C",
		map[string][]int***REMOVED***"seq": []int***REMOVED***1***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"seq:\n - A\n - 1\n - C",
		map[string]interface***REMOVED******REMOVED******REMOVED***"seq": []interface***REMOVED******REMOVED******REMOVED***"A", 1, "C"***REMOVED******REMOVED***,
	***REMOVED***,

	// Literal block scalar
	***REMOVED***
		"scalar: | # Comment\n\n literal\n\n \ttext\n\n",
		map[string]string***REMOVED***"scalar": "\nliteral\n\n\ttext\n"***REMOVED***,
	***REMOVED***,

	// Folded block scalar
	***REMOVED***
		"scalar: > # Comment\n\n folded\n line\n \n next\n line\n  * one\n  * two\n\n last\n line\n\n",
		map[string]string***REMOVED***"scalar": "\nfolded line\nnext line\n * one\n * two\n\nlast line\n"***REMOVED***,
	***REMOVED***,

	// Map inside interface with no type hints.
	***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"a": map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"b": "c"***REMOVED******REMOVED***,
	***REMOVED***,

	// Structs and type conversions.
	***REMOVED***
		"hello: world",
		&struct***REMOVED*** Hello string ***REMOVED******REMOVED***"world"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		&struct***REMOVED*** A struct***REMOVED*** B string ***REMOVED*** ***REMOVED******REMOVED***struct***REMOVED*** B string ***REMOVED******REMOVED***"c"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		&struct***REMOVED*** A *struct***REMOVED*** B string ***REMOVED*** ***REMOVED******REMOVED***&struct***REMOVED*** B string ***REMOVED******REMOVED***"c"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		&struct***REMOVED*** A map[string]string ***REMOVED******REMOVED***map[string]string***REMOVED***"b": "c"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		&struct***REMOVED*** A *map[string]string ***REMOVED******REMOVED***&map[string]string***REMOVED***"b": "c"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a:",
		&struct***REMOVED*** A map[string]string ***REMOVED******REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1",
		&struct***REMOVED*** A int ***REMOVED******REMOVED***1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1",
		&struct***REMOVED*** A float64 ***REMOVED******REMOVED***1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1.0",
		&struct***REMOVED*** A int ***REMOVED******REMOVED***1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1.0",
		&struct***REMOVED*** A uint ***REMOVED******REMOVED***1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: [1, 2]",
		&struct***REMOVED*** A []int ***REMOVED******REMOVED***[]int***REMOVED***1, 2***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1",
		&struct***REMOVED*** B int ***REMOVED******REMOVED***0***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 1",
		&struct ***REMOVED***
			B int "a"
		***REMOVED******REMOVED***1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: y",
		&struct***REMOVED*** A bool ***REMOVED******REMOVED***true***REMOVED***,
	***REMOVED***,

	// Some cross type conversions
	***REMOVED***
		"v: 42",
		map[string]uint***REMOVED***"v": 42***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: -42",
		map[string]uint***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 4294967296",
		map[string]uint64***REMOVED***"v": 4294967296***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: -4294967296",
		map[string]uint64***REMOVED******REMOVED***,
	***REMOVED***,

	// int
	***REMOVED***
		"int_max: 2147483647",
		map[string]int***REMOVED***"int_max": math.MaxInt32***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int_min: -2147483648",
		map[string]int***REMOVED***"int_min": math.MinInt32***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int_overflow: 9223372036854775808", // math.MaxInt64 + 1
		map[string]int***REMOVED******REMOVED***,
	***REMOVED***,

	// int64
	***REMOVED***
		"int64_max: 9223372036854775807",
		map[string]int64***REMOVED***"int64_max": math.MaxInt64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int64_max_base2: 0b111111111111111111111111111111111111111111111111111111111111111",
		map[string]int64***REMOVED***"int64_max_base2": math.MaxInt64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int64_min: -9223372036854775808",
		map[string]int64***REMOVED***"int64_min": math.MinInt64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int64_neg_base2: -0b111111111111111111111111111111111111111111111111111111111111111",
		map[string]int64***REMOVED***"int64_neg_base2": -math.MaxInt64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"int64_overflow: 9223372036854775808", // math.MaxInt64 + 1
		map[string]int64***REMOVED******REMOVED***,
	***REMOVED***,

	// uint
	***REMOVED***
		"uint_min: 0",
		map[string]uint***REMOVED***"uint_min": 0***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint_max: 4294967295",
		map[string]uint***REMOVED***"uint_max": math.MaxUint32***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint_underflow: -1",
		map[string]uint***REMOVED******REMOVED***,
	***REMOVED***,

	// uint64
	***REMOVED***
		"uint64_min: 0",
		map[string]uint***REMOVED***"uint64_min": 0***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint64_max: 18446744073709551615",
		map[string]uint64***REMOVED***"uint64_max": math.MaxUint64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint64_max_base2: 0b1111111111111111111111111111111111111111111111111111111111111111",
		map[string]uint64***REMOVED***"uint64_max_base2": math.MaxUint64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint64_maxint64: 9223372036854775807",
		map[string]uint64***REMOVED***"uint64_maxint64": math.MaxInt64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"uint64_underflow: -1",
		map[string]uint64***REMOVED******REMOVED***,
	***REMOVED***,

	// float32
	***REMOVED***
		"float32_max: 3.40282346638528859811704183484516925440e+38",
		map[string]float32***REMOVED***"float32_max": math.MaxFloat32***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float32_nonzero: 1.401298464324817070923729583289916131280e-45",
		map[string]float32***REMOVED***"float32_nonzero": math.SmallestNonzeroFloat32***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float32_maxuint64: 18446744073709551615",
		map[string]float32***REMOVED***"float32_maxuint64": float32(math.MaxUint64)***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float32_maxuint64+1: 18446744073709551616",
		map[string]float32***REMOVED***"float32_maxuint64+1": float32(math.MaxUint64 + 1)***REMOVED***,
	***REMOVED***,

	// float64
	***REMOVED***
		"float64_max: 1.797693134862315708145274237317043567981e+308",
		map[string]float64***REMOVED***"float64_max": math.MaxFloat64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float64_nonzero: 4.940656458412465441765687928682213723651e-324",
		map[string]float64***REMOVED***"float64_nonzero": math.SmallestNonzeroFloat64***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float64_maxuint64: 18446744073709551615",
		map[string]float64***REMOVED***"float64_maxuint64": float64(math.MaxUint64)***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"float64_maxuint64+1: 18446744073709551616",
		map[string]float64***REMOVED***"float64_maxuint64+1": float64(math.MaxUint64 + 1)***REMOVED***,
	***REMOVED***,

	// Overflow cases.
	***REMOVED***
		"v: 4294967297",
		map[string]int32***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: 128",
		map[string]int8***REMOVED******REMOVED***,
	***REMOVED***,

	// Quoted values.
	***REMOVED***
		"'1': '\"2\"'",
		map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"1": "\"2\""***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v:\n- A\n- 'B\n\n  C'\n",
		map[string][]string***REMOVED***"v": []string***REMOVED***"A", "B\nC"***REMOVED******REMOVED***,
	***REMOVED***,

	// Explicit tags.
	***REMOVED***
		"v: !!float '1.1'",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 1.1***REMOVED***,
	***REMOVED***, ***REMOVED***
		"v: !!null ''",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"%TAG !y! tag:yaml.org,2002:\n---\nv: !y!int '1'",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": 1***REMOVED***,
	***REMOVED***,

	// Non-specific tag (Issue #75)
	***REMOVED***
		"v: ! test",
		map[string]interface***REMOVED******REMOVED******REMOVED***"v": "test"***REMOVED***,
	***REMOVED***,

	// Anchors and aliases.
	***REMOVED***
		"a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
		&struct***REMOVED*** A, B, C, D int ***REMOVED******REMOVED***1, 2, 1, 2***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: &a ***REMOVED***c: 1***REMOVED***\nb: *a",
		&struct ***REMOVED***
			A, B struct ***REMOVED***
				C int
			***REMOVED***
		***REMOVED******REMOVED***struct***REMOVED*** C int ***REMOVED******REMOVED***1***REMOVED***, struct***REMOVED*** C int ***REMOVED******REMOVED***1***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: &a [1, 2]\nb: *a",
		&struct***REMOVED*** B []int ***REMOVED******REMOVED***[]int***REMOVED***1, 2***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		"b: *a\na: &a ***REMOVED***c: 1***REMOVED***",
		&struct ***REMOVED***
			A, B struct ***REMOVED***
				C int
			***REMOVED***
		***REMOVED******REMOVED***struct***REMOVED*** C int ***REMOVED******REMOVED***1***REMOVED***, struct***REMOVED*** C int ***REMOVED******REMOVED***1***REMOVED******REMOVED***,
	***REMOVED***,

	// Bug #1133337
	***REMOVED***
		"foo: ''",
		map[string]*string***REMOVED***"foo": new(string)***REMOVED***,
	***REMOVED***, ***REMOVED***
		"foo: null",
		map[string]*string***REMOVED***"foo": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"foo: null",
		map[string]string***REMOVED***"foo": ""***REMOVED***,
	***REMOVED***, ***REMOVED***
		"foo: null",
		map[string]interface***REMOVED******REMOVED******REMOVED***"foo": nil***REMOVED***,
	***REMOVED***,

	// Support for ~
	***REMOVED***
		"foo: ~",
		map[string]*string***REMOVED***"foo": nil***REMOVED***,
	***REMOVED***, ***REMOVED***
		"foo: ~",
		map[string]string***REMOVED***"foo": ""***REMOVED***,
	***REMOVED***, ***REMOVED***
		"foo: ~",
		map[string]interface***REMOVED******REMOVED******REMOVED***"foo": nil***REMOVED***,
	***REMOVED***,

	// Ignored field
	***REMOVED***
		"a: 1\nb: 2\n",
		&struct ***REMOVED***
			A int
			B int "-"
		***REMOVED******REMOVED***1, 0***REMOVED***,
	***REMOVED***,

	// Bug #1191981
	***REMOVED***
		"" +
			"%YAML 1.1\n" +
			"--- !!str\n" +
			`"Generic line break (no glyph)\n\` + "\n" +
			` Generic line break (glyphed)\n\` + "\n" +
			` Line separator\u2028\` + "\n" +
			` Paragraph separator\u2029"` + "\n",
		"" +
			"Generic line break (no glyph)\n" +
			"Generic line break (glyphed)\n" +
			"Line separator\u2028Paragraph separator\u2029",
	***REMOVED***,

	// Struct inlining
	***REMOVED***
		"a: 1\nb: 2\nc: 3\n",
		&struct ***REMOVED***
			A int
			C inlineB `yaml:",inline"`
		***REMOVED******REMOVED***1, inlineB***REMOVED***2, inlineC***REMOVED***3***REMOVED******REMOVED******REMOVED***,
	***REMOVED***,

	// Map inlining
	***REMOVED***
		"a: 1\nb: 2\nc: 3\n",
		&struct ***REMOVED***
			A int
			C map[string]int `yaml:",inline"`
		***REMOVED******REMOVED***1, map[string]int***REMOVED***"b": 2, "c": 3***REMOVED******REMOVED***,
	***REMOVED***,

	// bug 1243827
	***REMOVED***
		"a: -b_c",
		map[string]interface***REMOVED******REMOVED******REMOVED***"a": "-b_c"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"a: +b_c",
		map[string]interface***REMOVED******REMOVED******REMOVED***"a": "+b_c"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"a: 50cent_of_dollar",
		map[string]interface***REMOVED******REMOVED******REMOVED***"a": "50cent_of_dollar"***REMOVED***,
	***REMOVED***,

	// Duration
	***REMOVED***
		"a: 3s",
		map[string]time.Duration***REMOVED***"a": 3 * time.Second***REMOVED***,
	***REMOVED***,

	// Issue #24.
	***REMOVED***
		"a: <foo>",
		map[string]string***REMOVED***"a": "<foo>"***REMOVED***,
	***REMOVED***,

	// Base 60 floats are obsolete and unsupported.
	***REMOVED***
		"a: 1:1\n",
		map[string]string***REMOVED***"a": "1:1"***REMOVED***,
	***REMOVED***,

	// Binary data.
	***REMOVED***
		"a: !!binary gIGC\n",
		map[string]string***REMOVED***"a": "\x80\x81\x82"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
		map[string]string***REMOVED***"a": strings.Repeat("\x90", 54)***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: !!binary |\n  " + strings.Repeat("A", 70) + "\n  ==\n",
		map[string]string***REMOVED***"a": strings.Repeat("\x00", 52)***REMOVED***,
	***REMOVED***,

	// Ordered maps.
	***REMOVED***
		"***REMOVED***b: 2, a: 1, d: 4, c: 3, sub: ***REMOVED***e: 5***REMOVED******REMOVED***",
		&yaml.MapSlice***REMOVED******REMOVED***"b", 2***REMOVED***, ***REMOVED***"a", 1***REMOVED***, ***REMOVED***"d", 4***REMOVED***, ***REMOVED***"c", 3***REMOVED***, ***REMOVED***"sub", yaml.MapSlice***REMOVED******REMOVED***"e", 5***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***,

	// Issue #39.
	***REMOVED***
		"a:\n b:\n  c: d\n",
		map[string]struct***REMOVED*** B interface***REMOVED******REMOVED*** ***REMOVED******REMOVED***"a": ***REMOVED***map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"c": "d"***REMOVED******REMOVED******REMOVED***,
	***REMOVED***,

	// Custom map type.
	***REMOVED***
		"a: ***REMOVED***b: c***REMOVED***",
		M***REMOVED***"a": M***REMOVED***"b": "c"***REMOVED******REMOVED***,
	***REMOVED***,

	// Support encoding.TextUnmarshaler.
	***REMOVED***
		"a: 1.2.3.4\n",
		map[string]net.IP***REMOVED***"a": net.IPv4(1, 2, 3, 4)***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"a: 2015-02-24T18:19:39Z\n",
		map[string]time.Time***REMOVED***"a": time.Unix(1424801979, 0).In(time.UTC)***REMOVED***,
	***REMOVED***,

	// Encode empty lists as zero-length slices.
	***REMOVED***
		"a: []",
		&struct***REMOVED*** A []int ***REMOVED******REMOVED***[]int***REMOVED******REMOVED******REMOVED***,
	***REMOVED***,

	// UTF-16-LE
	***REMOVED***
		"\xff\xfe\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00\n\x00",
		M***REMOVED***"ñoño": "very yes"***REMOVED***,
	***REMOVED***,
	// UTF-16-LE with surrogate.
	***REMOVED***
		"\xff\xfe\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00 \x00=\xd8\xd4\xdf\n\x00",
		M***REMOVED***"ñoño": "very yes 🟔"***REMOVED***,
	***REMOVED***,

	// UTF-16-BE
	***REMOVED***
		"\xfe\xff\x00\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00\n",
		M***REMOVED***"ñoño": "very yes"***REMOVED***,
	***REMOVED***,
	// UTF-16-BE with surrogate.
	***REMOVED***
		"\xfe\xff\x00\xf1\x00o\x00\xf1\x00o\x00:\x00 \x00v\x00e\x00r\x00y\x00 \x00y\x00e\x00s\x00 \xd8=\xdf\xd4\x00\n",
		M***REMOVED***"ñoño": "very yes 🟔"***REMOVED***,
	***REMOVED***,

	// YAML Float regex shouldn't match this
	***REMOVED***
		"a: 123456e1\n",
		M***REMOVED***"a": "123456e1"***REMOVED***,
	***REMOVED***, ***REMOVED***
		"a: 123456E1\n",
		M***REMOVED***"a": "123456E1"***REMOVED***,
	***REMOVED***,
***REMOVED***

type M map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***

type inlineB struct ***REMOVED***
	B       int
	inlineC `yaml:",inline"`
***REMOVED***

type inlineC struct ***REMOVED***
	C int
***REMOVED***

func (s *S) TestUnmarshal(c *C) ***REMOVED***
	for i, item := range unmarshalTests ***REMOVED***
		c.Logf("test %d: %q", i, item.data)
		t := reflect.ValueOf(item.value).Type()
		var value interface***REMOVED******REMOVED***
		switch t.Kind() ***REMOVED***
		case reflect.Map:
			value = reflect.MakeMap(t).Interface()
		case reflect.String:
			value = reflect.New(t).Interface()
		case reflect.Ptr:
			value = reflect.New(t.Elem()).Interface()
		default:
			c.Fatalf("missing case for %s", t)
		***REMOVED***
		err := yaml.Unmarshal([]byte(item.data), value)
		if _, ok := err.(*yaml.TypeError); !ok ***REMOVED***
			c.Assert(err, IsNil)
		***REMOVED***
		if t.Kind() == reflect.String ***REMOVED***
			c.Assert(*value.(*string), Equals, item.value)
		***REMOVED*** else ***REMOVED***
			c.Assert(value, DeepEquals, item.value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *S) TestUnmarshalNaN(c *C) ***REMOVED***
	value := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err := yaml.Unmarshal([]byte("notanum: .NaN"), &value)
	c.Assert(err, IsNil)
	c.Assert(math.IsNaN(value["notanum"].(float64)), Equals, true)
***REMOVED***

var unmarshalErrorTests = []struct ***REMOVED***
	data, error string
***REMOVED******REMOVED***
	***REMOVED***"v: !!float 'error'", "yaml: cannot decode !!str `error` as a !!float"***REMOVED***,
	***REMOVED***"v: [A,", "yaml: line 1: did not find expected node content"***REMOVED***,
	***REMOVED***"v:\n- [A,", "yaml: line 2: did not find expected node content"***REMOVED***,
	***REMOVED***"a: *b\n", "yaml: unknown anchor 'b' referenced"***REMOVED***,
	***REMOVED***"a: &a\n  b: *a\n", "yaml: anchor 'a' value contains itself"***REMOVED***,
	***REMOVED***"value: -", "yaml: block sequence entries are not allowed in this context"***REMOVED***,
	***REMOVED***"a: !!binary ==", "yaml: !!binary value contains invalid base64 data"***REMOVED***,
	***REMOVED***"***REMOVED***[.]***REMOVED***", `yaml: invalid map key: \[\]interface \***REMOVED***\***REMOVED***\***REMOVED***"\."\***REMOVED***`***REMOVED***,
	***REMOVED***"***REMOVED******REMOVED***.***REMOVED******REMOVED***", `yaml: invalid map key: map\[interface\ \***REMOVED***\***REMOVED***\]interface \***REMOVED***\***REMOVED***\***REMOVED***".":interface \***REMOVED***\***REMOVED***\(nil\)\***REMOVED***`***REMOVED***,
	***REMOVED***"%TAG !%79! tag:yaml.org,2002:\n---\nv: !%79!int '1'", "yaml: did not find expected whitespace"***REMOVED***,
***REMOVED***

func (s *S) TestUnmarshalErrors(c *C) ***REMOVED***
	for _, item := range unmarshalErrorTests ***REMOVED***
		var value interface***REMOVED******REMOVED***
		err := yaml.Unmarshal([]byte(item.data), &value)
		c.Assert(err, ErrorMatches, item.error, Commentf("Partial unmarshal: %#v", value))
	***REMOVED***
***REMOVED***

var unmarshalerTests = []struct ***REMOVED***
	data, tag string
	value     interface***REMOVED******REMOVED***
***REMOVED******REMOVED***
	***REMOVED***"_: ***REMOVED***hi: there***REMOVED***", "!!map", map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***"hi": "there"***REMOVED******REMOVED***,
	***REMOVED***"_: [1,A]", "!!seq", []interface***REMOVED******REMOVED******REMOVED***1, "A"***REMOVED******REMOVED***,
	***REMOVED***"_: 10", "!!int", 10***REMOVED***,
	***REMOVED***"_: null", "!!null", nil***REMOVED***,
	***REMOVED***`_: BAR!`, "!!str", "BAR!"***REMOVED***,
	***REMOVED***`_: "BAR!"`, "!!str", "BAR!"***REMOVED***,
	***REMOVED***"_: !!foo 'BAR!'", "!!foo", "BAR!"***REMOVED***,
	***REMOVED***`_: ""`, "!!str", ""***REMOVED***,
***REMOVED***

var unmarshalerResult = map[int]error***REMOVED******REMOVED***

type unmarshalerType struct ***REMOVED***
	value interface***REMOVED******REMOVED***
***REMOVED***

func (o *unmarshalerType) UnmarshalYAML(unmarshal func(v interface***REMOVED******REMOVED***) error) error ***REMOVED***
	if err := unmarshal(&o.value); err != nil ***REMOVED***
		return err
	***REMOVED***
	if i, ok := o.value.(int); ok ***REMOVED***
		if result, ok := unmarshalerResult[i]; ok ***REMOVED***
			return result
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type unmarshalerPointer struct ***REMOVED***
	Field *unmarshalerType "_"
***REMOVED***

type unmarshalerValue struct ***REMOVED***
	Field unmarshalerType "_"
***REMOVED***

func (s *S) TestUnmarshalerPointerField(c *C) ***REMOVED***
	for _, item := range unmarshalerTests ***REMOVED***
		obj := &unmarshalerPointer***REMOVED******REMOVED***
		err := yaml.Unmarshal([]byte(item.data), obj)
		c.Assert(err, IsNil)
		if item.value == nil ***REMOVED***
			c.Assert(obj.Field, IsNil)
		***REMOVED*** else ***REMOVED***
			c.Assert(obj.Field, NotNil, Commentf("Pointer not initialized (%#v)", item.value))
			c.Assert(obj.Field.value, DeepEquals, item.value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *S) TestUnmarshalerValueField(c *C) ***REMOVED***
	for _, item := range unmarshalerTests ***REMOVED***
		obj := &unmarshalerValue***REMOVED******REMOVED***
		err := yaml.Unmarshal([]byte(item.data), obj)
		c.Assert(err, IsNil)
		c.Assert(obj.Field, NotNil, Commentf("Pointer not initialized (%#v)", item.value))
		c.Assert(obj.Field.value, DeepEquals, item.value)
	***REMOVED***
***REMOVED***

func (s *S) TestUnmarshalerWholeDocument(c *C) ***REMOVED***
	obj := &unmarshalerType***REMOVED******REMOVED***
	err := yaml.Unmarshal([]byte(unmarshalerTests[0].data), obj)
	c.Assert(err, IsNil)
	value, ok := obj.value.(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
	c.Assert(ok, Equals, true, Commentf("value: %#v", obj.value))
	c.Assert(value["_"], DeepEquals, unmarshalerTests[0].value)
***REMOVED***

func (s *S) TestUnmarshalerTypeError(c *C) ***REMOVED***
	unmarshalerResult[2] = &yaml.TypeError***REMOVED***[]string***REMOVED***"foo"***REMOVED******REMOVED***
	unmarshalerResult[4] = &yaml.TypeError***REMOVED***[]string***REMOVED***"bar"***REMOVED******REMOVED***
	defer func() ***REMOVED***
		delete(unmarshalerResult, 2)
		delete(unmarshalerResult, 4)
	***REMOVED***()

	type T struct ***REMOVED***
		Before int
		After  int
		M      map[string]*unmarshalerType
	***REMOVED***
	var v T
	data := `***REMOVED***before: A, m: ***REMOVED***abc: 1, def: 2, ghi: 3, jkl: 4***REMOVED***, after: B***REMOVED***`
	err := yaml.Unmarshal([]byte(data), &v)
	c.Assert(err, ErrorMatches, ""+
		"yaml: unmarshal errors:\n"+
		"  line 1: cannot unmarshal !!str `A` into int\n"+
		"  foo\n"+
		"  bar\n"+
		"  line 1: cannot unmarshal !!str `B` into int")
	c.Assert(v.M["abc"], NotNil)
	c.Assert(v.M["def"], IsNil)
	c.Assert(v.M["ghi"], NotNil)
	c.Assert(v.M["jkl"], IsNil)

	c.Assert(v.M["abc"].value, Equals, 1)
	c.Assert(v.M["ghi"].value, Equals, 3)
***REMOVED***

type proxyTypeError struct***REMOVED******REMOVED***

func (v *proxyTypeError) UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error ***REMOVED***
	var s string
	var a int32
	var b int64
	if err := unmarshal(&s); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if s == "a" ***REMOVED***
		if err := unmarshal(&b); err == nil ***REMOVED***
			panic("should have failed")
		***REMOVED***
		return unmarshal(&a)
	***REMOVED***
	if err := unmarshal(&a); err == nil ***REMOVED***
		panic("should have failed")
	***REMOVED***
	return unmarshal(&b)
***REMOVED***

func (s *S) TestUnmarshalerTypeErrorProxying(c *C) ***REMOVED***
	type T struct ***REMOVED***
		Before int
		After  int
		M      map[string]*proxyTypeError
	***REMOVED***
	var v T
	data := `***REMOVED***before: A, m: ***REMOVED***abc: a, def: b***REMOVED***, after: B***REMOVED***`
	err := yaml.Unmarshal([]byte(data), &v)
	c.Assert(err, ErrorMatches, ""+
		"yaml: unmarshal errors:\n"+
		"  line 1: cannot unmarshal !!str `A` into int\n"+
		"  line 1: cannot unmarshal !!str `a` into int32\n"+
		"  line 1: cannot unmarshal !!str `b` into int64\n"+
		"  line 1: cannot unmarshal !!str `B` into int")
***REMOVED***

type failingUnmarshaler struct***REMOVED******REMOVED***

var failingErr = errors.New("failingErr")

func (ft *failingUnmarshaler) UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error ***REMOVED***
	return failingErr
***REMOVED***

func (s *S) TestUnmarshalerError(c *C) ***REMOVED***
	err := yaml.Unmarshal([]byte("a: b"), &failingUnmarshaler***REMOVED******REMOVED***)
	c.Assert(err, Equals, failingErr)
***REMOVED***

type sliceUnmarshaler []int

func (su *sliceUnmarshaler) UnmarshalYAML(unmarshal func(interface***REMOVED******REMOVED***) error) error ***REMOVED***
	var slice []int
	err := unmarshal(&slice)
	if err == nil ***REMOVED***
		*su = slice
		return nil
	***REMOVED***

	var intVal int
	err = unmarshal(&intVal)
	if err == nil ***REMOVED***
		*su = []int***REMOVED***intVal***REMOVED***
		return nil
	***REMOVED***

	return err
***REMOVED***

func (s *S) TestUnmarshalerRetry(c *C) ***REMOVED***
	var su sliceUnmarshaler
	err := yaml.Unmarshal([]byte("[1, 2, 3]"), &su)
	c.Assert(err, IsNil)
	c.Assert(su, DeepEquals, sliceUnmarshaler([]int***REMOVED***1, 2, 3***REMOVED***))

	err = yaml.Unmarshal([]byte("1"), &su)
	c.Assert(err, IsNil)
	c.Assert(su, DeepEquals, sliceUnmarshaler([]int***REMOVED***1***REMOVED***))
***REMOVED***

// From http://yaml.org/type/merge.html
var mergeTests = `
anchors:
  list:
    - &CENTER ***REMOVED*** "x": 1, "y": 2 ***REMOVED***
    - &LEFT   ***REMOVED*** "x": 0, "y": 2 ***REMOVED***
    - &BIG    ***REMOVED*** "r": 10 ***REMOVED***
    - &SMALL  ***REMOVED*** "r": 1 ***REMOVED***

# All the following maps are equal:

plain:
  # Explicit keys
  "x": 1
  "y": 2
  "r": 10
  label: center/big

mergeOne:
  # Merge one map
  << : *CENTER
  "r": 10
  label: center/big

mergeMultiple:
  # Merge multiple maps
  << : [ *CENTER, *BIG ]
  label: center/big

override:
  # Override
  << : [ *BIG, *LEFT, *SMALL ]
  "x": 1
  label: center/big

shortTag:
  # Explicit short merge tag
  !!merge "<<" : [ *CENTER, *BIG ]
  label: center/big

longTag:
  # Explicit merge long tag
  !<tag:yaml.org,2002:merge> "<<" : [ *CENTER, *BIG ]
  label: center/big

inlineMap:
  # Inlined map 
  << : ***REMOVED***"x": 1, "y": 2, "r": 10***REMOVED***
  label: center/big

inlineSequenceMap:
  # Inlined map in sequence
  << : [ *CENTER, ***REMOVED***"r": 10***REMOVED*** ]
  label: center/big
`

func (s *S) TestMerge(c *C) ***REMOVED***
	var want = map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***
		"x":     1,
		"y":     2,
		"r":     10,
		"label": "center/big",
	***REMOVED***

	var m map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	err := yaml.Unmarshal([]byte(mergeTests), &m)
	c.Assert(err, IsNil)
	for name, test := range m ***REMOVED***
		if name == "anchors" ***REMOVED***
			continue
		***REMOVED***
		c.Assert(test, DeepEquals, want, Commentf("test %q failed", name))
	***REMOVED***
***REMOVED***

func (s *S) TestMergeStruct(c *C) ***REMOVED***
	type Data struct ***REMOVED***
		X, Y, R int
		Label   string
	***REMOVED***
	want := Data***REMOVED***1, 2, 10, "center/big"***REMOVED***

	var m map[string]Data
	err := yaml.Unmarshal([]byte(mergeTests), &m)
	c.Assert(err, IsNil)
	for name, test := range m ***REMOVED***
		if name == "anchors" ***REMOVED***
			continue
		***REMOVED***
		c.Assert(test, Equals, want, Commentf("test %q failed", name))
	***REMOVED***
***REMOVED***

var unmarshalNullTests = []func() interface***REMOVED******REMOVED******REMOVED***
	func() interface***REMOVED******REMOVED*** ***REMOVED*** var v interface***REMOVED******REMOVED***; v = "v"; return &v ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** var s = "s"; return &s ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** var s = "s"; sptr := &s; return &sptr ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** var i = 1; return &i ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** var i = 1; iptr := &i; return &iptr ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** m := map[string]int***REMOVED***"s": 1***REMOVED***; return &m ***REMOVED***,
	func() interface***REMOVED******REMOVED*** ***REMOVED*** m := map[string]int***REMOVED***"s": 1***REMOVED***; return m ***REMOVED***,
***REMOVED***

func (s *S) TestUnmarshalNull(c *C) ***REMOVED***
	for _, test := range unmarshalNullTests ***REMOVED***
		item := test()
		zero := reflect.Zero(reflect.TypeOf(item).Elem()).Interface()
		err := yaml.Unmarshal([]byte("null"), item)
		c.Assert(err, IsNil)
		if reflect.TypeOf(item).Kind() == reflect.Map ***REMOVED***
			c.Assert(reflect.ValueOf(item).Interface(), DeepEquals, reflect.MakeMap(reflect.TypeOf(item)).Interface())
		***REMOVED*** else ***REMOVED***
			c.Assert(reflect.ValueOf(item).Elem().Interface(), DeepEquals, zero)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *S) TestUnmarshalSliceOnPreset(c *C) ***REMOVED***
	// Issue #48.
	v := struct***REMOVED*** A []int ***REMOVED******REMOVED***[]int***REMOVED***1***REMOVED******REMOVED***
	yaml.Unmarshal([]byte("a: [2]"), &v)
	c.Assert(v.A, DeepEquals, []int***REMOVED***2***REMOVED***)
***REMOVED***

func (s *S) TestUnmarshalStrict(c *C) ***REMOVED***
	v := struct***REMOVED*** A, B int ***REMOVED******REMOVED******REMOVED***

	err := yaml.UnmarshalStrict([]byte("a: 1\nb: 2"), &v)
	c.Check(err, IsNil)
	err = yaml.Unmarshal([]byte("a: 1\nb: 2\nc: 3"), &v)
	c.Check(err, IsNil)
	err = yaml.UnmarshalStrict([]byte("a: 1\nb: 2\nc: 3"), &v)
	c.Check(err, ErrorMatches, "yaml: unmarshal errors:\n  line 3: field c not found in struct struct ***REMOVED*** A int; B int ***REMOVED***")
***REMOVED***

//var data []byte
//func init() ***REMOVED***
//	var err error
//	data, err = ioutil.ReadFile("/tmp/file.yaml")
//	if err != nil ***REMOVED***
//		panic(err)
//	***REMOVED***
//***REMOVED***
//
//func (s *S) BenchmarkUnmarshal(c *C) ***REMOVED***
//	var err error
//	for i := 0; i < c.N; i++ ***REMOVED***
//		var v map[string]interface***REMOVED******REMOVED***
//		err = yaml.Unmarshal(data, &v)
//	***REMOVED***
//	if err != nil ***REMOVED***
//		panic(err)
//	***REMOVED***
//***REMOVED***
//
//func (s *S) BenchmarkMarshal(c *C) ***REMOVED***
//	var v map[string]interface***REMOVED******REMOVED***
//	yaml.Unmarshal(data, &v)
//	c.ResetTimer()
//	for i := 0; i < c.N; i++ ***REMOVED***
//		yaml.Marshal(&v)
//	***REMOVED***
//***REMOVED***
