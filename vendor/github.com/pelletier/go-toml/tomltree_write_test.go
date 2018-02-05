package toml

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type failingWriter struct ***REMOVED***
	failAt  int
	written int
	buffer  bytes.Buffer
***REMOVED***

func (f *failingWriter) Write(p []byte) (n int, err error) ***REMOVED***
	count := len(p)
	toWrite := f.failAt - (count + f.written)
	if toWrite < 0 ***REMOVED***
		toWrite = 0
	***REMOVED***
	if toWrite > count ***REMOVED***
		f.written += count
		f.buffer.Write(p)
		return count, nil
	***REMOVED***

	f.buffer.Write(p[:toWrite])
	f.written = f.failAt
	return toWrite, fmt.Errorf("failingWriter failed after writing %d bytes", f.written)
***REMOVED***

func assertErrorString(t *testing.T, expected string, err error) ***REMOVED***
	expectedErr := errors.New(expected)
	if err == nil || err.Error() != expectedErr.Error() ***REMOVED***
		t.Errorf("expecting error %s, but got %s instead", expected, err)
	***REMOVED***
***REMOVED***

func TestTreeWriteToEmptyTable(t *testing.T) ***REMOVED***
	doc := `[[empty-tables]]
[[empty-tables]]`

	toml, err := Load(doc)
	if err != nil ***REMOVED***
		t.Fatal("Unexpected Load error:", err)
	***REMOVED***
	tomlString, err := toml.ToTomlString()
	if err != nil ***REMOVED***
		t.Fatal("Unexpected ToTomlString error:", err)
	***REMOVED***

	expected := `
[[empty-tables]]

[[empty-tables]]
`

	if tomlString != expected ***REMOVED***
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, tomlString)
	***REMOVED***
***REMOVED***

func TestTreeWriteToTomlString(t *testing.T) ***REMOVED***
	toml, err := Load(`name = ***REMOVED*** first = "Tom", last = "Preston-Werner" ***REMOVED***
points = ***REMOVED*** x = 1, y = 2 ***REMOVED***`)

	if err != nil ***REMOVED***
		t.Fatal("Unexpected error:", err)
	***REMOVED***

	tomlString, _ := toml.ToTomlString()
	reparsedTree, err := Load(tomlString)

	assertTree(t, reparsedTree, err, map[string]interface***REMOVED******REMOVED******REMOVED***
		"name": map[string]interface***REMOVED******REMOVED******REMOVED***
			"first": "Tom",
			"last":  "Preston-Werner",
		***REMOVED***,
		"points": map[string]interface***REMOVED******REMOVED******REMOVED***
			"x": int64(1),
			"y": int64(2),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestTreeWriteToTomlStringSimple(t *testing.T) ***REMOVED***
	tree, err := Load("[foo]\n\n[[foo.bar]]\na = 42\n\n[[foo.bar]]\na = 69\n")
	if err != nil ***REMOVED***
		t.Errorf("Test failed to parse: %v", err)
		return
	***REMOVED***
	result, err := tree.ToTomlString()
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %s", err)
	***REMOVED***
	expected := "\n[foo]\n\n  [[foo.bar]]\n    a = 42\n\n  [[foo.bar]]\n    a = 69\n"
	if result != expected ***REMOVED***
		t.Errorf("Expected got '%s', expected '%s'", result, expected)
	***REMOVED***
***REMOVED***

func TestTreeWriteToTomlStringKeysOrders(t *testing.T) ***REMOVED***
	for i := 0; i < 100; i++ ***REMOVED***
		tree, _ := Load(`
		foobar = true
		bar = "baz"
		foo = 1
		[qux]
		  foo = 1
		  bar = "baz2"`)

		stringRepr, _ := tree.ToTomlString()

		t.Log("Intermediate string representation:")
		t.Log(stringRepr)

		r := strings.NewReader(stringRepr)
		toml, err := LoadReader(r)

		if err != nil ***REMOVED***
			t.Fatal("Unexpected error:", err)
		***REMOVED***

		assertTree(t, toml, err, map[string]interface***REMOVED******REMOVED******REMOVED***
			"foobar": true,
			"bar":    "baz",
			"foo":    1,
			"qux": map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": 1,
				"bar": "baz2",
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func testMaps(t *testing.T, actual, expected map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	if !reflect.DeepEqual(actual, expected) ***REMOVED***
		t.Fatal("trees aren't equal.\n", "Expected:\n", expected, "\nActual:\n", actual)
	***REMOVED***
***REMOVED***

func TestTreeWriteToMapSimple(t *testing.T) ***REMOVED***
	tree, _ := Load("a = 42\nb = 17")

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"a": int64(42),
		"b": int64(17),
	***REMOVED***

	testMaps(t, tree.ToMap(), expected)
***REMOVED***

func TestTreeWriteToInvalidTreeSimpleValue(t *testing.T) ***REMOVED***
	tree := Tree***REMOVED***values: map[string]interface***REMOVED******REMOVED******REMOVED***"foo": int8(1)***REMOVED******REMOVED***
	_, err := tree.ToTomlString()
	assertErrorString(t, "invalid value type at foo: int8", err)
***REMOVED***

func TestTreeWriteToInvalidTreeTomlValue(t *testing.T) ***REMOVED***
	tree := Tree***REMOVED***values: map[string]interface***REMOVED******REMOVED******REMOVED***"foo": &tomlValue***REMOVED***value: int8(1), comment: "", position: Position***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***
	_, err := tree.ToTomlString()
	assertErrorString(t, "unsupported value type int8: 1", err)
***REMOVED***

func TestTreeWriteToInvalidTreeTomlValueArray(t *testing.T) ***REMOVED***
	tree := Tree***REMOVED***values: map[string]interface***REMOVED******REMOVED******REMOVED***"foo": &tomlValue***REMOVED***value: int8(1), comment: "", position: Position***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***
	_, err := tree.ToTomlString()
	assertErrorString(t, "unsupported value type int8: 1", err)
***REMOVED***

func TestTreeWriteToFailingWriterInSimpleValue(t *testing.T) ***REMOVED***
	toml, _ := Load(`a = 2`)
	writer := failingWriter***REMOVED***failAt: 0, written: 0***REMOVED***
	_, err := toml.WriteTo(&writer)
	assertErrorString(t, "failingWriter failed after writing 0 bytes", err)
***REMOVED***

func TestTreeWriteToFailingWriterInTable(t *testing.T) ***REMOVED***
	toml, _ := Load(`
[b]
a = 2`)
	writer := failingWriter***REMOVED***failAt: 2, written: 0***REMOVED***
	_, err := toml.WriteTo(&writer)
	assertErrorString(t, "failingWriter failed after writing 2 bytes", err)

	writer = failingWriter***REMOVED***failAt: 13, written: 0***REMOVED***
	_, err = toml.WriteTo(&writer)
	assertErrorString(t, "failingWriter failed after writing 13 bytes", err)
***REMOVED***

func TestTreeWriteToFailingWriterInArray(t *testing.T) ***REMOVED***
	toml, _ := Load(`
[[b]]
a = 2`)
	writer := failingWriter***REMOVED***failAt: 2, written: 0***REMOVED***
	_, err := toml.WriteTo(&writer)
	assertErrorString(t, "failingWriter failed after writing 2 bytes", err)

	writer = failingWriter***REMOVED***failAt: 15, written: 0***REMOVED***
	_, err = toml.WriteTo(&writer)
	assertErrorString(t, "failingWriter failed after writing 15 bytes", err)
***REMOVED***

func TestTreeWriteToMapExampleFile(t *testing.T) ***REMOVED***
	tree, _ := LoadFile("example.toml")
	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"title": "TOML Example",
		"owner": map[string]interface***REMOVED******REMOVED******REMOVED***
			"name":         "Tom Preston-Werner",
			"organization": "GitHub",
			"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
			"dob":          time.Date(1979, time.May, 27, 7, 32, 0, 0, time.UTC),
		***REMOVED***,
		"database": map[string]interface***REMOVED******REMOVED******REMOVED***
			"server":         "192.168.1.1",
			"ports":          []interface***REMOVED******REMOVED******REMOVED***int64(8001), int64(8001), int64(8002)***REMOVED***,
			"connection_max": int64(5000),
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
				[]interface***REMOVED******REMOVED******REMOVED***"gamma", "delta"***REMOVED***,
				[]interface***REMOVED******REMOVED******REMOVED***int64(1), int64(2)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	testMaps(t, tree.ToMap(), expected)
***REMOVED***

func TestTreeWriteToMapWithTablesInMultipleChunks(t *testing.T) ***REMOVED***
	tree, _ := Load(`
	[[menu.main]]
        a = "menu 1"
        b = "menu 2"
        [[menu.main]]
        c = "menu 3"
        d = "menu 4"`)
	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"menu": map[string]interface***REMOVED******REMOVED******REMOVED***
			"main": []interface***REMOVED******REMOVED******REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***"a": "menu 1", "b": "menu 2"***REMOVED***,
				map[string]interface***REMOVED******REMOVED******REMOVED***"c": "menu 3", "d": "menu 4"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	treeMap := tree.ToMap()

	testMaps(t, treeMap, expected)
***REMOVED***

func TestTreeWriteToMapWithArrayOfInlineTables(t *testing.T) ***REMOVED***
	tree, _ := Load(`
    	[params]
	language_tabs = [
    		***REMOVED*** key = "shell", name = "Shell" ***REMOVED***,
    		***REMOVED*** key = "ruby", name = "Ruby" ***REMOVED***,
    		***REMOVED*** key = "python", name = "Python" ***REMOVED***
	]`)

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"params": map[string]interface***REMOVED******REMOVED******REMOVED***
			"language_tabs": []interface***REMOVED******REMOVED******REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"key":  "shell",
					"name": "Shell",
				***REMOVED***,
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"key":  "ruby",
					"name": "Ruby",
				***REMOVED***,
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"key":  "python",
					"name": "Python",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	treeMap := tree.ToMap()
	testMaps(t, treeMap, expected)
***REMOVED***

func TestTreeWriteToFloat(t *testing.T) ***REMOVED***
	tree, err := Load(`a = 3.0`)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	str, err := tree.ToTomlString()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := `a = 3.0`
	if strings.TrimSpace(str) != strings.TrimSpace(expected) ***REMOVED***
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, str)
	***REMOVED***
***REMOVED***

func TestTreeWriteToSpecialFloat(t *testing.T) ***REMOVED***
	expected := `a = +inf
b = -inf
c = nan`

	tree, err := Load(expected)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	str, err := tree.ToTomlString()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if strings.TrimSpace(str) != strings.TrimSpace(expected) ***REMOVED***
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, str)
	***REMOVED***
***REMOVED***

func BenchmarkTreeToTomlString(b *testing.B) ***REMOVED***
	toml, err := Load(sampleHard)
	if err != nil ***REMOVED***
		b.Fatal("Unexpected error:", err)
	***REMOVED***

	for i := 0; i < b.N; i++ ***REMOVED***
		_, err := toml.ToTomlString()
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

var sampleHard = `# Test file for TOML
# Only this one tries to emulate a TOML file written by a user of the kind of parser writers probably hate
# This part you'll really hate

[the]
test_string = "You'll hate me after this - #"          # " Annoying, isn't it?

    [the.hard]
    test_array = [ "] ", " # "]      # ] There you go, parse this!
    test_array2 = [ "Test #11 ]proved that", "Experiment #9 was a success" ]
    # You didn't think it'd as easy as chucking out the last #, did you?
    another_test_string = " Same thing, but with a string #"
    harder_test_string = " And when \"'s are in the string, along with # \""   # "and comments are there too"
    # Things will get harder

        [the.hard."bit#"]
        "what?" = "You don't think some user won't do that?"
        multi_line_array = [
            "]",
            # ] Oh yes I did
            ]

# Each of the following keygroups/key value pairs should produce an error. Uncomment to them to test

#[error]   if you didn't catch this, your parser is broken
#string = "Anything other than tabs, spaces and newline after a keygroup or key value pair has ended should produce an error unless it is a comment"   like this
#array = [
#         "This might most likely happen in multiline arrays",
#         Like here,
#         "or here,
#         and here"
#         ]     End of array comment, forgot the #
#number = 3.14  pi <--again forgot the #         `
