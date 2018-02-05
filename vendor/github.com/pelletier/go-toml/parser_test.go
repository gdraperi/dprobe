package toml

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func assertSubTree(t *testing.T, path []string, tree *Tree, err error, ref map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	if err != nil ***REMOVED***
		t.Error("Non-nil error:", err.Error())
		return
	***REMOVED***
	for k, v := range ref ***REMOVED***
		nextPath := append(path, k)
		t.Log("asserting path", nextPath)
		// NOTE: directly access key instead of resolve by path
		// NOTE: see TestSpecialKV
		switch node := tree.GetPath([]string***REMOVED***k***REMOVED***).(type) ***REMOVED***
		case []*Tree:
			t.Log("\tcomparing key", nextPath, "by array iteration")
			for idx, item := range node ***REMOVED***
				assertSubTree(t, nextPath, item, err, v.([]map[string]interface***REMOVED******REMOVED***)[idx])
			***REMOVED***
		case *Tree:
			t.Log("\tcomparing key", nextPath, "by subtree assestion")
			assertSubTree(t, nextPath, node, err, v.(map[string]interface***REMOVED******REMOVED***))
		default:
			t.Log("\tcomparing key", nextPath, "by string representation because it's of type", reflect.TypeOf(node))
			if fmt.Sprintf("%v", node) != fmt.Sprintf("%v", v) ***REMOVED***
				t.Errorf("was expecting %v at %v but got %v", v, k, node)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func assertTree(t *testing.T, tree *Tree, err error, ref map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	t.Log("Asserting tree:\n", spew.Sdump(tree))
	assertSubTree(t, []string***REMOVED******REMOVED***, tree, err, ref)
	t.Log("Finished tree assertion.")
***REMOVED***

func TestCreateSubTree(t *testing.T) ***REMOVED***
	tree := newTree()
	tree.createSubTree([]string***REMOVED***"a", "b", "c"***REMOVED***, Position***REMOVED******REMOVED***)
	tree.Set("a.b.c", 42)
	if tree.Get("a.b.c") != 42 ***REMOVED***
		t.Fail()
	***REMOVED***
***REMOVED***

func TestSimpleKV(t *testing.T) ***REMOVED***
	tree, err := Load("a = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(42),
	***REMOVED***)

	tree, _ = Load("a = 42\nb = 21")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(42),
		"b": int64(21),
	***REMOVED***)
***REMOVED***

func TestNumberInKey(t *testing.T) ***REMOVED***
	tree, err := Load("hello2 = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"hello2": int64(42),
	***REMOVED***)
***REMOVED***

func TestIncorrectKeyExtraSquareBracket(t *testing.T) ***REMOVED***
	_, err := Load(`[a]b]
zyx = 42`)
	if err == nil ***REMOVED***
		t.Error("Error should have been returned.")
	***REMOVED***
	if err.Error() != "(1, 4): unexpected token" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestSimpleNumbers(t *testing.T) ***REMOVED***
	tree, err := Load("a = +42\nb = -21\nc = +4.2\nd = -2.1")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(42),
		"b": int64(-21),
		"c": float64(4.2),
		"d": float64(-2.1),
	***REMOVED***)
***REMOVED***

func TestSpecialFloats(t *testing.T) ***REMOVED***
	tree, err := Load(`
normalinf = inf
plusinf = +inf
minusinf = -inf
normalnan = nan
plusnan = +nan
minusnan = -nan
`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"normalinf": math.Inf(1),
		"plusinf":   math.Inf(1),
		"minusinf":  math.Inf(-1),
		"normalnan": math.NaN(),
		"plusnan":   math.NaN(),
		"minusnan":  math.NaN(),
	***REMOVED***)
***REMOVED***

func TestHexIntegers(t *testing.T) ***REMOVED***
	tree, err := Load(`a = 0xDEADBEEF`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(3735928559)***REMOVED***)

	tree, err = Load(`a = 0xdeadbeef`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(3735928559)***REMOVED***)

	tree, err = Load(`a = 0xdead_beef`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(3735928559)***REMOVED***)

	_, err = Load(`a = 0x_1`)
	if err.Error() != "(1, 5): invalid use of _ in hex number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestOctIntegers(t *testing.T) ***REMOVED***
	tree, err := Load(`a = 0o01234567`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(342391)***REMOVED***)

	tree, err = Load(`a = 0o755`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(493)***REMOVED***)

	_, err = Load(`a = 0o_1`)
	if err.Error() != "(1, 5): invalid use of _ in number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestBinIntegers(t *testing.T) ***REMOVED***
	tree, err := Load(`a = 0b11010110`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***"a": int64(214)***REMOVED***)

	_, err = Load(`a = 0b_1`)
	if err.Error() != "(1, 5): invalid use of _ in number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestBadIntegerBase(t *testing.T) ***REMOVED***
	_, err := Load(`a = 0k1`)
	if err.Error() != "(1, 5): unknown number base: k. possible options are x (hex) o (octal) b (binary)" ***REMOVED***
		t.Error("Error should have been returned.")
	***REMOVED***
***REMOVED***

func TestIntegerNoDigit(t *testing.T) ***REMOVED***
	_, err := Load(`a = 0b`)
	if err.Error() != "(1, 5): number needs at least one digit" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestNumbersWithUnderscores(t *testing.T) ***REMOVED***
	tree, err := Load("a = 1_000")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(1000),
	***REMOVED***)

	tree, err = Load("a = 5_349_221")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(5349221),
	***REMOVED***)

	tree, err = Load("a = 1_2_3_4_5")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(12345),
	***REMOVED***)

	tree, err = Load("flt8 = 9_224_617.445_991_228_313")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"flt8": float64(9224617.445991228313),
	***REMOVED***)

	tree, err = Load("flt9 = 1e1_00")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"flt9": float64(1e100),
	***REMOVED***)
***REMOVED***

func TestFloatsWithExponents(t *testing.T) ***REMOVED***
	tree, err := Load("a = 5e+22\nb = 5E+22\nc = -5e+22\nd = -5e-22\ne = 6.626e-34")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": float64(5e+22),
		"b": float64(5E+22),
		"c": float64(-5e+22),
		"d": float64(-5e-22),
		"e": float64(6.626e-34),
	***REMOVED***)
***REMOVED***

func TestSimpleDate(t *testing.T) ***REMOVED***
	tree, err := Load("a = 1979-05-27T07:32:00Z")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
	***REMOVED***)
***REMOVED***

func TestDateOffset(t *testing.T) ***REMOVED***
	tree, err := Load("a = 1979-05-27T00:32:00-07:00")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": time.Date(1979, time.May, 27, 0, 32, 0, 0, time.FixedZone("", -7*60*60)),
	***REMOVED***)
***REMOVED***

func TestDateNano(t *testing.T) ***REMOVED***
	tree, err := Load("a = 1979-05-27T00:32:00.999999999-07:00")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": time.Date(1979, time.May, 27, 0, 32, 0, 999999999, time.FixedZone("", -7*60*60)),
	***REMOVED***)
***REMOVED***

func TestSimpleString(t *testing.T) ***REMOVED***
	tree, err := Load("a = \"hello world\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": "hello world",
	***REMOVED***)
***REMOVED***

func TestSpaceKey(t *testing.T) ***REMOVED***
	tree, err := Load("\"a b\" = \"hello world\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a b": "hello world",
	***REMOVED***)
***REMOVED***

func TestDoubleQuotedKey(t *testing.T) ***REMOVED***
	tree, err := Load(`
	"key"        = "a"
	"\t"         = "b"
	"\U0001F914" = "c"
	"\u2764"     = "d"
	`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"key":        "a",
		"\t":         "b",
		"\U0001F914": "c",
		"\u2764":     "d",
	***REMOVED***)
***REMOVED***

func TestSingleQuotedKey(t *testing.T) ***REMOVED***
	tree, err := Load(`
	'key'        = "a"
	'\t'         = "b"
	'\U0001F914' = "c"
	'\u2764'     = "d"
	`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		`key`:        "a",
		`\t`:         "b",
		`\U0001F914`: "c",
		`\u2764`:     "d",
	***REMOVED***)
***REMOVED***

func TestStringEscapables(t *testing.T) ***REMOVED***
	tree, err := Load("a = \"a \\n b\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": "a \n b",
	***REMOVED***)

	tree, err = Load("a = \"a \\t b\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": "a \t b",
	***REMOVED***)

	tree, err = Load("a = \"a \\r b\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": "a \r b",
	***REMOVED***)

	tree, err = Load("a = \"a \\\\ b\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": "a \\ b",
	***REMOVED***)
***REMOVED***

func TestEmptyQuotedString(t *testing.T) ***REMOVED***
	tree, err := Load(`[""]
"" = 1`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"": map[string]interface***REMOVED******REMOVED******REMOVED***
			"": int64(1),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestBools(t *testing.T) ***REMOVED***
	tree, err := Load("a = true\nb = false")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": true,
		"b": false,
	***REMOVED***)
***REMOVED***

func TestNestedKeys(t *testing.T) ***REMOVED***
	tree, err := Load("[a.b.c]\nd = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": map[string]interface***REMOVED******REMOVED******REMOVED***
			"b": map[string]interface***REMOVED******REMOVED******REMOVED***
				"c": map[string]interface***REMOVED******REMOVED******REMOVED***
					"d": int64(42),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedQuotedUnicodeKeys(t *testing.T) ***REMOVED***
	tree, err := Load("[ j . \"ʞ\" . l ]\nd = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"j": map[string]interface***REMOVED******REMOVED******REMOVED***
			"ʞ": map[string]interface***REMOVED******REMOVED******REMOVED***
				"l": map[string]interface***REMOVED******REMOVED******REMOVED***
					"d": int64(42),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)

	tree, err = Load("[ g . h . i ]\nd = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"g": map[string]interface***REMOVED******REMOVED******REMOVED***
			"h": map[string]interface***REMOVED******REMOVED******REMOVED***
				"i": map[string]interface***REMOVED******REMOVED******REMOVED***
					"d": int64(42),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)

	tree, err = Load("[ d.e.f ]\nk = 42")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"d": map[string]interface***REMOVED******REMOVED******REMOVED***
			"e": map[string]interface***REMOVED******REMOVED******REMOVED***
				"f": map[string]interface***REMOVED******REMOVED******REMOVED***
					"k": int64(42),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayOne(t *testing.T) ***REMOVED***
	tree, err := Load("a = [1]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(1)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayZero(t *testing.T) ***REMOVED***
	tree, err := Load("a = []")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArraySimple(t *testing.T) ***REMOVED***
	tree, err := Load("a = [42, 21, 10]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(42), int64(21), int64(10)***REMOVED***,
	***REMOVED***)

	tree, _ = Load("a = [42, 21, 10,]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(42), int64(21), int64(10)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayMultiline(t *testing.T) ***REMOVED***
	tree, err := Load("a = [42,\n21, 10,]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(42), int64(21), int64(10)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayNested(t *testing.T) ***REMOVED***
	tree, err := Load("a = [[42, 21], [10]]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": [][]int64***REMOVED******REMOVED***int64(42), int64(21)***REMOVED***, ***REMOVED***int64(10)***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedArrayComment(t *testing.T) ***REMOVED***
	tree, err := Load(`
someArray = [
# does not work
["entry1"]
]`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"someArray": [][]string***REMOVED******REMOVED***"entry1"***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

func TestNestedEmptyArrays(t *testing.T) ***REMOVED***
	tree, err := Load("a = [[[]]]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": [][][]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayMixedTypes(t *testing.T) ***REMOVED***
	_, err := Load("a = [42, 16.0]")
	if err.Error() != "(1, 10): mixed types in array" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a = [42, \"hello\"]")
	if err.Error() != "(1, 11): mixed types in array" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestArrayNestedStrings(t *testing.T) ***REMOVED***
	tree, err := Load("data = [ [\"gamma\", \"delta\"], [\"Foo\"] ]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"data": [][]string***REMOVED******REMOVED***"gamma", "delta"***REMOVED***, ***REMOVED***"Foo"***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

func TestParseUnknownRvalue(t *testing.T) ***REMOVED***
	_, err := Load("a = !bssss")
	if err == nil ***REMOVED***
		t.Error("Expecting a parse error")
	***REMOVED***

	_, err = Load("a = /b")
	if err == nil ***REMOVED***
		t.Error("Expecting a parse error")
	***REMOVED***
***REMOVED***

func TestMissingValue(t *testing.T) ***REMOVED***
	_, err := Load("a = ")
	if err.Error() != "(1, 5): expecting a value" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestUnterminatedArray(t *testing.T) ***REMOVED***
	_, err := Load("a = [1,")
	if err.Error() != "(1, 8): unterminated array" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a = [1")
	if err.Error() != "(1, 7): unterminated array" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a = [1 2")
	if err.Error() != "(1, 8): missing comma" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestNewlinesInArrays(t *testing.T) ***REMOVED***
	tree, err := Load("a = [1,\n2,\n3]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(1), int64(2), int64(3)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayWithExtraComma(t *testing.T) ***REMOVED***
	tree, err := Load("a = [1,\n2,\n3,\n]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(1), int64(2), int64(3)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestArrayWithExtraCommaComment(t *testing.T) ***REMOVED***
	tree, err := Load("a = [1, # wow\n2, # such items\n3, # so array\n]")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": []int64***REMOVED***int64(1), int64(2), int64(3)***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestSimpleInlineGroup(t *testing.T) ***REMOVED***
	tree, err := Load("key = ***REMOVED***a = 42***REMOVED***")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"key": map[string]interface***REMOVED******REMOVED******REMOVED***
			"a": int64(42),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestDoubleInlineGroup(t *testing.T) ***REMOVED***
	tree, err := Load("key = ***REMOVED***a = 42, b = \"foo\"***REMOVED***")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"key": map[string]interface***REMOVED******REMOVED******REMOVED***
			"a": int64(42),
			"b": "foo",
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestExampleInlineGroup(t *testing.T) ***REMOVED***
	tree, err := Load(`name = ***REMOVED*** first = "Tom", last = "Preston-Werner" ***REMOVED***
point = ***REMOVED*** x = 1, y = 2 ***REMOVED***`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"name": map[string]interface***REMOVED******REMOVED******REMOVED***
			"first": "Tom",
			"last":  "Preston-Werner",
		***REMOVED***,
		"point": map[string]interface***REMOVED******REMOVED******REMOVED***
			"x": int64(1),
			"y": int64(2),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestExampleInlineGroupInArray(t *testing.T) ***REMOVED***
	tree, err := Load(`points = [***REMOVED*** x = 1, y = 2 ***REMOVED***]`)
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"points": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***
				"x": int64(1),
				"y": int64(2),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestInlineTableUnterminated(t *testing.T) ***REMOVED***
	_, err := Load("foo = ***REMOVED***")
	if err.Error() != "(1, 8): unterminated inline table" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestInlineTableCommaExpected(t *testing.T) ***REMOVED***
	_, err := Load("foo = ***REMOVED***hello = 53 test = foo***REMOVED***")
	if err.Error() != "(1, 19): comma expected between fields in inline table" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestInlineTableCommaStart(t *testing.T) ***REMOVED***
	_, err := Load("foo = ***REMOVED***, hello = 53***REMOVED***")
	if err.Error() != "(1, 8): inline table cannot start with a comma" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestInlineTableDoubleComma(t *testing.T) ***REMOVED***
	_, err := Load("foo = ***REMOVED***hello = 53,, foo = 17***REMOVED***")
	if err.Error() != "(1, 19): need field between two commas in inline table" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestDuplicateGroups(t *testing.T) ***REMOVED***
	_, err := Load("[foo]\na=2\n[foo]b=3")
	if err.Error() != "(3, 2): duplicated tables" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestDuplicateKeys(t *testing.T) ***REMOVED***
	_, err := Load("foo = 2\nfoo = 3")
	if err.Error() != "(2, 1): The following key was defined twice: foo" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestEmptyIntermediateTable(t *testing.T) ***REMOVED***
	_, err := Load("[foo..bar]")
	if err.Error() != "(1, 2): invalid table array key: empty table key" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestImplicitDeclarationBefore(t *testing.T) ***REMOVED***
	tree, err := Load("[a.b.c]\nanswer = 42\n[a]\nbetter = 43")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": map[string]interface***REMOVED******REMOVED******REMOVED***
			"b": map[string]interface***REMOVED******REMOVED******REMOVED***
				"c": map[string]interface***REMOVED******REMOVED******REMOVED***
					"answer": int64(42),
				***REMOVED***,
			***REMOVED***,
			"better": int64(43),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestFloatsWithoutLeadingZeros(t *testing.T) ***REMOVED***
	_, err := Load("a = .42")
	if err.Error() != "(1, 5): cannot start float with a dot" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a = -.42")
	if err.Error() != "(1, 5): cannot start float with a dot" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestMissingFile(t *testing.T) ***REMOVED***
	_, err := LoadFile("foo.toml")
	if err.Error() != "open foo.toml: no such file or directory" &&
		err.Error() != "open foo.toml: The system cannot find the file specified." ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestParseFile(t *testing.T) ***REMOVED***
	tree, err := LoadFile("example.toml")

	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"title": "TOML Example",
		"owner": map[string]interface***REMOVED******REMOVED******REMOVED***
			"name":         "Tom Preston-Werner",
			"organization": "GitHub",
			"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
			"dob":          time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
		***REMOVED***,
		"database": map[string]interface***REMOVED******REMOVED******REMOVED***
			"server":         "192.168.1.1",
			"ports":          []int64***REMOVED***8001, 8001, 8002***REMOVED***,
			"connection_max": 5000,
			"enabled":        true,
		***REMOVED***,
		"servers": map[string]interface***REMOVED******REMOVED******REMOVED***
			"alpha": map[string]interface***REMOVED******REMOVED******REMOVED***
				"ip": "10.0.0.1",
				"dc": "eqdc10",
			***REMOVED***,
			"beta": map[string]interface***REMOVED******REMOVED******REMOVED***
				"ip": "10.0.0.2",
				"dc": "eqdc10",
			***REMOVED***,
		***REMOVED***,
		"clients": map[string]interface***REMOVED******REMOVED******REMOVED***
			"data": []interface***REMOVED******REMOVED******REMOVED***
				[]string***REMOVED***"gamma", "delta"***REMOVED***,
				[]int64***REMOVED***1, 2***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestParseFileCRLF(t *testing.T) ***REMOVED***
	tree, err := LoadFile("example-crlf.toml")

	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"title": "TOML Example",
		"owner": map[string]interface***REMOVED******REMOVED******REMOVED***
			"name":         "Tom Preston-Werner",
			"organization": "GitHub",
			"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
			"dob":          time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
		***REMOVED***,
		"database": map[string]interface***REMOVED******REMOVED******REMOVED***
			"server":         "192.168.1.1",
			"ports":          []int64***REMOVED***8001, 8001, 8002***REMOVED***,
			"connection_max": 5000,
			"enabled":        true,
		***REMOVED***,
		"servers": map[string]interface***REMOVED******REMOVED******REMOVED***
			"alpha": map[string]interface***REMOVED******REMOVED******REMOVED***
				"ip": "10.0.0.1",
				"dc": "eqdc10",
			***REMOVED***,
			"beta": map[string]interface***REMOVED******REMOVED******REMOVED***
				"ip": "10.0.0.2",
				"dc": "eqdc10",
			***REMOVED***,
		***REMOVED***,
		"clients": map[string]interface***REMOVED******REMOVED******REMOVED***
			"data": []interface***REMOVED******REMOVED******REMOVED***
				[]string***REMOVED***"gamma", "delta"***REMOVED***,
				[]int64***REMOVED***1, 2***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestParseKeyGroupArray(t *testing.T) ***REMOVED***
	tree, err := Load("[[foo.bar]] a = 42\n[[foo.bar]] a = 69")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
			"bar": []map[string]interface***REMOVED******REMOVED******REMOVED***
				***REMOVED***"a": int64(42)***REMOVED***,
				***REMOVED***"a": int64(69)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestParseKeyGroupArrayUnfinished(t *testing.T) ***REMOVED***
	_, err := Load("[[foo.bar]\na = 42")
	if err.Error() != "(1, 10): was expecting token [[, but got unclosed table array key instead" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("[[foo.[bar]\na = 42")
	if err.Error() != "(1, 3): unexpected token table array key cannot contain ']', was expecting a table array key" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestParseKeyGroupArrayQueryExample(t *testing.T) ***REMOVED***
	tree, err := Load(`
      [[book]]
      title = "The Stand"
      author = "Stephen King"
      [[book]]
      title = "For Whom the Bell Tolls"
      author = "Ernest Hemmingway"
      [[book]]
      title = "Neuromancer"
      author = "William Gibson"
    `)

	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"book": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"title": "The Stand", "author": "Stephen King"***REMOVED***,
			***REMOVED***"title": "For Whom the Bell Tolls", "author": "Ernest Hemmingway"***REMOVED***,
			***REMOVED***"title": "Neuromancer", "author": "William Gibson"***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestParseKeyGroupArraySpec(t *testing.T) ***REMOVED***
	tree, err := Load("[[fruit]]\n name=\"apple\"\n [fruit.physical]\n color=\"red\"\n shape=\"round\"\n [[fruit]]\n name=\"banana\"")
	assertTree(t, tree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"fruit": []map[string]interface***REMOVED******REMOVED******REMOVED***
			***REMOVED***"name": "apple", "physical": map[string]interface***REMOVED******REMOVED******REMOVED***"color": "red", "shape": "round"***REMOVED******REMOVED***,
			***REMOVED***"name": "banana"***REMOVED***,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestTomlValueStringRepresentation(t *testing.T) ***REMOVED***
	for idx, item := range []struct ***REMOVED***
		Value  interface***REMOVED******REMOVED***
		Expect string
	***REMOVED******REMOVED***
		***REMOVED***int64(12345), "12345"***REMOVED***,
		***REMOVED***uint64(50), "50"***REMOVED***,
		***REMOVED***float64(123.45), "123.45"***REMOVED***,
		***REMOVED***true, "true"***REMOVED***,
		***REMOVED***"hello world", "\"hello world\""***REMOVED***,
		***REMOVED***"\b\t\n\f\r\"\\", "\"\\b\\t\\n\\f\\r\\\"\\\\\""***REMOVED***,
		***REMOVED***"\x05", "\"\\u0005\""***REMOVED***,
		***REMOVED***time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
			"1979-05-27T07:32:00Z"***REMOVED***,
		***REMOVED***[]interface***REMOVED******REMOVED******REMOVED***"gamma", "delta"***REMOVED***,
			"[\"gamma\",\"delta\"]"***REMOVED***,
		***REMOVED***nil, ""***REMOVED***,
	***REMOVED*** ***REMOVED***
		result, err := tomlValueStringRepresentation(item.Value, "", false)
		if err != nil ***REMOVED***
			t.Errorf("Test %d - unexpected error: %s", idx, err)
		***REMOVED***
		if result != item.Expect ***REMOVED***
			t.Errorf("Test %d - got '%s', expected '%s'", idx, result, item.Expect)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestToStringMapStringString(t *testing.T) ***REMOVED***
	tree, err := TreeFromMap(map[string]interface***REMOVED******REMOVED******REMOVED***"m": map[string]interface***REMOVED******REMOVED******REMOVED***"v": "abc"***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("unexpected error: %s", err)
	***REMOVED***
	want := "\n[m]\n  v = \"abc\"\n"
	got := tree.String()

	if got != want ***REMOVED***
		t.Errorf("want:\n%q\ngot:\n%q", want, got)
	***REMOVED***
***REMOVED***

func assertPosition(t *testing.T, text string, ref map[string]Position) ***REMOVED***
	tree, err := Load(text)
	if err != nil ***REMOVED***
		t.Errorf("Error loading document text: `%v`", text)
		t.Errorf("Error: %v", err)
	***REMOVED***
	for path, pos := range ref ***REMOVED***
		testPos := tree.GetPosition(path)
		if testPos.Invalid() ***REMOVED***
			t.Errorf("Failed to query tree path or path has invalid position: %s", path)
		***REMOVED*** else if pos != testPos ***REMOVED***
			t.Errorf("Expected position %v, got %v instead", pos, testPos)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDocumentPositions(t *testing.T) ***REMOVED***
	assertPosition(t,
		"[foo]\nbar=42\nbaz=69",
		map[string]Position***REMOVED***
			"":        ***REMOVED***1, 1***REMOVED***,
			"foo":     ***REMOVED***1, 1***REMOVED***,
			"foo.bar": ***REMOVED***2, 1***REMOVED***,
			"foo.baz": ***REMOVED***3, 1***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestDocumentPositionsWithSpaces(t *testing.T) ***REMOVED***
	assertPosition(t,
		"  [foo]\n  bar=42\n  baz=69",
		map[string]Position***REMOVED***
			"":        ***REMOVED***1, 1***REMOVED***,
			"foo":     ***REMOVED***1, 3***REMOVED***,
			"foo.bar": ***REMOVED***2, 3***REMOVED***,
			"foo.baz": ***REMOVED***3, 3***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestDocumentPositionsWithGroupArray(t *testing.T) ***REMOVED***
	assertPosition(t,
		"[[foo]]\nbar=42\nbaz=69",
		map[string]Position***REMOVED***
			"":        ***REMOVED***1, 1***REMOVED***,
			"foo":     ***REMOVED***1, 1***REMOVED***,
			"foo.bar": ***REMOVED***2, 1***REMOVED***,
			"foo.baz": ***REMOVED***3, 1***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestNestedTreePosition(t *testing.T) ***REMOVED***
	assertPosition(t,
		"[foo.bar]\na=42\nb=69",
		map[string]Position***REMOVED***
			"":          ***REMOVED***1, 1***REMOVED***,
			"foo":       ***REMOVED***1, 1***REMOVED***,
			"foo.bar":   ***REMOVED***1, 1***REMOVED***,
			"foo.bar.a": ***REMOVED***2, 1***REMOVED***,
			"foo.bar.b": ***REMOVED***3, 1***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestInvalidGroupArray(t *testing.T) ***REMOVED***
	_, err := Load("[table#key]\nanswer = 42")
	if err == nil ***REMOVED***
		t.Error("Should error")
	***REMOVED***

	_, err = Load("[foo.[bar]\na = 42")
	if err.Error() != "(1, 2): unexpected token table key cannot contain ']', was expecting a table key" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestDoubleEqual(t *testing.T) ***REMOVED***
	_, err := Load("foo= = 2")
	if err.Error() != "(1, 6): cannot have multiple equals for the same key" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestGroupArrayReassign(t *testing.T) ***REMOVED***
	_, err := Load("[hello]\n[[hello]]")
	if err.Error() != "(2, 3): key \"hello\" is already assigned and not of type table array" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***

func TestInvalidFloatParsing(t *testing.T) ***REMOVED***
	_, err := Load("a=1e_2")
	if err.Error() != "(1, 3): invalid use of _ in number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a=1e2_")
	if err.Error() != "(1, 3): invalid use of _ in number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a=1__2")
	if err.Error() != "(1, 3): invalid use of _ in number" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***

	_, err = Load("a=_1_2")
	if err.Error() != "(1, 3): cannot start number with underscore" ***REMOVED***
		t.Error("Bad error message:", err.Error())
	***REMOVED***
***REMOVED***
