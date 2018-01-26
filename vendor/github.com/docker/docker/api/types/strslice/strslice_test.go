package strslice

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStrSliceMarshalJSON(t *testing.T) ***REMOVED***
	for _, testcase := range []struct ***REMOVED***
		input    StrSlice
		expected string
	***REMOVED******REMOVED***
		// MADNESS(stevvooe): No clue why nil would be "" but empty would be
		// "null". Had to make a change here that may affect compatibility.
		***REMOVED***input: nil, expected: "null"***REMOVED***,
		***REMOVED***StrSlice***REMOVED******REMOVED***, "[]"***REMOVED***,
		***REMOVED***StrSlice***REMOVED***"/bin/sh", "-c", "echo"***REMOVED***, `["/bin/sh","-c","echo"]`***REMOVED***,
	***REMOVED*** ***REMOVED***
		data, err := json.Marshal(testcase.input)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if string(data) != testcase.expected ***REMOVED***
			t.Fatalf("%#v: expected %v, got %v", testcase.input, testcase.expected, string(data))
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestStrSliceUnmarshalJSON(t *testing.T) ***REMOVED***
	parts := map[string][]string***REMOVED***
		"":   ***REMOVED***"default", "values"***REMOVED***,
		"[]": ***REMOVED******REMOVED***,
		`["/bin/sh","-c","echo"]`: ***REMOVED***"/bin/sh", "-c", "echo"***REMOVED***,
	***REMOVED***
	for json, expectedParts := range parts ***REMOVED***
		strs := StrSlice***REMOVED***"default", "values"***REMOVED***
		if err := strs.UnmarshalJSON([]byte(json)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		actualParts := []string(strs)
		if !reflect.DeepEqual(actualParts, expectedParts) ***REMOVED***
			t.Fatalf("%#v: expected %v, got %v", json, expectedParts, actualParts)
		***REMOVED***

	***REMOVED***
***REMOVED***

func TestStrSliceUnmarshalString(t *testing.T) ***REMOVED***
	var e StrSlice
	echo, err := json.Marshal("echo")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := json.Unmarshal(echo, &e); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(e) != 1 ***REMOVED***
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	***REMOVED***

	if e[0] != "echo" ***REMOVED***
		t.Fatalf("expected `echo`, got: %q", e[0])
	***REMOVED***
***REMOVED***

func TestStrSliceUnmarshalSlice(t *testing.T) ***REMOVED***
	var e StrSlice
	echo, err := json.Marshal([]string***REMOVED***"echo"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := json.Unmarshal(echo, &e); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(e) != 1 ***REMOVED***
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	***REMOVED***

	if e[0] != "echo" ***REMOVED***
		t.Fatalf("expected `echo`, got: %q", e[0])
	***REMOVED***
***REMOVED***
